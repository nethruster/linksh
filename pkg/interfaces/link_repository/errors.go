package link_repository

import "errors"

var (
	//ErrInvalidID is returned when the provided username doesn't accomplish the requirements of models.Link.ID
	ErrInvalidID = errors.New("Invalid ID")
	//ErrInvalidContent is returned when the provided content doesn't accomplish the requirements of models.Link.Content
	ErrInvalidContent = errors.New("Invalid content")
	//ErrForbidden is returned when an ser user request to perform an action without enough privileges
	ErrForbidden = errors.New("Forbidden")
)
