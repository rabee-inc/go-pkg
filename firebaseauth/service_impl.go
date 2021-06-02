package firebaseauth

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"

	"github.com/rabee-inc/go-pkg/log"
)

type service struct {
	cli *auth.Client
}

// Authentication ... 認証を行う
func (s *service) Authentication(ctx context.Context, ah string) (string, map[string]interface{}, error) {
	token := getTokenByAuthHeader(ah)
	if token == "" {
		err := log.Warninge(ctx, "token empty error")
		return "", nil, err
	}

	t, err := s.cli.VerifyIDToken(ctx, token)
	if err != nil {
		msg := fmt.Sprintf("c.VerifyIDToken: %s", token)
		log.Warningm(ctx, msg, err)
		return "", nil, err
	}
	return t.UID, t.Claims, nil
}

// SetCustomClaims ... カスタムClaimsを設定
func (s *service) SetCustomClaims(ctx context.Context, userID string, claims map[string]interface{}) error {
	err := s.cli.SetCustomUserClaims(ctx, userID, claims)
	if err != nil {
		log.Errorm(ctx, "c.SetCustomUserClaims", err)
		return err
	}
	return nil
}

func (s *service) GetEmail(ctx context.Context, userID string) (string, error) {
	user, err := s.cli.GetUser(ctx, userID)
	if err != nil {
		log.Errorm(ctx, "s.cli.GetUser", err)
		return "", err
	}
	if user == nil {
		return "", nil
	}
	return user.Email, nil
}

func (s *service) GetTwitterID(ctx context.Context, userID string) (string, error) {
	user, err := s.cli.GetUser(ctx, userID)
	if err != nil {
		log.Errorm(ctx, "s.cli.GetUser", err)
		return "", err
	}
	if user == nil {
		return "", nil
	}

	dst := ""
	for _, userInfo := range user.ProviderUserInfo {
		if userInfo != nil && userInfo.ProviderID == "twitter.com" {
			dst = userInfo.UID
			break
		}
	}
	return dst, nil
}

func (s *service) ExistUser(ctx context.Context, userID string) (bool, error) {
	user, err := s.cli.GetUser(ctx, userID)
	if err != nil {
		log.Errorm(ctx, "s.cli.GetUser", err)
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return true, nil
}

func (s *service) IsEmailVerified(ctx context.Context, userID string) (bool, error) {
	user, err := s.cli.GetUser(ctx, userID)
	if err != nil {
		log.Errorm(ctx, "s.cli.GetUser", err)
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return user.EmailVerified, nil
}

func (s *service) IsLinkedProviders(ctx context.Context, userID string, providers []Provider) (bool, error) {
	user, err := s.cli.GetUser(ctx, userID)
	if err != nil {
		log.Errorm(ctx, "s.cli.GetUser", err)
		return false, err
	}
	if user == nil {
		return false, nil
	}

	for _, userInfo := range user.ProviderUserInfo {
		if userInfo != nil {
			providerID := Provider(userInfo.ProviderID)
			for _, provider := range providers {
				if provider == providerID {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (s *service) CreateUser(ctx context.Context, email string, password string, displayName string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(displayName)
	user, err := s.cli.CreateUser(ctx, params)
	if err != nil {
		log.Errorm(ctx, "s.cli.CreateUser", err)
		return "", err
	}
	return user.UID, nil
}

func (s *service) UpdateUser(ctx context.Context, userID string, email *string, password *string, displayName *string) error {
	params := (&auth.UserToUpdate{})
	if email != nil {
		params = params.Email(*email)
	}
	if password != nil {
		params = params.Password(*password)
	}
	if displayName != nil {
		params = params.DisplayName(*displayName)
	}
	_, err := s.cli.UpdateUser(ctx, userID, params)
	if err != nil {
		log.Errorm(ctx, "s.cli.UpdateUser", err)
		return err
	}
	return nil
}

func (s *service) DeleteUser(ctx context.Context, userID string) error {
	params := (&auth.UserToUpdate{}).Disabled(true)
	_, err := s.cli.UpdateUser(ctx, userID, params)
	if err != nil {
		log.Errorm(ctx, "s.cli.UpdateUser", err)
		return err
	}
	return nil
}

func (s *service) GeneratePasswordRemindURL(ctx context.Context, userID string, email string, setting *auth.ActionCodeSettings) (string, error) {
	url, err := s.cli.PasswordResetLinkWithSettings(ctx, email, setting)
	if err != nil {
		log.Errorm(ctx, "s.cli.PasswordResetLinkWithSettings", err)
		return "", err
	}
	return url, err
}

// NewService ... Serviceを作成する
func NewService(cli *auth.Client) Service {
	return &service{
		cli: cli,
	}
}
