// Package database provides database connection pool management and migrations
package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type DB struct {
	pool *pgxpool.Pool
}

type Config struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}
}

func New(ctx context.Context, cfg *Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := db.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) RunMigrationsFromEmbed(ctx context.Context, embedFS embed.FS, dir string) error {
	sqlDB, err := db.toSQLDB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(sqlDB, dir); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

func (db *DB) MigrateDown(ctx context.Context, embedFS embed.FS, dir string) error {
	sqlDB, err := db.toSQLDB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Down(sqlDB, dir); err != nil {
		return fmt.Errorf("rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}

func (db *DB) MigrationStatus(ctx context.Context, embedFS embed.FS, dir string) error {
	sqlDB, err := db.toSQLDB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	return goose.Status(sqlDB, dir)
}

func (db *DB) toSQLDB() (*sql.DB, error) {
	if db.pool == nil {
		return nil, fmt.Errorf("pgx pool is not initialized")
	}

	dsn := db.pool.Config().ConnConfig.ConnString()
	if dsn == "" {
		return nil, fmt.Errorf("empty DSN")
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sql.DB failed: %w", err)
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return sqlDB, nil
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return db.pool.Query(ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return db.pool.QueryRow(ctx, query, args...)
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err := db.pool.Exec(ctx, query, args...)
	return err
}

func (db *DB) Transaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (db *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var result int
	err := db.pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected health check result: %d", result)
	}

	return nil
}

func (db *DB) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}
