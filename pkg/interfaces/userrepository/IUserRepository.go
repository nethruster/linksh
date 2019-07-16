package userrepository

import (
	"github.com/nethruster/linksh/pkg/models"
)

//IUserRepository represents all the possible actions performed over the users
//The implementations of this interface will not be attached to an specific storage
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
	//If the user doesn't exists in the storage an error will be returned
	//This methods will permorn validations over the provided data
	//The data validations in this method can produce an ErrInvalidName or an ErrInvalidPassword
	Update(user models.User) error
	//Delete deletes an user from the storage
	Delete(id string) error
}
