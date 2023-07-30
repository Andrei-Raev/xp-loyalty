package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Andrei-Raev/xp-loyalty/internal/model"
	"github.com/golang-jwt/jwt"
)

type CredentialsRepository interface {
	Create(ctx context.Context, credentials model.Credentials) (string, error)
	GetByUsername(ctx context.Context, username string) (model.Credentials, error)
	GetByCredentials(ctx context.Context, username, password string) (model.Credentials, error)
}

type AuthService struct {
	repo              CredentialsRepository
	moderatorUsername string
	moderatorPassword string
	salt              string
	signingKey        []byte
	expireDuration    time.Duration
}

type AuthClaims struct {
	jwt.StandardClaims
	Credentials model.Credentials
}

func NewAuthService(credentialsRepo CredentialsRepository, moderatorUsername, moderatorPassword string) *AuthService {
	return &AuthService{
		repo:              credentialsRepo,
		moderatorUsername: moderatorUsername,
		moderatorPassword: moderatorPassword,
		salt:              "authSecertSalt",
		signingKey:        []byte("authSecretSigningKey"),
		expireDuration:    24 * time.Hour,
	}
}

func (s *AuthService) SignUp(ctx context.Context, credentials model.Credentials) (string, error) {
	if credentials.Username == s.moderatorUsername {
		return "", model.ErrUserExists
	}

	pwd := sha1.New()
	pwd.Write([]byte(credentials.Password))
	pwd.Write([]byte(s.salt))

	_, err := s.repo.GetByUsername(ctx, credentials.Username)
	switch {
	case errors.Is(err, model.ErrUserNotFound):
	case err == nil:
		return "", model.ErrUserExists
	default:
		return "", err
	}

	credentials.Password = fmt.Sprintf("%x", pwd.Sum(nil))
	return s.repo.Create(ctx, credentials)
}

func (s *AuthService) SignIn(ctx context.Context, username, password string) (string, error) {
	var err error
	credentials := model.Credentials{
		Password: password,
		CredentialsSecure: model.CredentialsSecure{
			Username: username,
			Role:     model.RoleModerator,
		},
	}

	if username != s.moderatorUsername || password != s.moderatorPassword {
		pwd := sha1.New()
		pwd.Write([]byte(password))
		pwd.Write([]byte(s.salt))
		password = fmt.Sprintf("%x", pwd.Sum(nil))

		credentials, err = s.repo.GetByCredentials(ctx, username, password)
		if err != nil {
			return "", err
		}
	}

	credentials.Password = ""
	claims := AuthClaims{
		Credentials: credentials,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.expireDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.signingKey)
}

func (s *AuthService) CheckAccess(ctx context.Context, token string, role model.Role) (model.Credentials, error) {
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return model.Credentials{}, model.ErrInvalidAccessToken
	}

	credentials, err := s.parseToken(tokenParts[1])
	if err != nil {
		return model.Credentials{}, err
	}

	if credentials.Role < role {
		return model.Credentials{}, model.ErrWrongRole
	}

	return credentials, nil
}

func (s *AuthService) parseToken(accessToken string) (model.Credentials, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return model.Credentials{}, model.ErrUnexpectedSigningMethod
		}
		return s.signingKey, nil
	})
	if err != nil {
		return model.Credentials{}, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.Credentials, nil
	}

	return model.Credentials{}, model.ErrInvalidAccessToken
}
