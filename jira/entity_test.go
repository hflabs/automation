package jira

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseIssue(t *testing.T) {
	tests := []struct {
		name    string
		srcPath string
		want    IssueJira
	}{
		{name: "1. Задачка из проекта INNA", srcPath: "./test_data/issue-inna.json",
			want: IssueJira{Id: "100842", Key: "TEST-842", Fields: FieldsIssue{
				Summary:             "Minor digest improvements",
				Description:         "(?)",
				BusinessValue:       5,
				StoryPoints:         8,
				WeightedJob:         5,
				Status:              IssueField{ID: "10425", Name: "Selected"},
				IssueType:           IssueField{ID: "3", Name: "Task"},
				Priority:            IssueField{ID: "4", Name: "Minor"},
				Resolution:          IssueField{ID: "10300", Name: "Bugfix required"},
				Assignee:            JiraUser{Name: "testuser", Key: "JIRAUSER45200", Email: "test@example.com", DisplayName: "Test User", Active: true},
				Creator:             JiraUser{Name: "testuser", Key: "JIRAUSER45200", Email: "test@example.com", DisplayName: "Test User", Active: true},
				Reporter:            JiraUser{Name: "testuser", Key: "JIRAUSER45200", Email: "test@example.com", DisplayName: "Test User", Active: true},
				Participants:        []JiraUser{{Name: "testuser", Key: "JIRAUSER45200", Email: "test@example.com", DisplayName: "Test User", Active: true}},
				Project:             IssueField{ID: "10001", Name: "Test Automation", Key: "TEST"},
				Components:          []IssueField{{ID: "48150", Name: "Jira"}},
				Created:             JiraTime{time.Date(2025, time.January, 17, 12, 3, 44, 0, time.Local)},
				Updated:             JiraTime{time.Date(2025, time.April, 1, 19, 55, 33, 0, time.Local)},
				BusinessDescription: "Make better",
				WhoWillGetBetter:    []IssueField{{ID: "11854", Value: "All Staff"}},
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
