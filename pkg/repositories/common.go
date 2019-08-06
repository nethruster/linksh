package repositories

import (
	errors "golang.org/x/xerrors"
	sto "github.com/nethruster/linksh/pkg/interfaces/storage"
	"github.com/nethruster/linksh/pkg/interfaces/userrepository"
)

func checkIfRequesterIsAdmin(storage sto.IStorage, requesterID string) (err error) {
	requester, err := storage.GetUser(requesterID)
	if err != nil {
		err = errors.Errorf("Error checking the requester %w", err)
		return
	}

	if !requester.IsAdmin {
		err = userrepository.ErrForbidden
	}
	return
}
