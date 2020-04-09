package mongo

import (
	"context"
	"fmt"
	istorage "github.com/nethruster/linksh/pkg/interfaces/storage"
	"github.com/nethruster/linksh/pkg/interfaces/user_repository"
	"github.com/nethruster/linksh/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	userCollectionName = "users"
	linksCollectionName = "links"
)

var (
	appName = "linksh"
)

type Storage struct {
	client *mongo.Client
	databaseName string
	DefaultTimeout time.Duration
}

func New(connectionString string, databaseName string, defaultTimeout time.Duration) (*Storage, error) {
	options := mongoOptions.Client().ApplyURI(connectionString)
	options.AppName = &appName
	client, err := mongo.NewClient(options)
	if err != nil {
		return nil, err
	}
	sto := Storage{
		client:         client,
		databaseName:   databaseName,
		DefaultTimeout: defaultTimeout,
	}

	err = client.Connect(sto.newTimeoutContext())
	if err != nil {
		return nil, err
	}

	return &sto, nil
}

func (sto *Storage) newTimeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(),sto.DefaultTimeout)
	return ctx
}

func (sto *Storage) close() error {
	return sto.client.Disconnect(sto.newTimeoutContext())
}

func (sto *Storage) db() *mongo.Database {
	return sto.client.Database(sto.databaseName)
}


//User related methods

func (sto *Storage) SaveUser(user models.User) error {
	_, err := sto.db().Collection(userCollectionName).InsertOne(sto.newTimeoutContext(), &user)
	if err != nil {
		//TODO test for already exists error
		return err
	}

	return nil
}

func (sto *Storage) GetUser(id string) (user models.User, err error) {
	result := sto.db().Collection(userCollectionName).FindOne(sto.newTimeoutContext(), bson.M{"_id": id})
	err = result.Err()

	if err == mongo.ErrNoDocuments {
		err = istorage.NewNotFoundError("users", "ID", id)
	}
	if  err != nil {
		err = fmt.Errorf("error searching user with id \"%s\":%w", id, err)
		return
	}
	if err = result.Decode(&user); err != nil {
		err = fmt.Errorf("error deconding user with id \"%s\":%w", id, err)
		return
	}
	return
}

func (sto *Storage) GetUserByName(name string) (user models.User, err error) {
	result := sto.db().Collection(userCollectionName).FindOne(sto.newTimeoutContext(), bson.M{"name": name})
	err = result.Err()

	if err == mongo.ErrNoDocuments {
		err = istorage.NewNotFoundError("users", "Name", name)
	}
	if err != nil {
		err = fmt.Errorf("error searching user with name \"%s\":%w", name, err)
		return
	}
	if err = result.Decode(&user); err != nil {
		err = fmt.Errorf("error deconding user with name \"%s\":%w", name, err)
		return
	}
	return
}

func (sto *Storage) ListUsers(limit, offset uint) ([]models.User, error) {
	options := mongoOptions.Find()
	options.SetSort(bson.D{{"name", -1}})
	if limit != 0 {
		options.SetLimit(int64(limit))
	}
	if offset != 0 {
		options.SetSkip(int64(offset))
	}
	ctx := sto.newTimeoutContext()
	cursor, err := sto.db().Collection(userCollectionName).Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []models.User
	err = cursor.All(ctx, &users)
	return users, err
}

func (sto *Storage) UpdateUser(user user_repository.UpdatePayload) error {
	var set bson.D
	if user.ID == "" {
		return istorage.NewNotFoundError("users", "id", "")
	}
	if user.Name != nil {
		set = append(set, bson.E{"name", user.Name})
	}
	if user.Password != nil {
		set = append(set, bson.E{"password", user.Password})
	}
	if user.IsAdmin != nil {
		set = append(set, bson.E{"isAdmin", user.IsAdmin})
	}
	if len(set) == 0 {
		return nil
	}
	result, err := sto.db().Collection(userCollectionName).
		UpdateOne(sto.newTimeoutContext(),
			bson.M{"_id": user.ID},
			bson.D{bson.E{"$set", set},})
	if err != nil {
		return fmt.Errorf("error updating  user with id \"%s\":%w", user.ID, err)
	}

	if result.MatchedCount == 0 {
		return  istorage.NewNotFoundError("users", "id", user.ID)
	}

	return nil
}

func (sto *Storage) DeleteUser(id string) error {
	if id == "" {
		return istorage.NewNotFoundError("users", "id", "")
	}
	result, err := sto.db().Collection(userCollectionName).DeleteOne(sto.newTimeoutContext(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("error removing  user with id \"%s\":%w", id, err)
	}
	if result.DeletedCount == 0 {
		return  istorage.NewNotFoundError("users", "id", id)
	}

	return nil
}

//Link related methods

func (sto *Storage) SaveLink(link models.Link) error {
	_, err := sto.db().Collection(linksCollectionName).InsertOne(sto.newTimeoutContext(), &link)
	if err != nil {
		//TODO test for already exists error
		return err
	}

	return nil
}

func (sto *Storage) GetLink(id string) (link models.Link, err error) {
	result := sto.db().Collection(linksCollectionName).FindOne(sto.newTimeoutContext(), bson.M{"_id": id})
	err = result.Err()

	if err == mongo.ErrNoDocuments {
		err = istorage.NewNotFoundError("link", "ID", id)
	}
	if  err != nil {
		err = fmt.Errorf("error searching link with id \"%s\":%w", id, err)
		return
	}
	if err = result.Decode(&link); err != nil {
		err = fmt.Errorf("error deconding user with id \"%s\":%w", id, err)
		return
	}
	return
}

func (sto *Storage) ListLinks(ownerID string, limit, offset uint) ([]models.Link, error) {
	filter :=  make(bson.M)
	options := mongoOptions.Find()
	options.SetSort(bson.D{{"createdAt", -1}})
	if limit != 0 {
		options.SetLimit(int64(limit))
	}
	if offset != 0 {
		options.SetSkip(int64(offset))
	}
	if ownerID != "" {
		filter["ownerId"] = ownerID
	}
	ctx := sto.newTimeoutContext()
	cursor, err := sto.db().Collection(linksCollectionName).Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var links []models.Link
	err = cursor.All(ctx, &links)
	return links, err
}

func (sto *Storage) UpdateLinkContent(id, content string) error {
	if id == "" {
		return istorage.NewNotFoundError("links", "id", "")
	}
	if content == "" {
		return nil
	}

	result, err := sto.db().Collection(linksCollectionName).
		UpdateOne(sto.newTimeoutContext(),
			bson.M{"_id": id},
			bson.D{bson.E{"$set", bson.D{{"content", content}}},})
	if err != nil {
		return fmt.Errorf("error updating link with id \"%s\":%w", id, err)
	}

	if result.MatchedCount == 0 {
		return  istorage.NewNotFoundError("links", "id", id)
	}

	return nil
}

func (sto *Storage) DeleteLink(id string) error {
	if id == "" {
		return istorage.NewNotFoundError("links", "id", "")
	}
	result, err := sto.db().Collection(linksCollectionName).DeleteOne(sto.newTimeoutContext(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("error removing link with id \"%s\":%w", id, err)
	}
	if result.DeletedCount == 0 {
		return  istorage.NewNotFoundError("links", "id", id)
	}

	return nil
}

func (sto *Storage) IncreaseLinkHitCount(id string) error {
	if id == "" {
		return istorage.NewNotFoundError("links", "id", "")
	}

	result, err := sto.db().Collection(linksCollectionName).
		UpdateOne(sto.newTimeoutContext(),
			bson.M{"_id": id},
			bson.D{bson.E{"$inc", bson.D{{"hits", 1}}},})
	if err != nil {
		return fmt.Errorf("error updating link with id \"%s\":%w", id, err)
	}

	if result.MatchedCount == 0 {
		return  istorage.NewNotFoundError("links", "id", id)
	}

	return nil
}

// Session related methods

func (sto *Storage) SaveSession(session models.Session) error {
	panic("not implemented")
}

func (sto *Storage) GetSession(id string) (models.Session, error) {
	panic("not implemented")
}


func (sto *Storage) ListSessions(ownerID string, limit, offset uint) ([]models.Session, error) {
	panic("not implemented")
}


func (sto *Storage) UpdateSessionToken(id string, tokenID string) error {
	panic("not implemented")
}


func (sto *Storage) DeleteSession(id string) error {
	panic("not implemented")
}
