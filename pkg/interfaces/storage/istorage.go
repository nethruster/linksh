package istorage

import (
	"github.com/nethruster/linksh/pkg/interfaces/user_repository"
	"github.com/nethruster/linksh/pkg/models"
)

//IStorage represents the storage functionality.
//No method of this interface will perform data validations as it's a job of the repositories, although it will check for uniqueness in the storage
type IStorage interface {
	//User related methods

	//SaveUser save the user in the storage
	//If there is a conflicting unique field this method will return an AlreadyExistsError
	SaveUser(user models.User) error
	//GetUser returns the user with specified ID from the storage
	//If the user does not exists in the storage an NotFoundError would be returned
	GetUser(id string) (models.User, error)
	//GetUserByName returns the user with specified name from the storage
	//If the user does not exists in the storage an NotFoundError would be returned
	GetUserByName(name string) (models.User, error)
	//ListUsers list the users in the storage with a limit and an offset
	//If the limit is set to 0, no limit will be established, the same applies to the offset
	ListUsers(limit, offset uint) ([]models.User, error)
	//UpdateUser replaces the values of the user in the storage with the non empty ones of the provided user
	//If the user does not exists in the storage a NotFoundError would be returned
	//If there is a conflicting unique field this method will return an AlreadyExistsError
	UpdateUser(user user_repository.UpdatePayload) error
	//DeleteUser deletes the user specified user from the storage
	//If the user does not exists in the storage an NotFoundError would be returned
	DeleteUser(id string) error

	//Link related methods

	//SaveLink save the link in the storage
	//If there is a conflicting unique field this method will return an AlreadyExistsError
	SaveLink(link models.Link) error
	//GetLink returns the link with specified ID from the storage
	//If the link does not exists in the storage an NotFoundError would be returned
	GetLink(id string) (models.Link, error)
	//ListLinks list the links in the storage with a limit and an offset
	//if the ownerID is not empty the search would be limited to the ones owned by the specified user
	//If the limit is set to 0, no limit will be established, the same applies to the offset
	ListLinks(ownerID string, limit, offset uint) ([]models.Link, error)
	//UpdateLinkContent replaces the values of the user in the storage with the non empty ones of the provided user
	//If the link does not exists in the storage an NotFoundError would be returned
	//If there is a conflicting unique field this method will return an AlreadyExistsError
	UpdateLinkContent(id, content string) error
	//DeleteLink deletes the link specified user from the storage
	//If the link does not exists in the storage an NotFoundError would be returned
	DeleteLink(id string) error
	//IncreaseLinkHitCount increases the hits number of a link in the storage
	//If the user does not exists in the storage an NotFoundError would be returned
	IncreaseLinkHitCount(id string) error

	// Session related methods

	// SaveSession saves the session into the storage
	// In case of unique field conflict this method will return an AlreadyExistsError
	SaveSession(session models.Session) error
	// GetSession gets the session with the specified ID from the storage
	// If the session does not exists in the storage a NotFoundError will be returned
	GetSession(id string) (models.Session, error)
	// ListSessions lists the sessions in the storage
	// if the ownerID is not empty the search will be limited to the ones owned by the specified user
	// if the limit is set to 0, no limit will be established, the same applies to the offset
	ListSessions(ownerID string, limit, offset uint) ([]models.Session, error)
	// UpdateSessionToken updates the lastToken of a session in the storage
	// if the session does not exists in the storage an NotFoundError would be returned
	UpdateSessionToken(id string, tokenID string) error
	// DeleteSession deletes a session
	// If the session does not exists in the storage a NotFoundError will be returned
	DeleteSession(id string) error
}
