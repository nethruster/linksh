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
	sto,  err := newStorage()
	if err != nil {
		t.Errorf("Connection failed: %+v", err)
	}
	sto.close()
}


func TestUserRelatedMethods(t *testing.T) {
	var sto istorage.IStorage
	var err error

	mongoSto, err := newStorage()
	if err != nil {
		panic("CDatabase connection failed: " + err.Error())
	}
	defer mongoSto.close()

	if err = mongoSto.client.Database(mongoSto.databaseName).Collection(userCollectionName).Drop(mongoSto.newTimeoutContext()); err != nil {
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

		//TODO re-enable when implemented
		/*t.Run("conflict", func (t *testing.T) {
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
		*/
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
				t.Errorf("Expected NotFound got %v: %v", reflect.TypeOf(err), err.Error())
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
		user, err := sto.GetUser("abc")
		if err != nil {
			panic(err)
		}
		if user.Name != name {
			t.Errorf("The value was not updated, expected %s, got %s", name, user.Name)
		}

		t.Run("not found", func(t *testing.T) {
			payload.ID = "404"
			err = sto.UpdateUser(payload)
			if !errors.As(err, &istorage.NotFoundError{}) {
				t.Errorf("Expected NotFound got %v: %v", reflect.TypeOf(err), err.Error())
			}
		})
	})

	t.Run("delete", func(t *testing.T) {
		err = sto.DeleteUser("abc")
		if err != nil {
			t.Error(err)
		}

		if _, err = sto.GetUser("abc"); !errors.As(err, &istorage.NotFoundError{}) {
			t.Error("The user was not deleted from the database")
		}

		t.Run("not found", func(t *testing.T) {
			err = sto.DeleteUser("404")
			if !errors.As(err, &istorage.NotFoundError{}) {
				t.Errorf("Expected NotFound got %v: %v", reflect.TypeOf(err), err.Error())
			}
		})
	})
}

func TestLinkRelatedMethods(t *testing.T) {
	var sto istorage.IStorage
	var err error

	mongoSto, err := newStorage()
	if err != nil {
		panic("CDatabase connection failed: " + err.Error())
	}
	defer mongoSto.close()
	sto = mongoSto

	if err = mongoSto.client.Database(mongoSto.databaseName).Collection(linksCollectionName).Drop(mongoSto.newTimeoutContext()); err != nil {
		t.Errorf("Error reseting the collection: %+v", err)
	}

	t.Run("save", func(t *testing.T) {
		link := models.Link{
			ID: "abc",
			Content: "example.tld",
			CreatedAt: 100,
			Hits: 0,
			OwnerID: "abc",
		}
		err = sto.SaveLink(link)
		if err != nil {
			t.Error(err)
		}
		link.ID += "d"
		link.CreatedAt++

		err = sto.SaveLink(link)
		if err != nil {
			t.Error(err)
		}
		link.ID += "e"
		link.CreatedAt++
		link.OwnerID = "abcd"

		err = sto.SaveLink(link)
		if err != nil {
			t.Error(err)
		}

		//TODO check for collision
	})

	t.Run("get", func(t *testing.T) {
		_, err = sto.GetLink("abc")
		if err != nil {
			t.Error(err)
		}

		_, err = sto.GetUser("404")
		if !errors.As(err, &istorage.NotFoundError{}) {
			t.Errorf("Expected NotFound got %v: %v", reflect.TypeOf(err), err.Error())
		}
	})

	t.Run("list", func(t *testing.T) {
		var links []models.Link
		t.Run("no restrictions", func (t *testing.T) {
			links, err = sto.ListLinks("",0,0)
			if err != nil {
				t.Error(err)
			}
			if len(links) == 0 {
				t.Error("No results were returned")
			}
			if links[0].ID == "" {
				t.Error("The results were empty")
			}
		})

		t.Run("limit set", func(t *testing.T) {
			var links2 []models.Link
			links2, err = sto.ListLinks("", 1,0)
			if err != nil {
				t.Error(err)
			}
			if len(links2) != 1 {
				t.Errorf("The limit was set to 1 but %v results were returned", len(links))
			}
		})

		t.Run("offset set", func(t *testing.T) {
			var links3 []models.Link
			links3, err = sto.ListLinks("", 0,2)
			if err != nil {
				t.Error(err)
			}
			if links3[0].ID != links[2].ID {
				t.Errorf("The user was not the expected %+v", links3[0])
			}
		})

		t.Run("ownerID set", func(t *testing.T) {
				links, err = sto.ListLinks("abc",0,0)
				if err != nil {
					t.Error(err)
				}
				if len(links) == 0 {
					t.Error("No results were returned")
				}
				if links[0].ID == "" {
					t.Error("The results were empty")
				}
				for i, link := range links {
					if link.OwnerID != "abc" {
						t.Errorf("All links should be owned by the user with ID \"abc\", but links[%v] is owned by the user with ID \"%v\"", i, link.OwnerID)
					}
				}

			t.Run("limit set", func(t *testing.T) {
				var links2 []models.Link
				links2, err = sto.ListLinks("abc", 1,0)
				if err != nil {
					t.Error(err)
				}
				if len(links2) != 1 {
					t.Errorf("The limit was set to 1 but %v results were returned", len(links))
				}
				if links2[0].OwnerID != "abc" {
					t.Errorf("The link should be owned by the user with ID \"abc\", but was owned by the user with ID \"%v\" instead", links[0].OwnerID)
				}
			})

			t.Run("offset set", func(t *testing.T) {
				var links3 []models.Link
				links3, err = sto.ListLinks("abc", 0,1)
				if err != nil {
					t.Error(err)
				}
				if links3[0].ID != links[1].ID {
					t.Errorf("The user was not the expected %+v", links3[0])
				}
				for i, link := range links3 {
					if link.OwnerID != "abc" {
						t.Errorf("All links should be owned by the user with ID \"abc\", but links[%v] is owned by the user with ID \"%v\"", i, link.OwnerID)
					}
				}
			})
		})
	})

	t.Run("update", func(t *testing.T) {
		err = sto.UpdateLinkContent("abc", "example2.tld")
		if err != nil {
			t.Error(err)
		}
		link, err := sto.GetLink("abc")
		if err != nil {
			panic(err)
		}
		if link.Content != "example2.tld" {
			t.Errorf("The link was notupdated, content expected to be example2.tld, but was %s instead", link.Content)
		}

		t.Run("not found", func(t *testing.T) {
			_, err = sto.GetUser("404")
			if !errors.As(err, &istorage.NotFoundError{}) {
				t.Errorf("Expected NotFound got %v: %v", reflect.TypeOf(err), err.Error())
			}
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("delete", func(t *testing.T) {
			err = sto.DeleteLink("abc")
			if err != nil {
				t.Error(err)
			}

			if _, err = sto.GetLink("abc"); !errors.As(err, &istorage.NotFoundError{}) {
				t.Error("The link was not deleted from the database")
			}

			t.Run("not found", func(t *testing.T) {
				err = sto.DeleteLink("404")
				if !errors.As(err, &istorage.NotFoundError{}) {
					t.Errorf("Expected NotFound got %v: %v", reflect.TypeOf(err), err.Error())
				}
			})
		})
	})
}
