package repositories

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	sto "github.com/nethruster/linksh/pkg/interfaces/storage"
	"github.com/nethruster/linksh/pkg/models"
	"time"
)

type SessionRepository struct {
	Storage sto.IStorage
}

func (sr *SessionRepository) Create(userID string, expireDate int64) (models.Session, error) {
	if userID == "" {
		return models.Session{}, fmt.Errorf("UserID can not be empty")
	}


	id, err := generateSessionID()
	if err != nil {
		return models.Session{}, fmt.Errorf("error creating session ID%w", err)
	}
	session := models.Session{
		ID:         id,
		UserID:     userID,
		LastToken:  "",
		CreatedAt: time.Now().Unix(),
		ExpireDate: expireDate,
	}

	err = sr.Storage.SaveSession(session)
	if err != nil {
		return session, fmt.Errorf("error creating the session%w", err)
	}

	return session, nil
}

func (sr *SessionRepository) List(userID string, limit, offset uint) ([]models.Session, error) {
	panic("not implemented")
}

func (sr *SessionRepository) ValidateToken(sessionToken string) (string, error) {
	panic("not implemented")
}

func (sr *SessionRepository) GenerateToken(sessionID string) (string, error) {
	panic("not implemented")
}

func (sr *SessionRepository) ValidateAndRenew(sessionToken string) (string, error) {
	panic("not implemented")
}

func (sr *SessionRepository) Delete(id string) error {
	panic("not implemented")
}


func (sr *SessionRepository) DeleteByUser(userID, id string) error {
	panic("not implemented")
}

func generateSessionID() (string, error) {
	return gonanoid.Nanoid()
}
