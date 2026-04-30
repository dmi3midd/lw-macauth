package models

import "github.com/golang-jwt/jwt/v5"

type UserDto struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

type TokensPair struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

type AccessClaims struct {
	Username string
	Email    string
	IsAdmin  bool
	jwt.RegisteredClaims
}
