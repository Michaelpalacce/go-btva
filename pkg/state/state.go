package state

// SetStateOption follows the famous functional options pattern. It returns an error if there was an issue performing an operation on the state
type (
	SetStateOption func(*State) error

	GetStateOption        func(*State) *internalState
	GetSuccessStateOption func(*State) bool
	GetMsgStateOption     func(*State) string
	GetErrStateOption     func(*State) error
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
func (s *State) Get(option GetStateOption) *internalState  { return option(s) }
func (s *State) GetDone(option GetSuccessStateOption) bool { return option(s) }
func (s *State) GetMsg(option GetMsgStateOption) string    { return option(s) }
func (s *State) GetErr(option GetErrStateOption) error     { return option(s) }

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
	s.State[key] = internalState{Done: done, Msg: msg, Step: step, Err: err}
}

// Init can be used to initialize the internals of the State object
func (s *State) Init() {
	if s.State == nil {
		s.State = make(map[string]internalState)
	}
}
