package args

import (
	"flag"
)

// Args will parse the CLI arguments once and return the parsed options from then on
func Args() *Options {
	if options.parsed {
		return options
	}

	// Software
	flag.BoolVar(&options.Software.InstallJava, "installJava", true, "Flag that will specify if you want to install the correct java version. If java already exists, regardless of which version, this is ignored")
	flag.StringVar(&options.Software.LinuxJavaVersion, "linuxJavaVersion", "17", "Which version of java to install?")

	flag.BoolVar(&options.Software.InstallMvn, "installMvn", true, "Flag that will specify if you want to install the correct mvn version. If mvn already exists, regardless of which version, this is ignored")
	flag.StringVar(&options.Software.LinuxMvnVersion, "linuxMvnVersion", "3.9.9", "Which version of mvn to install?")

	flag.BoolVar(&options.Software.InstallNode, "installNode", true, "Flag that will specify if you want to install the correct node version. This will work by installing fnm locally.")
	flag.StringVar(&options.Software.LinuxNodeVersion, "linuxNodeVersion", "22", "Which version of node to install?")

	// Local
	flag.BoolVar(&options.Local.SetupM2, "setupM2", true, "Do you want to overwrite your current ~/.m2/settings.xml file with a proposed configuration from the tool?")
	flag.BoolVar(&options.Local.SaveState, "saveState", true, "Do you want to preserve the state of the machine execution? Allows you to resume in case of failures.")
	flag.StringVar(&options.Local.StateJson, "stateJson", "state.json5", "The file to store the current state in")

	// Infra
	flag.BoolVar(&options.Infra.MinimalInfrastructure, "minimalInfrastructure", true, "Do you want to spin up a mininmal infrastructure example?")

	flag.StringVar(&options.Infra.SSHVMIP, "sshVmIp", "", "IP of the VM where to setup the minimal infrastructure example.")
	flag.StringVar(&options.Infra.SSHUsername, "sshUsername", "root", "Username of the user to ssh with. This MUST be a root user or a user that can sudo without a password.")
	flag.StringVar(&options.Infra.SSHPassword, "sshPassword", "", "Password of the user to ssh with. Either this or sshPrivateKey must be provided.")
	flag.StringVar(&options.Infra.SSHPrivateKey, "sshPrivateKey", "", "Private key to use for authentication. Either this or sshPassword must be provided.")
	flag.StringVar(&options.Infra.SSHPrivateKeyPassphrase, "sshPrivateKeyPassphrase", "", "Passphrase for the private key if any. Optional")

	flag.Parse()

	if options.Infra.MinimalInfrastructure {
		if options.Infra.SSHPrivateKey == "" && options.Infra.SSHPassword == "" {
			panic("Either sshPrivateKey or sshPassword must be provided")
		}

		if options.Infra.SSHVMIP == "" {
			panic("sshVmIp must be provided")
		}

		if options.Infra.SSHUsername == "" {
			panic("sshUsername must be provided")
		}
	}

	options.parsed = true

	return options
}
