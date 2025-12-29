package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"golang.org/x/oauth2"

	oauth2github "golang.org/x/oauth2/github"
)

var GithubClientID = os.Getenv("GITHUB_CLIENT_ID")
var GithubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
var BackendURL = os.Getenv("BACKEND_URL")

func isProduction() bool {
	return os.Getenv("ENVIRONMENT") == "production"
}

func getCookieSameSite() http.SameSite {
	if isProduction() {
		return http.SameSiteNoneMode // Required for cross-origin in production
	}
	return http.SameSiteLaxMode // Works for localhost in dev
}

func getCookieSecure() bool {
	return isProduction() // true in production (HTTPS), false in dev (HTTP)
}

var conf *oauth2.Config = &oauth2.Config{
	ClientID:     GithubClientID,
	ClientSecret: GithubClientSecret,
	RedirectURL:  BackendURL + "/auth/github-callback",
	Scopes:       []string{"user:email", "read:user"},
	Endpoint:     oauth2github.Endpoint,
}

// Cookie configuration for OAuth state parameter
// Automatically adjusts based on ENVIRONMENT variable
var cookieConf gologin.CookieConfig = gologin.CookieConfig{
	Name:     "oauth_state",
	Path:     "/",
	MaxAge:   300, // 5 minutes
	HTTPOnly: true,
	Secure:   getCookieSecure(),
	SameSite: getCookieSameSite(),
}

type Handler struct {
	authService *AuthService
}

func NewHandler(authService *AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) Login() http.Handler {
	return github.StateHandler(cookieConf, github.LoginHandler(conf, nil))
}

func (h *Handler) Callback() http.Handler {
	return github.StateHandler(cookieConf, github.CallbackHandler(conf, http.HandlerFunc(h.githubCallbackHandler), nil))
}

func (h *Handler) Logout() http.Handler {
	return http.HandlerFunc(h.logoutHandler)
}

func (h *Handler) RefreshToken() http.Handler {
	return http.HandlerFunc(h.refreshTokenHandler)
}

func (h *Handler) Me() http.Handler {
	return http.HandlerFunc(h.meHandler)
}

func (h *Handler) meHandler(w http.ResponseWriter, r *http.Request) {
	// Get JWT from cookie and validate
	cookie, err := r.Cookie("cinema-log-access-token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Return user data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear cookies - settings must match how they were set
	http.SetCookie(w, &http.Cookie{
		Name:     "cinema-log-access-token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "cinema-log-refresh-token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
	})
}

func (h *Handler) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	githubUser, err := github.UserFromContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	jwtResponse, err := h.authService.HandleGithubCallback(ctx, githubUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	h.setCookies(w, jwtResponse.Jwt, jwtResponse.RefreshToken)

	// Redirect to user profile: http://localhost:4200/profile/{userId}
	frontendURL := os.Getenv("FRONTEND_URL")
	http.Redirect(w, r, fmt.Sprintf("%s/profile/%s", frontendURL, jwtResponse.User.ID), http.StatusTemporaryRedirect)
}

func (h *Handler) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cinema-log-refresh-token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateRefreshToken(cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	jwt, refreshToken, err := h.authService.GenerateJWT(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.setCookies(w, jwt, refreshToken)
	w.WriteHeader(http.StatusOK)
}

func (*Handler) setCookies(w http.ResponseWriter, jwt string, refreshToken string) {
	// Cookie settings automatically adjust based on ENVIRONMENT variable:
	// - Development: SameSite=Lax, Secure=false (works on localhost)
	// - Production: SameSite=None, Secure=true (works cross-origin with HTTPS)
	http.SetCookie(w, &http.Cookie{
		Name:     "cinema-log-access-token",
		Value:    jwt,
		Path:     "/",
		HttpOnly: true,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
		MaxAge:   86400, // 24 hours to match JWT expiration
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "cinema-log-refresh-token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   getCookieSecure(),
		SameSite: getCookieSameSite(),
		MaxAge:   604800,
	})
}
