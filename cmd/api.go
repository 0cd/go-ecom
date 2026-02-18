package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	repo "github.com/0cd/go-ecom/internal/adapters/sqlc"
	"github.com/0cd/go-ecom/internal/auth"
	"github.com/0cd/go-ecom/internal/middleware"
	"github.com/0cd/go-ecom/internal/orders"
	"github.com/0cd/go-ecom/internal/products"
	"github.com/0cd/go-ecom/internal/users"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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

	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	userService := users.NewService(repo.New(a.db))
	userHandler := users.NewHandler(userService)

	authService := auth.NewService(userService)
	authHandler := auth.NewHandler(authService)

	productService := products.NewService(repo.New(a.db))
	productHandler := products.NewHandler(productService)

	ordersService := orders.NewService(repo.New(a.db), a.db)
	ordersHandler := orders.NewHandler(ordersService)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/logout", authHandler.Logout)
		r.Post("/refresh", authHandler.Refresh)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Route("/me", func(r chi.Router) {
			r.Get("/", userHandler.GetMe)
			r.Patch("/password", userHandler.UpdateUserPassword)
			r.Patch("/email", userHandler.UpdateUserEmail)
			r.Patch("/verify", userHandler.VerifyUser)
		})
	})

	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.ListProducts)
		r.Get("/{id}", productHandler.FindProductByID)
	})

	r.Route("/orders", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Post("/", ordersHandler.PlaceOrder)
		// TODO: get user's orders
	})

	// admin only endpoints
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.AdminMiddleware(userService))

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.ListUsers)
			r.Get("/search", userHandler.SearchUsers)
			r.Get("/{id}", userHandler.FindUserByID)
			r.Patch("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHandler.CreateProduct)
			r.Patch("/{id}", productHandler.UpdateProduct)
			r.Put("/{id}", productHandler.ReplaceProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
		})

		r.Route("/orders", func(r chi.Router) {
			// TODO: get all orders endpoint
			r.Get("/{id}", ordersHandler.FindOrderByID)
			// TODO: update order (maybe)
			// TODO: delete order
		})
	})

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
