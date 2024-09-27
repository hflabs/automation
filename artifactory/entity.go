package artifactory

import "time"

type artifactory struct {
	baseUrl    string
	repo       string
	aqlBaseUrl string
	username   string
	password   string
}

type Item struct {
	Name    string    `json:"name"`
	Repo    string    `json:"repo"`
	Path    string    `json:"path"`
	Type    string    `json:"type"`
	Size    int       `json:"size"`
	Created time.Time `json:"created"`
}

type Result struct {
	Items []Item `json:"results"`
	Count struct {
		Total int `json:"total"`
	} `json:"range"`
}
