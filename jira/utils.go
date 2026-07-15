package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const TimeFormatJira = "2006-01-02T15:04:05.000-0700"

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
	stamp := j.Time.Format(TimeFormatJira)
	return []byte(`"` + stamp + `"`), nil
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

func findTransitionToStatus(transitions []Transition, targetStatusId string) (Transition, bool) {
	for _, transition := range transitions {
		if transition.To.ID == targetStatusId {
			return transition, true
		}
	}
	return Transition{}, false
}

func findStatusRoute(currentStatusId, targetStatusId string, currentTransitions []Transition) []string {
	if currentStatusId == targetStatusId {
		return []string{currentStatusId}
	}

	graph := knownTransitionStatusGraph()
	if _, ok := graph[currentStatusId]; !ok {
		graph[currentStatusId] = nil
	}
	for _, transition := range currentTransitions {
		graph[currentStatusId] = appendUniqueString(graph[currentStatusId], transition.To.ID)
	}

	visited := map[string]bool{currentStatusId: true}
	queue := [][]string{{currentStatusId}}

	for len(queue) > 0 {
		route := queue[0]
		queue = queue[1:]
		statusId := route[len(route)-1]

		for _, nextStatusId := range graph[statusId] {
			if visited[nextStatusId] {
				continue
			}

			nextRoute := append(append([]string{}, route...), nextStatusId)
			if nextStatusId == targetStatusId {
				return nextRoute
			}

			visited[nextStatusId] = true
			queue = append(queue, nextRoute)
		}
	}

	return nil
}

func knownTransitionStatusGraph() map[string][]string {
	status := Issue.Status
	return map[string][]string{
		status.New:        {status.Assigned, status.NoNeedReaction},
		status.Assigned:   {status.InProgressHRP},
		status.Backlog:    {status.Rated, status.Done},
		status.Rated:      {status.Selected, status.Backlog, status.Done},
		status.Selected:   {status.InProgress, status.Backlog, status.Done},
		status.InProgress: {status.Resolved, status.Backlog, status.Delay},
		status.Resolved:   {status.Done, status.ToRelease, status.Selected, status.Delay},
		status.Delay:      {status.InProgress},
		status.CodeReview: {status.Resolved},
		status.Closed:     {status.Reopened},
	}
}

func appendUniqueString(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}
