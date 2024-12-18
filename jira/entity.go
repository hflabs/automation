package jira

type jira struct {
	BaseUrl  string
	Username string
	Password string
}

type WebhookIssue struct {
	Timestamp      Timestamp    `json:"timestamp,omitempty"`
	WebhookEvent   string       `json:"webhookEvent,omitempty"`
	IssueEventType string       `json:"issue_event_type_name,omitempty"`
	UserEvent      JiraUser     `json:"user,omitempty"`
	Issue          IssueJira    `json:"issue,omitempty"`
	Comment        IssueComment `json:"comment,omitempty"`
	Changelog      struct {
		Id    string          `json:"id,omitempty"`
		Items []ChangelogItem `json:"items,omitempty"`
	} `json:"changelog,omitempty"`
}

type WebhookComment struct {
	Timestamp    Timestamp    `json:"timestamp,omitempty"`
	WebhookEvent string       `json:"webhookEvent,omitempty"`
	Comment      IssueComment `json:"comment,omitempty"`
}

type IssueJira struct {
	Id     string      `json:"id,omitempty"`
	Key    string      `json:"key,omitempty"`
	Fields FieldsIssue `json:"fields,omitempty"`
}

type ChangelogItem struct {
	Field      string `json:"field,omitempty,omitempty"`
	FieldType  string `json:"fieldType,omitempty"`
	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
	FromString string `json:"fromString,omitempty"`
	ToString   string `json:"toString,omitempty"`
}

type FieldsIssue struct {
	Summary       string       `json:"summary,omitempty"`
	BusinessValue float64      `json:"customfield_10084,omitempty"`
	StoryPoints   float64      `json:"customfield_10083,omitempty"`
	WeightedJob   float64      `json:"customfield_12580,omitempty"`
	Status        IssueField   `json:"status,omitempty"`
	IssueType     IssueIdField `json:"issuetype,omitempty"`
	Priority      IssueIdField `json:"priority,omitempty"`
	Resolution    IssueIdField `json:"resolution,omitempty"`
	Assignee      JiraUser     `json:"assignee,omitempty"`
	Creator       JiraUser     `json:"creator,omitempty"`
	Reporter      JiraUser     `json:"reporter,omitempty"`
	Project       IssueField   `json:"project,omitempty"`
}

type IssueIdField struct {
	ID string `json:"id,omitempty"`
}

type IssueNameField struct {
	Name string `json:"name,omitempty"`
}

type IssueField struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

type JiraUser struct {
	Name  string `json:"name,omitempty"`
	Key   string `json:"key,omitempty"`
	Email string `json:"emailAddress,omitempty"`
}

type IssueWatchersResponse struct {
	Watchers []JiraUser `json:"watchers,omitempty"`
}

type IssueCommentsResponse struct {
	Total    int            `json:"total,omitempty"`
	Comments []IssueComment `json:"comments,omitempty"`
}

type IssueComment struct {
	Author       JiraUser `json:"author,omitempty"`
	Body         string   `json:"body,omitempty"`
	UpdateAuthor JiraUser `json:"update_author,omitempty"`
	Created      JiraTime `json:"created,omitempty"`
	Updated      JiraTime `json:"updated,omitempty"`
}

type SearchResponse struct {
	StartAt    int         `json:"startAt,omitempty"`
	MaxResults int         `json:"maxResults,omitempty"`
	Total      int         `json:"total,omitempty"`
	Issues     []IssueJira `json:"issues,omitempty"`
}

type UpdateIssueRequest struct {
	Fields map[string]interface{} `json:"fields"`
}

type TransitionIssueRequest struct {
	Transition IssueIdField `json:"transition,omitempty"`
}
