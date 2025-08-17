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
	// TODO: Clear cookies
	return nil
}

func (h *Handler) RefreshToken() http.Handler {
	// TODO: Implement token refresh logic
	return nil
}

func (h *Handler) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
  ctx := r.Context()
  githubUser, err := github.UserFromContext(ctx)
  if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }

  createdOrReturnedUser, jwt, refreshToken, err := h.authService.HandleGithubCallback(ctx, githubUser)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-type", "application/json")
  http.SetCookie(w, &http.Cookie{
    Name:  "cinema-log-access-token",
    Value: jwt,
    Path:  "/",
	HttpOnly: true,  
	Secure:   false, // will set to true in production
	SameSite: http.SameSiteStrictMode, 
	MaxAge:  3600,  
  })
  http.SetCookie(w, &http.Cookie{
	Name:  "cinema-log-refresh-token",
	Value: refreshToken,
	Path:  "/",
	HttpOnly: true,  
	Secure:   false,
	SameSite: http.SameSiteStrictMode, 
	MaxAge:  3600,  
  })

  // Redirect to user profile: http://localhost:4200/profile/{newOrCurrentUser.Username}
  frontendURL := os.Getenv("FRONTEND_URL")
  http.Redirect(w, r, fmt.Sprintf("%s/profile/%s", frontendURL, createdOrReturnedUser.Username), http.StatusTemporaryRedirect)

}