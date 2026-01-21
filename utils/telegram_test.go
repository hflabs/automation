package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mymmrac/telego"
	"github.com/stretchr/testify/require"
)

func TestConvertTgLinks(t *testing.T) {
	tests := []struct {
		name        string
		msgText     string
		msgEntities []telego.MessageEntity
		want        string
	}{
		{name: "1. Нет никаких ссылок", msgText: "Привет!\n\nТут не ок, что почему-то не попали версии раньше 2.1.1\nВ РН и инструкции 2.0.0 всё это есть",
			want: "Привет!\n\nТут не ок, что почему-то не попали версии раньше 2.1.1\nВ РН и инструкции 2.0.0 всё это есть"},
		{name: "2. Есть одна ссылка", msgText: "Привет!\n\nТут не ок, что почему-то не попали версии раньше 2.1.1\nВ РН и инструкции 2.0.0 всё это есть",
			msgEntities: []telego.MessageEntity{{Type: MsgEntityLinkType, Offset: 71, Length: 16, URL: "https://confluence.ru/pages/viewpage.action?pageId=1724023564"}},
			want:        "Привет!\n\nТут не ок, что почему-то не попали версии раньше 2.1.1\nВ РН и <a href='https://confluence.ru/pages/viewpage.action?pageId=1724023564'>инструкции 2.0.0</a> всё это есть"},
		{name: "3. Есть несколько ссылок", msgText: "очень странно (в спэйсе апдейтера это есть). \nПросто пример вот. Это ведь не ок?",
			msgEntities: []telego.MessageEntity{
				{Type: MsgEntityLinkType, Offset: 17, Length: 16, URL: "https://confluence.ru/pages/viewpage.action?pageId=1784283523"},
				{Type: MsgEntityLinkType, Offset: 53, Length: 6, URL: "https://confluence.ru/pages/viewpage.action?pageId=1787790238"}},
			want: "очень странно (в <a href='https://confluence.ru/pages/viewpage.action?pageId=1784283523'>спэйсе апдейтера</a> это есть). \nПросто <a href='https://confluence.ru/pages/viewpage.action?pageId=1787790238'>пример</a> вот. Это ведь не ок?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertTgLinks(tt.msgText, tt.msgEntities)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_SplitTextIntoChunksWithSize(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		chunkSize  int
		wantChunks int
	}{
		{name: "1. Делит эмоджи в разные чанки",
			input:      "<b>Привет 👋</b>",
			chunkSize:  11,
			wantChunks: 9999},
		{name: "2. бесконечный цикл",
			input:      "<a href=\"https://example.com/very/long/url\">Link</a>",
			chunkSize:  10,
			wantChunks: 9999},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SmartSplitTextIntoChunks(tt.input, tt.chunkSize)
			require.Len(t, got, tt.wantChunks)
			for _, chunk := range got {
				err := validateChunkHTML(chunk)
				require.NoError(t, err)
			}
		})
	}
}

func Test_SplitTextIntoChunks(t *testing.T) {
	tests := []struct {
		name       string
		sourceFile string
		wantPrefix string
		wantChunks int
	}{
		{name: "1. Длинный кусок текста с HTML тэгами посреди которых может порезаться сообщение", sourceFile: "long_text_with_html_markdown.html",
			wantPrefix: "Борис, привет!", wantChunks: 5},
		{name: "2. Код-ревью от LLM с HTML разметкой", sourceFile: "code_review_ai_html_markdown.html",
			wantPrefix: "## PR Reviewer Guide", wantChunks: 2},
		{name: "3. Оповещение о комментарии в МР с код-ревью от LLM с HTML разметкой", sourceFile: "notification_comment_with_code_review_ai_html_markdown.html",
			wantPrefix: "Петр(@petr) оставил комментарий в твоём Merge Request Тестовый", wantChunks: 2},
		{name: "4. Багфикс дайджеста №1", sourceFile: "digest_short_1.html",
			wantPrefix: "Екатерина, привет!", wantChunks: 1},
		{name: "5. Багфикс дайджеста №2", sourceFile: "digest_short_2.html",
			wantPrefix: "Максим, привет!", wantChunks: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.ReadFile(fmt.Sprintf("./test_data/%s", tt.sourceFile))
			require.NoError(t, err)

			got := SmartSplitTextIntoChunks(string(file), 4096)
			require.Len(t, got, tt.wantChunks)
			if tt.wantChunks > 0 && !strings.HasPrefix(got[0], tt.wantPrefix) {
				t.Errorf("want prefix `%s`, got `%s`", got, tt.wantPrefix)
			}
			for _, chunk := range got {
				err = validateChunkHTML(chunk)
				require.NoError(t, err)
			}
		})
	}
}
