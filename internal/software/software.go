package software

// Software is a contianer for operations upon a OS software
type Software interface {
	Install() error
	Remove() error
	Exists() bool
}
