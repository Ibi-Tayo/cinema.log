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
	"cinema.log.server.golang/internal/users"
)

type Server struct {
	port int

	db          *sql.DB
	userHandler *users.Handler
	authHandler *auth.Handler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	
	// Initialize database
	db := database.New()
	
	// Wire up dependencies: Database -> Store -> Service -> Handler
	userStore := users.NewStore(db)
	userService := users.NewService(userStore)
	userHandler := users.NewHandler(userService)

	authService := auth.NewService(userService)
	authHandler := auth.NewHandler(authService)

	
	NewServer := &Server{
		port:        port,
		db:          db,
		userHandler: userHandler,
		authHandler: authHandler,
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
