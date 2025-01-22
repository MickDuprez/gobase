package auth

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
}

func (a *AuthDB) CreateUser(email, password, name string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := a.db.Exec(
		`INSERT INTO users (email, password_hash, name) VALUES (?, ?, ?)`,
		email, string(hash), name,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		CreatedAt:    time.Now(),
	}, nil
}

func (a *AuthDB) GetUserByEmail(email string) (*User, error) {
	var user User
	err := a.db.QueryRow(
		`SELECT id, email, password_hash, name, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *AuthDB) ValidateUser(email, password string) (*User, error) {
	user, err := a.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (a *AuthDB) GetUserByID(id int64) (*User, error) {
	var user User
	err := a.db.QueryRow(
		`SELECT id, email, password_hash, name, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
