package tenant

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func generateInvitationToken() (string, string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", "", fmt.Errorf("failde to generate random bytes: %w", err)
	}

	rawToken := base64.URLEncoding.EncodeToString(b)
	tokenHash := hashToken(rawToken)

	return rawToken, tokenHash, nil
}

func hashToken(token string) string {
	hashBytes := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hashBytes[:])
}
