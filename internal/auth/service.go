package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"uuid"

	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/routes"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service struct {
	cfg               *config.Config
	googleOauthConfig *oauth2.Config
	authRepo          *repository
}

func NewService(cfg *config.Config, authRepo *repository) *Service {
	redirectURL := fmt.Sprintf("http://%s%s", cfg.Host, routes.ILoginGoogleCallback)

	googleOauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	return &Service{
		cfg:               cfg,
		googleOauthConfig: googleOauthConfig,
		authRepo:          authRepo,
	}
}

func (u *Service) GetOrCreateUser(
	ctx context.Context,
	email string,
	googleID string,
) (*User, error) {
	id := uuid.NewV7()

	user, err := u.authRepo.GetOrCreateUserWithGoogleID(ctx, id, email, googleID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GenerateGoogleLoginURL(ctx context.Context) string {
	return s.googleOauthConfig.AuthCodeURL("todo")
}

type googleUserInfo struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	HD            string `json:"hd"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (s *Service) ProcessGoogleLogin(
	ctx context.Context,
	code, state string,
) (*Session, error) {
	token, err := s.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	url := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var userInfo googleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("unmarshal user info: %w", err)
	}

	if !userInfo.VerifiedEmail {
		return nil, fmt.Errorf("user email is not verified")
	}

	user, err := s.GetOrCreateUser(ctx, userInfo.Email, userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("get or create user: %w", err)
	}

	sessionID := uuid.NewV7()
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 14)
	session, err := s.authRepo.CreateSession(ctx, sessionID, user.ID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return session, nil
}

func (s *Service) VerifySession(ctx context.Context, sessionID string) (*Session, error) {
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("parsing session id: %w", err)
	}

	session, err := s.authRepo.GetSession(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("retrieve session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired fof user with id: %d", session.UserID)
	}

	return session, err
}

func (s *Service) Logout(ctx context.Context, sessionIDStr string) error {
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return err
	}

	err = s.authRepo.DeleteSession(ctx, sessionID)
	return err
}
