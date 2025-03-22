package state

// internalState is a struct that contains the internal state of a state key
type internalState struct {
	// Done signifies that the current state is complete
	Done bool `json:"successful"`
	// Msg is any human readable message that was or was not added
	Msg string `json:"msg"`
	// Step signifies at which step of the execution we are
	Step int `json:"step"`
	// Err is the State-Step error if any
	Err error
}

// State contains the current state of the application
type State struct {
	_state map[string]internalState
}

// NewState will return a new state object, ready to be used
func NewState() *State {
	state := &State{}

	state.Init()

	return state
}

// SetStateOption follows the famous functional options pattern. It returns an error if there was an issue performing an operation on the state
type (
	SetStateOption func(*State) error

	GetStateOption        func(*State) *internalState
	GetSuccessStateOption func(*State) bool
	GetMsgStateOption     func(*State) string
	GetErrStateOption     func(*State) error
)

func (s *State) Set(option SetStateOption) error           { return option(s) }
func (s *State) Get(option GetStateOption) *internalState  { return option(s) }
func (s *State) GetDone(option GetSuccessStateOption) bool { return option(s) }
func (s *State) GetMsg(option GetMsgStateOption) string    { return option(s) }
func (s *State) GetErr(option GetErrStateOption) error     { return option(s) }

func (s *State) GetValue(key string) *internalState {
	if value, ok := s._state[key]; !ok {
		return nil
	} else {
		return &value
	}
}

func (s *State) SetValue(key string, done bool, msg string, step int, err error) {
	s._state[key] = internalState{Done: done, Msg: msg, Step: step, Err: err}
}

// Init can be used to initialize the internals of the State object
func (s *State) Init() {
	if s._state == nil {
		s._state = make(map[string]internalState)
	}
}
