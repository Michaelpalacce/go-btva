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
	flag.StringVar(&options.Software.JavaLinuxPackage, "javaLinuxPackage", "openjdk-17-jdk", "Which package to install with apt?")
	flag.BoolVar(&options.Software.InstallMvn, "installMvn", true, "Flag that will specify if you want to install the correct mvn version. If mvn already exists, regardless of which version, this is ignored")
	flag.BoolVar(&options.Software.InstallNode, "installNode", true, "Flag that will specify if you want to install the correct node version. This will work by installing fnm locally.")

	// Local
	flag.BoolVar(&options.Local.SetupM2, "setupM2", true, "Do you want to overwrite your current ~/.m2/settings.xml file with a proposed configuration from the tool?")
	flag.BoolVar(&options.Local.SaveState, "saveState", true, "Do you want to preserve the state of the machine execution? Allows you to resume in case of failures.")
	flag.StringVar(&options.Local.StateJson, "stateJson", "state.json5", "The file to store the current state in")

	// Infra
	flag.BoolVar(&options.Infra.MinimalInfrastructure, "minimalInfrastructure", true, "Do you want to spin up a mininmal infrastructure example?")

	flag.Parse()

	options.parsed = true

	return options
}
