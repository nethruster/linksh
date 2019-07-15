package interfaces

import (
	"github.com/nethruster/linksh/pkg/models"
)

//IUserRepository represents all the possible actions performed over the users
//The implementations of this interface will not be attached to an specific storage
type IUserRepository interface {
	CheckLoginCredentials(username, password string) (bool, error)
	Create(username, password string) (models.User, error)
	Get(id string) (models.User, error)
	List(limit, offset uint) ([]models.User, error)
	Update(user models.User) error
	Delete(userID string) error
}
