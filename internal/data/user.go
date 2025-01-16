package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID             int64          `json:"id"`
	Username       string         `json:"username"`
	Password       string         `json:"-"`
	Role           string         `json:"role"`
	Token          sql.NullString `json:"access_token"`
	TokenExpiredAt int64          `json:"token_expired_at"`
	Avatar         sql.NullString `json:"avatar"`
}

var ANONYMOUS_USER = &User{}

func (u *User) IsAnonymous() bool {
	return u == ANONYMOUS_USER
}

func (m UserModel) Insert(user *User) error {
	stmt := `INSERT INTO users (username, password, role) VALUES(?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, stmt, user.Username, user.Password, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) UpdateAvatar(username string, avatar string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE users SET avatar = ? WHERE username = ?"
	_, err := m.DB.ExecContext(ctx, query, avatar, username)
	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) GetByUsername(username string) (*User, error) {
	if username == "" {
		return nil, ErrRecordNotFound
	}

	query := "SELECT username, password, role, token, avatar FROM users WHERE username = ?"

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, username).Scan(
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Token,
		&user.Avatar,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
