package models

//Link represents a link in the core logic
type Link struct {
	//ID must be unique and no longer that 100 characters
	ID        string   `json:"id" bson:"_id"`
	//Content must be no longer that 2000 characters
	Content   string   `json:"content" bson:"content"`
	Hits      uint     `json:"hits" bson:"hits"`
	//CreatedAt must be an Unix EPOCH
	CreatedAt int64 `json:"createdAt" bson:"createdAt"`
	OwnerID   string    `json:"ownerId" bson:"ownerId"`
}
