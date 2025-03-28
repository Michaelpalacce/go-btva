package args

import (
	"flag"
	"fmt"

	"github.com/Michaelpalacce/go-btva/pkg/prompt"
)

// Args will parse the CLI arguments once and return the parsed options from then on
// This will panic if there are any validation issues
func Args() *Options {
	if options.parsed {
		return options
	}

	// Software
	flag.BoolVar(&options.Software.InstallJava, "installJava", true, "Flag that will specify if you want to install the correct java version. If java already exists, regardless of which version, this is ignored")
	flag.StringVar(&options.Software.JavaVersion, "javaVersion", "17", "Which version of java to install?")

	flag.BoolVar(&options.Software.InstallMvn, "installMvn", true, "Flag that will specify if you want to install the correct mvn version. If mvn already exists, regardless of which version, this is ignored")
	flag.StringVar(&options.Software.MvnVersion, "mvnVersion", "3.9.9", "Which version of mvn to install?")

	flag.BoolVar(&options.Software.InstallNode, "installNode", true, "Flag that will specify if you want to install the correct node version. This will work by installing fnm locally.")
	flag.StringVar(&options.Software.NodeVersion, "nodeVersion", "22", "Which version of node to install?")

	// Local
	flag.BoolVar(&options.Local.SetupM2, "setupM2", true, "Do you want to overwrite your current ~/.m2/settings.xml file with a proposed configuration from the tool?")

	// Infra
	flag.BoolVar(&options.Infra.MinimalInfrastructure, "minimalInfrastructure", true, "Do you want to spin up a mininmal infrastructure example?")

	flag.StringVar(&options.Infra.SSHVMIP, "sshVmIp", "", "IP of the VM where to setup the minimal infrastructure example.")
	flag.StringVar(&options.Infra.SSHUsername, "sshUsername", "root", "Username of the user to ssh with. This MUST be a root user or a user that can sudo without a password.")
	flag.StringVar(&options.Infra.SSHPassword, "sshPassword", "", "Password of the user to ssh with. Either this or sshPrivateKey must be provided.")
	flag.StringVar(&options.Infra.SSHPrivateKey, "sshPrivateKey", "", "Private key to use for authentication. Either this or sshPassword must be provided.")
	flag.StringVar(&options.Infra.SSHPrivateKeyPassphrase, "sshPrivateKeyPassphrase", "", "Passphrase for the private key if any. Optional")

	flag.Parse()

	if err := validate(options); err != nil {
		panic(err)
	}

	options.parsed = true

	return options
}

// validate will validate the options and return an error if something is wrong
func validate(options *Options) error {
	if options.Infra.MinimalInfrastructure {
		var err error

		if options.Infra.SSHVMIP == "" {
			fmt.Print("MinimalInfrastructure selected, but you did not provide sshVmIp, please type in the IP: ")
			if options.Infra.SSHVMIP, err = prompt.AskText(); err != nil {
				return fmt.Errorf("sshVmIp must be provided. Err: %w", err)
			}
		}

		if options.Infra.SSHPrivateKey == "" && options.Infra.SSHPassword == "" {
			fmt.Print("MinimalInfrastructure selected, but you did not provide sshPassword or sshPrivateKey, please type in password: ")
			if options.Infra.SSHPassword, err = prompt.AskPass(); err != nil {
				return fmt.Errorf("sshPassword must be provided. Err: %w", err)
			}
			fmt.Println("")
		}

		if options.Infra.SSHUsername == "" {
			fmt.Print("MinimalInfrastructure selected, but you did not provide sshUsername, please type in the username of root or a passwordless sudo user")
			if options.Infra.SSHUsername, err = prompt.AskText(); err != nil {
				return fmt.Errorf("sshUsername must be provided. Err: %w", err)
			}
		}
	}

	return nil
}
