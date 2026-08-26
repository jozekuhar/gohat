package tenant

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func generateInviteToken() (string, string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", "", fmt.Errorf("failde to generate random bytes: %w", err)
	}

	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	return token, tokenHash, nil
}
