package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SearchSuite struct {
	suite.Suite
	srv   *httptest.Server
	mux   *http.ServeMux
	total int // общий параметр, который тесты меняют перед вызовом
}

// makeIssues формирует список тестовых задач Jira: ID=1..total
func makeIssues(total int) []IssueJira {
	res := make([]IssueJira, 0, total)
	for i := 1; i <= total; i++ {
		res = append(res, IssueJira{Id: strconv.Itoa(i), Key: "KEY-" + strconv.Itoa(i)})
	}
	return res
}

// SetupSuite поднимает общий httptest.Server и настраивает обработчик /search
func (s *SearchSuite) SetupSuite() {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		q := r.URL.Query()
		startAt, _ := strconv.Atoi(q.Get("startAt"))
		maxResults, _ := strconv.Atoi(q.Get("maxResults"))

		issues := makeIssues(s.total)
		start := startAt
		if start > len(issues) {
			start = len(issues)
		}
		end := startAt + maxResults
		if end > len(issues) {
			end = len(issues)
		}
		page := issues[start:end]

		resp := SearchResponse{StartAt: startAt, MaxResults: len(page), Total: s.total, Issues: page}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(resp)
		s.NoError(err, "encode response")
	})

	s.srv = httptest.NewServer(s.mux)
}

// TearDownSuite завершает работу сервера
func (s *SearchSuite) TearDownSuite() {
	if s.srv != nil {
		s.srv.Close()
	}
}

// Точка входа для запуска сьюта
func TestSearchSuite(t *testing.T) {
	suite.Run(t, new(SearchSuite))
}
