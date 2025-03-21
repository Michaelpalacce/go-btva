package args

import "flag"

func ParseArguments() Options {
	options := Options{
		Software: Software{},
		Mvn:      Mvn{},
	}

	// Software
	flag.BoolVar(&options.Software.InstallJava, "installJava", true, "Flag that will specify if you want to install the correct java version. Note: This will replace your current installation, so be carefull")
	flag.BoolVar(&options.Software.InstallMvn, "installMvn", true, "Flag that will specify if you want to install the correct mvn version. Note: This will replace your current installation, so be carefull")
	flag.BoolVar(&options.Software.InstallNode, "installNode", true, "Flag that will specify if you want to install the correct node version. This will work by installing fnm locally")

	// MVN
	flag.BoolVar(&options.Mvn.SetupM2, "setupM2", true, "Do you want to overwrite your current ~/.m2/settings.xml file with a proposed configuration from the tool?")

	// Infra
	flag.BoolVar(&options.Infra.MinimalInfrastructure, "minimalInfrastructure", true, "Do you want to spin up a mininmal infrastructure example?")

	flag.Parse()

	return options
}
