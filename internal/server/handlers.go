package server

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dmi3midd/lw-macauth/internal/models"
	"github.com/dmi3midd/lw-macauth/internal/shared/apierror"
	"github.com/dmi3midd/lw-macauth/internal/shared/utils"
)

type GenerateTokensRequest struct {
	User      models.UserDto `json:"user"`
	ServiceId string         `json:"serviceId"`
}

func (r GenerateTokensRequest) Validate() error {
	if strings.TrimSpace(r.ServiceId) == "" {
		return errors.New("serviceId is required")
	}
	if strings.TrimSpace(r.User.UserId) == "" {
		return errors.New("user.userId is required")
	}
	if strings.TrimSpace(r.User.Username) == "" {
		return errors.New("user.username is required")
	}
	if strings.TrimSpace(r.User.Email) == "" {
		return errors.New("user.email is required")
	}
	return nil
}

type GenerateTokensResponse struct {
	Tokens  models.TokensPair `json:"tokens"`
	TokenId string            `json:"tokenId"`
}

func (s *Server) GenerateTokens(w http.ResponseWriter, r *http.Request) error {
	req, err := utils.BindAndValidate[GenerateTokensRequest](r)
	if err != nil {
		return err
	}

	tokensPair, tokenId, err := s.tokenService.GenerateTokens(&req.User, req.ServiceId)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(GenerateTokensResponse{
		Tokens:  *tokensPair,
		TokenId: tokenId,
	})
}

type ValidateAccessTokenRequest struct {
	AccessToken string `json:"accessToken"`
}

func (r ValidateAccessTokenRequest) Validate() error {
	if strings.TrimSpace(r.AccessToken) == "" {
		return errors.New("accessToken is required")
	}
	return nil
}

type ValidateAccessTokenResponse struct {
	User models.UserDto `json:"user"`
}

func (s *Server) ValidateAccessToken(w http.ResponseWriter, r *http.Request) error {
	req, err := utils.BindAndValidate[ValidateAccessTokenRequest](r)
	if err != nil {
		return err
	}

	userData, _, err := s.tokenService.ValidateAccessToken(req.AccessToken)
	if err != nil {
		return apierror.NewUnauthorizedError(err, "Invalid access token")
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(ValidateAccessTokenResponse{
		User: *userData,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (r RefreshTokenRequest) Validate() error {
	if strings.TrimSpace(r.RefreshToken) == "" {
		return errors.New("refreshToken is required")
	}
	return nil
}

type RefreshTokenResponse struct {
	UserId  string `json:"userId"`
	TokenId string `json:"tokenId"`
}

func (s *Server) ValidateRefreshToken(w http.ResponseWriter, r *http.Request) error {
	req, err := utils.BindAndValidate[RefreshTokenRequest](r)
	if err != nil {
		return err
	}

	tokenId, userId, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return apierror.NewUnauthorizedError(err, "Invalid refresh token")
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(RefreshTokenResponse{
		UserId:  userId,
		TokenId: tokenId,
	})
}

type GetPublicKeyResponse struct {
	Key rsa.PublicKey `json:"key"`
}

func (s *Server) GetPublicKey(w http.ResponseWriter, r *http.Request) error {
	key := s.tokenService.GetPublicKey()

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(GetPublicKeyResponse{
		Key: key,
	})
}
