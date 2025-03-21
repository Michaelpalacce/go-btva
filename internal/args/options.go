package args

// Software holds a list of all the software that needs to be installed, so that means that this struct should contain all the possible
// software the tool supports.
type Software struct {
	InstallJava bool `json:"installJava"`
	InstallMvn  bool `json:"installMvn"`
	InstallNode bool `json:"installNode"`
}

type Mvn struct {
	SetupM2 bool `json:"setupM2"`
}

type Infra struct {
	MinimalInfrastructure bool `json:"minimalInfrastructure"`
}

// Options is a struct for options that the tool can accept.
type Options struct {
	Software Software `json:"software"`
	Mvn      Mvn      `json:"mvn"`
	Infra    Infra    `json:"infra"`
}
