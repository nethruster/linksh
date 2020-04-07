package models

//Session describes a session
type Session struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	LastToken  string `json:"last_token" bson:"last_token"`
	//ExpireDate must be an Unix EPOCH
	ExpireDate uint64 `json:"expire_date" bson:"expire_date"`
}
