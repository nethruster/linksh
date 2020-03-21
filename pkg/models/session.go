package models

//Session describes a session
//ExpireDate is an Unix EPOCH
type Session struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	LastToken  string `json:"last_token" bson:"last_token"`
	ExpireDate uint64 `json:"expire_date" bson:"expire_date"`
}
