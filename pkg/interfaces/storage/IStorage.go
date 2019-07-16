package istorage

import "github.com/nethruster/linksh/pkg/models"

//IStorage represents the storage functionality.
//No method of this interface will perform data validations as it's a job of the repositories, although it will check for uniqueness in the storage
type IStorage interface {
	//User related methos

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
	ListUsers(limit, offset uint) ([]models.User, error)
	//UpdateUser replaces the values of the user in the storage with the non empty ones of the provided user
	//If the user does not exists in the storage an NotFoundError would be returned
	//If there is a conflicting unique field this method will return an AlreadyExistsError
	UpdateUser(user models.User) error
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
	ListLinks(limit, offset uint) ([]models.Link, error)
	//UpdateUser replaces the values of the user in the storage with the non empty ones of the provided user
	//If the link does not exists in the storage an NotFoundError would be returned
	//If there is a conflicting unique field this method will return an AlreadyExistsError
	UpdateLinkContent(id, content string) error
	//DeleteLink deletes the link specified user from the storage
	//If the link does not exists in the storage an NotFoundError would be returned
	DeleteLink(id string) error
}
