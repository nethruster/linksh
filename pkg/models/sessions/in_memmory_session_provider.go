package sessions

import (
	"bufio"
	"encoding/gob"
	"errors"
	"os"
	"sync"
	"time"
)

//InMemmorySessionProvider is valid provider which stores the sessions on RAM
type InMemmorySessionProvider struct {
	mutex   sync.RWMutex
    storage map[string]*Session
    autoSaveShouldBeRunning bool
}

//NewInMemmorySessionProvider returns a new instance of InMemmorySessionProvider
func NewInMemmorySessionProvider() *InMemmorySessionProvider {
	provider := InMemmorySessionProvider{
		storage: make(map[string]*Session),
	}

	return &provider
}

//Add a session to the storage
func (provider *InMemmorySessionProvider) Add(session Session) error {
	provider.mutex.Lock()
	defer provider.mutex.Unlock()
	provider.storage[session.ID] = &session

	return nil
}

//Get returns the requested session, if not found it returns an error
func (provider *InMemmorySessionProvider) Get(id string) (Session, error) {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()

	session, ok := provider.storage[id]
	if !ok {
		return *session, errors.New("Session not found")
	}
	return *session, nil
}

//GetByOwnerID Returns the sessions which belongs to the selected user
func (provider *InMemmorySessionProvider) GetByOwnerID(ownerID string) (map[string]Session, error) {
	provider.mutex.RLock()
	defer provider.mutex.RUnlock()

	result := make(map[string]Session)
	for i, session := range provider.storage {
		if session.OwnerID == ownerID {
			result[i] = *session
		}
	}

	return result, nil
}

//Update a session
func (provider *InMemmorySessionProvider) Update(session Session) error {
	provider.mutex.Lock()
	defer provider.mutex.Unlock()
	provider.storage[session.ID] = &session

	return nil
}

//Delete the session with selected id
func (provider *InMemmorySessionProvider) Delete(id string) error {
	provider.mutex.Lock()
	defer provider.mutex.Unlock()

	delete(provider.storage, id)

	return nil
}

//GC deletes expired entries from the storage
func (provider *InMemmorySessionProvider) GC() error {
	var targetedSessionsID []string
	now := time.Now().Unix()
	provider.mutex.RLock()
	for id, session := range provider.storage {
		if session.ExpiresOn < now {
			targetedSessionsID = append(targetedSessionsID, id)
		}
	}
	provider.mutex.RUnlock()

	for _, id := range targetedSessionsID {
		provider.Delete(id)
	}

	return nil
}

//DumpToDisk creates a copy of the stored sessions in the specified path
func (provider *InMemmorySessionProvider) DumpToDisk(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	bufferedWriter := bufio.NewWriter(file)
	encoder := gob.NewEncoder(bufferedWriter)

	provider.mutex.RLock()
	err = encoder.Encode(provider.storage)
	provider.mutex.RUnlock()
	if err != nil {
		return err
	}

	err = bufferedWriter.Flush()
	if err != nil {
		return err
	}

	return nil
}

//RecoveFromDisk recoves the data from a previus dump in the spcified path
func (provider *InMemmorySessionProvider) RecoveFromDisk(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)

	err = decoder.Decode(&provider.storage)
	if err != nil {
		return err
	}

	return nil
}

//EnableAutoSave dump the session storage to disk once every time the provided duration has passed in the desired path
// Only a job of autosave can be running at the same time for each InMemmorySessionProvider instance.
func (provider *InMemmorySessionProvider) EnableAutoSave(path string, x time.Duration) error {
    if provider.autoSaveShouldBeRunning {
        return errors.New("The autosave job is already running")
    }
    provider.autoSaveShouldBeRunning = true
    go func() {
        for provider.autoSaveShouldBeRunning {
            provider.DumpToDisk(path)
            time.Sleep(x)
        }
    }()

    return nil
}

//DisableAutoSave disable the autosave functionality
func (provider *InMemmorySessionProvider) DisableAutoSave() error {
    if !provider.autoSaveShouldBeRunning {
        return errors.New("The autosave job was not running")
    }
    provider.autoSaveShouldBeRunning = false

    return nil
}
