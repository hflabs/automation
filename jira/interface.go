package jira

import "context"

type ApiJira interface {
	QueryTasks(ctx context.Context, query string, pageSize int) ([]IssueJira, error)
	GetIssueById(ctx context.Context, issueId string) (IssueJira, error)

	GetIssueComments(ctx context.Context, issueKey string) ([]IssueComment, error)
	GetIssueWatchers(ctx context.Context, issueKey string) ([]JiraUser, error)
	GetIssueChangelog(ctx context.Context, issueId string) ([]ChangeLog, error)
	GetProjectVersions(ctx context.Context, projectKey string) ([]ProjectVersion, error)
	GetUserByKey(ctx context.Context, userKey string) (JiraUser, error)

	UpdateIssue(ctx context.Context, issueKey string, req UpdateIssueRequest) error
	CommentIssue(ctx context.Context, issueKey, comment string) error
	UpdateIssueAssignee(ctx context.Context, issueKey string, assigneeName string) error

	// TransitionIssue - низкоуровневый метод (по ID перехода)
	TransitionIssue(ctx context.Context, issueKey, transitionID string) error
	// TransitionToStatus - высокоуровневый метод (по Названию статуса назначения)
	TransitionToStatus(ctx context.Context, issueKey, targetStatusName string) error

	// Утилиты
	GetFieldID(name string) (string, bool) // Получить ID поля по имени (например, "Story Points" -> "customfield_10083")
	RefreshFields(ctx context.Context) error
}
