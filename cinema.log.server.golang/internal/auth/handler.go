package auth

import (
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

var conf *oauth2.Config = &oauth2.Config{
        ClientID:     GithubClientID,
        ClientSecret: GithubClientSecret,
		RedirectURL: "http://localhost:8080/auth/github-callback",
        Scopes:       []string{"user:email", "read:user"},
        Endpoint:     oauth2github.Endpoint,
}

// gologin has a default cookie configuration for debug deployments (no TLS).
var cookieConf gologin.CookieConfig = gologin.DebugOnlyCookieConfig

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

func (h *Handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:   "cinema-log-access-token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "cinema-log-refresh-token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
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

  // Redirect to user profile: http://localhost:4200/profile/{newOrCurrentUser.Username}
  frontendURL := os.Getenv("FRONTEND_URL")
  http.Redirect(w, r, fmt.Sprintf("%s/profile/%s", frontendURL, jwtResponse.User.Username), http.StatusTemporaryRedirect)
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
	http.SetCookie(w, &http.Cookie{
		Name:     "cinema-log-access-token",
		Value:    jwt,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   600,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "cinema-log-refresh-token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800,
	})
}