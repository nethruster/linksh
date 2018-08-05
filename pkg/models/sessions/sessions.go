package sessions

import (
	"fmt"
)

//Session is the datatype which contains all the relevant information about a session
type Session struct {
	ID        string
	OwnerID   string
	ExpiresOn int64
	CreatedAt int64
}

//Encode return a string with the basic session data in the spected format for the its usage as a auth method
func (session *Session) Encode() string {
	return fmt.Sprintf("%v%v", session.ID, session.OwnerID)
}

//SessionProvider is the interface for a valid session storage
type SessionProvider interface {
	Add(session Session) error
	Get(id string) (Session, error)
	GetByOwnerID(ownerID string) (map[string]Session, error)
	Update(session Session) error
	Delete(id string) error
	GC() error
}
