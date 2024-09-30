package jira

import "time"

type jira struct {
	BaseUrl     string
	Username    string
	Password    string
	Suggestions Fields
}

type WebhookIssue struct {
	IssueEventType string    `json:"issue_event_type_name,omitempty"`
	Issue          IssueJira `json:"issue,omitempty"`
}

type IssueJira struct {
	Key       string      `json:"key,omitempty"`
	Fields    FieldsIssue `json:"fields,omitempty"`
	Changelog struct {
		Items []IssueItemChangelog `json:"items,omitempty"`
	} `json:"changelog,omitempty"`
}

type IssueItemChangelog struct {
	Field string `json:"field,omitempty,omitempty"`
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
}

type FieldsIssue struct {
	BusinessValue int          `json:"customfield_10084,omitempty"`
	StoryPoints   int          `json:"customfield_10083,omitempty"`
	WeightedJob   float64      `json:"customfield_12580,omitempty"`
	Status        IssueField   `json:"status,omitempty"`
	Issuetype     IssueIdField `json:"issuetype,omitempty"`
	Priority      IssueIdField `json:"priority,omitempty"`
	Resolution    IssueIdField `json:"resolution,omitempty"`
	Assignee      JiraUser     `json:"assignee,omitempty"`
	Creator       JiraUser     `json:"creator,omitempty"`
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
	Author  JiraUser  `json:"author,omitempty"`
	Body    string    `json:"body,omitempty,omitempty"`
	Created time.Time `json:"created,omitempty,omitempty"`
}

type SearchResponse struct {
	StartAt    int         `json:"startAt,omitempty"`
	MaxResults int         `json:"maxResults,omitempty"`
	Total      int         `json:"total,omitempty"`
	Issues     []IssueJira `json:"issues,omitempty"`
}
