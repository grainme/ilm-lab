package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grainme/movie-api/internal/cache"
	handlers "github.com/grainme/movie-api/internal/handler"
	"github.com/grainme/movie-api/internal/middleware"
	"github.com/grainme/movie-api/internal/repository/postgres"
	"github.com/grainme/movie-api/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbDSN := os.Getenv("DB_DSN")
	redisAddr := os.Getenv("REDIS_ADDR")

	// Migrations (tables creation)
	log.Println("Running database migrations...")
	m, err := migrate.New("file://db/migrations", dbDSN)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// ErrNoChanges means DB is already up-to-date
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	// setup postgresDB
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
		os.Exit(1)
	}

	// setup redis
	rdb, err := cache.NewRedisClient(redisAddr)
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
		os.Exit(1)
	}
	defer rdb.Close()

	movieRepo := postgres.NewPostgresMovieRepository(db)
	reviewRepo := postgres.NewPostgresReviewRepository(db)
	userRepo := postgres.NewPostgresUserRepository(db)

	movieService := service.NewMovieService(movieRepo, rdb)
	reviewService := service.NewReviewService(reviewRepo, rdb)
	userService := service.NewUserService(userRepo, rdb)

	movieHandler := handlers.NewMovieHandler(movieService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	userHandler := handlers.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(chimw.Logger)

	// movie routes
	r.Get("/movies", movieHandler.GetAllMovies)
	r.Get("/movies/{id}", movieHandler.GetMovieById)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Post("/movies", movieHandler.AddMovie)
		r.Put("/movies/{id}", movieHandler.UpdateMovieTitleById)
		r.With(middleware.Authorize).Delete("/movies/{id}", movieHandler.DeleteById)
	})
	r.Get("/movies/{id}/reviews", movieHandler.GetMovieWithReviews)

	// reviews routes
	r.Get("/reviews", reviewHandler.GetAllReviews)
	r.Post("/reviews", reviewHandler.AddReview)
	r.Get("/reviews/{id}", reviewHandler.GetAllReviewsByMovieId)

	// auth routes
	r.Post("/auth/login", userHandler.Login)
	r.Post("/auth/register", userHandler.Register)
	r.Post("/auth/refresh", userHandler.RefreshToken)
	r.Post("/auth/logout", userHandler.Logout)

	port := os.Getenv("PORT")
	if port == "" {
		// default value
		port = "3000"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
