package infra

import (
	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type InfraComponent struct {
	os      *os.OS
	state   *state.State
	options *options.RunOptions
}

func NewInfraComponent(os *os.OS, state *state.State) *InfraComponent {
	return &InfraComponent{os: os, state: state, options: state.Options}
}

const (
	INFRA_STATE = "Infra"

	// Public
	INFRA_GITLAB_ADMIN_PASSWORD_KEY = "gitlabPassword"
	INFRA_GITLAB_ADMIN_PAT_KEY      = "gitlabPat"
	INFRA_NEXUS_PASSWORD_KEY        = "nexusPassword"
)
