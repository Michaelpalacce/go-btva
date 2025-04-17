package run

import (
	"flag"
	"os"

	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/pkg/args"
)

// This is a single instance of the runOptions. We don't want to parse them more than once
var runOptions = &options.RunOptions{
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
func (c *RunCommand) Args() *options.RunOptions {
	if runOptions.Parsed {
		return runOptions
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

	// Interactive
	args.AddVar(&runOptions.Prompt, "prompt", "p", true, "The tool will prompt you for most actions that will be taken (mainly things that affect your hardware). Actions can be Yes-ed, No-ed or Aborted.")

	// Software
	args.AddVar(&runOptions.Software.JavaVersion, "javaVersion", "jv", "17", "Which version of java to install? If not set, will skip installation.")
	args.AddVar(&runOptions.Software.NodeVersion, "nodeVersion", "nv", "22", "Which version of node to install? If not set, will skip installation.")
	args.AddVar(&runOptions.Software.MvnVersion, "mvnVersion", "mv", "3.9.9", "Which version of mvn to install? If not set, will skip installation.")
	args.AddVar(&runOptions.Software.VsCodeVersion, "vsCodeVersion", "vv", "latest", "Which version of VSCode to install? Supports `latest` or empty. If not set, will skip installation.")
	args.AddVar(&runOptions.Software.ChocolateyVersion, "chocolateyVersion", "cv", "latest", "(WINDOWS ONLY) Which version of chocolatey to install? Supports `latest` or empty. If not set, will skip installation.")

	// Minimal Infra

	args.AddVar(&runOptions.MinimalInfra.MinimalInfrastructureGitlab, "minimalInfrastructureGitlab", "mig", false, "Do you want to spin up minimal infrastructure Gitlab?")
	args.AddVar(&runOptions.MinimalInfra.MinimalInfrastructureNexus, "minimalInfrastructureNexus", "min", false, "Do you want to spin up only minimal infrastructure nexus? Modifies settings.xml with Nexus settings.")

	args.AddVar(&runOptions.MinimalInfra.SSHVMIP, "sshVmIp", "", "", "IP of the VM where to setup the minimal infrastructure example.")

	args.AddVar(&runOptions.MinimalInfra.SSHUsername, "sshUsername", "", "root", "Username of the user to ssh with. This MUST be a root user or a user that can sudo without a password.")
	args.AddVar(&runOptions.MinimalInfra.SSHPassword, "sshPassword", "", "", "Password of the user to ssh with. Either this or sshPrivateKey must be provided.")
	args.AddVar(&runOptions.MinimalInfra.SSHPrivateKey, "sshPrivateKey", "", "", "Private key to use for authentication. Either this or sshPassword must be provided.")
	args.AddVar(&runOptions.MinimalInfra.SSHPrivateKeyPassphrase, "sshPrivateKeyPassphrase", "", "", "Passphrase for the private key if any. Optional")
	args.AddVar(&runOptions.MinimalInfra.DockerUsername, "dockerUsername", "", "", "Docker username to use when setting up the minimal infra.")
	args.AddVar(&runOptions.MinimalInfra.DockerPAT, "dockerPat", "", "", "Docker Public Access Token to use when setting up the minimal infra. If dockerUsername is provided and this isn't, you will be prompted.")

	// Aria Automation
	args.AddVar(&runOptions.Aria.Automation.Port, "ariaAutomationFqdn", "aaFQDN", "vra-l-01a.corp.local", "Fully Qualified Domain Name for Aria Automation without the protocol (https://) and port (:443).")
	args.AddVar(&runOptions.Aria.Automation.FQDN, "ariaAutomationPort", "aaPort", "443", "Aria Automation port")
	args.AddVar(&runOptions.Aria.Automation.Username, "ariaAutomationUsername", "aaUsername", "configurationadmin", "Username to use for authentication to Aria Automation.")
	args.AddVar(&runOptions.Aria.Automation.Password, "ariaAutomationPassword", "aaPassword", "", "Password to use for authentication to Aria Automation.")
	args.AddVar(&runOptions.Aria.Automation.OrgName, "ariaAutomationOrgName", "aaOrgName", "vidm-l-01a", "Aria Automation organization name. Can be found in the dropdown at the top.")
	args.AddVar(&runOptions.Aria.Automation.ProjectName, "ariaAutomationProjectName", "aaProjectName", "Development", "Aria Automation default project name to push to. Used mainly for vra-ng archetype.")

	// Artifactory
	args.AddVar(&runOptions.ArtifactManager.Password, "artifactManagerPassword", "amPassword", "", "Password for existing ArtifactManager.")
	args.AddVar(&runOptions.ArtifactManager.GroupRepo, "artifactManagerGroupRepo", "amGroupRepo", "", "Group repository to use for fetching both snapshots and release artifacts.")
	args.AddVar(&runOptions.ArtifactManager.ReleaseRepo, "artifactManagerReleaseRepo", "amReleaseRepo", "", "Release repository to use for fetching and uploading release artifacts.")
	args.AddVar(&runOptions.ArtifactManager.SnapshotRepo, "artifactManagerSnapshotRepo", "amSnapshotRepo", "", "Snapshot repository to use for fetching and uploading sanpshot artifacts.")

	args.Parse()

	runOptions.Parsed = true

	return runOptions
}
