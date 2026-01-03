package server

import (
	"context"
	"net/http"
	"os"

	"cinema.log.server.golang/internal/middleware"
)

// isAuthExempt checks if a path should bypass authentication
func isAuthExempt(path string) bool {
	exemptPaths := []string{
		"/auth/github-login",
		"/auth/github-callback",
		"/auth/refresh-token",
		"/auth/dev/login",
	}
	for _, exemptPath := range exemptPaths {
		if path == exemptPath {
			return true
		}
	}
	return false
}

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
	mux.Handle("GET /auth/me", s.authHandler.Me())
	mux.Handle("GET /auth/dev/login", s.authHandler.DevLogin()) // only in dev environment

	// Film routes
	mux.HandleFunc("GET /films/{id}", s.filmHandler.GetFilmById)
	mux.HandleFunc("POST /films", s.filmHandler.CreateFilm)
	mux.HandleFunc("GET /films/search", s.filmHandler.GetFilmsFromExternal) // query param name = "f"
	mux.HandleFunc("GET /films/for-comparison", s.filmHandler.GetFilmsForComparison) // query params: userId, filmId

	// Review routes
	mux.HandleFunc("GET /reviews/{userId}", s.reviewHandler.GetAllReviews)
	mux.HandleFunc("POST /reviews", s.reviewHandler.CreateReview)
	mux.HandleFunc("PUT /reviews/{id}", s.reviewHandler.UpdateReview)
	mux.HandleFunc("DELETE /reviews", s.reviewHandler.DeleteReview)

	// Rating routes
	mux.HandleFunc("GET /ratings/{userId}", s.ratingHandler.GetRatingsByUserId)
	mux.HandleFunc("GET /ratings", s.ratingHandler.GetRating) // query params: userId, filmId
	mux.HandleFunc("POST /ratings/compare-films", s.ratingHandler.CompareFilms)

	// Wrap the mux with middleware
	return s.corsMiddleware(s.authMiddleware(mux))
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get frontend URL from environment, fallback to localhost for development
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:4200"
		}

		// Set CORS headers - must use specific origin with credentials
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

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
		// Allow certain auth paths to bypass token validation
		if isAuthExempt(r.URL.Path) {
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
		ctx := context.WithValue(r.Context(), middleware.KeyUser, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
