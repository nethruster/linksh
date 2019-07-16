package repositories

import (
	gonanoid "github.com/matoous/go-nanoid"
	sto "github.com/nethruster/linksh/pkg/interfaces/storage"
	"github.com/nethruster/linksh/pkg/interfaces/userrepository"
	"github.com/nethruster/linksh/pkg/models"
	"golang.org/x/crypto/bcrypt"
	errors "golang.org/x/xerrors"
)

//UserRepository implements IUserRepository
type UserRepository struct {
	Storage sto.IStorage
}

//CheckLoginCredentials checks if the provided credentials are valid to perform a login
func (ur *UserRepository) CheckLoginCredentials(name string, password []byte) (bool, error) {
	user, err := ur.Storage.GetUserByName(name)
	if err != nil {
		var notFoundError *sto.NotFoundError
		if errors.As(err, notFoundError) {
			return false, nil
		}
		return false, errors.Errorf("Error checking the login credentials %w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, password)

	return err == nil, nil
}

//Create creates an user and save it to the storage
//This methods will permorn validations over the provided data
func (ur *UserRepository) Create(name string, password []byte, isAdmin bool) (user models.User, err error) {
	//Check if the user's name is valid
	if length := len(name); length == 0 || length > 100 {
		err = userrepository.ErrInvalidName
		return
	}
	//Check if the user's password is valid
	if len(password) < 6 {
		err = userrepository.ErrInvalidPassword
		return
	}

	pwHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	id, err := generateUserID()
	if err != nil {
		return
	}
	user = models.User{
		ID:       id,
		Name:     name,
		Password: pwHash,
		IsAdmin:  isAdmin,
	}

	err = ur.Storage.SaveUser(user)
	return
}

//Get returns an user from the storage
func (ur *UserRepository) Get(id string) (models.User, error) {
	return ur.Storage.GetUser(id)
}

//List lits the users
//If limit is set to 0, no limit will be established
func (ur *UserRepository) List(limit, offset uint) ([]models.User, error) {
	panic(errors.New("*UserRepository.List not implemented"))
}

//Update replaces the values of the user in the storage with the values of the user provided by parameter
//If the user doesn't exists in the storage an error will be returned
//This methods will permorn validations over the provided data
func (ur *UserRepository) Update(u1 models.User) error {
	panic(errors.New("*UserRepository.Update not implemented"))
}

//Delete deletes an user from the storage
func (ur *UserRepository) Delete(id string) error {
	panic(errors.New("*UserRepository.Delete not implemented"))
}

func generateUserID() (string, error) {
	return gonanoid.Nanoid()
}
