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
	flag.StringVar(&options.Software.JavaVersion, "javaVersion", "17", "Which version of java to install? If not set, will skip installation.")
	flag.StringVar(&options.Software.MvnVersion, "mvnVersion", "3.9.9", "Which version of mvn to install? If not set, will skip installation.")
	flag.StringVar(&options.Software.NodeVersion, "nodeVersion", "22", "Which version of node to install? If not set, will skip installation.")

	// Infra
	flag.BoolVar(&options.MinimalInfra.MinimalInfrastructureGitlab, "minimalInfrastructureGitlab", true, "Do you want to spin up minimal infrastructure Gitlab?")
	flag.BoolVar(&options.MinimalInfra.MinimalInfrastructureNexus, "minimalInfrastructureNexus", false, "Do you want to spin up only minimal infrastructure nexus?")

	flag.StringVar(&options.MinimalInfra.SSHVMIP, "sshVmIp", "", "IP of the VM where to setup the minimal infrastructure example.")
	flag.StringVar(&options.MinimalInfra.SSHUsername, "sshUsername", "root", "Username of the user to ssh with. This MUST be a root user or a user that can sudo without a password.")
	flag.StringVar(&options.MinimalInfra.SSHPassword, "sshPassword", "", "Password of the user to ssh with. Either this or sshPrivateKey must be provided.")
	flag.StringVar(&options.MinimalInfra.SSHPrivateKey, "sshPrivateKey", "", "Private key to use for authentication. Either this or sshPassword must be provided.")
	flag.StringVar(&options.MinimalInfra.SSHPrivateKeyPassphrase, "sshPrivateKeyPassphrase", "", "Passphrase for the private key if any. Optional")
	flag.StringVar(&options.MinimalInfra.DockerUsername, "dockerUsername", "", "Docker username to use when setting up the minimal infra.")
	flag.StringVar(&options.MinimalInfra.DockerPAT, "dockerPat", "", "Docker Public Access Token to use when setting up the minimal infra. If dockerUsername is provided and this isn't, you will be prompted.")

	// Aria

	// Automation
	flag.StringVar(&options.Aria.Automation.Port, "ariaAutomationFqdn", "vra-l-01a.corp.local", "Fully Qualified Domain Name for Aria Automation without the protocol (https://) and port (:443).")
	flag.StringVar(&options.Aria.Automation.FQDN, "ariaAutomationPort", "443", "Aria Automation port")
	flag.StringVar(&options.Aria.Automation.Username, "ariaAutomationUsername", "configurationadmin", "Username to use for authentication to Aria Automation.")
	flag.StringVar(&options.Aria.Automation.Password, "ariaAutomationPassword", "", "Password to use for authentication to Aria Automation.")
	flag.StringVar(&options.Aria.Automation.OrgName, "ariaAutomationOrgName", "vidm-l-01a", "Aria Automation organization name. Can be found in the dropdown at the top.")
	flag.StringVar(&options.Aria.Automation.ProjectName, "ariaAutomationProjectName", "Development", "Aria Automation default project name to push to. Used mainly for vra-ng archetype.")

	flag.Parse()

	options.parsed = true

	return options
}
