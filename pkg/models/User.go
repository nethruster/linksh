package models

//User represents a user in the core logic
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
}
