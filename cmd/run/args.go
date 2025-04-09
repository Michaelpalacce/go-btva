package run

import (
	"flag"
	"os"

	"github.com/Michaelpalacce/go-btva/cmd/run/run_options"
	"github.com/Michaelpalacce/go-btva/internal/args"
)

// This is a single instance of the options. We don't want to parse them more than once
var options = &run_options.RunOptions{
	Parsed: false,
}

var usage = `go-btva run is used for:
- [x] Setup local environment
- [x] Install needed software
- [x] Setup minimal infrastructure on a linux machine
- [x] Connect local env with existing Artifact Manager
`

var examples = `
# Basic Usage

go-btva comes with a lot of defaults. Information will be asked of you if it's needed.

'''
go-btva run # This will run the tool with all the defaults.
'''

# Skipping software installation

'''
go-btva run -jv="" -mv="" -vv="" -nv=""
'''

# With Minimal Infrastructure

'''
go-btva run -mig -min
'''

# With Custom Artifact Manager

'''
go-btva run --artifactManagerGroupRepo=REPO --artifactManagerReleaseRepo=REPO --artifactManagerSnapshotRepo=REPO
'''
`

// Args will parse the CLI arguments once and return the parsed options from then on
// This will panic if there are any validation issues
func (c *RunCommand) Args() *run_options.RunOptions {
	if options.Parsed {
		return options
	}

	args, err := args.NewArgs(
		os.Args[2:],
		args.WithExamples(examples),
		args.WithUsage(usage),
		args.WithFs(flag.NewFlagSet("run", flag.ExitOnError)),
	)
	if err != nil {
		panic(err)
	}

	// Software
	args.AddVar(&options.Software.JavaVersion, "javaVersion", "jv", "17", "Which version of java to install? If not set, will skip installation.")
	args.AddVar(&options.Software.NodeVersion, "nodeVersion", "nv", "22", "Which version of node to install? If not set, will skip installation.")
	args.AddVar(&options.Software.MvnVersion, "mvnVersion", "mv", "3.9.9", "Which version of mvn to install? If not set, will skip installation.")
	args.AddVar(&options.Software.VsCodeVersion, "vsCodeVersion", "vv", "latest", "Which version of node to install? Supports `latest` or empty. If not set, will skip installation.")

	// Minimal Infra

	args.AddVar(&options.MinimalInfra.MinimalInfrastructureGitlab, "minimalInfrastructureGitlab", "mig", false, "Do you want to spin up minimal infrastructure Gitlab?")
	args.AddVar(&options.MinimalInfra.MinimalInfrastructureNexus, "minimalInfrastructureNexus", "min", false, "Do you want to spin up only minimal infrastructure nexus? Modifies settings.xml with Nexus settings.")

	args.AddVar(&options.MinimalInfra.SSHVMIP, "sshVmIp", "", "", "IP of the VM where to setup the minimal infrastructure example.")

	args.AddVar(&options.MinimalInfra.SSHUsername, "sshUsername", "", "root", "Username of the user to ssh with. This MUST be a root user or a user that can sudo without a password.")
	args.AddVar(&options.MinimalInfra.SSHPassword, "sshPassword", "", "", "Password of the user to ssh with. Either this or sshPrivateKey must be provided.")
	args.AddVar(&options.MinimalInfra.SSHPrivateKey, "sshPrivateKey", "", "", "Private key to use for authentication. Either this or sshPassword must be provided.")
	args.AddVar(&options.MinimalInfra.SSHPrivateKeyPassphrase, "sshPrivateKeyPassphrase", "", "", "Passphrase for the private key if any. Optional")
	args.AddVar(&options.MinimalInfra.DockerUsername, "dockerUsername", "", "", "Docker username to use when setting up the minimal infra.")
	args.AddVar(&options.MinimalInfra.DockerPAT, "dockerPat", "", "", "Docker Public Access Token to use when setting up the minimal infra. If dockerUsername is provided and this isn't, you will be prompted.")

	// Aria Automation
	args.AddVar(&options.Aria.Automation.Port, "ariaAutomationFqdn", "aaFQDN", "vra-l-01a.corp.local", "Fully Qualified Domain Name for Aria Automation without the protocol (https://) and port (:443).")
	args.AddVar(&options.Aria.Automation.FQDN, "ariaAutomationPort", "aaPort", "443", "Aria Automation port")
	args.AddVar(&options.Aria.Automation.Username, "ariaAutomationUsername", "aaUsername", "configurationadmin", "Username to use for authentication to Aria Automation.")
	args.AddVar(&options.Aria.Automation.Password, "ariaAutomationPassword", "aaPassword", "", "Password to use for authentication to Aria Automation.")
	args.AddVar(&options.Aria.Automation.OrgName, "ariaAutomationOrgName", "aaOrgName", "vidm-l-01a", "Aria Automation organization name. Can be found in the dropdown at the top.")
	args.AddVar(&options.Aria.Automation.ProjectName, "ariaAutomationProjectName", "aaProjectName", "Development", "Aria Automation default project name to push to. Used mainly for vra-ng archetype.")

	// Artifactory
	args.AddVar(&options.ArtifactManager.Password, "artifactManagerPassword", "amPassword", "", "Password for existing ArtifactManager.")
	args.AddVar(&options.ArtifactManager.GroupRepo, "artifactManagerGroupRepo", "amGroupRepo", "", "Group repository to use for fetching both snapshots and release artifacts.")
	args.AddVar(&options.ArtifactManager.ReleaseRepo, "artifactManagerReleaseRepo", "amReleaseRepo", "", "Release repository to use for fetching and uploading release artifacts.")
	args.AddVar(&options.ArtifactManager.SnapshotRepo, "artifactManagerSnapshotRepo", "amSnapshotRepo", "", "Snapshot repository to use for fetching and uploading sanpshot artifacts.")

	args.Parse()

	options.Parsed = true

	return options
}
