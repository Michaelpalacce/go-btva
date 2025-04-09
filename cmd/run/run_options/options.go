package run_options

// Software holds a list of all the software that needs to be installed, so that means that this struct should contain all the possible
// software the tool supports.
type Software struct {
	// Versions
	JavaVersion   string `json:"javaVersion"`
	MvnVersion    string `json:"mvnVersion"`
	NodeVersion   string `json:"nodeVersion"`
	VsCodeVersion string `json:"vsCodeVersion"`
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

// HasMinimalInfra returns true if either gitlab or nexus minimal infra is required
func (m *MinimalInfra) HasMinimalInfra() bool {
	return m.MinimalInfrastructureGitlab || m.MinimalInfrastructureNexus
}

func (m *MinimalInfra) HasNexus() bool {
	return m.MinimalInfrastructureNexus
}

func (m *MinimalInfra) HasGitlab() bool {
	return m.MinimalInfrastructureGitlab
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

// IsPartial can be used to check if at least one of the keys is set... This can indicate a partial setup and may be good to prompt the user
// for the rest.
func (a *ArtifactManager) IsPartial() bool {
	return a.Password != "" || a.ReleaseRepo != "" || a.GroupRepo != "" || a.SnapshotRepo != ""
}

// RunOptions is the the spec from the user what is wanted.
type RunOptions struct {
	Software        Software        `json:"software"`
	MinimalInfra    MinimalInfra    `json:"mininalInfra"`
	Aria            Aria            `json:"aria"`
	ArtifactManager ArtifactManager `json:"artifactManager"`

	// Parsed is an internal variable that tells us that the options have already been Parsed and don't need a second go
	Parsed bool
}
