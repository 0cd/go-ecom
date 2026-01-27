package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/0cd/go-ecom/internal/env"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	cfg := config{
		address: ":1337",
		db: dbConfig{
			dsn: env.GetString(
				"GOOSE_DBSTRING",
				"host=localhost user=postgres password=postgres dbname=postgres sslmode=disable",
			),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Connected to database", "dsn", cfg.db.dsn)

	api := app{
		config: cfg,
		db:     conn,
	}

	if err := api.start(api.mount()); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
