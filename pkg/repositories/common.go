package repositories

import (
	sto "github.com/nethruster/linksh/pkg/interfaces/storage"
	"github.com/nethruster/linksh/pkg/interfaces/user_repository"
	errors "golang.org/x/xerrors"
)

func checkIfRequesterIsAdmin(storage sto.IStorage, requesterID string) (err error) {
	requester, err := storage.GetUser(requesterID)
	if err != nil {
		err = errors.Errorf("Error checking the requester %w", err)
		return
	}

	if !requester.IsAdmin {
		err = user_repository.ErrForbidden
	}
	return
}
