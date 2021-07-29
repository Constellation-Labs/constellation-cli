package gh

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Release struct {
	Url string `json:"url"`
	Name string `json:"name"`
	TagName string `json:"tag_name"`
	IsDraft bool `json:"draft"`
	IsPrerelease bool `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
}

type Github interface {
	GetReleases() (*[]Release)
	GetLatestCliRelease() (*Release)
	GetLatestUpdaterRelease (*Release)
}

type github struct {
	organization string
	repository string
}

func NewGithub(organization string, repository string) *github {
	return & github { organization, repository }
}

func (gh *github) GetReleases() (*[]Release, error) {
	endpointUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", gh.organization, gh.repository)

	resp, err := http.Get(endpointUrl)

	if err != nil {
		log.Fatalf("Cannot execute request err=%s", err.Error())

		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Invalid status code")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	releases := make([]Release, 0)

	err = json.Unmarshal(body, &releases)

	if err != nil {
		return nil, err
	}

	return &releases, nil
}

// TODO
func (gh *github) GetLatestCliRelease() (*Release, error) {
	return nil, nil
}

// TODO
func (gh *github) GetLatestUpdaterRelease() (*Release, error) {
	return nil, nil
}