package mongo

import (
	"context"
	"github.com/nethruster/linksh/pkg/interfaces/user_repository"
	"github.com/nethruster/linksh/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoStorage struct {
	client *mongo.Client
	DefaultTimeout time.Duration
}

func (sto *MongoStorage) Connect(connectionString string) error {
	options := mongoOptions.Client().ApplyURI(connectionString)
	*options.AppName = "linksh"
	client, err := mongo.NewClient(options)
	if err != nil {
		return err
	}

	err = client.Connect(sto.newTimeoutContext())
	if err != nil {
		return err
	}
	
	sto.client = client
	return nil
}

func (sto *MongoStorage) newTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(),sto.DefaultTimeout)
	return ctx
}


//User related methods

//SaveUser save the user in the storage
//If there is a conflicting unique field this method will return an AlreadyExistsError
func (sto *MongoStorage) SaveUser(user models.User) error {

}

//GetUser returns the user with specified ID from the storage
//If the user does not exists in the storage an NotFoundError would be returned
func (sto *MongoStorage) GetUser(id string) (models.User, error) {

}

//GetUserByName returns the user with specified name from the storage
//If the user does not exists in the storage an NotFoundError would be returned
func (sto *MongoStorage) GetUserByName(name string) (models.User, error) {

}

//ListUsers list the users in the storage with a limit and an offset
//If the limit is set to 0, no limit will be established, the same applies to the offset
func (sto *MongoStorage) ListUsers(limit, offset uint) ([]models.User, error) {

}

//UpdateUser replaces the values of the user in the storage with the non empty ones of the provided user
//If the user does not exists in the storage a NotFoundError would be returned
//If there is a conflicting unique field this method will return an AlreadyExistsError
func (sto *MongoStorage) UpdateUser(user user_repository.UpdatePayload) error {

}

//DeleteUser deletes the user specified user from the storage
//If the user does not exists in the storage an NotFoundError would be returned
func (sto *MongoStorage) DeleteUser(id string) error {

}

//Link related methods

//SaveLink save the link in the storage
//If there is a conflicting unique field this method will return an AlreadyExistsError
func (sto *MongoStorage) SaveLink(link models.Link) error {

}

//GetLink returns the link with specified ID from the storage
//If the link does not exists in the storage an NotFoundError would be returned
func (sto *MongoStorage) GetLink(id string) (models.Link, error) {

}

//ListLinks list the links in the storage with a limit and an offset
//if the ownerID is not empty the search would be limited to the ones owned by the specified user
//If the limit is set to 0, no limit will be established, the same applies to the offset
func (sto *MongoStorage) ListLinks(ownerID string, limit, offset uint) ([]models.Link, error) {

}

//UpdateUser replaces the values of the user in the storage with the non empty ones of the provided user
//If the link does not exists in the storage an NotFoundError would be returned
//If there is a conflicting unique field this method will return an AlreadyExistsError
func (sto *MongoStorage) UpdateLinkContent(id, content string) error {

}

//DeleteLink deletes the link specified user from the storage
//If the link does not exists in the storage an NotFoundError would be returned
func (sto *MongoStorage) DeleteLink(id string) error {

}

//IncreaseLinkHitCount increases the hits number of a link in the storage
//If the user does not exists in the storage an NotFoundError would be returned
func (sto *MongoStorage) IncreaseLinkHitCount(id string) error {

}
