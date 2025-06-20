package jira

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestParseIssue(t *testing.T) {
	tests := []struct {
		name    string
		srcPath string
		want    IssueJira
	}{
		{name: "1. Задачка из проекта INNA", srcPath: "./test_data/issue-inna.json",
			want: IssueJira{Id: "418688", Key: "INNA-842", Fields: FieldsIssue{
				Summary:       "Мелкие улучшения дайджеста",
				Description:   "(?)",
				BusinessValue: 5,
				StoryPoints:   8,
				WeightedJob:   5,
				Status:        IssueField{ID: "10425", Name: "Выбрано"},
				IssueType:     IssueField{ID: "3", Name: "Задача"},
				Priority:      IssueField{ID: "4", Name: "Незначительный"},
				Resolution:    IssueField{ID: "10300", Name: "Нужен багфикс"},
				Assignee:      JiraUser{Name: "ilyavas", Key: "JIRAUSER45200", Email: "test@yandex.ru", DisplayName: "Илья Васильев", Active: true},
				Creator:       JiraUser{Name: "ilyavas", Key: "JIRAUSER45200", Email: "test@yandex.ru", DisplayName: "Илья Васильев", Active: true},
				Reporter:      JiraUser{Name: "ilyavas", Key: "JIRAUSER45200", Email: "test@yandex.ru", DisplayName: "Илья Васильев", Active: true},
				Participants:  []JiraUser{{Name: "ilyavas", Key: "JIRAUSER45200", Email: "test@yandex.ru", DisplayName: "Илья Васильев", Active: true}},
				Project:       IssueField{ID: "18510", Name: "Inner Automation", Key: "INNA"},
				Components:    []IssueField{{ID: "48150", Name: "Jira (Джира)"}},
				Created:       JiraTime{time.Date(2025, time.January, 17, 12, 3, 44, 0, time.Local)},
				Updated:       JiraTime{time.Date(2025, time.April, 1, 19, 55, 33, 0, time.Local)},
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.srcPath)
			require.NoError(t, err)
			var got IssueJira
			err = json.Unmarshal(data, &got)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}

}
