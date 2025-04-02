package args

// Software holds a list of all the software that needs to be installed, so that means that this struct should contain all the possible
// software the tool supports.
type Software struct {
	// Versions
	JavaVersion string `json:"javaVersion"`
	MvnVersion  string `json:"mvnVersion"`
	NodeVersion string `json:"nodeVersion"`
}

// MinimalInfra will hold different minimal infra decisions that need to be taken
type MinimalInfra struct {
	MinimalInfrastructureGitlab bool `json:"minimalInfrastructureGitlab"`
	MinimalInfrastructureNexus  bool `json:"minimalInfrastructureNexus"`

	SSHVMIP                 string `json:"sshVmIP"`
	SSHUsername             string `json:"sshUsername"`
	SSHPassword             string `json:"sshPassword"`
	SSHPrivateKey           string `json:"sshPrivateKey"`
	SSHPrivateKeyPassphrase string `json:"sshPrivateKeyPassphrase"`

	DockerUsername string `json:"dockerUsername"`
	DockerPAT      string `json:"dockerPat"`
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

// Holds ArtifactManager (jFrog Artifactory, Nexus, Azure Packages, etc) details that will be used to generate the settings.xml file
type ArtifactManager struct {
	ReleaseRepo  string `json:"releaseRepo"`
	SnapshotRepo string `json:"snapshotRepo"`
	GroupRepo    string `json:"groupRepo"`
	Password     string `json:"password"`
}

// Options is the the spec from the user what is wanted.
type Options struct {
	Software        Software        `json:"software"`
	MinimalInfra    MinimalInfra    `json:"mininalInfra"`
	Aria            Aria            `json:"aria"`
	ArtifactManager ArtifactManager `json:"artifactManager"`

	// parsed is an internal variable that tells us that the options have already been parsed and don't need a second go
	parsed bool
}

// This is a single instance of the options. We don't want to parse them more than once
var options = &Options{
	Software: Software{},

	parsed: false,
}
