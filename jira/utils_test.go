package jira

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"
)

func TestParseJiraDate(t *testing.T) {
	tests := []struct {
		name         string
		respFilepath string
		createdDate  JiraTime
	}{
		{"01. INNA-565 standard format", "comment-response.json", JiraTime{time.Date(2024, 4, 9, 11, 35, 17, 0, time.Local)}}, //"2024-04-09T11:35:17.000+0300"
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respContent, errT := os.ReadFile(path.Join("test_data", tt.respFilepath))
			if errT != nil {
				t.Fatal(errT)
			}
			var want IssueCommentsResponse
			err := json.Unmarshal(respContent, &want)
			require.NoError(t, err)
			require.Equal(t, tt.createdDate, want.Comments[0].Created)
		})
	}
}

func TestParseJiraTimestamp(t *testing.T) {
	tests := []struct {
		name         string
		respFilepath string
		timestamp    Timestamp
	}{
		{"01. create issue webhook", "created-issue-webhook.json", Timestamp{time.Date(2024, time.November, 6, 20, 15, 55, 829000000, time.Local)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respContent, errT := os.ReadFile(path.Join("test_data", tt.respFilepath))
			if errT != nil {
				t.Fatal(errT)
			}
			var want WebhookIssue
			err := json.Unmarshal(respContent, &want)
			require.NoError(t, err)
			require.Equal(t, tt.timestamp, want.Timestamp)
		})
	}
}
