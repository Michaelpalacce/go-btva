package args

import (
	"flag"
)

// Args will parse the CLI arguments once and return the parsed options from then on
// This will panic if there are any validation issues
func Args() *Options {
	if options.parsed {
		return options
	}

	// Software
	AddNewVar(&options.Software.JavaVersion, "javaVersion", "jv", "17", "Which version of java to install? If not set, will skip installation.")
	AddNewVar(&options.Software.NodeVersion, "nodeVersion", "nv", "22", "Which version of node to install? If not set, will skip installation.")
	AddNewVar(&options.Software.VsCodeVersion, "vsCodeVersion", "vv", "latest", "Which version of node to install? Supports `latest` or empty. If not set, will skip installation.")

	// Minimal Infra

	AddNewVar(&options.MinimalInfra.MinimalInfrastructureGitlab, "minimalInfrastructureGitlab", "mig", false, "Do you want to spin up minimal infrastructure Gitlab?")
	AddNewVar(&options.MinimalInfra.MinimalInfrastructureNexus, "minimalInfrastructureNexus", "min", false, "Do you want to spin up only minimal infrastructure nexus? Modifies settings.xml with Nexus settings.")

	AddNewVar(&options.MinimalInfra.SSHVMIP, "sshVmIp", "", "", "IP of the VM where to setup the minimal infrastructure example.")

	AddNewVar(&options.MinimalInfra.SSHUsername, "sshUsername", "", "root", "Username of the user to ssh with. This MUST be a root user or a user that can sudo without a password.")
	AddNewVar(&options.MinimalInfra.SSHPassword, "sshPassword", "", "", "Password of the user to ssh with. Either this or sshPrivateKey must be provided.")
	AddNewVar(&options.MinimalInfra.SSHPrivateKey, "sshPrivateKey", "", "", "Private key to use for authentication. Either this or sshPassword must be provided.")
	AddNewVar(&options.MinimalInfra.SSHPrivateKeyPassphrase, "sshPrivateKeyPassphrase", "", "", "Passphrase for the private key if any. Optional")
	AddNewVar(&options.MinimalInfra.DockerUsername, "dockerUsername", "", "", "Docker username to use when setting up the minimal infra.")
	AddNewVar(&options.MinimalInfra.DockerPAT, "dockerPat", "", "", "Docker Public Access Token to use when setting up the minimal infra. If dockerUsername is provided and this isn't, you will be prompted.")

	// Aria Automation
	AddNewVar(&options.Aria.Automation.Port, "ariaAutomationFqdn", "aaFQDN", "vra-l-01a.corp.local", "Fully Qualified Domain Name for Aria Automation without the protocol (https://) and port (:443).")
	AddNewVar(&options.Aria.Automation.FQDN, "ariaAutomationPort", "aaPort", "443", "Aria Automation port")
	AddNewVar(&options.Aria.Automation.Username, "ariaAutomationUsername", "aaUsername", "configurationadmin", "Username to use for authentication to Aria Automation.")
	AddNewVar(&options.Aria.Automation.Password, "ariaAutomationPassword", "aaPassword", "", "Password to use for authentication to Aria Automation.")
	AddNewVar(&options.Aria.Automation.OrgName, "ariaAutomationOrgName", "aaOrgName", "vidm-l-01a", "Aria Automation organization name. Can be found in the dropdown at the top.")
	AddNewVar(&options.Aria.Automation.ProjectName, "ariaAutomationProjectName", "aaProjectName", "Development", "Aria Automation default project name to push to. Used mainly for vra-ng archetype.")

	// Artifactory
	AddNewVar(&options.ArtifactManager.Password, "artifactManagerPassword", "amPassword", "", "Password for existing ArtifactManager.")
	AddNewVar(&options.ArtifactManager.GroupRepo, "artifactManagerGroupRepo", "amGroupRepo", "", "Group repository to use for fetching both snapshots and release artifacts.")
	AddNewVar(&options.ArtifactManager.ReleaseRepo, "artifactManagerReleaseRepo", "amReleaseRepo", "", "Release repository to use for fetching and uploading release artifacts.")
	AddNewVar(&options.ArtifactManager.SnapshotRepo, "artifactManagerSnapshotRepo", "amSnapshotRepo", "", "Snapshot repository to use for fetching and uploading sanpshot artifacts.")

	flag.Usage = Usage

	flag.Parse()

	options.parsed = true

	return options
}
