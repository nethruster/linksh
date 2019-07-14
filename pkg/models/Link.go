package models

//Link represents a link in the core logic
type Link struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Hits    uint   `json:"hits"`
	OwnerID string `json:"owner_id"`
}
