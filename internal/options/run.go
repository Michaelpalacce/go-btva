package options

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/pkg/prompt"
)

// Software holds a list of all the software that needs to be installed, so that means that this struct should contain all the possible
// software the tool supports.
type Software struct {
	// Versions
	JavaVersion       string `json:"javaVersion"`
	MvnVersion        string `json:"mvnVersion"`
	NodeVersion       string `json:"nodeVersion"`
	VsCodeVersion     string `json:"vsCodeVersion"`
	ChocolateyVersion string `json:"chocolateyVersion"`
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
	Prompt          bool            `json:"prompt"`

	// Parsed is an internal variable that tells us that the options have already been Parsed and don't need a second go
	Parsed bool
}

// @NOTE
// Validations should be ran when they are actually needed. Used to fill missing content.

// ValidateMinimalInfra runs validation if all the inputs for running minimal infrastructure are set
func (options *RunOptions) ValidateMinimalInfra() error {
	var err error
	if options.MinimalInfra.MinimalInfrastructureNexus || options.MinimalInfra.MinimalInfrastructureGitlab {
		if options.MinimalInfra.SSHVMIP == "" {
			if options.MinimalInfra.SSHVMIP, err = prompt.AskText("MinimalInfrastructure selected, but you did not provide sshVmIp, please type in the IP: "); err != nil {
				return fmt.Errorf("sshVmIp must be provided. Err: %w", err)
			}
		}

		if options.MinimalInfra.SSHPrivateKey == "" && options.MinimalInfra.SSHPassword == "" {
			if options.MinimalInfra.SSHPassword, err = prompt.AskPass("MinimalInfrastructure selected, but you did not provide sshPassword or sshPrivateKey, please type in password: "); err != nil {
				return fmt.Errorf("sshPassword must be provided. Err: %w", err)
			}
		}

		if options.MinimalInfra.SSHUsername == "" {
			if options.MinimalInfra.SSHUsername, err = prompt.AskText("MinimalInfrastructure selected, but you did not provide sshUsername, please type in the username of root or a passwordless sudo user"); err != nil {
				return fmt.Errorf("sshUsername must be provided. Err: %w", err)
			}
		}
	}

	if options.MinimalInfra.DockerUsername != "" && options.MinimalInfra.DockerPAT == "" {
		if options.MinimalInfra.DockerPAT, err = prompt.AskPass("dockerUsername passed, but you did not provide dockerPat, please type in public access token: "); err != nil {
			return fmt.Errorf("dockerPat must be provided with dockerUsername. Err: %w", err)
		}
	}

	return nil
}

// ValidateAriaAutomation will prompt the user a series of question needed to build the aria inventory if the settings are missing
func (options *RunOptions) ValidateAriaAutomation() error {
	var err error
	if options.Aria.Automation.FQDN == "" {
		if options.Aria.Automation.FQDN, err = prompt.AskText(fmt.Sprintf("What is Aria Automation's FQDN without `https://`. Current (%s)", options.Aria.Automation.FQDN)); err != nil {
			return err
		}
	}

	if options.Aria.Automation.Port == "" {
		if options.Aria.Automation.Port, err = prompt.AskText(fmt.Sprintf("What is Aria Automation's port? Current (%s)", options.Aria.Automation.Port)); err != nil {
			return err
		}
	}

	if options.Aria.Automation.Username == "" {
		if options.Aria.Automation.Username, err = prompt.AskText(fmt.Sprintf("What is the username of the account for Aria Automation? Current (%s)", options.Aria.Automation.Username)); err != nil {
			return err
		}
	}

	if options.Aria.Automation.Password == "" {
		if options.Aria.Automation.Password, err = prompt.AskPass("What is the password of the account for Aria Automation?"); err != nil {
			return err
		}
	}

	if options.Aria.Automation.OrgName == "" {
		if options.Aria.Automation.OrgName, err = prompt.AskText(fmt.Sprintf("What is the org name used in Aria Automation? Current (%s)", options.Aria.Automation.OrgName)); err != nil {
			return err
		}
	}

	if options.Aria.Automation.ProjectName == "" {
		if options.Aria.Automation.ProjectName, err = prompt.AskText(fmt.Sprintf("What is the default project name in Aria Automation you want to push automation code to? Current (%s)", options.Aria.Automation.ProjectName)); err != nil {
			return err
		}
	}

	return nil
}

// ValidateArtifactManagerArguments will prompt for missing artifact manager options
func (options *RunOptions) ValidateArtifactManagerArguments() error {
	var err error
	if options.ArtifactManager.ReleaseRepo == "" {
		if options.ArtifactManager.ReleaseRepo, err = prompt.AskText("Artifact Manager setup partially passed. What is the release repository:"); err != nil {
			return err
		}
	}

	if options.ArtifactManager.SnapshotRepo == "" {
		if options.ArtifactManager.SnapshotRepo, err = prompt.AskText("Artifact Manager setup partially passed. What is the snapshot repository:"); err != nil {
			return err
		}
	}

	if options.ArtifactManager.GroupRepo == "" {
		if options.ArtifactManager.GroupRepo, err = prompt.AskText("Artifact Manager setup partially passed. What is the group repository:"); err != nil {
			return err
		}
	}

	if options.ArtifactManager.Password == "" {
		if options.ArtifactManager.Password, err = prompt.AskPass("Artifact Manager setup partially passed. What is the password used to authenticate to the artifact manager?"); err != nil {
			return err
		}
	}

	return nil
}
