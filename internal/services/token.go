package services

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"lw-macauth/internal/config"
	"lw-macauth/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/xid"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
	ErrInvalidAccessToken      = errors.New("invalid access token")
	ErrSubjectAndIDNotFound    = errors.New("subject and id not found")
)

type TokenService interface {
	// GenerateTokens generates pair with access and refresh tokens.
	// It returns tokens pair, tokenId (tokensPair, tokenId, error).
	GenerateTokens(user *models.UserDto, serviceId string) (*models.TokensPair, string, error)
	// ValidateRefreshToken validates refresh token.
	// It returns tokenId and userId (tokenId, userId, error).
	// It returns ("", "", error) if validation go wrong.
	// It returns ErrUnexpectedSigningMethod if the token uses an unexpected signing method.
	// It returns ErrInvalidRefreshToken if the token is invalid.
	// It returns ErrSubjectAndIDNotFound if subject or token ID are not found in claims.
	ValidateRefreshToken(refreshToken string) (string, string, error)
	// ValidateAccessToken validates access token.
	// It returns userDto and tokenId (userDto, tokenId, error).
	// It returns (nil, "", error) if validation go wrong.
	// It returns ErrUnexpectedSigningMethod if the token uses an unexpected signing method.
	// It returns ErrInvalidAccessToken if the token is invalid.
	// It returns ErrSubjectAndIDNotFound if subject or token ID are not found in claims.
	ValidateAccessToken(accessToken string) (*models.UserDto, string, error)
	// GetPublicKey returns public rsa keys
	GetPublicKey() rsa.PublicKey
}

type tokenService struct {
	keys config.KeysPair
}

func NewTokenService(keys *config.KeysPair) TokenService {
	return &tokenService{
		keys: *keys,
	}
}

func (s *tokenService) GenerateTokens(user *models.UserDto, serviceId string) (*models.TokensPair, string, error) {
	op := "tokenService.GenerateTokens"
	accessExpiry, _ := time.ParseDuration("30m")
	refreshExpiry, _ := time.ParseDuration("336h")
	now := time.Now()
	id := xid.New().String()

	// Access token
	accessClaims := models.AccessClaims{
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			Issuer:    "lw-macauth",
			Subject:   user.UserId,
			Audience:  jwt.ClaimStrings{serviceId},
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims).SignedString(s.keys.PrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	// Refresh token
	refreshClaims := jwt.RegisteredClaims{
		ID:        id,
		Issuer:    "lw-macauth",
		Subject:   user.UserId,
		Audience:  jwt.ClaimStrings{serviceId},
		ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims).SignedString(s.keys.PrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return &models.TokensPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, id, nil
}

func (s *tokenService) ValidateRefreshToken(refreshToken string) (string, string, error) {
	op := "tokenService.ValidateRefreshToken"
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("%s: %w %v", op, ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return s.keys.PublicKey, nil
	})

	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid {
		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidRefreshToken)
	}

	userId := claims.Subject
	tokenId := claims.ID

	if userId == "" || tokenId == "" {
		return "", "", fmt.Errorf("%s: %w", op, ErrSubjectAndIDNotFound)
	}

	return tokenId, userId, nil
}

func (s *tokenService) ValidateAccessToken(accessToken string) (*models.UserDto, string, error) {
	op := "tokenService.ValidateAccessToken"
	claims := &models.AccessClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("%s: %w %v", op, ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return s.keys.PublicKey, nil
	})

	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	if !token.Valid {
		return nil, "", fmt.Errorf("%s: %w", op, ErrInvalidAccessToken)
	}

	userId := claims.Subject
	tokenId := claims.ID

	if userId == "" || tokenId == "" {
		return nil, "", fmt.Errorf("%s: %w", op, ErrSubjectAndIDNotFound)
	}

	return &models.UserDto{
		UserId:   userId,
		Username: claims.Username,
		Email:    claims.Email,
		IsAdmin:  claims.IsAdmin,
	}, tokenId, nil
}

func (s *tokenService) GetPublicKey() rsa.PublicKey {
	return *s.keys.PublicKey
}
