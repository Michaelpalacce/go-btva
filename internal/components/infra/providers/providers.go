package providers

type Provider interface {
	// Init will handle any initializing tasks
	// This can include creating new clients and connections, spinning up containers, executing terraform to provisioni infra, etc
	Init() error

	// Artifactory will handle setting up the artifact management tool
	// This can include: fetching passwords, setting rest endpoints, creating repositories, etc
	Artifactory() error
	// Build will handle the build environment together with CI/CD setup
	// This cna include registering a new runner, fetcinh passwords, setting endpoints, etc
	Build() error

	// Disable will handle finalizing tasks
	// Closing connections, etc
	Disable() error
}
