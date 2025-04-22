package infra

const (
	_BTVA_INSTALL_DIR_INFRA         = "/opt/build-tools-for-vmware-aria/infrastructure"
	_BTVA_MINIMAL_INFRA_INSTALL_URL = "https://raw.githubusercontent.com/vmware/build-tools-for-vmware-aria/refs/heads/main/infrastructure/install.sh"
)

// Scripts that will be executed at different stages for the minimal infra setup
const (
	_MINIMAL_INFRA_INSTALL_SCRIPT = "curl -o- %s | bash -s -- %s %q"
	_FIX_GITLAB_EXTERNAL_URL      = "sed -i \"s|external_url 'http://infra.corp.local/gitlab'|external_url '%q'|\" %s/docker-compose.yml"
	_WAIT_FOR_MINIMAL_INFRA_UP    = "docker compose -f %s/docker-compose.yml up -d --wait"

	_GET_GITLAB_PASSWORD = "docker exec gitlab-ce test -f /etc/gitlab/initial_root_password && docker exec gitlab-ce grep 'Password:' /etc/gitlab/initial_root_password | awk '{print $2}'"
	_GET_GITLAB_PAT      = "docker exec gitlab-ce gitlab-rails runner 'token = User.find_by_username(\"root\").personal_access_tokens.create(scopes: [:read_user, :read_repository, :api, :create_runner, :manage_runner, :sudo, :admin_mode], name: \"Automation token\", expires_at: 356.days.from_now); token.set_token(\"%s\"); token.save! '"

	_REGISTER_GITLAB_RUNNER = "docker exec gitlab-runner gitlab-runner register --non-interactive --url \"%s\" --token \"%s\" --executor \"docker\" --docker-image alpine:latest --description \"docker-runner\""

	_GET_NEXUS_PASSWORD = "docker exec nexus cat /nexus-data/admin.password"
)
