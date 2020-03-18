package models

//Link represents a link in the core logic
type Link struct {
	ID      string `json:"id"`      //must be unique and no longer that 100 characters
	Content string `json:"content"` //must be no longer that 2000 characters
	Hits    uint   `json:"hits"`
	OwnerID string `json:"ownerId"`
}
