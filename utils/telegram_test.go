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
		wantChunks []string
	}{
		{name: "1. Не должен делить эмоджи в разные чанки при границе в начале эмоджи",
			input:      "<b>Привет 👋</b>",
			chunkSize:  20,
			wantChunks: []string{"<b>Привет </b>", "<b>👋</b>"}},
		{name: "2. Не должен делить эмоджи в разные чанки при границе в середине эмоджи",
			input:      "<b>Привет 👋</b>",
			chunkSize:  21,
			wantChunks: []string{"<b>Привет </b>", "<b>👋</b>"}},
		{name: "3. Не должен делить эмоджи в разные чанки при границе в конце эмоджи",
			input:      "<b>Привет 👋</b>",
			chunkSize:  22,
			wantChunks: []string{"<b>Привет </b>", "<b>👋</b>"}},
		{name: "4. Не должно быть бесконечного цикла на экстремально маленьком размере чанка и длинной ссылке (по одному символу помимо тэгов)",
			input:     "<a href=\"https://example.com/very/long/url\">Link</a>",
			chunkSize: 10,
			wantChunks: []string{
				"<a href=\"https://example.com/very/long/url\">L</a>",
				"<a href=\"https://example.com/very/long/url\">i</a>",
				"<a href=\"https://example.com/very/long/url\">n</a>",
				"<a href=\"https://example.com/very/long/url\">k</a>",
			}},
		{name: "5. Не должно быть бесконечного цикла на маленьком размере чанка и длинной ссылке, фактический контент корректно режется пополам",
			input:     "<a href=\"https://example.com/very/long/url\">Link</a>",
			chunkSize: 50,
			wantChunks: []string{
				"<a href=\"https://example.com/very/long/url\">Li</a>",
				"<a href=\"https://example.com/very/long/url\">nk</a>",
			}},
		{
			name:      "6. Вложенные теги: корректный порядок закрытия и открытия",
			input:     "<b><i>Ну очень жирный и курсивный текст</i></b>",
			chunkSize: 45,
			// Ожидаем, что теги закроются и откроются в правильном порядке
			wantChunks: []string{
				"<b><i>Ну очень жирный и</i></b>",
				"<b><i> курсивный текст</i></b>",
			},
		},
		{
			name:      "7. Одиночные теги (br): не должны попадать в стек",
			input:     "Строка 1<br/>Строка 2",
			chunkSize: 15,
			wantChunks: []string{
				"Строка 1<br/>",
				"Строка 2",
			},
		},
		{
			name:      "8. Резка очень длинного слова без пробелов",
			input:     "ОченьДлинноеСловоКотороеНеВлезает",
			chunkSize: 20,
			wantChunks: []string{
				"ОченьДлинн",
				"оеСловоКот",
				"ороеНеВлез",
				"ает",
			},
		},
		{
			name:      "9. Текст с атрибутами (сложные теги)",
			input:     "<a href=\"http://example.com\" title=\"test\">Link text here</a>",
			chunkSize: 55,
			wantChunks: []string{
				"<a href=\"http://example.com\" title=\"test\">Link text</a>",
				"<a href=\"http://example.com\" title=\"test\"> here</a>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SmartSplitTextIntoChunks(tt.input, tt.chunkSize)
			require.Equal(t, tt.wantChunks, got)
			for _, chunk := range got {
				err := ValidateChunkHTML(chunk)
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
				err = ValidateChunkHTML(chunk)
				require.NoError(t, err)
			}
		})
	}
}

func Test_SanitizeForTelegram(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "1. Простые разрешенные теги",
			input: "<b>Bold</b> <i>Italic</i> <u>Underline</u> <s>Strike</s>",
			want:  "<b>Bold</b> <i>Italic</i> <u>Underline</u> <s>Strike</s>",
		},
		{
			name:  "2. Очистка запрещенных структурных тегов (413 Error)",
			input: "<html><body><center><h1>413 Large</h1></center><hr></body></html>",
			want:  "413 Large", // Теги удалены
		},
		{
			name:  "3. Ссылки: сохранение href и удаление опасных атрибутов",
			input: `<a href="https://t.me/test" onclick="alert('xss')" style="color:red">Link</a>`,
			want:  `<a href="https://t.me/test">Link</a>`,
		},
		{
			name:  "4. Спойлеры: поддержка тега ",
			input: `<tg-spoiler>Secret</tg-spoiler> <span class="bad">Public</span>`,
			want:  `<tg-spoiler>Secret</tg-spoiler> Public`,
		},
		{
			name:  "5. Форматирование кода и цитат",
			input: "<blockquote>Quote</blockquote> <pre><code>fmt.Print()</code></pre>",
			want:  "<blockquote>Quote</blockquote> <pre><code>fmt.Print()</code></pre>",
		},
		{
			name:  "6. Кастомные эмодзи",
			input: `<tg-emoji emoji-id="53212345">👋</tg-emoji>`,
			want:  `<tg-emoji emoji-id="53212345">👋</tg-emoji>`,
		},
		{
			name:  "7. Экранирование спецсимволов (автоматически)",
			input: "1 < 2 & 3 > 2",
			want:  "1 &lt; 2 &amp; 3 &gt; 2",
		},
		{
			name:  "8. Вложенные запрещенные теги",
			input: "<div><p>Параграф <b>жирный</b></p></div>",
			want:  "Параграф <b>жирный</b>", // div и p удалены как элементы, но контент сохранен
		},
		{
			name:  "9. Протоколы ссылок",
			input: `<a href="https://t.me">Safe</a> <a href="javascript:alert(1)">Unsafe</a>`,
			want:  `<a href="https://t.me">Safe</a> Unsafe`, // unsafe ссылка стала просто текстом (без тега <a>)
		},
		{
			name: "10. Есть заголовок HTML вмест с телом",
			input: `Не удалось вкинуть файлик sql party: ErrValidator: status code 413.
Body:<html>
<head><title>413 Request Entity Too Large</title></head>
<body>
<center><h1>413 Request Entity Too Large</h1></center>
<hr><center>nginx</center>
</body>
</html>`,
			want: "Не удалось вкинуть файлик sql party: ErrValidator: status code 413." +
				"\nBody:" +
				"\n" +
				"\n413 Request Entity Too Large" +
				"\nnginx",
		},
		{
			name:  "11. Пробелы между тегами и переносами строк: должны сохраняться, но не создавать лишних переносов",
			input: "<p>Мне хочется пробелов везде</p><br/><p>везде</p> ",
			want:  "Мне хочется пробелов везде везде",
		},

		{
			name: "12. Нет HTML форматирования, но есть JSON, должен остаться без изменений",
			input: `Не удалось вкинуть файлик sql party: ErrValidator: status code 413.
Body:{
	"errorMessages": [],
	"errors": {
		"summary": "field summary can not be longer than 255 characters",
		"priority": "The priority field is required",
	}
}`,
			want: `Не удалось вкинуть файлик sql party: ErrValidator: status code 413.
Body:{
	"errorMessages": [],
	"errors": {
		"summary": "field summary can not be longer than 255 characters",
		"priority": "The priority field is required",
	}
}`},
		{
			name: "13. Есть заголовок HTML вместе с телом, но они экранированы. Поэтому ничего не делаем",
			input: `Не удалось вкинуть файлик sql party: ErrValidator: status code 413.
Body:&lt;html&gt;
&lt;head&gt;&lt;title&gt;413 Request Entity Too Large&lt;/title&lt;/head&gt;
&lt;body&gt;
&lt;center&gt;&lt;h1&gt;413 Request Entity Too Large&lt;/h1&gt;&lt;/center&gt;
&lt;hr&gt;&lt;center&gt;nginx&lt;/center&gt;
&lt;/body&gt;
&lt;/html&gt;`,
			want: `Не удалось вкинуть файлик sql party: ErrValidator: status code 413.
Body:&lt;html&gt;
&lt;head&gt;&lt;title&gt;413 Request Entity Too Large&lt;/title&lt;/head&gt;
&lt;body&gt;
&lt;center&gt;&lt;h1&gt;413 Request Entity Too Large&lt;/h1&gt;&lt;/center&gt;
&lt;hr&gt;&lt;center&gt;nginx&lt;/center&gt;
&lt;/body&gt;
&lt;/html&gt;`,
		},
		{
			name: "14. Markdown форматирование в тексте должно сохраняться",
			input: `*Тестовое сообщение*, Это *жирный*, а это _курсив_, __Подчеркнутый__ и ~зачеркнутый~, ||Скрытый спойлер||
* Список возможностей:*, • Поддержка [ссылок](https://google.com)`,
			want: `*Тестовое сообщение*, Это *жирный*, а это _курсив_, __Подчеркнутый__ и ~зачеркнутый~, ||Скрытый спойлер||
* Список возможностей:*, • Поддержка [ссылок](https://google.com)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeForTelegram(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_normalizeNewLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "1. Обычные Unix переносы (3 в 2)",
			input: "Line 1\n\n\nLine 2",
			want:  "Line 1\n\nLine 2",
		},
		{
			name:  "2. Windows переносы (3 в 2)",
			input: "Line 1\r\n\r\n\r\nLine 2",
			want:  "Line 1\n\nLine 2",
		},
		{
			name:  "3. Смешанные переносы и пробелы между ними",
			input: "Line 1 \r\n  \n \t\r\nLine 2",
			want:  "Line 1\n\nLine 2",
		},
		{
			name:  "4. Много пустых строк (экстремальный случай)",
			input: "Start\n\n\n\n\n\n\n\nEnd",
			want:  "Start\n\nEnd",
		},
		{
			name:  "5. Пробелы в конце строки не затрагивают переносы",
			input: "Text with spaces    \nNext line",
			want:  "Text with spaces\nNext line",
		},
		{
			name:  "6. Текст без лишних переносов (не должен измениться)",
			input: "Line 1\n\nLine 2",
			want:  "Line 1\n\nLine 2",
		},
		{
			name:  "7. Переносы в самом начале и конце (должны удалиться TrimSpace)",
			input: "\n\n\nHeader\n\nFooter\n\n\n",
			want:  "Header\n\nFooter",
		},
		{
			name:  "8. Только пробельные символы",
			input: "  \n  \r\n  \t  ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeNewLines(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
