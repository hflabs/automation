package jira

import "context"

type ApiJira interface {
	SearchTasks(ctx context.Context, query string, pageSize, offset int, fields ...string) (SearchResponse, error)
	SearchAllTasks(ctx context.Context, query string, fields ...string) ([]IssueJira, error)
	GetIssueById(ctx context.Context, issueId string, fields ...string) (IssueJira, error)

	GetIssueComments(ctx context.Context, issueKey string) ([]IssueComment, error)
	GetIssueWatchers(ctx context.Context, issueKey string) ([]JiraUser, error)
	GetIssueChangelog(ctx context.Context, issueId string) ([]ChangeLog, error)
	GetProjectVersions(ctx context.Context, projectKey string) ([]ProjectVersion, error)
	GetUserByKey(ctx context.Context, userKey string) (JiraUser, error)
	GetFields(ctx context.Context) ([]IssueField, error)
	GetIssueTypeMeta(ctx context.Context, projectKey, issueTypeId string) (*IssueTypeMeta, error)

	CreateIssueFromMap(ctx context.Context, req map[string]any) (CreatedIssueResponse, error)
	CreateIssue(ctx context.Context, req FieldsIssue) (CreatedIssueResponse, error)

	UpdateIssueFromMap(ctx context.Context, issueKey string, req map[string]any) error
	UpdateIssue(ctx context.Context, issueKey string, req FieldsIssue) error
	UpdateIssueAssignee(ctx context.Context, issueKey string, assigneeName string) error

	AddLabel(ctx context.Context, issueKey string, label string) error

	CommentIssue(ctx context.Context, issueKey, comment string) error

	// TransitionIssue - низкоуровневый метод (по ID перехода)
	TransitionIssue(ctx context.Context, issueKey, transitionID string) error
	// TransitionToStatus - высокоуровневый метод (по ID статуса назначения)
	TransitionToStatus(ctx context.Context, issueKey, targetStatusId string) error
}
