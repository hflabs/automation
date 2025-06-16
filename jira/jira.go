package jira

import (
	"cmp"
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
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s/comment", j.BaseUrl, issueKey)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Comments, nil
}

func (j *jira) GetIssueWatchers(issueKey string) ([]JiraUser, error) {
	var resp IssueWatchersResponse
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s/watchers", j.BaseUrl, issueKey)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Watchers, nil
}

func (j *jira) GetProjectVersions(projectKey string) ([]ProjectVersion, error) {
	var resp []ProjectVersion
	err := requests.
		URL(fmt.Sprintf("%s/project/%s/versions", j.BaseUrl, projectKey)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *jira) GetIssueById(issueId string) (IssueJira, error) {
	var resp IssueJira
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueId)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return IssueJira{}, err
	}
	return resp, nil
}

func (j *jira) GetUserByKey(userKey string) (JiraUser, error) {
	var resp JiraUser
	err := requests.
		URL(fmt.Sprintf("%s/user?key=%s", j.BaseUrl, url.QueryEscape(userKey))).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return JiraUser{}, err
	}
	return resp, nil
}

func (j *jira) GetIssueChangelog(issueId string) ([]ChangeLog, error) {
	var resp IssueJira
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s?expand=changelog", j.BaseUrl, issueId)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Changelog.Histories, nil
}

func (j *jira) UpdateIssue(issueKey string, req UpdateIssueRequest) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueKey)).
		Put().
		BasicAuth(j.Username, j.Password).
		BodyJSON(req).
		AddValidator(validateStatus).
		Fetch(context.Background())
}

func (j *jira) UpdateIssueAssignee(issueKey, assigneeName string) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s/assignee", j.BaseUrl, issueKey)).
		Put().
		BasicAuth(j.Username, j.Password).
		BodyJSON(JiraUser{Name: assigneeName}).
		AddValidator(validateStatus).
		Fetch(context.Background())
}

func (j *jira) TransitionIssue(issueKey, transition string) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s/transitions", j.BaseUrl, issueKey)).
		Post().
		BasicAuth(j.Username, j.Password).
		BodyJSON(TransitionIssueRequest{IssueField{ID: transition}}).
		AddValidator(validateStatus).
		Fetch(context.Background())
}

func (j *jira) CommentIssue(issueKey, comment string) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s/comment", j.BaseUrl, issueKey)).
		Post().
		BasicAuth(j.Username, j.Password).
		BodyJSON(IssueComment{Body: comment}).
		AddValidator(validateStatus).
		Fetch(context.Background())
}

func (j *jira) QueryTasks(query string, pageSize int) ([]IssueJira, error) {
	var tasks SearchResponse
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}
	err := requests.
		URL(fmt.Sprintf("%s/search?jql=%s&maxResults=%d", j.BaseUrl, url.QueryEscape(query), cmp.Or(pageSize, 50))).
		BasicAuth(j.Username, j.Password).
		ToJSON(&tasks).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return tasks.Issues, nil
}
