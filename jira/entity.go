package jira

type WebhookIssue struct {
	Timestamp      Timestamp    `json:"timestamp,omitzero"`
	WebhookEvent   string       `json:"webhookEvent,omitzero"`
	IssueEventType string       `json:"issue_event_type_name,omitzero"`
	UserEvent      JiraUser     `json:"user,omitzero"`
	Issue          IssueJira    `json:"issue,omitzero"`
	Comment        IssueComment `json:"comment,omitzero"`
	Changelog      ChangeLog    `json:"changelog,omitzero"`
	Version        Version      `json:"version,omitzero"`
}

type ChangeLog struct {
	Id      string          `json:"id,omitzero"`
	Author  JiraUser        `json:"author,omitzero"`
	Created JiraTime        `json:"created,omitzero"`
	Items   []ChangelogItem `json:"items,omitzero"`
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
	Timestamp    Timestamp    `json:"timestamp,omitzero"`
	WebhookEvent string       `json:"webhookEvent,omitzero"`
	Comment      IssueComment `json:"comment,omitzero"`
}

type IssueJira struct {
	Id        string       `json:"id,omitzero"`
	Key       string       `json:"key,omitzero"`
	Fields    FieldsIssue  `json:"fields,omitzero"`
	Changelog IssueHistory `json:"changelog,omitzero"`
}

type IssueHistory struct {
	Total     int         `json:"total,omitzero"`
	Histories []ChangeLog `json:"histories,omitzero"`
}

type ChangelogItem struct {
	Field      string `json:"field,omitzero,omitzero"`
	FieldType  string `json:"fieldType,omitzero"`
	From       string `json:"from,omitzero"`
	To         string `json:"to,omitzero"`
	FromString string `json:"fromString,omitzero"`
	ToString   string `json:"toString,omitzero"`
}

type FieldsIssue struct {
	Summary             string       `json:"summary,omitzero"`
	Description         string       `json:"description,omitzero"`
	BusinessValue       float64      `json:"customfield_10084,omitzero"`
	StoryPoints         float64      `json:"customfield_10083,omitzero"`
	WeightedJob         float64      `json:"customfield_12580,omitzero"`
	ReleaseNotes        string       `json:"customfield_13082,omitzero"`
	ReleaseInstruction  string       `json:"customfield_13081,omitzero"`
	DueDate             string       `json:"duedate,omitzero"`
	Status              IssueField   `json:"status,omitzero"`
	IssueType           IssueField   `json:"issuetype,omitzero"`
	Priority            IssueField   `json:"priority,omitzero"`
	Resolution          IssueField   `json:"resolution,omitzero"`
	Assignee            JiraUser     `json:"assignee,omitzero"`
	Creator             JiraUser     `json:"creator,omitzero"`
	Reporter            JiraUser     `json:"reporter,omitzero"`
	Participants        []JiraUser   `json:"customfield_10380,omitzero"`
	Project             IssueField   `json:"project,omitzero"`
	Components          []IssueField `json:"components,omitzero"`
	LearnTime           string       `json:"customfield_14481,omitzero"`
	LearnForWho         string       `json:"customfield_13881,omitzero"`
	LearnWhatLike       string       `json:"customfield_14483,omitzero"`
	LearnWhatUseful     string       `json:"customfield_13784,omitzero"`
	LearnWhatBad        string       `json:"customfield_14484,omitzero"`
	LearnWhatLearned    string       `json:"customfield_13880,omitzero"`
	LearnWillRecommend  string       `json:"customfield_13882,omitzero"`
	LearnPeople         []JiraUser   `json:"customfield_14480,omitzero"`
	LearnField          IssueField   `json:"customfield_14380,omitzero"`
	LearnLink           string       `json:"customfield_13782,omitzero"`
	Created             JiraTime     `json:"created,omitzero"`
	Updated             JiraTime     `json:"updated,omitzero"`
	FreeStringValue     string       `json:"freeValue,omitzero"`
	BusinessDescription string       `json:"customfield_10000,omitzero"`
	WhoWillGetBetter    []IssueField `json:"customfield_12680,omitzero"`
}

type IssueField struct {
	ID    string `json:"id,omitzero"`
	Name  string `json:"name,omitzero"`
	Key   string `json:"key,omitzero"`
	Value string `json:"value,omitzero"`
}

type JiraUser struct {
	Name        string `json:"name,omitzero"`
	Key         string `json:"key,omitzero"`
	Email       string `json:"emailAddress,omitzero"`
	DisplayName string `json:"displayName,omitzero"`
	Active      bool   `json:"active,omitzero"`
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
	Watchers []JiraUser `json:"watchers,omitzero"`
}

type IssueCommentsResponse struct {
	Total    int            `json:"total,omitzero"`
	Comments []IssueComment `json:"comments,omitzero"`
}

type IssueComment struct {
	Author       JiraUser `json:"author,omitzero"`
	Body         string   `json:"body,omitzero"`
	UpdateAuthor JiraUser `json:"update_author,omitzero"`
	Created      JiraTime `json:"created,omitzero"`
	Updated      JiraTime `json:"updated,omitzero"`
}

type Version struct {
	Id              string `json:"id,omitzero"`
	Name            string `json:"name,omitzero"`
	Description     string `json:"description,omitzero"`
	Archived        bool   `json:"archived,omitzero"`
	Released        bool   `json:"released,omitzero"`
	Overdue         bool   `json:"overdue,omitzero"`
	UserReleaseDate string `json:"userReleaseDate,omitzero"`
	ProjectId       int    `json:"projectId,omitzero"`
	Self            string `json:"self,omitzero"`
}

type SearchResponse struct {
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Issues     []IssueJira `json:"issues"`
}

type UpsertIssueRequestFromMap struct {
	Fields map[string]interface{} `json:"fields"`
}

type UpsertIssueRequest struct {
	Fields FieldsIssue `json:"fields"`
}

// CreatedIssueResponse — ответ Jira на создание задачи (POST /issue)
// Содержит минимальные поля, которые Jira возвращает по умолчанию.
type CreatedIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

type TransitionIssueRequest struct {
	Transition IssueField `json:"transition,omitzero"`
}
type TransitionsResponse struct {
	Expand      string       `json:"expand"`
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	ID   string     `json:"id"`
	Name string     `json:"name"` // Название перехода (например "Start Progress")
	To   IssueField `json:"to"`   // Целевой статус
}
