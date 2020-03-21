package models

//User represents a user in the core logic
type User struct {
	ID       string `json:"id" bson:"_id"`     //must be unique
	Name     string `json:"name" bson:"name"`  //must be unique and no longer that 100 characters
	Password []byte `json:"-" bson:"password"` //must be at least than six characters long
	IsAdmin  bool   `json:"isAdmin"`
}
