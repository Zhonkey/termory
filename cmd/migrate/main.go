package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"trainer/internal/infrastructure/database"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <up|down|status>")
		os.Exit(1)
	}

	cmd := os.Args[1]

	cfg := database.DefaultConfig()
	cfg.DSN = os.Getenv("DB_DSN")

	ctx := context.Background()
	db, err := database.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch cmd {
	case "up":
		if err := db.RunMigrationsFromEmbed(ctx, migrationsFS, "migrations"); err != nil {
			log.Fatal(err)
		}
		log.Println("✅ Migrations applied successfully")
	case "down":
		if err := db.MigrateDown(ctx, migrationsFS, "migrations"); err != nil {
			log.Fatal(err)
		}
		log.Println("✅ Migration rolled back successfully")
	case "status":
		if err := db.MigrationStatus(ctx, migrationsFS, "migrations"); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Unknown command. Usage: migrate <up|down|status>")
		os.Exit(1)
	}
}
