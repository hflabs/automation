package jira

import (
	"cmp"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/carlmjohnson/requests"
)

type jira struct {
	BaseUrl  string
	Username string
	Password string
}

func NewJira(baseUrl, user, password string) ApiJira {
	return &jira{BaseUrl: strings.TrimRight(baseUrl, "/"), Username: user, Password: password}
}

// GetFields — возвращает полный список полей в Jira
func (j *jira) GetFields(ctx context.Context) ([]IssueField, error) {
	var fields []IssueField
	err := requests.
		URL(fmt.Sprintf("%s/field", j.BaseUrl)).
		BasicAuth(j.Username, j.Password).
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
		BasicAuth(j.Username, j.Password).
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
		BasicAuth(j.Username, j.Password).
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
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *jira) GetIssueById(ctx context.Context, issueId string) (IssueJira, error) {
	var resp IssueJira
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueId)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return IssueJira{}, err
	}
	return resp, nil
}

func (j *jira) GetUserByKey(ctx context.Context, userKey string) (JiraUser, error) {
	var resp JiraUser
	err := requests.
		URL(fmt.Sprintf("%s/user?key=%s", j.BaseUrl, url.QueryEscape(userKey))).
		BasicAuth(j.Username, j.Password).
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
		BasicAuth(j.Username, j.Password).
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
		BasicAuth(j.Username, j.Password).
		BodyJSON(JiraUser{Name: assigneeName}).
		AddValidator(validateStatus).
		Fetch(ctx)
}

// TransitionIssue - низкоуровневый метод, принимает ID перехода
func (j *jira) TransitionIssue(ctx context.Context, issueKey, transitionID string) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s/transitions", j.BaseUrl, issueKey)).
		Post().
		BasicAuth(j.Username, j.Password).
		BodyJSON(TransitionIssueRequest{Transition: IssueField{ID: transitionID}}).
		AddValidator(validateStatus).
		Fetch(ctx)
}

// TransitionToStatus - высокоуровневый метод.
// 1. Получает доступные переходы.
// 2. Ищет переход в статус targetStatusId
// 3. Если найден - выполняет. Если нет - возвращает ошибку со списком доступных.
func (j *jira) TransitionToStatus(ctx context.Context, issueKey, targetStatusId string) error {
	// 1. Получаем возможные переходы для этой задачи
	var meta TransitionsResponse
	err := requests.
		URL(fmt.Sprintf("%s/issue/%s/transitions", j.BaseUrl, issueKey)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&meta).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transitions: %w", err)
	}

	// 2. Ищем нужный переход
	var targetTransitionID string
	// Мапа, чтобы в ответе был не только ID, но и название статуса
	availableStatuses := make(map[string]string, len(meta.Transitions))
	for _, t := range meta.Transitions {
		availableStatuses[t.To.ID] = t.To.Name
		if t.To.ID == targetStatusId {
			targetTransitionID = t.ID
			break
		}
	}
	if targetTransitionID == "" {
		return fmt.Errorf("cannot transition issue %s to status '%s'. Available statuses: %v",
			issueKey, targetStatusId, formatAvailableStatuses(availableStatuses))
	}
	// 3. Выполняем переход
	return j.TransitionIssue(ctx, issueKey, targetTransitionID)
}

func (j *jira) CommentIssue(ctx context.Context, issueKey, comment string) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s/comment", j.BaseUrl, issueKey)).
		Post().
		BasicAuth(j.Username, j.Password).
		BodyJSON(IssueComment{Body: comment}).
		AddValidator(validateStatus).
		Fetch(ctx)
}

// SearchTasks — постраничный поиск задач по JQL запросу, pageSize по умолчанию 50
func (j *jira) SearchTasks(ctx context.Context, query string, pageSize, offset int) (SearchResponse, error) {
	var resp SearchResponse
	if query == "" {
		return SearchResponse{}, fmt.Errorf("query is empty")
	}
	err := requests.
		URL(fmt.Sprintf("%s/search", j.BaseUrl)).
		Param("jql", query).
		Param("maxResults", strconv.Itoa(cmp.Or(pageSize, 50))).
		Param("startAt", strconv.Itoa(offset)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return SearchResponse{}, err
	}
	return resp, nil
}

// SearchAllTasks — поиск всех задач по JQL запросу
func (j *jira) SearchAllTasks(ctx context.Context, query string) ([]IssueJira, error) {
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}
	const pageSize = 1000
	offset := 0
	var all []IssueJira
	for {
		resp, err := j.SearchTasks(ctx, query, pageSize, offset)
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
		BasicAuth(j.Username, j.Password).
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
		BasicAuth(j.Username, j.Password).
		BodyJSON(UpsertIssueRequest{Fields: req}).
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
		BasicAuth(j.Username, j.Password).
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
		BasicAuth(j.Username, j.Password).
		BodyJSON(UpsertIssueRequest{Fields: req}).
		ToJSON(&created).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return CreatedIssueResponse{}, err
	}
	return created, nil
}
