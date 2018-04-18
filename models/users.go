package models

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"errors"
)

type User struct {
    Id       string    `gorm:"primary_key; type:char(36)" json:"id"`
    Username string    `json:"username"`
    Email    string    `gorm:"UNIQUE_INDEX" json:"email"`
    Password []byte    `gorm:"type:binary(60)" json:"-"`
    Apikey   string    `gorm:"UNIQUE_INDEX; type:char(36)" json:"apikey"`
    IsAdmin  bool      `gorm:"DEFAULT:false"`
    Links    []Link    `json:"links"`
	Sessions []Session `json:"sessions"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (u *User) SaveToDatabase(db *gorm.DB) error {
	id, err := GenerateUserId()
	apikey,err := GenerateUserApiKey()
	hash, err := HashPassword(u.Password)
	if err != nil {return err}

	u.Id = id
	u.Links = nil
	u.Apikey = apikey
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Password = hash

	return db.Create(&u).Error
}

func GenerateUserId() (string, error) {
	return gonanoid.Nanoid(36)
}

func GenerateUserApiKey() (string, error) {
	return gonanoid.Nanoid(36)
}
func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (u *User) CheckIfCorrectPassword(plainPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, plainPassword)
	if err != nil {
		return false
	}

	return true
}

// Validations
func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("Missing username")
	} else if len(username) > 255 {
		return errors.New("Username must not be longer than 255 characters")
	} else {
		return nil
	}
}
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("Missing email")
	} else if len(email) > 255 {
		return errors.New("Email is to long")
	} else if !emailRegex.MatchString(email) {
		return errors.New("Invalid email format")
	} else {
		return nil
	}
}
func ValidatePassword(password []byte) error {
	if len(password) == 0  {
		return errors.New("Missing password")
	} else {
		return nil
	}
}

func (u *User) ValidateUser() []error {
	var errs []error

	err := ValidateUsername(u.Username)
	if err != nil {
		errs = append(errs, err)
	}

	err = ValidateEmail(u.Email)
	if err != nil {
		errs = append(errs, err)
	}

	err = ValidatePassword(u.Password)
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
