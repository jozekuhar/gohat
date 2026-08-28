package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
	"uuid"

	"mimokocke/internal/shared/clock"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/routes"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Service struct {
	cfg               *config.Config
	logger            *slog.Logger
	clock             clock.Clock
	googleOauthConfig *oauth2.Config
	authRepo          *repository
}

func NewService(
	cfg *config.Config,
	logger *slog.Logger,
	clock clock.Clock,
	authRepo *repository,
) *Service {
	redirectURL := fmt.Sprintf("http://%s%s", cfg.Host, routes.CallbackSignInGoogle)

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
		logger:            logger,
		clock:             clock,
		googleOauthConfig: googleOauthConfig,
		authRepo:          authRepo,
	}
}

func (s *Service) LoginUserWithPassword(
	ctx context.Context,
	email, password string,
) (Session, error) {
	authentication, err := s.authRepo.GetAuthenticationByEmail(ctx, email, AuthProviderPassword)
	if err != nil {
		return Session{}, err
	}

	if !CheckPasswordHash(password, *authentication.PasswordHash) {
		return Session{}, fmt.Errorf("invalid password")
	}

	session, err := s.createSession(ctx, nil, authentication.UserID)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

type RegisterUserWithPasswordParams struct {
	Email    string
	Password string
}

func (s *Service) RegisterUserWithPassword(
	ctx context.Context,
	params RegisterUserWithPasswordParams,
) (Session, error) {
	var zero Session

	tx, err := s.authRepo.pool.Begin(ctx)
	if err != nil {
		return zero, nil
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	user, err := s.authRepo.CreateUser(ctx, tx, User{
		ID:    uuid.NewV7(),
		Email: params.Email,
	})
	if err != nil {
		return zero, err
	}

	passwordHash, err := HashPassword(params.Password)
	if err != nil {
		return zero, err
	}

	_, err = s.authRepo.CreateAuthentication(ctx, tx, Authentication{
		ID:           uuid.NewV7(),
		UserID:       user.ID,
		Provider:     AuthProviderPassword,
		PasswordHash: &passwordHash,
	})
	if err != nil {
		return zero, err
	}

	session, err := s.createSession(ctx, tx, user.ID)
	if err != nil {
		return zero, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return zero, nil
	}

	return session, nil
}

func (s *Service) GenerateGoogleSignInURLAndState(ctx context.Context) (string, string) {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	url := s.googleOauthConfig.AuthCodeURL(state)

	return url, state
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

func (s *Service) ProcessGoogleSignIn(ctx context.Context, code string) (Session, error) {
	var zero Session

	token, err := s.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return zero, err
	}

	url := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken
	res, err := http.Get(url)
	if err != nil {
		return zero, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return zero, err
	}

	var userInfo googleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return zero, fmt.Errorf("unmarshal user info: %w", err)
	}

	if !userInfo.VerifiedEmail {
		return zero, fmt.Errorf("user email is not verified")
	}

	// Ce ni authentication za ta mail potem naredimo userja in
	// ce dobimo error to pomeni da user ze obstaja torej bo moral
	// se prvo prijavit in potem linkat.

	authentication, err := s.authRepo.GetAuthenticationByProvider(
		ctx,
		AuthProviderGoogle,
		userInfo.ID,
	)
	if err != nil {
		tx, err := s.authRepo.pool.Begin(ctx)
		if err != nil {
			return zero, err
		}
		defer func() {
			err := tx.Rollback(ctx)
			if err != nil {
				s.logger.Error("rollback failed", "err", err)
				return
			}
		}()

		newUser, err := s.authRepo.CreateUser(ctx, nil, User{
			ID:    uuid.NewV7(),
			Email: userInfo.Email,
		})
		if err != nil {
			return zero, err
		}

		authentication, err = s.authRepo.CreateAuthentication(ctx, tx, Authentication{
			ID:         uuid.NewV7(),
			UserID:     newUser.ID,
			Provider:   AuthProviderGoogle,
			ProviderID: &userInfo.ID,
		})
		if err != nil {
			return zero, err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return zero, err
		}
	}

	session, err := s.createSession(ctx, nil, authentication.UserID)
	if err != nil {
		return zero, err
	}

	return session, nil
}

func (s *Service) VerifySession(ctx context.Context, sessionID string) (Session, error) {
	var zero Session

	id, err := uuid.Parse(sessionID)
	if err != nil {
		return zero, fmt.Errorf("parsing session id: %w", err)
	}

	session, err := s.authRepo.GetSession(ctx, id)
	if err != nil {
		return zero, fmt.Errorf("retrieve session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return zero, fmt.Errorf("session expired: user %s", session.UserID)
	}

	return session, nil
}

func (s *Service) LogoutUser(ctx context.Context, sessionIDStr string) error {
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return err
	}

	err = s.authRepo.DeleteSession(ctx, sessionID)
	return err
}

func (s *Service) createSession(ctx context.Context, tx pgx.Tx, userID uuid.UUID) (Session, error) {
	session, err := s.authRepo.CreateSession(ctx, tx, Session{
		ID:        uuid.NewV7(),
		UserID:    userID,
		ExpiresAt: s.clock.NowUTC().Add(time.Hour * 24 * 14), // 14 days
	})
	if err != nil {
		return Session{}, err
	}
	return session, nil
}
