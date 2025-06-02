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
	Changelog      ChangeLog    `json:"changelog,omitempty"`
	Version        Version      `json:"version,omitempty"`
}

type ChangeLog struct {
	Id    string          `json:"id,omitempty"`
	Items []ChangelogItem `json:"items,omitempty"`
}

func (c ChangeLog) FindItemByField(field string) ChangelogItem {
	for _, item := range c.Items {
		if field == item.Field {
			return item
		}
	}
	return ChangelogItem{}
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
	Summary            string       `json:"summary,omitempty"`
	Description        string       `json:"description,omitempty"`
	BusinessValue      float64      `json:"customfield_10084,omitempty"`
	StoryPoints        float64      `json:"customfield_10083,omitempty"`
	WeightedJob        float64      `json:"customfield_12580,omitempty"`
	ReleaseNotes       string       `json:"customfield_13082,omitempty"`
	ReleaseInstruction string       `json:"customfield_13081,omitempty"`
	DueDate            string       `json:"duedate,omitempty"`
	Status             IssueField   `json:"status,omitempty"`
	IssueType          IssueField   `json:"issuetype,omitempty"`
	Priority           IssueField   `json:"priority,omitempty"`
	Resolution         IssueField   `json:"resolution,omitempty"`
	Assignee           JiraUser     `json:"assignee,omitempty"`
	Creator            JiraUser     `json:"creator,omitempty"`
	Reporter           JiraUser     `json:"reporter,omitempty"`
	Participants       []JiraUser   `json:"customfield_10380,omitempty"`
	Project            IssueField   `json:"project,omitempty"`
	Components         []IssueField `json:"components,omitempty"`
	LearnTime          string       `json:"customfield_14481,omitempty"`
	LearnForWho        string       `json:"customfield_13881,omitempty"`
	LearnWhatLike      string       `json:"customfield_14483,omitempty"`
	LearnWhatUseful    string       `json:"customfield_13784,omitempty"`
	LearnWhatBad       string       `json:"customfield_14484,omitempty"`
	LearnWhatLearned   string       `json:"customfield_13880,omitempty"`
	LearnWillRecommend string       `json:"customfield_13882,omitempty"`
	LearnPeoples       []JiraUser   `json:"customfield_14480,omitempty"`
	LearnField         IssueField   `json:"customfield_14380,omitempty"`
	LearnLink          string       `json:"customfield_13782,omitempty"`
}

type IssueField struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type JiraUser struct {
	Name        string `json:"name,omitempty"`
	Key         string `json:"key,omitempty"`
	Email       string `json:"emailAddress,omitempty"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active,omitempty"`
}

type ProjectVersion struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Archived        bool   `json:"archived"`
	Released        bool   `json:"released"`
	ReleaseDate     string `json:"releaseDate"`
	UserReleaseDate string `json:"userReleaseDate"`
	ProjectID       int    `json:"projectId"`
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

type Version struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	Archived        bool   `json:"archived,omitempty"`
	Released        bool   `json:"released,omitempty"`
	Overdue         bool   `json:"overdue,omitempty"`
	UserReleaseDate string `json:"userReleaseDate,omitempty"`
	ProjectId       int    `json:"projectId,omitempty"`
	Self            string `json:"self,omitempty"`
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
	Transition IssueField `json:"transition,omitempty"`
}
