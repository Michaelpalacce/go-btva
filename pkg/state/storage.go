package state

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

// Storage is an interface that defines how to store and load the state
type Storage interface {
	// Commit will store the state
	Commit(state *State) error
	// Load will load the state
	Load(state *State) error
}

// JsonStorage is a struct that will store the state in JSON
type JsonStorage struct {
	Filepath string
}

// Commit will save the current State object to file
// @NOTE: On windows the file permissions are ignored
func (s *JsonStorage) Commit(state *State) error {
	bytes, err := json.MarshalIndent(state, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(s.Filepath, bytes, 0o640)
}

// Load will load the state from a file if it exists.
func (s *JsonStorage) Load(state *State) error {
	bytes, err := os.ReadFile(s.Filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		bytes = []byte("{}")
	}

	if err = json.Unmarshal(bytes, state); err != nil {
		return fmt.Errorf("error while loading previous state. Error was %w", err)
	} else {
		slog.Info("Loaded previous state", "state", state.State)
	}

	return nil
}

// WithJsonStorage will make the State store data in JSON
func WithJsonStorage(filepath string, load bool) SetStateOption {
	return func(s *State) error {
		storage := &JsonStorage{Filepath: filepath}
		s.storage = append(s.storage, storage)

		if load {
			return storage.Load(s)
		}

		return nil
	}
}
