package jira

import (
	"cmp"
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/carlmjohnson/requests"
)

type jira struct {
	BaseUrl string
	Token   string
}

func NewJira(baseUrl, token string) ApiJira {
	return &jira{BaseUrl: strings.TrimRight(baseUrl, "/"), Token: token}
}

// GetFields — возвращает полный список полей в Jira
func (j *jira) GetFields(ctx context.Context) ([]IssueField, error) {
	var fields []IssueField
	err := requests.
		URL(fmt.Sprintf("%s/field", j.BaseUrl)).
		Bearer(j.Token).
		ToJSON(&fields).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func (j *jira) GetIssueComments(ctx context.Context, issueKey string) ([]IssueComment, error) {
	var resp IssueCommentsResponse
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s/comment", j.BaseUrl, issueKey)).
		Bearer(j.Token).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Comments, nil
}

func (j *jira) GetIssueWatchers(ctx context.Context, issueKey string) ([]JiraUser, error) {
	var resp IssueWatchersResponse
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s/watchers", j.BaseUrl, issueKey)).
		Bearer(j.Token).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Watchers, nil
}

func (j *jira) GetProjectVersions(ctx context.Context, projectKey string) ([]ProjectVersion, error) {
	var resp []ProjectVersion
	err := requests.
		URL(fmt.Sprintf("%s/project/%s/versions", j.BaseUrl, projectKey)).
		Bearer(j.Token).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *jira) GetIssueById(ctx context.Context, issueId string, fields ...string) (IssueJira, error) {
	if issueId == "" {
		return IssueJira{}, fmt.Errorf("issueId is empty")
	}
	var resp IssueJira
	req := requests.
		URL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueId)).
		Bearer(j.Token).
		ToJSON(&resp).
		AddValidator(validateStatus)
	// Если поля указаны, добавляем их в URL через запятую
	if len(fields) > 0 {
		req.Param("fields", strings.Join(fields, ","))
	}
	err := req.Fetch(ctx)
	if err != nil {
		return IssueJira{}, err
	}
	return resp, nil
}

func (j *jira) GetUserByKey(ctx context.Context, userKey string) (JiraUser, error) {
	var resp JiraUser
	err := requests.
		URL(fmt.Sprintf("%s/user?key=%s", j.BaseUrl, url.QueryEscape(userKey))).
		Bearer(j.Token).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return JiraUser{}, err
	}
	return resp, nil
}

func (j *jira) GetIssueChangelog(ctx context.Context, issueId string) ([]ChangeLog, error) {
	var resp IssueJira
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s?expand=changelog", j.BaseUrl, issueId)).
		Bearer(j.Token).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return resp.Changelog.Histories, nil
}

func (j *jira) UpdateIssueAssignee(ctx context.Context, issueKey, assigneeName string) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s/assignee", j.BaseUrl, issueKey)).
		Put().
		Bearer(j.Token).
		BodyJSON(JiraUser{Name: assigneeName}).
		AddValidator(validateStatus).
		Fetch(ctx)
}

// TransitionIssue - низкоуровневый метод, принимает ID перехода
func (j *jira) TransitionIssue(ctx context.Context, issueKey, transitionID string) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s/transitions", j.BaseUrl, issueKey)).
		Post().
		Bearer(j.Token).
		BodyJSON(TransitionIssueRequest{Transition: IssueField{ID: transitionID}}).
		AddValidator(validateStatus).
		Fetch(ctx)
}

// TransitionToStatus is a high-level transition method by target status ID.
// It can move through known intermediate statuses when Jira does not expose a
// direct transition from the issue's current status.
func (j *jira) TransitionToStatus(ctx context.Context, issueKey, targetStatusId string) error {
	if strings.TrimSpace(issueKey) == "" {
		return fmt.Errorf("issueKey is empty")
	}
	if strings.TrimSpace(targetStatusId) == "" {
		return fmt.Errorf("targetStatusId is empty")
	}

	issue, err := j.GetIssueById(ctx, issueKey, Issue.Fields.Status)
	if err != nil {
		return fmt.Errorf("failed to get issue status: %w", err)
	}
	currentStatusId := issue.Fields.Status.ID
	if currentStatusId == targetStatusId {
		return nil
	}

	const maxTransitionToStatusSteps = 20
	for step := 0; step < maxTransitionToStatusSteps; step++ {
		transitions, err := j.getIssueTransitions(ctx, issueKey)
		if err != nil {
			return fmt.Errorf("failed to get transitions: %w", err)
		}
		if len(transitions) == 0 {
			return fmt.Errorf("cannot transition issue %s from status '%s' to status '%s': no available transitions",
				issueKey, currentStatusId, targetStatusId)
		}

		if transition, ok := findTransitionToStatus(transitions, targetStatusId); ok {
			return j.TransitionIssue(ctx, issueKey, transition.ID)
		}

		route := findStatusRoute(currentStatusId, targetStatusId, transitions)
		if len(route) < 2 {
			return fmt.Errorf("cannot transition issue %s from status '%s' to status '%s'. Available statuses: %v",
				issueKey, currentStatusId, targetStatusId, formatAvailableStatuses(transitions))
		}

		nextStatusId := route[1]
		transition, ok := findTransitionToStatus(transitions, nextStatusId)
		if !ok {
			return fmt.Errorf("cannot transition issue %s from status '%s' to next route status '%s'. Available statuses: %v",
				issueKey, currentStatusId, nextStatusId, formatAvailableStatuses(transitions))
		}
		if err := j.TransitionIssue(ctx, issueKey, transition.ID); err != nil {
			return fmt.Errorf("failed to transition issue %s from status '%s' to status '%s': %w",
				issueKey, currentStatusId, nextStatusId, err)
		}
		currentStatusId = nextStatusId
	}

	return fmt.Errorf("cannot transition issue %s to status '%s': exceeded %d transition steps",
		issueKey, targetStatusId, maxTransitionToStatusSteps)
}

func (j *jira) getIssueTransitions(ctx context.Context, issueKey string) ([]Transition, error) {
	var meta TransitionsResponse
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s/transitions", j.BaseUrl, issueKey)).
		Bearer(j.Token).
		ToJSON(&meta).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return meta.Transitions, nil
}

func (j *jira) CommentIssue(ctx context.Context, issueKey, comment string) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s/comment", j.BaseUrl, issueKey)).
		Post().
		Bearer(j.Token).
		BodyJSON(IssueComment{Body: comment}).
		AddValidator(validateStatus).
		Fetch(ctx)
}

// SearchTasks — постраничный поиск задач по JQL запросу, pageSize по умолчанию 50
func (j *jira) SearchTasks(ctx context.Context, query string, pageSize, offset int, fields ...string) (SearchResponse, error) {
	req := SearchRequest{
		Jql:        query,
		StartAt:    offset,
		MaxResults: cmp.Or(pageSize, 50),
		Fields:     fields,
	}
	var resp SearchResponse
	if query == "" {
		return SearchResponse{}, fmt.Errorf("query is empty")
	}
	err := requests.
		URL(fmt.Sprintf("%s/search", j.BaseUrl)).
		BodyJSON(&req).
		Bearer(j.Token).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return SearchResponse{}, err
	}
	return resp, nil
}

// SearchAllTasks — поиск всех задач по JQL запросу
func (j *jira) SearchAllTasks(ctx context.Context, query string, fields ...string) ([]IssueJira, error) {
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}
	const pageSize = 1000
	offset := 0
	var all []IssueJira
	for {
		resp, err := j.SearchTasks(ctx, query, pageSize, offset, fields...)
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Issues...)
		offset += len(resp.Issues)
		if len(resp.Issues) == 0 || offset >= resp.Total {
			break
		}
	}
	return all, nil
}

func (j *jira) UpdateIssueFromMap(ctx context.Context, issueKey string, req map[string]any) error {
	if strings.TrimSpace(issueKey) == "" {
		return fmt.Errorf("issueKey is empty")
	}
	if req == nil || len(req) == 0 {
		return fmt.Errorf("request fields are empty")
	}
	return requests.
		URL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueKey)).
		Put().
		Bearer(j.Token).
		BodyJSON(UpsertIssueRequestFromMap{Fields: req}).
		AddValidator(validateStatus).
		Fetch(ctx)
}

// UpdateIssue — Обновить задачу по её ключу
func (j *jira) UpdateIssue(ctx context.Context, issueKey string, req FieldsIssue) error {
	if strings.TrimSpace(issueKey) == "" {
		return fmt.Errorf("issueKey is empty")
	}
	// Для обновления Jira допускает частичный набор полей, поэтому дополнительных проверок не делаем
	return requests.
		URL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueKey)).
		Put().
		Bearer(j.Token).
		BodyJSON(UpsertIssueRequest{Fields: req}).
		AddValidator(validateStatus).
		Fetch(ctx)
}

func (j *jira) AddLabel(ctx context.Context, issueKey string, label string) error {
	if strings.TrimSpace(issueKey) == "" {
		return fmt.Errorf("issueKey is empty")
	}
	req := UpdateIssueRequest{Update: UpdateIssue{Labels: []UpdateField{{Add: label}}}}
	return requests.
		URL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueKey)).
		Put().
		Bearer(j.Token).
		BodyJSON(req).
		AddValidator(validateStatus).
		Fetch(ctx)
}

func (j *jira) CreateIssueFromMap(ctx context.Context, req map[string]any) (CreatedIssueResponse, error) {
	if req == nil || len(req) == 0 {
		return CreatedIssueResponse{}, fmt.Errorf("request fields are empty")
	}
	var created CreatedIssueResponse
	err := requests.
		URL(fmt.Sprintf("%s/issue", j.BaseUrl)).
		Post().
		Bearer(j.Token).
		BodyJSON(UpsertIssueRequestFromMap{Fields: req}).
		ToJSON(&created).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return CreatedIssueResponse{}, err
	}
	return created, nil
}

// CreateIssue — создать задачу, поля Summary, Project и IssueType обязательны для любых проектов
func (j *jira) CreateIssue(ctx context.Context, req FieldsIssue) (CreatedIssueResponse, error) {
	// Базовая валидация обязательных полей для создания задачи в Jira
	if strings.TrimSpace(req.Summary) == "" {
		return CreatedIssueResponse{}, fmt.Errorf("summary is empty")
	}
	// Тут в целом по проекту и типу задачи можно узнать какие поля ещё нужны и их проверять, но как будто нет смысла,
	// жира в любом случае их вернет в ошибке
	if req.Project.ID == "" && strings.TrimSpace(req.Project.Key) == "" {
		return CreatedIssueResponse{}, fmt.Errorf("project is empty (need id or key)")
	}
	if req.IssueType.ID == "" && strings.TrimSpace(req.IssueType.Name) == "" {
		return CreatedIssueResponse{}, fmt.Errorf("issuetype is empty (need id or name)")
	}

	var created CreatedIssueResponse
	err := requests.
		URL(fmt.Sprintf("%s/issue", j.BaseUrl)).
		Post().
		Bearer(j.Token).
		BodyJSON(UpsertIssueRequest{Fields: req}).
		ToJSON(&created).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return CreatedIssueResponse{}, err
	}
	return created, nil
}

func (j *jira) GetIssueTypeMeta(ctx context.Context, projectKey, issueTypeId string) (*IssueTypeMeta, error) {
	resp := &IssueTypeMeta{}
	err := requests.
		URL(fmt.Sprintf("%s/issue/createmeta/%s/issuetypes/%s", j.BaseUrl, projectKey, issueTypeId)).
		Bearer(j.Token).
		ToJSON(resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *jira) GetJiraProjects(ctx context.Context) ([]JiraProject, error) {
	var projects []JiraProject
	err := requests.
		URL(fmt.Sprintf("%s/project", j.BaseUrl)).
		Bearer(j.Token).
		ToJSON(&projects).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (j *jira) GetJiraProjectComponents(ctx context.Context, projectKey string) ([]JiraComponent, error) {
	var components []JiraComponent
	err := requests.
		URL(fmt.Sprintf("%s/project/%s/components", j.BaseUrl, projectKey)).
		Bearer(j.Token).
		ToJSON(&components).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return components, nil
}
