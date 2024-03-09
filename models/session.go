package models

import "database/sql"

type Session struct {
	ID int
	// Token is only set when creating a new session. When look un a session
	// this will be left empty. as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token
	Token     string
	UserId    int
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	// TODO: Create the session token
	// TODO: Implement SessionService.Create
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: Implement SessionService.User
	return nil, nil
}
