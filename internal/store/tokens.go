package store

import (
	"database/sql"
	"time"

	"github.com/cykj40/beginner_go/internal/store/tokens"
)

type PostgresTokenStore struct {
	DB *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		DB: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int64, scope string) error
}

func (t *PostgresTokenStore) CreateNewToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`
	_, err := t.DB.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int64, scope string) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2
   `

	_, err := t.DB.Exec(query, scope, userID)
	return err
}
