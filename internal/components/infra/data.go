package infra

import "github.com/Michaelpalacce/go-btva/internal/state"

func GitlabAdminPassword(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY))
}

func GitlabPat(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PAT_KEY))
}

func GitlabRunnerAuthToken(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, _INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY))
}

func NexusAdminPassword(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_NEXUS_PASSWORD_KEY))
}
