package session_repository

import (
	"github.com/nethruster/linksh/pkg/models"
)

//ISessionRepository represents all the possible actions performed over the sessions
type ISessionRepository interface {
	// Create creates a session and save it to the storage
	// if the expire date is set 0 the session will not expire
	Create(userID string, expireDate int64) (models.Session, error)
	// List the sessions
	// If the limit is set to 0, no limit will be established, the same applies to the offset
	// if the userID is not empty the search will be limited to the ones with the specified userID
	List(userID string, limit, offset uint) ([]models.Session, error)
	// ValidateToken validates a JWT
	// Return the userID and an error if necessary
	// If the token is invalid ErrInvalidToken will be returned
	// If the token is valid but expired ErrExpiredToken will be returned
	ValidateToken(sessionToken string) (string, error)
	// GenerateToken generates a JWT
	// If the session does not exists in the storage an error pkg/interfaces/storage.NotFoundError will be returned
	GenerateToken(sessionID string) (string, error)
	// ValidateAndRenew validates a JWT, if it is valid but the token has expired (and the session does not) and it is the last issued token for the session a new one would be generated
	// if the token has expired and it is not the last issued token for the session, the session will be deleted
	// If the session does not exists in the storage an error pkg/interfaces/storage.NotFoundError will be returned
	ValidateAndRenew(sessionToken string) (string, error)
	// Delete deletes a session
	// If the session does not exists in the storage an error pkg/interfaces/storage.NotFoundError will be returned
	Delete(id string) error
	// Delete deletes a session
	// The requester must own the session to perform this action, otherwise an pkg/interfaces/user_repository.ErrForbidden will be returned
	DeleteByUser(userID, id string) error
}
