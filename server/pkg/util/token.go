package util

import (
	"time"

	"github.com/google/uuid"
)

const (
	TokenCookieKey   = "Ttttt" // token
	RoleCookieKey    = "Rrrr"  // role
	AdminIdCookieKey = "Ddd"   // id
)

func NewToken() string {
	uid := uuid.New()
	return uid.String()
}

func GetExpireAt() *time.Time {
	t := time.Now().Add(7 * 24 * time.Hour)
	return &t
}
