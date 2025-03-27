package software

// Software is a contianer for operations upon a OS software
type Software interface {
	Install() error
	Exists() bool

	GetName() string
	GetVersion() string
}

const (
	JavaSoftwareKey = "Java"
	MvnSoftwareKey  = "Maven"
	NodeSoftwareKey = "NodeJs"
)
