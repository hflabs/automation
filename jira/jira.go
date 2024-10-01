package jira

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/url"
)

func NewJira(baseUrl, user, password string) ApiJira {
	return &jira{baseUrl, user, password}
}

func (j *jira) GetIssueComments(issueKey string) ([]IssueComment, error) {
	var resp IssueCommentsResponse
	err := requests.New().
		BaseURL(fmt.Sprintf("%s/issue/%s/comment", j.BaseUrl, issueKey)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Comments, nil
}

func (j *jira) GetIssueWatchers(issueKey string) ([]JiraUser, error) {
	var resp IssueWatchersResponse
	err := requests.New().
		BaseURL(fmt.Sprintf("%s/issue/%s/watchers", j.BaseUrl, issueKey)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Watchers, nil
}

func (j *jira) UpdateIssue(issueKey string, req UpdateIssueRequest) error {
	return requests.New().
		BaseURL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueKey)).
		Put().
		BasicAuth(j.Username, j.Password).
		BodyJSON(req).
		Fetch(context.Background())
}

func (j *jira) TransitionIssue(issueKey, transition string) error {
	return requests.New().
		BaseURL(fmt.Sprintf("%s/issue/%s/transitions", j.BaseUrl, issueKey)).
		Post().
		BasicAuth(j.Username, j.Password).
		BodyJSON(TransitionIssueRequest{IssueIdField{ID: transition}}).
		Fetch(context.Background())
}

func (j *jira) CommentIssue(issueKey, comment string) error {
	return requests.New().
		BaseURL(fmt.Sprintf("%s/issue/%s/comment", j.BaseUrl, issueKey)).
		Post().
		BasicAuth(j.Username, j.Password).
		BodyJSON(IssueComment{Body: comment}).
		Fetch(context.Background())
}

func (j *jira) QueryTasks(query string) ([]IssueJira, error) {
	var tasks SearchResponse
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}
	err := requests.New().
		BaseURL(fmt.Sprintf("%s/search?jql=%s", j.BaseUrl, url.QueryEscape(query))).
		BasicAuth(j.Username, j.Password).
		ToJSON(&tasks).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return tasks.Issues, nil
}
