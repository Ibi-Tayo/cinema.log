package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"cinema.log.server.golang/internal/auth"
	"cinema.log.server.golang/internal/database"
	"cinema.log.server.golang/internal/films"
	"cinema.log.server.golang/internal/ratings"
	"cinema.log.server.golang/internal/reviews"
	"cinema.log.server.golang/internal/users"
)

type Server struct {
	port int

	db            *sql.DB
	userHandler   *users.Handler
	authHandler   *auth.Handler
	authService   *auth.AuthService
	filmHandler   *films.Handler
	reviewHandler *reviews.Handler
	ratingHandler *ratings.Handler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Initialize database with migrations // change to just database.New() if not needing auto migrations
	db := database.NewWithMigrations()

	// Wire up dependencies: Database -> Store -> Service -> Handler
	userStore := users.NewStore(db)
	userService := users.NewService(userStore)
	userHandler := users.NewHandler(userService)

	authService := auth.NewService(userService)
	authHandler := auth.NewHandler(authService)

	// filmStore := films.NewStore(db)
	// filmService := films.NewService(filmStore)
	// filmHandler := films.NewHandler(filmService)

	// reviewStore := reviews.NewStore(db)
	// reviewService := reviews.NewService(reviewStore)
	// reviewHandler := reviews.NewHandler(reviewService)

	// ratingStore := ratings.NewStore(db)
	// ratingService := ratings.NewService(ratingStore)
	// ratingHandler := ratings.NewHandler(ratingService)

	NewServer := &Server{
		port:          port,
		db:            db,
		userHandler:   userHandler,
		authHandler:   authHandler,
		authService:   authService,
		// filmHandler:   filmHandler,
		// reviewHandler: reviewHandler,
		// ratingHandler: ratingHandler,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
