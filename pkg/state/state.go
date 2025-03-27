package state

import "log/slog"

// SetStateOption follows the famous functional options pattern. It returns an error if there was an issue performing an operation on the state
type (
	SetStateOption func(*State) error

	GetStateOption            func(*State) *internalState
	GetSuccessStateOption     func(*State) bool
	GetStepStateOption        func(*State) int
	GetMsgStateOption         func(*State) string
	GetErrStateOption         func(*State) error
	GetContextPropStateOption func(*State) string
)

// internalState is a struct that contains the internal state of a state key
type internalState struct {
	// Done signifies that the current state is complete
	Done bool `json:"done"`
	// Msg is any human readable message that was or was not added
	Msg string `json:"msg"`
	// Step signifies at which step of the execution we are
	Step int `json:"step"`
	// Err is the State-Step error if any
	Err error `json:"-"`
	// Context is a container for additional data that must be stored for the state
	Context map[string]string `json:"data"`
}

// State contains the current state of the application
type State struct {
	State map[string]internalState `json:"state"`

	// storage is an array, as you may want to store the state in multiple places
	storage []Storage
}

// NewState will return a new state object, ready to be used
func NewState() *State {
	state := &State{}

	state.Init()

	return state
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
		storage.Commit(s)
	}

	return nil
}
func (s *State) Get(option GetStateOption) *internalState              { return option(s) }
func (s *State) GetDone(option GetSuccessStateOption) bool             { return option(s) }
func (s *State) GetMsg(option GetMsgStateOption) string                { return option(s) }
func (s *State) GetStep(option GetStepStateOption) int                 { return option(s) }
func (s *State) GetErr(option GetErrStateOption) error                 { return option(s) }
func (s *State) GetContextKey(option GetContextPropStateOption) string { return option(s) }

// GetValue is used to get the internal state of a key... not recommended for direct use
func (s *State) GetValue(key string) *internalState {
	if value, ok := s.State[key]; !ok {
		return nil
	} else {
		return &value
	}
}

// SetValue is used to directly set a key in the internal storage... not recommended for direct use
func (s *State) SetValue(key string, done bool, msg string, step int, err error) {
	s.State[key] = internalState{Done: done, Msg: msg, Step: step, Err: err, Context: make(map[string]string)}
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

func GetStep(key string) GetStepStateOption {
	return func(s *State) int {
		value := s.GetValue(key)
		if value == nil {
			return 0
		}

		return value.Step
	}
}

func GetDone(key string) GetSuccessStateOption {
	return func(s *State) bool {
		value := s.GetValue(key)
		if value == nil {
			return false
		}

		return value.Done
	}
}

// State Options

func WithDone(key string, done bool) SetStateOption {
	return func(s *State) error {
		if _, ok := s.State[key]; !ok {
			s.State[key] = internalState{}
		}

		value := s.State[key]
		value.Done = done
		s.State[key] = value

		return nil
	}
}

// WithMsg sets the message of the state.
// @NOTE: It will also log the message, as setting a message is for human consumption
func WithMsg(key string, msg string) SetStateOption {
	return func(s *State) error {
		if _, ok := s.State[key]; !ok {
			s.State[key] = internalState{}
		}

		value := s.State[key]
		value.Msg = msg
		s.State[key] = value

		slog.Info(msg)

		return nil
	}
}

// WithStep sets the step, however it will NOT decrement a step
func WithStep(key string, step int) SetStateOption {
	return func(s *State) error {
		if _, ok := s.State[key]; !ok {
			s.State[key] = internalState{}
		}

		value := s.State[key]
		if value.Step < step {
			value.Step = step
		}
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
