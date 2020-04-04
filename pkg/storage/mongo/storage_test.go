package mongo

import (
	"errors"
	istorage "github.com/nethruster/linksh/pkg/interfaces/storage"
	"github.com/nethruster/linksh/pkg/interfaces/user_repository"
	"github.com/nethruster/linksh/pkg/models"
	"os"
	"reflect"
	"testing"
	"time"
)

func newStorage() (*Storage, error) {
	dbString, ok := os.LookupEnv("LINKSH_TEST_MONGOSTRING")
	if !ok {
		panic("LINKSH_TEST_MONGOSTRING is not set")
	}
	dbName, ok := os.LookupEnv("LINKSH_TEST_DB_NAME")
	if !ok {
		dbName = "linksh_test"
	}

	return New(dbString, dbName, 10 * time.Second)
}

func TestMongoConnection(t *testing.T) {
	_,  err := newStorage()
	if err != nil {
		t.Errorf("Connection failed: %+v", err)
	}
}


func TestUser(t *testing.T) {
	var sto *Storage
	var err error

	sto, err = newStorage()
	if err != nil {
		panic("CDatabase connection failed: " + err.Error())
	}

	if err = sto.client.Database(sto.databaseName).Collection(userCollectionName).Drop(sto.newTimeoutContext()); err != nil {
		t.Errorf("Error reseting the collection: %+v", err)
	}

	t.Run("save", func(t *testing.T) {
		user := models.User{
			ID: "abc",
			Name: "testUser",
			Password: []byte("1234"),
			IsAdmin: true,
		}
		err = sto.SaveUser(user)
		if err != nil {
			t.Error(err)
		}

		user.ID += "d"
		user.Name += "2"
		err = sto.SaveUser(user)
		if err != nil {
			t.Error(err)
		}

		user.ID += "e"
		user.Name = "testUser3"
		err = sto.SaveUser(user)
		if err != nil {
			t.Error(err)
		}

		t.Run("conflict", func (t *testing.T) {
			oldName := user.Name

			t.Run("id", func(t *testing.T) {
				user.Name = "otherName"
				err = sto.SaveUser(user)
				var conflictErr *istorage.AlreadyExistsError
				if !errors.As(err, &conflictErr) {
					t.Errorf("Expected conflic error, got %v: %v", reflect.TypeOf(err), err.Error())
				}
				if conflictErr.Field != "ID" {
					t.Errorf("Expected conflic in field ID but it was on field %s instead", conflictErr.Field)
				}
			})

			t.Run("name", func(t *testing.T) {
				user.Name = oldName
				err = sto.SaveUser(user)
				var conflictErr *istorage.AlreadyExistsError
				if !errors.As(err, &conflictErr) {
					t.Errorf("Expected conflic error, got %v: %v", reflect.TypeOf(err), err.Error())
				}
				if conflictErr.Field != "ID" {
					t.Errorf("Expected conflic in field Name but it was on field %s instead", conflictErr.Field)
				}
			})
		})

	})

	t.Run("get", func (t *testing.T) {
		t.Run("byID", func(t *testing.T) {
			user, err := sto.GetUser("abc")
			if err != nil {
				t.Error(err)
			}
			if user.ID != "abc" {
				t.Errorf("UserID doesn't match, expected \"%v\" got º\"%v\"", "abc", user.ID)
			}

			_, err = sto.GetUser("404")
			if !errors.As(err, &istorage.NotFoundError{}) {
				t.Errorf("Expected NotFound got %v", err.Error())
			}
		})

		t.Run("byName", func(t *testing.T) {
			user, err := sto.GetUserByName("testUser")
			if err != nil {
				t.Error(err)
			}
			if user.ID != "abc" {
				t.Errorf("UserName doesn't match, expected \"%v\" got º\"%v\"", "abc", user.ID)
			}

			_, err = sto.GetUser("404")
			if !errors.As(err, &istorage.NotFoundError{}) {
				t.Errorf("Expected NotFound got %v", err.Error())
			}
		})
	})

	t.Run("list", func (t *testing.T) {
		var users []models.User
		t.Run("no restrictions", func (t *testing.T) {
			users, err = sto.ListUsers(0,0)
			if err != nil {
				t.Error(err)
			}
			if len(users) == 0 {
				t.Error("No results were returned")
			}
			if users[0].ID == "" {
				t.Error("The results were empty")
			}
		})

		t.Run("limit set", func(t *testing.T) {
			var users2 []models.User
			users2, err = sto.ListUsers(1,0)
			if err != nil {
				t.Error(err)
			}
			if len(users2) != 1 {
				t.Errorf("The limit was set to 1 but %v results were returned", len(users))
			}
		})

		t.Run("offset set", func(t *testing.T) {
			var users3 []models.User
			users3, err = sto.ListUsers(0,2)
			if err != nil {
				t.Error(err)
			}
			if users3[0].ID != users[2].ID {
				t.Errorf("The user was not the expected %+v", users[0])
			}
		})
	})

	t.Run("update", func (t *testing.T) {
		name := "Paco"
		payload := user_repository.UpdatePayload{ID: "abc", Name: &name}
		err = sto.UpdateUser(payload)
		if err != nil {
			t.Error(err)
		}

		payload.ID = "404"
		err = sto.UpdateUser(payload)
		if !errors.As(err, &istorage.NotFoundError{}) {
			t.Errorf("Expected NotFound got %v", err.Error())
		}
	})

	t.Run("delete", func(t *testing.T) {
		err = sto.DeleteUser("abc")
		if err != nil {
			t.Error(err)
		}

		err = sto.DeleteUser("404")
		if !errors.As(err, &istorage.NotFoundError{}) {
			t.Errorf("Expected NotFound got %v", err.Error())
		}
	})
}

