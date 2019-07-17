package linkrepository

import "github.com/nethruster/linksh/pkg/models"

//ILinkRepository represents all the possible actions performed over the links
//The implementations of this interface will not be attached to an specific storage
//The methods with the suffix 'ByUser' will only be perform if the requester has enough privileges, if not an ErrForbidden would be returned
type ILinkRepository interface {
	//Create creates a link and save it to the storage
	//This methods will permorn validations over the provided data
	//If the id is left blank, a random one would be assigned
	//The data validations in this method can produce an ErrInvalidID or an ErrInvalidContent
	Create(id, content, ownerID string) (models.Link, error)
	//GetLink returns the link with specified ID from the storage
	//If the user does not exists in the storage an NotFoundError would be returned
	GetLink(id string) (models.Link, error)
	//List lits the users
	//If limit is set to 0, no limit will be established, the same happens to the offset
	//if the ownerID is not empty the search would be limited to the owned owned by the specified user
	List(ownerID string, limit, offset uint) ([]models.Link, error)
	//UpdateContent replaces  the content of an existing link
	//If the user doesn't exists in the Link an error would be returned
	//This methods will permorn validations over the provided data
	//The data validations in this method can produce an ErrInvalidContent
	UpdateContent(id, content string)
	//Delete deletes a link from the storage
	Delete(id string) error
	//GetLinkByUser returns the link with specified ID from the storage
	//If the user does not exists in the storage an NotFoundError would be returned
	//The requester must own the link or be an admin to perform this action
	GetLinkByUser(requesterID, id string) models.Link
	//List lits the users
	//If limit is set to 0, no limit will be established, the same happens to the offset
	//if the ownerID is not empty the search would be limited to the owned owned by the specified user
	//The requester must be the owner of the links or an admin to perform this action
	ListByUser(requesterID, ownerID string, limit, offset uint) ([]models.Link, error)
	//IncreaseHitCount increases the hits number of a link in the storage
	//If the user does not exists in the storage an NotFoundError would be returned
	IncreaseHitCount(id string) error
}
