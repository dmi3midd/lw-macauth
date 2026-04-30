package models

import "github.com/golang-jwt/jwt/v5"

type UserDto struct {
	UserId   string
	Username string
	Email    string
	IsAdmin  bool
}

type TokensPair struct {
	RefreshToken string
	AccessToken  string
	TokenId      string
}

type AccessClaims struct {
	Username string
	Email    string
	IsAdmin  bool
	jwt.RegisteredClaims
}
