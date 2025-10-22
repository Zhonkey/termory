package database

import (
	"context"
	"fmt"
	"time"
	"trainer/internal/domain/user"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) user.Repository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT 
			u.id, u.role, u.email, u.first_name, u.last_name, u.password, 
			u.created_at, u.updated_at
		FROM users u
		WHERE id = $1
	`

	row := r.db.pool.QueryRow(ctx, query, id)

	return r.scanUser(ctx, row)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT 
			u.id, u.role, u.email, u.first_name, u.last_name, u.password, 
			u.created_at, u.updated_at
		FROM users u
		WHERE email = $1
	`

	row := r.db.pool.QueryRow(ctx, query, email)

	return r.scanUser(ctx, row)
}

func (r *UserRepository) FindByToken(ctx context.Context, token uuid.UUID) (*user.User, error) {
	query := `
		SELECT 
			u.id, u.role, u.email, u.first_name, u.last_name, u.password, 
			u.created_at, u.updated_at
		FROM users u
		LEFT JOIN refresh_tokens ON refresh_tokens.userid = u.id
		WHERE refresh_tokens.id = $1
	`

	row := r.db.pool.QueryRow(ctx, query, token)

	return r.scanUser(ctx, row)
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	_, err := r.trx(ctx, func(tx pgx.Tx) (any, error) {
		query := `
			INSERT INTO users (id, role, email, first_name, last_name, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err := r.db.pool.Exec(ctx, query, u.ID, u.Role, u.Email, u.FirstName, u.LastName, u.Password, u.CreatedAt, u.UpdatedAt)
		if err != nil {
			return nil, err
		}

		for _, t := range u.GetRefreshTokens() {
			query := `
				INSERT INTO refresh_tokens (id, user_id, expires_at, created_at)
				VALUES ($1, $2, $3, $4)
			`
			_, err := r.db.pool.Exec(ctx, query, t.ID, u.ID, t.ExpiresAt, t.CreatedAt)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	_, err := r.trx(ctx, func(tx pgx.Tx) (any, error) {
		query := `
			UPDATE users
			SET role=$2, email=$3, first_name=$4, last_name=$5, password=$6
			WHERE id=$1
		`
		_, err := r.db.pool.Exec(ctx, query, u.ID, u.Role, u.Email, u.FirstName, u.LastName, u.Password)
		if err != nil {
			return nil, err
		}

		existedTokens, err := r.findTokensByUserId(ctx, u.ID)
		if err != nil {
			return nil, err
		}

		mappedExistedTokens := make(map[uuid.UUID]*user.RefreshToken, len(existedTokens))
		for _, t := range existedTokens {
			mappedExistedTokens[t.ID] = t
		}

		for _, t := range u.GetRevokedTokens() {
			_, err = tx.Exec(ctx, `DELETE FROM refresh_tokens WHERE id=$1`, t.ID)
			if err != nil {
				return nil, err
			}
		}

		for _, t := range u.GetRefreshTokens() {
			if mappedExistedTokens[t.ID] == nil {
				query := `
				INSERT INTO refresh_tokens (id, user_id, expires_at, created_at)
				VALUES ($1, $2, $3, $4)
			`
				_, err := r.db.pool.Exec(ctx, query, t.ID, u.ID, t.ExpiresAt, t.CreatedAt)
				if err != nil {
					return nil, err
				}
			}
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.trx(ctx, func(tx pgx.Tx) (any, error) {
		queryTokens := `DELETE FROM refresh_tokens WHERE user_id=$1`
		_, err := tx.Exec(ctx, queryTokens, id)
		if err != nil {
			return nil, err
		}

		queryUser := `DELETE FROM users WHERE id=$1`
		_, err = tx.Exec(ctx, queryUser, id)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) scanUser(ctx context.Context, row pgx.Row) (*user.User, error) {
	var (
		id           uuid.UUID
		role         string
		email        string
		firstName    string
		lastName     string
		passwordHash string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := row.Scan(
		&id, &role, &email, &firstName, &lastName, &passwordHash,
		&createdAt, &updatedAt,
	)

	if err != nil {
		return nil, err
	}

	tokens, err := r.findTokensByUserId(ctx, id)

	if err != nil {
		return nil, err
	}

	return user.NewUserFromStorage(
		id,
		email,
		firstName,
		lastName,
		passwordHash,
		user.Role(role),
		createdAt,
		updatedAt,
		tokens,
	), nil
}

func (r *UserRepository) findTokensByUserId(ctx context.Context, userId uuid.UUID) ([]*user.RefreshToken, error) {
	query := `
		SELECT id, expires_at, created_at
		FROM refresh_tokens
		WHERE user_id=$1
	`

	rows, err := r.db.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tokens := make([]*user.RefreshToken, 0)

	for rows.Next() {
		var (
			tokenID      uuid.UUID
			tokenCreated time.Time
			tokenExpires time.Time
		)

		err := rows.Scan(
			&tokenID, &tokenExpires, &tokenCreated,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		tokens = append(tokens, user.NewRefreshTokenFromStorage(tokenID, tokenExpires, tokenCreated))
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return tokens, nil
}

func (r *UserRepository) trx(ctx context.Context, call func(tx pgx.Tx) (any, error)) (any, error) {
	tx, err := r.db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	res, err := call(tx)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return res, nil
}
