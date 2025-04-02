package args

// Software holds a list of all the software that needs to be installed, so that means that this struct should contain all the possible
// software the tool supports.
type Software struct {
	// Versions
	JavaVersion string `json:"javaVersion"`
	MvnVersion  string `json:"mvnVersion"`
	NodeVersion string `json:"nodeVersion"`
}

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

type AriaAutomation struct {
	FQDN        string `json:"fqdn"`
	Port        string `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	OrgName     string `json:"orgName"`
	ProjectName string `json:"projectName"`
}

type Aria struct {
	Automation AriaAutomation `json:"automation"`
}

// Holds Artifactory (jFrog Artifactory, Nexus, Azure Packages, etc) details that will be used to generate the settings.xml file
type Artifactory struct {
	ReleaseRepo  string
	SnapshotRepo string
	GroupRepo    string
	Password     string
}

// Options is a struct for options that the tool can accept.
type Options struct {
	Software    Software    `json:"software"`
	Infra       Infra       `json:"infra"`
	Aria        Aria        `json:"aria"`
	Artifactory Artifactory `json:"artifactory"`

	// parsed is an internal variable that tells us that the options have already been parsed and don't need a second go
	parsed bool
}

// This is a single instance of the options. We don't want to parse them more than once
var options = &Options{
	Software: Software{},
	// Local:    Local{},

	parsed: false,
}
