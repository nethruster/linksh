package models

//User represents a user in the core logic
type User struct {
	ID       string `json:"id"`   //must be unique
	Name     string `json:"name"` //must be unique and no longer that 100 characters
	Password []byte `json:"-"`    //must be atleast than six characters long
	IsAdmin  bool   `json:"isAdmin"`
}
