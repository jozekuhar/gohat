package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"tmpl/internal/model"
	"tmpl/internal/repository"
	"tmpl/internal/shared/config"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	cfg               *config.Config
	googleOauthConfig *oauth2.Config
	authRepo          *repository.Auth
	userSrv           *User
}

func NewAuth(cfg *config.Config, authRepo *repository.Auth, userSrv *User) *Auth {
	googleOauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
	return &Auth{
		cfg:               cfg,
		googleOauthConfig: googleOauthConfig,
		authRepo:          authRepo,
		userSrv:           userSrv,
	}
}

func (s *Auth) GenerateGoogleLoginURL(ctx context.Context) string {
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

func (s *Auth) ProcessGoogleLogin(ctx context.Context, code, state string) (*model.Session, error) {
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

	user, err := s.userSrv.GetOrCreateUser(ctx, userInfo.Email, userInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("get or create user: %w", err)
	}

	sessionID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 14)
	session, err := s.authRepo.CreateSession(ctx, sessionID, user.ID, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return session, nil
}

func (s *Auth) VerifySession(ctx context.Context, sessionID string) error {
	id, err := uuid.Parse(sessionID)
	if err != nil {
		return fmt.Errorf("parsing session id: %w", err)
	}

	session, err := s.authRepo.GetSession(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieve session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return fmt.Errorf("session expired fof user with id: %d", session.UserID)
	}

	return nil
}

func (s *Auth) Logout(ctx context.Context, sessionIDStr string) error {
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return err
	}

	err = s.authRepo.DeleteSession(ctx, sessionID)
	return err
}
