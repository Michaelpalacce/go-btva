package args

// Software holds a list of all the software that needs to be installed, so that means that this struct should contain all the possible
// software the tool supports.
type Software struct {
	InstallJava bool `json:"installJava"`
	InstallMvn  bool `json:"installMvn"`
	InstallNode bool `json:"installNode"`
}

// Mvn will hold different configurations for mvn that may be needed
type Mvn struct {
	SetupM2 bool `json:"setupM2"`
}

// Infra will hold different infra decisions that need to be taken
type Infra struct {
	MinimalInfrastructure bool `json:"minimalInfrastructure"`
}

// Options is a struct for options that the tool can accept.
type Options struct {
	Software Software `json:"software"`
	Mvn      Mvn      `json:"mvn"`
	Infra    Infra    `json:"infra"`

	// parsed is an internal variable that tells us that the options have already been parsed and don't need a second go
	parsed bool
}

// This is a single instance of the options. We don't want to parse them more than once
var options = &Options{
	Software: Software{},
	Mvn:      Mvn{},

	parsed: true,
}
