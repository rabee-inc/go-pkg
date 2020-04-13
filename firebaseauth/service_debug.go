package firebaseauth

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"

	"github.com/rabee-inc/go-pkg/log"
)

type serviceDebug struct {
	cli         *auth.Client
	dummyClaims map[string]interface{}
}

// Authentication ... 認証を行う
func (s *serviceDebug) Authentication(ctx context.Context, ah string) (string, map[string]interface{}, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	if user := getUserByAuthHeader(ah); user != "" {
		return user, s.dummyClaims, nil
	}

	// 通常の処理
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
func (s *serviceDebug) SetCustomClaims(ctx context.Context, userID string, claims map[string]interface{}) error {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	ah := getAuthHeader(ctx)
	if getUserByAuthHeader(ah) != "" {
		return nil
	}

	// 通常の処理
	err := s.cli.SetCustomUserClaims(ctx, userID, claims)
	if err != nil {
		log.Errorm(ctx, "c.SetCustomUserClaims", err)
		return err
	}
	return nil
}

func (s *serviceDebug) GetEmail(ctx context.Context, userID string) (string, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	ah := getAuthHeader(ctx)
	if user := getUserByAuthHeader(ah); user != "" {
		return "hirose.yuuki@rabee.jp", nil
	}

	// FirebaseAuthUserを取得
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

func (s *serviceDebug) GetTwitterID(ctx context.Context, userID string) (string, error) {
	// AuthorizationHeaderからUserが取得できたらデバッグリクエストと判定する
	ah := getAuthHeader(ctx)
	if getUserByAuthHeader(ah) != "" {
		return "", nil
	}

	// FirebaseAuthUserを取得
	user, err := s.cli.GetUser(ctx, userID)
	if err != nil {
		log.Errorm(ctx, "s.cli.GetUser", err)
		return "", err
	}
	if user == nil {
		return "", err
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

func (s *serviceDebug) ExistUser(ctx context.Context, userID string) (bool, error) {
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

func (s *serviceDebug) CreateUser(ctx context.Context, email string, password string, displayName string) (string, error) {
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

func (s *serviceDebug) UpdateUser(ctx context.Context, userID string, email *string, password *string, displayName *string) error {
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

func (s *serviceDebug) DeleteUser(ctx context.Context, userID string) error {
	params := (&auth.UserToUpdate{}).Disabled(true)
	_, err := s.cli.UpdateUser(ctx, userID, params)
	if err != nil {
		log.Errorm(ctx, "s.cli.UpdateUser", err)
		return err
	}
	return nil
}

// NewDebugService ... DebugServiceを作成する
func NewDebugService(cli *auth.Client, dummyClaims map[string]interface{}) Service {
	return &serviceDebug{
		cli:         cli,
		dummyClaims: dummyClaims,
	}
}
