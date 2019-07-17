package userrepository

import "errors"

var (
	//ErrInvalidName is returned when the provided username doesn't accomplish the requirements of models.User.Name
	ErrInvalidName = errors.New("Invalid username")
	//ErrInvalidPassword is returned when the provided password doesn't accomplish the requirements of models.User.Password
	ErrInvalidPassword = errors.New("Invalid password")
	//ErrForbidden is returned when an ser user request to perform an action without enough privileges
	ErrForbidden = errors.New("Forbidden")
)
