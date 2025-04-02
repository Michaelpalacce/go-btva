package orchestrator

// WithOptions will set the correct state based on the set Options.
// @WARN: THIS RESETS THE STATE OF THE ORCHESTRATOR
func WithOptions() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		o.Reset()

		// Software validation is inside. Will skip if needed
		o.Tasks(
			WithAllSoftware(),
		)

		if o.Options.MinimalInfra.MinimalInfrastructureGitlab || o.Options.MinimalInfra.MinimalInfrastructureNexus {
			o.Tasks(
				WithPartialMinimalInfrastructureSetup(),
			)

			if o.Options.MinimalInfra.MinimalInfrastructureNexus {
				o.Tasks(
					WithPartialMinimalInfrastructureNexus(),
					WithPartialMinimalInfrastructureSettingsXml(),
				)
			}

			if o.Options.MinimalInfra.MinimalInfrastructureGitlab {
				o.Tasks(
					WithPartialMinimalInfrastructureGitlab(),
				)
			}
		}

		return nil
	}
}
