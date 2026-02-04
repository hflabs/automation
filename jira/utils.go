package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const TimeFormatJira = "2006-01-02T15:04:05.000Z0700"

type JiraTime struct {
	time.Time
}

func (j *JiraTime) UnmarshalJSON(b []byte) error {
	// 1. Безопасное удаление кавычек
	s := strings.Trim(string(b), "\"")
	// 2. Проверка на null
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(TimeFormatJira, s)
	if err != nil {
		return err
	}
	j.Time = t
	return nil
}

func (j JiraTime) MarshalJSON() ([]byte, error) {
	return []byte(j.Format(TimeFormatJira)), nil
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

func formatAvailableStatuses(availableStatuses []Transition) string {
	pairs := strings.Builder{}
	for index, status := range availableStatuses {
		pairs.WriteString(fmt.Sprintf("%s:%v", status.ID, status.Name))
		if index != len(availableStatuses)-1 {
			pairs.WriteString(", ")
		}
	}
	return pairs.String()
}
