package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"
)

var expirationTime time.Duration = 30 * 24 * time.Hour

// PSQLProjectStore implements ProjectStore interface using an SQL database
type PSQLTokenStore struct {
	db *sql.DB
}

// NewSQLProjectStore creates a new SQLProjectStore
func NewPSQLTokenStore(db *sql.DB) *PSQLTokenStore {
	return &PSQLTokenStore{db: db}
}

func (store *PSQLTokenStore) ValidateRefreshToken(token string) (int, error) {
	var userID int
	err := store.db.QueryRow("SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW() AND revoked = FALSE", token).Scan(&userID)
	return userID, err
}

func (store *PSQLTokenStore) InvalidateRefreshToken(token string) error {
	query := "UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1"
	result, err := store.db.Exec(query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no token invalidated, token not found or already revoked")
	}
	return nil
}

func (store *PSQLTokenStore) GenerateRefreshToken(userID int) (string, error) {
	token, err := generateRandomToken()
	if err != nil {
		return "", err
	}

	_, err = store.db.Exec("INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)",
		token, userID, time.Now().Add(expirationTime)) // 30 days expiration
	return token, err
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

/*
func validateRefreshToken(token string) error {
	var userID int
err := db.QueryRow("SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW() AND revoked = FALSE", token).Scan(&userID)
if err != nil {
    log.Printf("Invalid token: %v", err)
    return false
}

func generateRefreshToken() error {
	token, err := generateRandomToken()
if err != nil {
    log.Fatalf("Error generating token: %v", err)
}

_, err = db.Exec("INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)",
    token, userID, time.Now().Add(30*24*time.Hour)) // 30 days expiration
	return err
}*/
