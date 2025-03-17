package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    int64
	CreatedAt time.Time
	ExpiresAt time.Time
	Data      map[string]interface{}
}

// Helper methods for working with session data
func (s *Session) SetValue(key string, value interface{}) {
	if s.Data == nil {
		s.Data = make(map[string]interface{})
	}
	s.Data[key] = value
}

func (s *Session) GetValue(key string) interface{} {
	if s.Data == nil {
		return nil
	}
	return s.Data[key]
}

func (s *Session) GetInt(key string) (int64, bool) {
	val := s.GetValue(key)
	if val == nil {
		return 0, false
	}

	switch v := val.(type) {
	case int64:
		return v, true
	case float64:
		return int64(v), true
	case int:
		return int64(v), true
	default:
		return 0, false
	}
}

func (s *Session) GetString(key string) (string, bool) {
	val := s.GetValue(key)
	if val == nil {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (a *AuthDB) CreateSession(userID int64, duration time.Duration) (*Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        id,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
		Data:      make(map[string]interface{}),
	}

	data, err := json.Marshal(session.Data)
	if err != nil {
		return nil, err
	}

	_, err = a.db.Exec(
		`INSERT INTO sessions (id, user_id, expires_at, data) VALUES (?, ?, ?, ?)`,
		session.ID, session.UserID, session.ExpiresAt, string(data),
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (a *AuthDB) GetSession(id string) (*Session, error) {
	var session Session
	var dataStr string

	err := a.db.QueryRow(
		`SELECT id, user_id, created_at, expires_at, data FROM sessions WHERE id = ?`,
		id,
	).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt, &dataStr)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(dataStr), &session.Data); err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		a.DeleteSession(id)
		return nil, errors.New("session expired")
	}

	return &session, nil
}

func (a *AuthDB) DeleteSession(id string) error {
	_, err := a.db.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}
