package jira

type ApiJira interface {
	QueryTasks(query string) ([]IssueJira, error)

	GetIssueComments(issueKey string) ([]IssueComment, error)
	GetIssueWatchers(issueKey string) ([]JiraUser, error)

	UpdateIssue(issueKey string, jsonBody interface{}) error
	CommentIssue(issueKey, comment string) error
	TransitionIssue(issueKey, transition string) error
}
