package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	repo "github.com/0cd/go-ecom/internal/adapters/sqlc"
	"github.com/0cd/go-ecom/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type app struct {
	config config
	db     *pgx.Conn
}

type config struct {
	address string
	db      dbConfig
}

type dbConfig struct {
	dsn string
}

func (a *app) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	productService := products.NewService(repo.New(a.db))
	productHandler := products.NewHandler(productService)

	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", productHandler.FindProductByID)

	return r
}

func (a *app) start(h http.Handler) error {
	srv := &http.Server{
		Addr:         a.config.address,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute * 1,
	}

	log.Printf("Server is running on port %s", strings.Split(a.config.address, ":")[1])

	return srv.ListenAndServe()
}
