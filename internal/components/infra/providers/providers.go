package providers

type Provider interface {
	Init() error

	Artifactory() error
	Build() error

	Disable() error
}
