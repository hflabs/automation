package jira

type ApiJira interface {
	QueryTasks(query string, pageSize int) ([]IssueJira, error)

	GetIssueComments(issueKey string) ([]IssueComment, error)
	GetIssueWatchers(issueKey string) ([]JiraUser, error)

	GetIssueById(issueId string) (IssueJira, error)

	UpdateIssue(issueKey string, req UpdateIssueRequest) error
	CommentIssue(issueKey, comment string) error
	TransitionIssue(issueKey, transition string) error

	GetProjectVersions(projectKey string) ([]ProjectVersion, error)
}
