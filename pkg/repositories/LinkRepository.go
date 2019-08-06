package repositories

import (
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

//Get returns the link with specified ID from the storage
//If the link does not exists in the storage an NotFoundError would be returned
func (lr *LinkRepository) Get(id string) (models.Link, error) {
	if id == "" {
		return models.Link{}, linkrepository.ErrInvalidID
	}

	return lr.Storage.GetLink(id)
}
//GetContentAndIncreaseHitCount return the link content and increases the hits number of a link in the storage
//If the link does not exists in the storage an NotFoundError would be returned
func (lr *LinkRepository) GetContentAndIncreaseHitCount(id string) (string, error) {
	link, err := lr.Get(id)
	if err != nil {
		return "", err
	}
	if err = lr.IncreaseHitCount(id); err != nil {
		return "", err
	}

	return link.Content, nil
}
//List lits the users
//If limit is set to 0, no limit will be established, the same happens to the offset
//if the ownerID is not empty the search would be limited to the owned by the specified user
func (lr *LinkRepository) List(ownerID string, limit, offset uint) ([]models.Link, error) {
	return lr.Storage.ListLinks(ownerID, limit, offset)
}
//UpdateContent replaces  the content of an existing link
//If the link doesn't exists in the Link an error would be returned
//This methods will permorn validations over the provided data
//The data validations in this method can produce an ErrInvalidContent
func (lr *LinkRepository) UpdateContent(id, content string) error {
	err := validateContent(content)
	if err != nil {
		return err
	}
	err = lr.Storage.UpdateLinkContent(id, content)
	if err != nil {
		return err
	}

	return nil
}
//Delete deletes a link from the storage
//If the link does not exists in the storage an NotFoundError would be returned
func (lr *LinkRepository) Delete(id string) error {
	return lr.Storage.DeleteLink(id)
}
//IncreaseHitCount increases the hits number of a link in the storage
//If the link does not exists in the storage an NotFoundError would be returned
func (lr *LinkRepository) IncreaseHitCount(id string) error {
	return lr.Storage.IncreaseLinkHitCount(id)
}

//GetByUser returns the link with specified ID from the storage
//If the link does not exists in the storage an NotFoundError would be returned
//The requester must own the link or be an admin to perform this action
func (lr *LinkRepository) GetByUser(requesterID, id string)  (models.Link, error) {
	link, err := lr.Get(id)
	if err != nil {
		return link, err
	}
	if link.OwnerID != requesterID {
		if err = checkIfRequesterIsAdmin(lr.Storage, requesterID); err != nil {
			return link, err
		}
	}

	return link, nil
}
//ListByUser lits the users
//If limit is set to 0, no limit will be established, the same happens to the offset
//if the ownerID is not empty the search would be limited to the owned owned by the specified user
//The requester must be the owner of the links or an admin to perform this action
func (lr *LinkRepository) ListByUser(requesterID, ownerID string, limit, offset uint) ([]models.Link, error) {
	var err error
	if requesterID != ownerID {
		if err = checkIfRequesterIsAdmin(lr.Storage, requesterID); err != nil {
			return nil, err
		}
	}

	return lr.List(ownerID, limit, offset)
}
//UpdateContentByUser replaces  the content of an existing link
//If the link doesn't exists in the Link an error would be returned
//This methods will permorn validations over the provided data
//The data validations in this method can produce an ErrInvalidContent
//The requester must own the link or be an admin to perform this action
func (lr *LinkRepository) UpdateContentByUser(requesterID, id, content string) error {
	link, err := lr.Get(id)
	if err != nil {
		return err
	}
	if link.OwnerID != requesterID {
		if err = checkIfRequesterIsAdmin(lr.Storage, requesterID); err != nil {
			return err
		}
	}

	return lr.UpdateContent(id, content)
}

//DeleteByUser deletes a link from the storage
//If the link does not exists in the storage an NotFoundError would be returned
//The requester must own the link or be an admin to perform this action
func (lr *LinkRepository) DeleteByUser(requesterID, id string) error {
	link, err := lr.Get(id)
	if err != nil {
		return err
	}
	if link.OwnerID != requesterID {
		if err = checkIfRequesterIsAdmin(lr.Storage, requesterID); err != nil {
			return err
		}
	}

	return lr.Delete(id)
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
