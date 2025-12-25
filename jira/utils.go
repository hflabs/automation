package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type JiraTime struct {
	time.Time
}

func (j *JiraTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse("2006-01-02T15:04:05.000Z0700", s)
	if err != nil {
		return err
	}
	j.Time = t
	return nil
}

func (j *JiraTime) MarshalJSON() ([]byte, error) {
	return []byte(j.Format("2006-01-02T15:04:05.000Z0700")), nil
}

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UnixMilli())
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var unix int64
	if err := json.Unmarshal(data, &unix); err != nil {
		return err
	}
	t.Time = time.UnixMilli(unix).Local()
	return nil
}

func validateStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("status code %v.\nBody:%s", resp.StatusCode, string(b))
}

func formatAvailableStatuses(availableStatuses map[string]string) string {
	pairs := make([]string, 0, len(availableStatuses))
	for id, name := range availableStatuses {
		pairs = append(pairs, fmt.Sprintf("%s:%v", id, name))
	}
	return strings.Join(pairs, ", ")
}
