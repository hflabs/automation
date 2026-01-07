package jira

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"
)

// Временные интеграционные тесты против реальной Jira.
// Запускаются ТОЛЬКО если заданы переменные окружения:
//  - JIRA_BASE_URL   — базовый URL, как ожидает клиент сейчас, например: https://host/rest/api/2
//  - JIRA_USER       — имя пользователя (basic auth)
//  - JIRA_PASS       — пароль (basic auth)
//  - JIRA_PROJECT_KEY или JIRA_PROJECT_ID — идентификация проекта
//  - JIRA_ISSUETYPE_NAME или JIRA_ISSUETYPE_ID — тип задачи
// Если переменные не заданы, тесты будут пропущены.

type realCfg struct {
	baseURL string
	user    string
	pass    string
	projKey string
	projID  string
	itName  string
	itID    string
}

// Если локально влом прописывать переменные в ENV на винде
func setEnvTest(url, user, pass, projKey, itName string) {
	os.Setenv("JIRA_BASE_URL", url)
	os.Setenv("JIRA_USER", user)
	os.Setenv("JIRA_PASS", pass)
	os.Setenv("JIRA_PROJECT_KEY", projKey)
	os.Setenv("JIRA_ISSUETYPE_NAME", itName)
}

func loadRealCfg(t *testing.T) (realCfg, bool) {
	t.Helper()
	cfg := realCfg{
		baseURL: os.Getenv("JIRA_BASE_URL"),
		user:    os.Getenv("JIRA_USER"),
		pass:    os.Getenv("JIRA_PASS"),
		projKey: os.Getenv("JIRA_PROJECT_KEY"),
		projID:  os.Getenv("JIRA_PROJECT_ID"),
		itName:  os.Getenv("JIRA_ISSUETYPE_NAME"),
		itID:    os.Getenv("JIRA_ISSUETYPE_ID"),
	}

	if cfg.baseURL == "" || cfg.user == "" || cfg.pass == "" {
		t.Skip("Пропуск: задайте JIRA_BASE_URL, JIRA_USER, JIRA_PASS для запуска интеграционных тестов")
		return realCfg{}, false
	}
	// Нужен проект и тип задачи — допускаем указание по key/id и name/id соответственно
	if cfg.projKey == "" && cfg.projID == "" {
		t.Skip("Пропуск: задайте JIRA_PROJECT_KEY или JIRA_PROJECT_ID")
		return realCfg{}, false
	}
	if cfg.itName == "" && cfg.itID == "" {
		t.Skip("Пропуск: задайте JIRA_ISSUETYPE_NAME или JIRA_ISSUETYPE_ID")
		return realCfg{}, false
	}
	return cfg, true
}

func uniqueSummary(prefix string) string {
	return prefix + " - AUTOTEST - " + strconv.FormatInt(time.Now().UnixNano(), 10)
}

// TestReal_CreateIssueFromMap_And_UpdateFromMap — проверяет создание и обновление через map
func TestReal_CreateIssueFromMap_And_UpdateFromMap(t *testing.T) {
	cfg, ok := loadRealCfg(t)
	if !ok {
		return
	}

	j := NewJira(cfg.baseURL, cfg.user, cfg.pass)
	ctx := context.Background()

	// Формируем тело создания через map
	fields := map[string]any{
		"summary":           uniqueSummary("Создание через map"),
		"description":       "created via map",
		"components":        []map[string]any{{"name": "Мониторинг"}},
		"customfield_10000": "test",
		"customfield_12680": []map[string]any{{"value": "Другое (напишу в задаче)"}},
	}
	// project
	if cfg.projID != "" {
		fields["project"] = map[string]any{"id": cfg.projID}
	} else {
		fields["project"] = map[string]any{"key": cfg.projKey}
	}
	// issuetype
	if cfg.itID != "" {
		fields["issuetype"] = map[string]any{"id": cfg.itID}
	} else {
		fields["issuetype"] = map[string]any{"name": cfg.itName}
	}

	created, err := j.CreateIssueFromMap(ctx, fields)
	if err != nil {
		t.Fatalf("CreateIssueFromMap error: %v", err)
	}
	if created.Key == "" || created.ID == "" {
		t.Fatalf("CreateIssueFromMap: пустые Key/ID в ответе: %+v", created)
	}

	// Обновляем summary и description через map
	upd := map[string]any{
		"summary":     uniqueSummary("Обновление через map"),
		"description": "updated via map",
	}
	if err := j.UpdateIssueFromMap(ctx, created.Key, upd); err != nil {
		t.Fatalf("UpdateIssueFromMap error: %v", err)
	}

	// Проверяем через GetIssueById (в Jira это idOrKey)
	got, err := j.GetIssueById(ctx, created.Key)
	if err != nil {
		t.Fatalf("GetIssueById error: %v", err)
	}
	if got.Fields.Summary != upd["summary"].(string) {
		t.Fatalf("summary не обновился: want %q, got %q", upd["summary"], got.Fields.Summary)
	}
}

// TestReal_CreateIssue_And_Update — проверяет создание и обновление через типизированные структуры
func TestReal_CreateIssue_And_Update(t *testing.T) {
	cfg, ok := loadRealCfg(t)
	if !ok {
		return
	}

	j := NewJira(cfg.baseURL, cfg.user, cfg.pass)
	ctx := context.Background()

	req := FieldsIssue{
		Summary:             uniqueSummary("Создание через структуру"),
		Description:         "created via struct",
		Components:          []IssueField{{Name: "Мониторинг"}},
		WhoWillGetBetter:    []IssueField{{Value: "Другое (напишу в задаче)"}},
		BusinessDescription: "тест",
	}
	if cfg.projID != "" {
		req.Project = IssueField{ID: cfg.projID}
	} else {
		req.Project = IssueField{Key: cfg.projKey}
	}
	if cfg.itID != "" {
		req.IssueType = IssueField{ID: cfg.itID}
	} else {
		req.IssueType = IssueField{Name: cfg.itName}
	}

	created, err := j.CreateIssue(ctx, req)
	if err != nil {
		t.Fatalf("CreateIssue error: %v", err)
	}
	if created.Key == "" || created.ID == "" {
		t.Fatalf("CreateIssue: пустые Key/ID в ответе: %+v", created)
	}

	// Обновление через типизированный запрос
	newSummary := uniqueSummary("Обновление через структуру")
	if err := j.UpdateIssue(ctx, created.Key, FieldsIssue{
		Summary:     newSummary,
		Description: "updated via struct",
	}); err != nil {
		t.Fatalf("UpdateIssue error: %v", err)
	}

	got, err := j.GetIssueById(ctx, created.Key)
	if err != nil {
		t.Fatalf("GetIssueById error: %v", err)
	}
	if got.Fields.Summary != newSummary {
		t.Fatalf("summary не обновился: want %q, got %q", newSummary, got.Fields.Summary)
	}
}
