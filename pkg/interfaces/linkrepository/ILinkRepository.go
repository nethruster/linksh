package linkrepository

import "github.com/nethruster/linksh/pkg/models"

//ILinkRepository represents all the possible actions performed over the links
//The implementations of this interface will not be attached to an specific storage
//The methods with the suffix 'ByUser' will only be perform if the requester has enough privileges, if not an pkg/interfaces/userrepository.ErrForbidden would be returned
type ILinkRepository interface {
	//Create creates a link and save it to the storage
	//This methods will permorn validations over the provided data
	//If the id is left blank, a random one would be assigned
	//The data validations in this method can produce an ErrInvalidID or an ErrInvalidContent
	Create(id, content, ownerID string) (models.Link, error)
	//Get returns the link with specified ID from the storage
	//If the link does not exists in the storage an error pkg/interfaces/storage.NotFoundError would be returned
	Get(id string) (models.Link, error)
	//GetContentAndIncreaseHitCount return the link content and increases the hits number of a link in the storage
	//If the link does not exists in the storage an error pkg/interfaces/storage.NotFoundError would be returned
	GetContentAndIncreaseHitCount(id string) (string, error)
	//List lits the users
	//If limit is set to 0, no limit will be established, the same happens to the offset
	//if the ownerID is not empty the search would be limited to the ones owned by the specified user
	List(ownerID string, limit, offset uint) ([]models.Link, error)
	//UpdateContent replaces  the content of an existing link
	//If the link doesn't exists in the Link an error would be returned
	//This methods will permorn validations over the provided data
	//The data validations in this method can produce an ErrInvalidContent
	UpdateContent(id, content string) error
	//Delete deletes a link from the storage
	//If the link does not exists in the storage an error pkg/interfaces/storage.NotFoundError would be returned
	Delete(id string) error
	//IncreaseHitCount increases the hits number of a link in the storage
	//If the link does not exists in the storage an error pkg/interfaces/storage.NotFoundError would be returned
	IncreaseHitCount(id string) error

	//GetByUser returns the link with specified ID from the storage
	//If the link does not exists in the storage an error pkg/interfaces/storage.NotFoundError would be returned
	//The requester must own the link or be an admin to perform this action
	GetByUser(requesterID, id string) (models.Link, error)
	//ListByUser lits the users
	//If limit is set to 0, no limit will be established, the same happens to the offset
	//if the ownerID is not empty the search would be limited to the ones owned by the specified user
	//The requester must be the owner of the links or an admin to perform this action
	ListByUser(requesterID, ownerID string, limit, offset uint) ([]models.Link, error)
	//UpdateContentByUser replaces  the content of an existing link
	//If the link doesn't exists in the Link an error would be returned
	//This methods will permorn validations over the provided data
	//The data validations in this method can produce an ErrInvalidContent
	//The requester must own the link or be an admin to perform this action
	UpdateContentByUser(requesterID, id, content string) error
	//DeleteByUser deletes a link from the storage
	//If the link does not exists in the storage an error pkg/interfaces/storage.NotFoundError would be returned
	//The requester must own the link or be an admin to perform this action
	DeleteByUser(requesterID, id string) error
}
