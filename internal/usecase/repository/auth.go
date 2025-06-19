package repository

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/bbrighter/dreams-api/internal/entity"
)

type AuthRepo struct {
	Auth entity.User
}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{}
}

const (
	validPassword string = "0803"
	validUser     string = "Benni"
)

func (r *AuthRepo) isValidUser() bool {
	return r.Auth.Name == validUser && r.Auth.Password == validPassword
}

func (r *AuthRepo) SetAuth(name, pw string) string {
	r.Auth.Name = name
	r.Auth.Password = pw
	if !r.isValidUser() {
		return ""
	}
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}
	r.Auth.Token = hex.EncodeToString(tokenBytes)
	return r.Auth.Token
}

func (r *AuthRepo) IsValidToken(token string) bool {
	return r.Auth.Token == token
}

func (r *AuthRepo) ClearToken() {
	r.Auth.Token = ""
	r.Auth.Password = ""
}
