package store

import (
	"database/sql"
	"time"

	"github.com/makhammatovb/Articles/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int64, scope string) error
	GetToken(hash []byte) (*tokens.Token, error)
}

func (t *PostgresTokenStore) CreateNewToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4);
	`
	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	if err != nil {
		return err
	}
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int64, scope string) error {
	query := `
	DELETE FROM tokens WHERE user_id = $1 AND scope = $2;
	`
	_, err := t.db.Exec(query, userID, scope)
	if err != nil {
		return err
	}
	return err
}

// bu xato (token hash byte qabul qiladigon bo'ldi)
func (t *PostgresTokenStore) GetToken(hash []byte) (*tokens.Token, error) {
	token := &tokens.Token{}
	query := `
	SELECT hash, user_id, expiry, scope FROM tokens WHERE hash = $1;
	`
	err := t.db.QueryRow(query, hash).Scan(&token.Hash, &token.UserID, &token.Expiry, &token.Scope)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	token.Hash = hash
	return token, nil
}

