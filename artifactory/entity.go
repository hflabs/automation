package artifactory

import (
	"encoding/json"
	"strings"
	"time"
)

type artifactory struct {
	baseUrl    string
	repo       string
	aqlBaseUrl string
	username   string
	password   string
}

type ListOptions struct {
	SortDirection string // "asc" or "desc"
	SortValue     string // e.g., "name", "size", "created"
	Limit         int    // maximum number of items to return
	Name          string // match item name to string with '*' or exact search without '*'
	Type          string // item type to search
}

type findQuery struct {
	Repo string `json:"repo"`
	Path string `json:"path"`
	Type string `json:"type"`
	Name *name  `json:"name,omitempty"`
}

type name struct {
	Value string
}

// MarshalJSON implements custom JSON marshaling for the Name field
func (n name) MarshalJSON() ([]byte, error) {
	if strings.Contains(n.Value, "*") {
		return json.Marshal(map[string]string{
			"$match": n.Value,
		})
	}
	return json.Marshal(n.Value)
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
