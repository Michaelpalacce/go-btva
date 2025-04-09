package orchestrator

// WithOptions will set the correct state based on the set Options.
// @WARN: THIS RESETS THE STATE OF THE ORCHESTRATOR
// @TODO: It'd be great if there is a better way to handle this.
func WithOptions() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		o.Reset()

		opts := o.State.Options

		// Software validation is inside. Will skip if needed
		o.Tasks(WithAllSoftware())

		if opts.MinimalInfra.HasMinimalInfra() {
			o.Tasks(WithPartialMinimalInfrastructureSetup())

			if opts.MinimalInfra.HasNexus() {
				o.Tasks(
					WithPartialMinimalInfrastructureNexus(),
					WithPartialMinimalInfrastructureSettingsXml(),
				)
			}

			if opts.MinimalInfra.HasGitlab() {
				o.Tasks(WithPartialMinimalInfrastructureGitlab())
			}
		}

		if opts.ArtifactManager.IsPartial() {
			o.Tasks(WithSettingsXml())
		}

		return nil
	}
}
