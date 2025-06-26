package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateInvitationToken generates a secure random token for workspace invitations
func GenerateInvitationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GetInvitationExpiration returns the expiration time for invitations (7 days from now)
func GetInvitationExpiration() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}

// IsInvitationExpired checks if an invitation has expired
func IsInvitationExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}