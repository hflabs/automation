package jira

type ApiJira interface {
	QueryTasks(query string, pageSize int) ([]IssueJira, error)

	GetIssueComments(issueKey string) ([]IssueComment, error)
	GetIssueWatchers(issueKey string) ([]JiraUser, error)

	GetIssueById(issueId string) (IssueJira, error)
	GetIssueChangelog(issueId string) ([]ChangeLog, error)

	UpdateIssue(issueKey string, req UpdateIssueRequest) error
	CommentIssue(issueKey, comment string) error
	TransitionIssue(issueKey, transition string) error
	UpdateIssueAssignee(issueKey string, assigneeName string) error

	GetProjectVersions(projectKey string) ([]ProjectVersion, error)

	GetUserByKey(userKey string) (JiraUser, error)
}
