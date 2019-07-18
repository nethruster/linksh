package repositories

import (
	"errors"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/nethruster/linksh/pkg/interfaces/linkrepository"
	sto "github.com/nethruster/linksh/pkg/interfaces/storage"
	"github.com/nethruster/linksh/pkg/models"
)

//LinkRepository implements ILinkRepository
type LinkRepository struct {
	Storage sto.IStorage
}

//Create creates a link and save it to the storage
//This methods will permorn validations over the provided data
//If the id is left blank, a random one would be assigned
//The data validations in this method can produce an ErrInvalidID or an ErrInvalidContent
func (lr *LinkRepository) Create(id, content, ownerID string) (link models.Link, err error) {
	mustGenerateID := id == ""
	if mustGenerateID {
		id, err = generateLinkID()
	} else {
		err = validateID(id)
	}
	if err != nil {
		return
	}
	err = validateContent(content)
	if err != nil {
		return
	}
	link = models.Link{
		ID:      id,
		Content: content,
		OwnerID: ownerID,
	}

	err = lr.Storage.SaveLink(link)
	return
}

//GetLink returns the link with specified ID from the storage
//If the user does not exists in the storage an NotFoundError would be returned
func (lr *LinkRepository) GetLink(id string) (models.Link, error) {
	panic(errors.New("Not implemented"))
}

//List lits the users
//If limit is set to 0, no limit will be established, the same happens to the offset
//if the ownerID is not empty the search would be limited to the owned owned by the specified user
func (lr *LinkRepository) List(ownerID string, limit, offset uint) ([]models.Link, error) {
	panic(errors.New("Not implemented"))
}

//UpdateContent replaces  the content of an existing link
//If the user doesn't exists in the Link an error would be returned
//This methods will permorn validations over the provided data
//The data validations in this method can produce an ErrInvalidContent
func (lr *LinkRepository) UpdateContent(id, content string) {
	panic(errors.New("Not implemented"))
}

//Delete deletes a link from the storage
func (lr *LinkRepository) Delete(id string) error {
	panic(errors.New("Not implemented"))
}

//GetLinkByUser returns the link with specified ID from the storage
//If the user does not exists in the storage an NotFoundError would be returned
//The requester must own the link or be an admin to perform this action
func (lr *LinkRepository) GetLinkByUser(requesterID, id string) models.Link {
	panic(errors.New("Not implemented"))
}

//ListByUser lits the users
//If limit is set to 0, no limit will be established, the same happens to the offset
//if the ownerID is not empty the search would be limited to the owned owned by the specified user
//The requester must be the owner of the links or an admin to perform this action
func (lr *LinkRepository) ListByUser(requesterID, ownerID string, limit, offset uint) ([]models.Link, error) {
	panic(errors.New("Not implemented"))
}

//IncreaseHitCount increases the hits number of a link in the storage
//If the user does not exists in the storage an NotFoundError would be returned
func (lr *LinkRepository) IncreaseHitCount(id string) error {
	panic(errors.New("Not implemented"))
}

func validateID(id string) error {
	if length := len(id); length == 0 || length > 100 {
		return linkrepository.ErrInvalidID
	}
	return nil
}

func validateContent(content string) error {
	if length := len(content); length == 0 || length > 2000 {
		return linkrepository.ErrInvalidContent
	}
	return nil
}

func generateLinkID() (string, error) {
	return gonanoid.Nanoid(7)
}
