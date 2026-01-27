package jira

import (
	"context"
)

// Табличные тесты для SearchAllTasks на базе общего TestSuite с одним сервером
func (s *SearchSuite) TestSearchAllTasks() {
	ctx := context.Background()

	tests := []struct {
		name    string // с нумерацией и описанием на русском
		jql     string
		total   int
		wantLen int
		wantErr bool
	}{
		{name: "01. Пустой запрос — ожидается ошибка", jql: "", total: 10, wantLen: 0, wantErr: true},
		{name: "02. total=0 — пустой результат без ошибок", jql: "project = TEST", total: 0, wantLen: 0, wantErr: false},
		{name: "03. total < pageSize — одна страница (500)", jql: "project = TEST", total: 500, wantLen: 500, wantErr: false},
		{name: "04. total == pageSize — одна полная страница (1000)", jql: "project = TEST", total: 1000, wantLen: 1000, wantErr: false},
		{name: "05. total > pageSize — несколько страниц (1500)", jql: "project = TEST", total: 1500, wantLen: 1500, wantErr: false},
	}

	for _, tc := range tests {
		tc := tc
		s.Run(tc.name, func() {
			// Настраиваем количество задач на этом кейсе
			s.total = tc.total

			j := &jira{BaseUrl: s.srv.URL, Username: "user", Password: "pass"}
			got, err := j.SearchAllTasks(ctx, tc.jql)
			if tc.wantErr {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().Len(got, tc.wantLen)
		})
	}
}
