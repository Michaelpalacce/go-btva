package args

// Software holds a list of all the software that needs to be installed, so that means that this struct should contain all the possible
// software the tool supports.
type Software struct {
	// Versions
	JavaVersion string `json:"javaVersion"`
	MvnVersion  string `json:"mvnVersion"`
	NodeVersion string `json:"nodeVersion"`
}

// Local will hold different configurations for mvn that may be needed
type Local struct{}

// Infra will hold different infra decisions that need to be taken
type Infra struct {
	// MinimalInfrastructure
	MinimalInfrastructure bool `json:"minimalInfrastructure"`

	SSHVMIP                 string `json:"sshVmIP"`
	SSHUsername             string `json:"sshUsername"`
	SSHPassword             string `json:"sshPassword"`
	SSHPrivateKey           string `json:"sshPrivateKey"`
	SSHPrivateKeyPassphrase string `json:"sshPrivateKeyPassphrase"`

	DockerUsername string `json:"dockerUsername"`
	DockerPAT      string `json:"dockerPat"`
	// MinimalInfrastructure
}

// Options is a struct for options that the tool can accept.
type Options struct {
	Software Software `json:"software"`
	Local    Local    `json:"local"`
	Infra    Infra    `json:"infra"`

	// parsed is an internal variable that tells us that the options have already been parsed and don't need a second go
	parsed bool
}

// This is a single instance of the options. We don't want to parse them more than once
var options = &Options{
	Software: Software{},
	Local:    Local{},

	parsed: false,
}
