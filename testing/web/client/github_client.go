package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Object struct {
	Sha  string `json:"sha"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Tag struct {
	Ref    string `json:"ref"`
	NodeId string `json:"node_id"`
	Url    string `json:"url"`
	Object Object `json:"object"`
}

type GitHubClient struct {
	baseUrl string
}

func New() GitHubClient {
	return NewWithBaseUrl("https://api.github.com")
}

func NewWithBaseUrl(baseUrl string) GitHubClient {
	return GitHubClient{
		baseUrl: baseUrl,
	}
}

func (g GitHubClient) GetTags(owner, repo string) ([]Tag, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/git/refs/tags", g.baseUrl, owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("network issue while getting tags: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("unable to ready response body: %s", err)
		}
		return nil, fmt.Errorf("%s returned %d: %s", url, resp.StatusCode, string(body))
	}
	var tags []Tag
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return nil, fmt.Errorf("unable to decode json response: %w", err)
	}
	return tags, nil
}
