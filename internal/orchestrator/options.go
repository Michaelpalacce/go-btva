package orchestrator

// WithOptions will set the correct state based on the set Options.
// @WARN: THIS RESETS THE STATE OF THE ORCHESTRATOR
func WithOptions() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		o.Reset()

		if o.Options.Infra.MinimalInfrastructure == true {
			o.Tasks(
				WithFullMinimalInfrastructure(),
			)
		}

		return nil
	}
}
