package jira

import (
	"cmp"
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/carlmjohnson/requests"
)

type jira struct {
	BaseUrl  string
	Username string
	Password string

	fieldIdMap map[string]string // Кэш для маппинга "Имя поля" -> "ID поля"
	mu         sync.RWMutex
}

// NewJira теперь возвращает ошибку, так как мы инициализируем карту полей при старте
func NewJira(ctx context.Context, baseUrl, user, password string) (ApiJira, error) {
	j := &jira{
		BaseUrl:    strings.TrimRight(baseUrl, "/"),
		Username:   user,
		Password:   password,
		fieldIdMap: make(map[string]string),
	}
	// Инициализируем карту полей при создании клиента
	err := j.RefreshFields(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to RefreshFieldMap: %w", err)
	}
	return j, nil
}

func (j *jira) RefreshFields(ctx context.Context) error {
	var fields []IssueField
	err := requests.
		URL(fmt.Sprintf("%s/field", j.BaseUrl)).
		BasicAuth(j.Username, j.Password).
		ToJSON(&fields).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return err
	}
	fieldMap := make(map[string]string, len(fields))
	for _, f := range fields {
		fieldMap[f.Name] = f.ID // Сохраняем mapping вида: "Story Points" -> "customfield_10083"
	}
	j.mu.Lock()
	defer j.mu.Unlock()

	j.fieldIdMap = fieldMap
	return nil
}

// GetFieldID возвращает ID поля по его имени (из кэша)
func (j *jira) GetFieldID(name string) (string, bool) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	id, ok := j.fieldIdMap[name]
	return id, ok
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

func (j *jira) UpdateIssue(ctx context.Context, issueKey string, req UpdateIssueRequest) error {
	return requests.
		URL(fmt.Sprintf("%s/issue/%s", j.BaseUrl, issueKey)).
		Put().
		BasicAuth(j.Username, j.Password).
		BodyJSON(req).
		AddValidator(validateStatus).
		Fetch(ctx)
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
// 2. Ищет переход в статус targetStatusName (например "Done" или "In Progress").
// 3. Если найден - выполняет. Если нет - возвращает ошибку со списком доступных.
func (j *jira) TransitionToStatus(ctx context.Context, issueKey, targetStatusName string) error {
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
	var availableStatuses []string
	for _, t := range meta.Transitions {
		availableStatuses = append(availableStatuses, t.To.Name)
		// Сравниваем case-insensitive, так надежнее
		if strings.EqualFold(t.To.Name, targetStatusName) {
			targetTransitionID = t.ID
			break
		}
	}
	if targetTransitionID == "" {
		return fmt.Errorf("cannot transition issue %s to status '%s'. Available statuses: %v",
			issueKey, targetStatusName, strings.Join(availableStatuses, ", "))
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

func (j *jira) QueryTasks(ctx context.Context, query string, pageSize int) ([]IssueJira, error) {
	var tasks SearchResponse
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}
	err := requests.
		URL(fmt.Sprintf("%s/search?jql=%s&maxResults=%d", j.BaseUrl, url.QueryEscape(query), cmp.Or(pageSize, 50))).
		BasicAuth(j.Username, j.Password).
		ToJSON(&tasks).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return tasks.Issues, nil
}
