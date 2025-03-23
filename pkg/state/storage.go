package state

import "fmt"

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

func (s *JsonStorage) Commit(state *State) error {
	fmt.Println("@TODO: Commit state to ", s.Filepath)
	return nil
}

// @TODO: FINISH
func (s *JsonStorage) Load(state *State) error {
	fmt.Println("@TODO: Load state from ", s.Filepath)
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
