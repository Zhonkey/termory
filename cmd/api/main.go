package main

import (
	"context"
	"log"
	"os"
	"trainer/internal/infrastructure/database"
	"trainer/internal/interfaces/http"
)

func main() {
	ctx := context.Background()

	cfg := database.DefaultConfig()
	cfg.DSN = os.Getenv("DB_DSN")

	db, err := database.New(ctx, cfg)

	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	srv := http.NewServer(db, os.Getenv("CORS_ALLOWED_ORIGINS"))
	if err := srv.Run(ctx); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
