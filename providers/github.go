package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubProvider struct {
	GenericProvider
}

type githubUser struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

/*
Github() creates a github oauth client.

The client ID and client secret are taken from the environment variables GITHUB_CLIENT_ID &
GITHUB_CLIENT_SECRET
*/
func Github(baseUrl string) Provider {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
		RedirectURL:  baseUrl + "callback/",
		Scopes:       []string{"user:email"},
	}

	return &GithubProvider{
		GenericProvider: GenericProvider{
			Name:   "github",
			Config: config,
		},
	}
}

func (p *GithubProvider) GetUserInfo(client *http.Client) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github return error.  %d %s", resp.StatusCode, resp.Status)
	}

	// fields we care about include id, login, name, email
	dec := json.NewDecoder(resp.Body)
	var user githubUser
	err = dec.Decode(&user)

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		PROVIDER_EMAIL:    user.Email,
		PROVIDER_USERNAME: user.Login,
		PROVIDER_ID:       user.Id,
		PROVIDER_NAME:     user.Name,
	}
	return result, nil
}
