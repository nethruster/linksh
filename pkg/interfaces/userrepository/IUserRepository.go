package userrepository

import (
	"github.com/nethruster/linksh/pkg/models"
)

//IUserRepository represents all the possible actions performed over the users
//The implementations of this interface will not be attached to an specific storage
//The methods with the suffix 'ByUser' will only be perform if the requester has enough privileges, if not an ErrForbidden would be returned
type IUserRepository interface {
	//CheckLoginCredentials checks if the provided credentials are valid to perform a login
	CheckLoginCredentials(name string, password []byte) (bool, error)
	//Create creates an user and save it to the storage
	//This methods will permorn validations over the provided data
	//The data validations in this method can produce an ErrInvalidName or an ErrInvalidPassword
	Create(name string, password []byte, isAdmin bool) (models.User, error)
	//Get returns an user from the storage
	Get(id string) (models.User, error)
	//List lits the users
	//If limit is set to 0, no limit will be established
	List(limit, offset uint) ([]models.User, error)
	//Update replaces the values of the user in the storage with the values of the user provided by parameter
	//If the user doesn't exists in the storage an error would be returned
	//This methods will permorn validations over the provided data
	//The data validations in this method can produce an ErrInvalidName or an ErrInvalidPassword
	Update(user UpdatePayload) error
	//Delete deletes an user from the storage
	Delete(id string) error

	//CreateByUser creates an user and save it to the storage
	//This methods will permorn validations over the provided data
	//The data validations in this method can produce an ErrInvalidName or an ErrInvalidPassword
	//The requester must be an admin to perform this action
	CreateByUser(requesterID string, name string, password []byte, isAdmin bool) (models.User, error)
	//GetByUser returns an user from the storage
	//The requester must only request information about himself or be an admin to perform this action
	GetByUser(requesterID, id string) (models.User, error)
	//ListByUser lits the users
	//If limit is set to 0, no limit will be established
	//The requester must be an admin to perform this action
	ListByUser(requesterID string, limit, offset uint) ([]models.User, error)
	//UpdateByUser replaces the values of the user in the storage with the values of the user provided by parameter
	//If the user doesn't exists in the storage an error would be returned
	//This methods will permorn validations over the provided data
	//The data validations in this method can produce an ErrInvalidName or an ErrInvalidPassword
	//The requestor can only modify information about himself or otherwise be an admin to perform this action. The isAdmin property can only be changed by other admins.
	UpdateByUser(requesterID string, user UpdatePayload) error
	//DeleteByUser deletes an user from the storage
	//The requester must only delete himself or be an admin to perform this action
	DeleteByUser(id string) error
}

//UpdatePayload is a clone of models.User with nullable fields used to perform operations as only the not null fields will be the ones updated
type UpdatePayload struct {
	ID       string  `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Password []byte  `json:"password,omitempty"`
	IsAdmin  *bool   `json:"isAdmin"`
}
