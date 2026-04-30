package server

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	errs "lw-macauth/internal/errors"
	"lw-macauth/internal/models"
	"net/http"
	"strings"
)

type GenerateTokensRequest struct {
	User      models.UserDto `json:"user"`
	ServiceId string         `json:"serviceId"`
}

type GenerateTokensResponse struct {
	Tokens  models.TokensPair `json:"tokens"`
	TokenId string            `json:"tokenId"`
}

func (s *Server) GenerateTokens(w http.ResponseWriter, r *http.Request) error {
	var req GenerateTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewBadRequestError(err, "Invalid request body")
	}

	tokensPair, tokenId, err := s.tokenService.GenerateTokens(&req.User, req.ServiceId)
	if err != nil {
		return errs.NewInternalServerError(err)
	}

	response := GenerateTokensResponse{
		Tokens:  *tokensPair,
		TokenId: tokenId,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}
	return nil
}

type ValidateAccessTokenResponse struct {
	User models.UserDto `json:"user"`
}

func (s *Server) ValidateAccessToken(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	token := ""
	if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
		token = after
	}
	if token == "" {
		return errs.NewUnauthorizedError(
			fmt.Errorf("Invalid or empty Authorization header"),
			"Invalid or empty Authorization header",
		)
	}
	userData, _, err := s.tokenService.ValidateAccessToken(token)
	if err != nil {
		return errs.NewUnauthorizedError(err, "Invalid access token")
	}

	w.Header().Set("Content-Type", "application/json")
	response := ValidateAccessTokenResponse{
		User: *userData,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}
	return nil
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	UserId  string `json:"userId"`
	TokenId string `json:"tokenId"`
}

func (s *Server) ValidateRefreshToken(w http.ResponseWriter, r *http.Request) error {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewBadRequestError(err, "Invalid request body")
	}
	refreshToken := req.RefreshToken
	if refreshToken == "" {
		return errs.NewUnauthorizedError(
			fmt.Errorf("Invalid or empty refresh token"),
			"Invalid or empty refresh token",
		)
	}
	tokenId, userId, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return errs.NewUnauthorizedError(err, "Invalid refresh token")
	}

	response := RefreshTokenResponse{
		UserId:  userId,
		TokenId: tokenId,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}
	return nil
}

type GetPublicKeyResponse struct {
	Key rsa.PublicKey `json:"key"`
}

func (s *Server) GetPublicKey(w http.ResponseWriter, r *http.Request) error {
	key := s.tokenService.GetPublicKey()

	w.Header().Set("Content-Type", "application/json")
	response := GetPublicKeyResponse{
		Key: key,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}
	return nil
}
