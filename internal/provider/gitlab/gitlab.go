package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	token      string
	user       string
	baseURL    string
	httpClient *http.Client
}

func New(token, user, baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	// Ensure it doesn't end with a trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		token:      token,
		user:       user,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) EnsureRepo(name string, private bool) (string, error) {
	visibility := "public"
	if private {
		visibility = "private"
	}

	body, _ := json.Marshal(map[string]any{
		"name":       name,
		"path":       name,
		"visibility": visibility,
	})

	apiURL := c.baseURL + "/api/v4/projects"
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("gitlab req create error: %w", err)
	}
	
	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gitlab create project request: %w", err)
	}
	defer resp.Body.Close()
	
	// Ignore errors that are not related to "already exists" conflicts.
	// HTTP 201 = Created, HTTP 400 = Bad Request (often implies it already exists in the API v4)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gitlab api exception: status=%d res=%s", resp.StatusCode, string(b))
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid gitlab url: %w", err)
	}

	pushURL := fmt.Sprintf("%s://oauth2:%s@%s/%s/%s.git", u.Scheme, c.token, u.Host, c.user, name)
	return pushURL, nil
}
