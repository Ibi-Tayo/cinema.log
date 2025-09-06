package server

import (
	"context"
	"net/http"
)

type key int

const (
	keyPrincipalID key = iota
	keyUser
	// ...
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes

	// User routes
	mux.HandleFunc("GET /users/{id}", s.userHandler.GetUserById)
	mux.HandleFunc("GET /users", s.userHandler.GetAllUsers)
	mux.HandleFunc("POST /users", s.userHandler.CreateUser)
	mux.HandleFunc("PUT /users", s.userHandler.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", s.userHandler.DeleteUser)

	// Auth routes
	mux.Handle("GET /auth/github-login", s.authHandler.Login())
	mux.Handle("GET /auth/logout", s.authHandler.Logout())
	mux.Handle("GET /auth/github-callback", s.authHandler.Callback())
	mux.Handle("GET /auth/refresh-token", s.authHandler.RefreshToken())

	// Film routes
	mux.HandleFunc("GET /films/{id}", s.filmHandler.GetFilmById)
	mux.HandleFunc("GET /films/search", s.filmHandler.GetFilmsFromExternal) // query param name = "f"

	// Wrap the mux with middleware
	return s.corsMiddleware(s.authMiddleware(mux))
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow github login path, otherwise validate auth token
		if r.URL.Path == "/auth/github-login" ||
			r.URL.Path == "/auth/github-callback" {
			next.ServeHTTP(w, r)
			return
		}
		// Check for authentication token in cookie
		authToken, err := r.Cookie("cinema-log-access-token")
		if err != nil {
			http.Error(w, http.ErrNoCookie.Error(), http.StatusUnauthorized)
			return
		}

		authTokenString := authToken.Value
		user, err := s.authService.ValidateJWT(authTokenString)
		if err != nil {
			http.Error(w, "jwt invalid", http.StatusUnauthorized)
			return
		}
		// so that downstream handlers can extract user from context
		ctx := context.WithValue(r.Context(), keyUser, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
