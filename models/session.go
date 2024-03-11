package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/sajadjalilian/lenslocked/rand"
)

type Session struct {
	ID int
	// Token is only set when creating a new session. When look un a session
	// this will be left empty. as we only store the hash of a session token
	// in our database, and we cannot reverse it into a raw token
	Token     string
	UserId    int
	TokenHash string
}

const (
	// The minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserId:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;`, session.UserId, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1,$2)
		RETURNING id;`, session.UserId, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT user_id
		FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	row = ss.DB.QueryRow(`
		SELECT email, password_hash
		FROM users WHERE id = $1;`, user.ID)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHask := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHask[:])
}
