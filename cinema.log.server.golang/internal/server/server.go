package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"cinema.log.server.golang/internal/database"
	"cinema.log.server.golang/internal/users"
)

type Server struct {
	port int

	db          database.Service
	userHandler *users.Handler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	
	// Initialize database
	db := database.New()
	
	// Wire up dependencies: Database -> Store -> Service -> Handler
	userStore := users.NewStore(db)
	userService := users.NewService(userStore)
	userHandler := users.NewHandler(userService)
	
	NewServer := &Server{
		port:        port,
		db:          db,
		userHandler: userHandler,
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
