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
	err = validateName(name)
	if err != nil {
		return
	}
	err = validatePassword(password)
	if err != nil {
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
	return ur.Storage.ListUsers(limit, offset)
}

//Update replaces the values of the user in the storage with the values of the user provided by parameter
//If the user doesn't exists in the storage an error will be returned
//This methods will permorn validations over the provided data
func (ur *UserRepository) Update(payload userrepository.UpdatePayload) (err error) {
	if payload.Name != nil {
		err = validateName(*payload.Name)
		if err != nil {
			return
		}
	}
	if payload.Password != nil {
		err = validatePassword(payload.Password)
		if err != nil {
			return
		}
	}

	return ur.Storage.UpdateUser(payload)
}

//Delete deletes an user from the storage
func (ur *UserRepository) Delete(id string) error {
	return ur.Storage.DeleteUser(id)
}

//CreateByUser creates an user and save it to the storage
//This methods will permorn validations over the provided data
//The data validations in this method can produce an ErrInvalidName or an ErrInvalidPassword
//The requester must be an admin to perform this action
func (ur *UserRepository) CreateByUser(requesterID string, name string, password []byte, isAdmin bool) (user models.User, err error) {
	err = ur.checkIfRequesterIsAdmin(requesterID)
	if err != nil {
		return
	}

	return ur.Create(name, password, isAdmin)
}

//GetByUser returns an user from the storage
//The requester must only request information about himself or be an admin to perform this action
func (ur *UserRepository) GetByUser(requesterID, id string) (user models.User, err error) {
	if requesterID != id {
		err = ur.checkIfRequesterIsAdmin(requesterID)
		if err != nil {
			return
		}
	}

	return ur.Get(id)
}

//ListByUser lits the users
//If limit is set to 0, no limit will be established
//The requester must be an admin to perform this action
func (ur *UserRepository) ListByUser(requesterID string, limit, offset uint) ([]models.User, error) {
	panic(errors.New("Not implemented"))
}

//UpdateByUser replaces the values of the user in the storage with the values of the user provided by parameter
//If the user doesn't exists in the storage an error would be returned
//This methods will permorn validations over the provided data
//The data validations in this method can produce an ErrInvalidName or an ErrInvalidPassword
//The requester must only modify information about himself or be an admin to perform this action
func (ur *UserRepository) UpdateByUser(requesterID string, user userrepository.UpdatePayload) error {
	panic(errors.New("Not implemented"))
}

//DeleteByUser deletes an user from the storage
//The requester must only delete himself or be an admin to perform this action
func (ur *UserRepository) DeleteByUser(id string) error {
	panic(errors.New("Not implemented"))
}

func (ur *UserRepository) checkIfRequesterIsAdmin(requesterID string) (err error) {
	requester, err := ur.Storage.GetUser(requesterID)
	if err != nil {
		err = errors.Errorf("Error checking the requester %w", err)
		return
	}

	if !requester.IsAdmin {
		err = userrepository.ErrForbidden
	}
	return
}

func generateUserID() (string, error) {
	return gonanoid.Nanoid()
}

func validateName(name string) error {
	if length := len(name); length == 0 || length > 100 {
		return userrepository.ErrInvalidName
	}
	return nil
}

func validatePassword(password []byte) error {
	if len(password) < 6 {
		return userrepository.ErrInvalidPassword
	}
	return nil
}
