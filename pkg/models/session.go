package models

//Session describes a session
//ExpireDate is an Unix EPOCH
type Session struct {
	ID string
	UserID string
	LastToken string
	ExpireDate uint64
}
