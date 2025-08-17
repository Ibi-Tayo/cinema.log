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

var conf *oauth2.Config = &oauth2.Config{
        ClientID:     GithubClientID,
        ClientSecret: GithubClientSecret,
        Scopes:       []string{},
        Endpoint:     oauth2github.Endpoint,
}

// gologin has a default cookie configuration for debug deployments (no TLS).
var cookieConf gologin.CookieConfig = gologin.DebugOnlyCookieConfig

type Handler struct {
	authService *AuthService
}

type AuthService struct {
	LoginHandler  http.Handler
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
	return github.StateHandler(cookieConf, github.CallbackHandler(conf, http.HandlerFunc(githubCallbackHandler), nil))
}

func (h *Handler) Logout() http.Handler {
	// TODO: Clear cookies
	return nil
}

func (h *Handler) RefreshToken() http.Handler {
	// TODO: Implement token refresh logic
	return nil
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
  ctx := r.Context()
  githubUser, err := github.UserFromContext(ctx)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // TODO: 1. Insert user into DB if doesn't exist
  // TODO: 2. Generate JWT
  // TODO: 3. Set JWT as cookie
  // TODO: 4. Redirect to user profile: http://localhost:4200/profile/{newOrCurrentUser.Username}

  w.Header().Set("Content-type", "application/json")
  buf, _ := json.Marshal(githubUser)
  fmt.Fprint(w, string(buf))
}