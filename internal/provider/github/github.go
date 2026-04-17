package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/isyuricunha/gitflect/internal/provider"
)

type Client struct {
	token      string
	user       string
	httpClient *http.Client
}

func New(token, user string) *Client {
	return &Client{
		token: token,
		user:  user,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type ghRepo struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func (c *Client) ListRepos(visibility string) ([]provider.Repo, error) {
	var all []provider.Repo
	page := 1

	for {
		url := fmt.Sprintf("https://api.github.com/user/repos?per_page=100&page=%d&affiliation=owner", page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("github req create: %w", err)
		}
		
		req.Header.Set("Authorization", "token "+c.token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("github list: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("github api error: status %d", resp.StatusCode)
		}

		var repos []ghRepo
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("github decode: %w", err)
		}
		resp.Body.Close()

		if len(repos) == 0 {
			break
		}

		for _, r := range repos {
			if visibility == "public" && r.Private {
				continue
			}
			if visibility == "private" && !r.Private {
				continue
			}
			all = append(all, provider.Repo{
				Name:    r.Name,
				Private: r.Private,
				CloneURL: fmt.Sprintf(
					"https://%s@github.com/%s/%s.git",
					c.token, c.user, r.Name,
				),
			})
		}
		page++
	}
	return all, nil
}
