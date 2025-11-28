package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

const (
	ScopeAuth = "authentication"
	ScopeResetPassword = "reset-password"
)

func GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	// creates a byte slice of 32 empty bytes
	emptyBytes := make([]byte, 32)
	// fills the slice with secure random data
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}
	// encodes the 32 random bytes into a Base32 String
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}
