package handlers

import (
	"crypto/rand"
	"fmt"
)

func GenerateRecoveryToken() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return fmt.Sprintf(
		"%s-%s-%s-%s",
		b[0:4],
		b[4:8],
		b[8:12],
		b[12:16],
	)
}
