// Package storage is responsible for storing data.
package storage

import (
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// FileStorage is a storage that uses plain text file for storing data.
type FileStorage struct {
	file  string
	state State
	mx    *sync.Mutex
}

// NewFileStorage creates new file storage.
func NewFileStorage(file string) (*FileStorage, error) {
	s := &FileStorage{
		file:  file,
		state: State{Feeds: map[string]time.Time{}},
		mx:    &sync.Mutex{},
	}

	// Read or init
	b, err := ioutil.ReadFile(s.file)
	if os.IsNotExist(err) {
		if err := s.save(); err != nil {
			return nil, errors.Wrap(err, "init file")
		}
		return s, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	// Unmarshal
	if err = yaml.Unmarshal(b, &s.state); err != nil {
		return nil, errors.Wrap(err, "decode data")
	}
	if s.state.Feeds == nil {
		s.state.Feeds = map[string]time.Time{}
	}
	return s, nil
}

// GetChats gets list of chat IDs.
func (s *FileStorage) GetChats() []int64 {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.state.Chats
}

// SaveChats saves list of chat IDs.
func (s *FileStorage) SaveChats(chats []int64) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.state.Chats = chats
	return s.save()
}

// GetLastUpdate gets last update time of the feed.
func (s *FileStorage) GetLastUpdate(feed string) time.Time {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.state.Feeds[feed]
}

// SaveLastUpdate saves last feed update.
func (s *FileStorage) SaveLastUpdate(feed string, t time.Time) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.state.Feeds[feed] = t
	return s.save()
}

// save rewrites whole current state in file.
func (s *FileStorage) save() error {
	b, err := yaml.Marshal(s.state)
	if err != nil {
		return errors.Wrap(err, "encode data")
	}
	err = ioutil.WriteFile(s.file, b, 0o600)
	if err != nil {
		return errors.Wrap(err, "write data to file")
	}
	return nil
}
