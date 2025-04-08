package orchestrator

// WithOptions will set the correct state based on the set Options.
// @WARN: THIS RESETS THE STATE OF THE ORCHESTRATOR
// @TODO: It'd be great if there is a better way to handle this.
func WithOptions() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		o.Reset()

		// Software validation is inside. Will skip if needed
		o.Tasks(
			WithAllSoftware(),
		)

		if o.State.Options.MinimalInfra.MinimalInfrastructureGitlab || o.State.Options.MinimalInfra.MinimalInfrastructureNexus {
			o.Tasks(
				WithPartialMinimalInfrastructureSetup(),
			)

			if o.State.Options.MinimalInfra.MinimalInfrastructureNexus {
				o.Tasks(
					WithPartialMinimalInfrastructureNexus(),
					WithPartialMinimalInfrastructureSettingsXml(),
				)
			}

			if o.State.Options.MinimalInfra.MinimalInfrastructureGitlab {
				o.Tasks(
					WithPartialMinimalInfrastructureGitlab(),
				)
			}
		}

		if o.State.Options.ArtifactManager.Password != "" || o.State.Options.ArtifactManager.ReleaseRepo != "" || o.State.Options.ArtifactManager.GroupRepo != "" || o.State.Options.ArtifactManager.SnapshotRepo != "" {
			o.Tasks(
				WithSettingsXml(),
			)
		}

		return nil
	}
}
