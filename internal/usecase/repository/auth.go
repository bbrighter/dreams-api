package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

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
	h := sha256.New()
	h.Write([]byte(name + pw + time.Now().String()))
	sum := h.Sum(nil)
	r.Auth.Token = hex.EncodeToString(sum)
	return r.Auth.Token
}

func (r *AuthRepo) IsValidToken(token string) bool {
	return r.Auth.Token == token
}
