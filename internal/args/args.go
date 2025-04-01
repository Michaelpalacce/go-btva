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
	flag.BoolVar(&options.Infra.MinimalInfrastructure, "minimalInfrastructure", true, "Do you want to spin up a mininmal infrastructure example?")

	flag.StringVar(&options.Infra.SSHVMIP, "sshVmIp", "", "IP of the VM where to setup the minimal infrastructure example.")
	flag.StringVar(&options.Infra.SSHUsername, "sshUsername", "root", "Username of the user to ssh with. This MUST be a root user or a user that can sudo without a password.")
	flag.StringVar(&options.Infra.SSHPassword, "sshPassword", "", "Password of the user to ssh with. Either this or sshPrivateKey must be provided.")
	flag.StringVar(&options.Infra.SSHPrivateKey, "sshPrivateKey", "", "Private key to use for authentication. Either this or sshPassword must be provided.")
	flag.StringVar(&options.Infra.SSHPrivateKeyPassphrase, "sshPrivateKeyPassphrase", "", "Passphrase for the private key if any. Optional")
	flag.StringVar(&options.Infra.DockerUsername, "dockerUsername", "", "Docker username to use when setting up the minimal infra.")
	flag.StringVar(&options.Infra.DockerPAT, "dockerPat", "", "Docker Public Access Token to use when setting up the minimal infra. If dockerUsername is provided and this isn't, you will be prompted.")

	flag.Parse()

	options.parsed = true

	return options
}
