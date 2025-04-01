package state

import (
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
)

// SetStateOption follows the famous functional options pattern. It returns an error if there was an issue performing an operation on the state
type (
	SetStateOption func(*State) error

	GetStateOption            func(*State) *internalState
	GetSuccessStateOption     func(*State) bool
	GetMsgStateOption         func(*State) string
	GetErrStateOption         func(*State) error
	GetContextPropStateOption func(*State) string
)

// internalState is a struct that contains the internal state of a state key
type internalState struct {
	// Msg is any human readable message that was or was not added
	Msg string `json:"msg"`
	// Err is the error if any
	Err error `json:"-"`
	// Context is a container for additional data that must be stored for the state
	Context map[string]string `json:"data,omitempty"`
}

// State contains the current state of the application
type State struct {
	State map[string]internalState `json:"state"`

	// Options are stored here so we can preserve them between multiple runs
	Options *args.Options `json:"options,omitempty"`

	// storage is an array, as you may want to store the state in multiple places
	storage []Storage
}

// NewState will return a new state object, ready to be used
func NewState(options ...SetStateOption) (*State, error) {
	state := &State{}

	state.Init()

	if err := state.Modify(options...); err != nil {
		return nil, err
	}

	return state, nil
}

// Modify is used to modify the internal configuration of the State object.
func (s *State) Modify(options ...SetStateOption) error {
	for _, option := range options {
		if err := option(s); err != nil {
			return err
		}
	}

	return nil
}

// Set is used to modify the internal storage of the State object.
func (s *State) Set(options ...SetStateOption) error {
	if err := s.Modify(options...); err != nil {
		return err
	}

	for _, storage := range s.storage {
		storage.Commit(*s)
	}

	return nil
}

// Get is used to get the internal state of a key
// You can pass any of the State Options provided above or write your own one
func Get[T any](s *State, option func(*State) T) T {
	return option(s)
}

// GetValue is used to get the internal state of a key... not recommended for direct use, unless writing your own GetOptions
func (s *State) GetValue(key string) *internalState {
	if value, ok := s.State[key]; !ok {
		return nil
	} else {
		return &value
	}
}

// Init can be used to initialize the internals of the State object
func (s *State) Init() {
	if s.State == nil {
		s.State = make(map[string]internalState)
	}
}

// Get Options

func GetState(key string) GetStateOption {
	return func(s *State) *internalState {
		if value, ok := s.State[key]; ok {
			return &value
		}

		return nil
	}
}

func GetContextProp(key, prop string) GetContextPropStateOption {
	return func(s *State) string {
		if value, ok := s.State[key]; ok {
			if val, ok := value.Context[prop]; ok {
				return val
			}
		}

		return ""
	}
}

// State Options

// WithMsg wraps WithQuietMsg but it will also log the message
func WithMsg(key string, msg string) SetStateOption {
	return func(s *State) error {
		slog.Info(msg)

		return WithQuietMsg(key, msg)(s)
	}
}

// WithQuietMsg sets the message of the state
func WithQuietMsg(key string, msg string) SetStateOption {
	return func(s *State) error {
		if _, ok := s.State[key]; !ok {
			s.State[key] = internalState{}
		}

		value := s.State[key]
		value.Msg = msg
		s.State[key] = value

		return nil
	}
}

func WithErr(key string, err error) SetStateOption {
	return func(s *State) error {
		if _, ok := s.State[key]; !ok {
			s.State[key] = internalState{}
		}

		value := s.State[key]
		value.Err = err

		if err != nil {
			value.Msg = err.Error()
		}
		s.State[key] = value

		return nil
	}
}

func WithContextProp(key, prop, value string) SetStateOption {
	return func(s *State) error {
		if _, ok := s.State[key]; !ok {
			s.State[key] = internalState{
				Context: make(map[string]string),
			}
		}

		state := s.State[key]
		if state.Context == nil {
			state.Context = make(map[string]string)
		}

		state.Context[prop] = value
		s.State[key] = state

		return nil
	}
}
