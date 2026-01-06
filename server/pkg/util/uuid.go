package util

import (
	"github.com/google/uuid"
	"strings"
)

func NewPrivateKey() string {
	privateKey := uuid.NewString()
	return strings.ReplaceAll(privateKey, "-", "")
}
