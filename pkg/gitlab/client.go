package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GitlabClient struct {
	privateToken string
	baseUrl      string
}

func NewGitlabClient(baseUrl, pat string) *GitlabClient {
	return &GitlabClient{
		baseUrl:      baseUrl,
		privateToken: pat,
	}
}

type getRunnerAuthTokenRequest struct {
	RunnerType string `json:"runner_type"`
}

type getRunnerAuthTokenResponse struct {
	Token string `json:"token"`
}

// GetRunnerAuthToken will give you a token you can use with a gitlab runner to register it
// runnerType is usually instance_type
func (c *GitlabClient) GetRunnerAuthToken(runnerType string) (string, error) {
	url := "/api/v4/user/runners"
	resp, err := c.postBody(url, getRunnerAuthTokenRequest{RunnerType: runnerType})
	if err != nil {
		return "", fmt.Errorf("erro submitting body for %s, err was: %w", url, err)
	}

	defer resp.Body.Close()

	var res getRunnerAuthTokenResponse

	json.NewDecoder(resp.Body).Decode(&res)

	return res.Token, nil
}

func (c *GitlabClient) postBody(url string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling json data. err was %w", err)
	}

	request, err := http.NewRequest("POST", c.url(url), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making a new request. Err was %w", err)
	}

	request.Header.Add("content-type", "application/json")
	request.Header.Add("accepts", "application/json")

	return http.DefaultClient.Do(c.authRequest(request))
}

func (c *GitlabClient) url(path string) string {
	return fmt.Sprintf("%s%s", c.baseUrl, path)
}

func (c *GitlabClient) authRequest(request *http.Request) *http.Request {
	request.Header.Add("Private-Token", c.privateToken)
	return request
}
