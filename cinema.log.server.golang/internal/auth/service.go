package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/users"
	"cinema.log.server.golang/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v52/github"
)

var tokenSecret = os.Getenv("TOKEN_SECRET")

type AuthService struct {
	userService users.UserService
}

type JwtResponse struct {
	User         *domain.User
	Jwt          string
	RefreshToken string
}

func NewService(userService users.UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) HandleGithubCallback(ctx context.Context, githubUser *github.User) (*JwtResponse, error) {
	user, err := s.userService.GetOrCreateUserByGithubId(ctx, githubUser.GetID(),
		githubUser.GetName(), githubUser.GetLogin(), githubUser.GetAvatarURL())

	if err != nil {
		return nil, fmt.Errorf("failed to get or create user: %w", err)
	}

	jwt, refreshToken, err := s.GenerateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return &JwtResponse{
		User:         user,
		Jwt:          jwt,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GenerateJWT(user *domain.User) (string, string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID.String(),
		"name":     user.Name,
		"username": user.Username,
		"iss":      "cinema.log.server.golang",
		"aud":      "cinema.log.client",
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 1 day expiration
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID.String(),
		"name":     user.Name,
		"username": user.Username,
		"iss":      "cinema.log.server.golang",
		"aud":      "cinema.log.client",
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 day expiration
	})

	// Sign and get the complete encoded token as a string using the secret
	jwtTokenString, err := jwtToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return jwtTokenString, refreshTokenString, nil
}

func (s *AuthService) ValidateJWT(tkn string) (*domain.User, error) {
	ctx := context.Background()
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token claims: missing user ID")
		}

		userUuid, err := utils.ParseUUID(userID)
		if err != nil {
			return nil, fmt.Errorf("invalid token claims: %w", err)
		}

		user, err := s.userService.GetUserById(ctx, userUuid)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by ID: %w", err)
		}

		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthService) ValidateRefreshToken(tkn string) (*domain.User, error) {
	ctx := context.Background()
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token claims: missing user ID")
		}

		userUuid, err := utils.ParseUUID(userID)
		if err != nil {
			return nil, fmt.Errorf("invalid token claims: %w", err)
		}

		user, err := s.userService.GetUserById(ctx, userUuid)
		if err != nil {
			return nil, fmt.Errorf("failed to get user by ID: %w", err)
		}

		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthService) HandleDevLogin(ctx context.Context) (*JwtResponse, error) {
	user, err := s.userService.GetOrCreateUserByGithubId(ctx, 0, "Dev User", "devuser", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get first user: %w", err)
	}

	jwt, refreshToken, err := s.GenerateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return &JwtResponse{
		User:         user,
		Jwt:          jwt,
		RefreshToken: refreshToken,
	}, nil
}