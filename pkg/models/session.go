package models

//Session describes a session
type Session struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	LastToken  string `json:"last_token" bson:"last_token"`
	// CreatedAt must be an Unix EPOCH
	CreatedAt int64 `json:"createdAt" bson:"createdAt"`
	//ExpireDate must be an Unix EPOCH
	ExpireDate int64 `json:"expire_date" bson:"expire_date"`
}
