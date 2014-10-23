package providers

import (
	"encoding/json"
	"fmt"
	"github.com/golang/oauth2"
	"net/http"
	"os"
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

func Github(baseUrl string) Provider {
	options := &oauth2.Options{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  baseUrl + "callback/",
		Scopes:       []string{"user"},
	}
	config, _ := oauth2.NewConfig(options, "https://github.com/login/oauth/authorize", "https://github.com/login/oauth/access_token")

	return &GithubProvider{
		GenericProvider: GenericProvider{
			Name:   "github",
			Config: config,
		},
	}
}

func (p *GithubProvider) GetUserInfo(t *oauth2.Transport) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := t.RoundTrip(req)

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
