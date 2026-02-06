package utils

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/carlmjohnson/requests"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mymmrac/telego"
	"golang.org/x/net/html"
)

const (
	FileEndpoint      = "https://api.telegram.org/file/bot%s/%s"
	MsgEntityLinkType = "text_link"
	PrivateChatType   = "private"
)

func DownloadFile(url, filepath string) error {
	return requests.URL(url).ToFile(filepath).Fetch(context.Background())
}

func CreateLink(token, filePath string) string {
	return fmt.Sprintf(FileEndpoint, token, filePath)
}

func IsTextOrCaption(update telego.Update) bool {
	if update.Message == nil {
		return false
	}
	if update.Message.Text != "" {
		return true
	}
	if update.Message.Caption != "" {
		return true
	}
	return false
}

func GetTextOrCaption(update telego.Message) string {
	if update.Text != "" {
		return update.Text
	}
	if update.Caption != "" {
		return update.Caption
	}
	return ""
}

func IsCommand(update telego.Update) bool {
	if update.Message == nil {
		return false
	}
	if strings.HasPrefix(update.Message.Text, "/") {
		return true
	}
	return false
}

func IsPrivateCommand(update telego.Update) bool {
	return IsCommand(update) && update.Message.Chat.Type == PrivateChatType
}

func IsMessage(update telego.Update) bool {
	if update.Message == nil {
		return false
	}
	if IsTextOrCaption(update) {
		return true
	}
	return false
}

func IsPrivateMessage(update telego.Update) bool {
	return IsMessage(update) && update.Message.Chat.Type == PrivateChatType
}

func IsPhoto(update telego.Update) bool {
	if update.Message.Photo != nil {
		return true
	}
	return false
}

func IsPrivatePhoto(update telego.Update) bool {
	return IsPhoto(update) && update.Message.Chat.Type == PrivateChatType
}

func IsDocument(update telego.Update) bool {
	if update.Message.Document != nil {
		return true
	}
	return false
}

func IsPrivateDocument(update telego.Update) bool {
	return IsDocument(update) && update.Message.Chat.Type == PrivateChatType
}

func IsCallback(update telego.Update) bool {
	if update.CallbackQuery != nil {
		return true
	}
	return false
}

func IsReaction(update telego.Update) bool {
	if update.MessageReaction != nil {
		return true
	}
	return false
}

func EscapeHtmlQuotes(text string) string {
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

func ConvertTgLinks(msgText string, msgEntities []telego.MessageEntity) string {
	if len(msgEntities) == 0 {
		return msgText
	}
	runes := []rune(msgText)
	result := strings.Builder{}
	offset := 0
	for _, msgEntity := range msgEntities {
		if msgEntity.Type == MsgEntityLinkType && msgEntity.URL != "" {
			result.WriteString(string(runes[offset:msgEntity.Offset]))
			result.WriteString(fmt.Sprintf("<a href='%s'>%s</a>",
				msgEntity.URL, string(runes[msgEntity.Offset:msgEntity.Offset+msgEntity.Length])))
			offset = msgEntity.Offset + msgEntity.Length
		}
	}
	result.WriteString(string(runes[offset:]))
	return result.String()
}

// Стек открытых тегов. Храним полный текст тега, например "<b class='x'>"
type tagStackItem struct {
	Name  string // например "b", "a"
	Token string // полный текст открывающего тега
}

func newTagStackItem(name, token string) tagStackItem {
	return tagStackItem{name, token}
}

// SmartSplitTextIntoChunks разбивает HTML текст на части, сохраняя валидность HTML.
func SmartSplitTextIntoChunks(s string, chunkSize int) []string {
	if len(s) <= chunkSize {
		return []string{s}
	}
	var chunks []string
	var currentChunk bytes.Buffer
	var stack []tagStackItem
	var hasContent bool // Флаг: добавлен ли в текущий чанк текст или void-теги

	z := html.NewTokenizer(strings.NewReader(s))

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		token := z.Token()
		tokenString := token.String()

		for len(tokenString) > 0 {
			closingOverhead := countClosingOverhead(stack)

			// Считаем оверхед на открытие тегов, если чанк будет новым
			openingOverhead := 0
			if currentChunk.Len() == 0 {
				for _, t := range stack {
					openingOverhead += len(t.Token)
				}
			}

			availableSpace := chunkSize - currentChunk.Len() - openingOverhead - closingOverhead

			// 1. Если места нет и в чанке уже есть полезный контент — сбрасываем
			if availableSpace <= 0 && hasContent {
				chunks = append(chunks, flushChunk(&currentChunk, stack))
				hasContent = false
				continue
			}

			// 2. Если чанк пустой, подготавливаем его (открываем теги из стека)
			if currentChunk.Len() == 0 && len(stack) > 0 {
				currentChunk.WriteString(reopenTags(stack))
				// Пересчитываем место
				availableSpace = chunkSize - currentChunk.Len() - closingOverhead
			}

			// 3. Обработка контента
			if tt == html.TextToken {
				head, tail := splitTextBySpace(tokenString, availableSpace)

				// Если не влезло ни слова и в чанке пусто — берем хоть что-то для прогресса
				if head == "" && !hasContent {
					head, tail = splitByRune(tokenString, 1)
				}

				if head != "" {
					currentChunk.WriteString(head)
					hasContent = true
					tokenString = tail
				} else {
					// Если head пустой и контент уже был — сбрасываем чанк
					chunks = append(chunks, flushChunk(&currentChunk, stack))
					hasContent = false
				}
			} else {
				// Это тег. Просто добавляем его в текущий чанк.
				// Теги не вызывают немедленный flush, чтобы не плодить пустые <a></a>
				currentChunk.WriteString(tokenString)
				updateStack(tt, token, tokenString, &stack)
				// Если это одиночный тег (br, img), это считается контентом
				if isVoidElement(token.Data) {
					hasContent = true
				}
				tokenString = "" // Тег обработан целиком
			}
		}
	}
	// Финальный сброс, если в последнем чанке было хоть что-то полезное
	if hasContent {
		chunks = append(chunks, flushChunk(&currentChunk, stack))
	}
	return chunks
}

// splitTextBySpace делит текст на две части:
// head - часть, которая влезает в limit (стараясь не резать слова и символы UTF-8)
// tail - остаток
func splitTextBySpace(text string, limit int) (head, tail string) {
	if limit <= 0 {
		return "", text
	}
	if len(text) <= limit {
		return text, ""
	}
	// Находим последнюю корректную границу UTF-8 символа внутри limit
	lastValidIdx := 0
	for i := range text {
		if i > limit {
			break
		}
		lastValidIdx = i
	}

	candidate := text[:lastValidIdx]
	lastSpaceIndex := strings.LastIndexFunc(candidate, unicode.IsSpace)
	if lastSpaceIndex > 0 {
		return text[:lastSpaceIndex+1], text[lastSpaceIndex+1:]
	}
	return text[:lastValidIdx], text[lastValidIdx:]
}

// flushChunk — Закрывает теги, возвращает строку чанка и очищает буфер
func flushChunk(buf *bytes.Buffer, stack []tagStackItem) string {
	// Закрываем теги в обратном порядке (LIFO)
	for i := len(stack) - 1; i >= 0; i-- {
		buf.WriteString("</" + stack[i].Name + ">")
	}
	res := buf.String()
	buf.Reset()
	return res
}

// reopenTags — Создает строку с открывающими тегами для нового чанка
func reopenTags(stack []tagStackItem) string {
	var sb strings.Builder
	// Открываем теги в прямом порядке (FIFO)
	for _, tag := range stack {
		sb.WriteString(tag.Token)
	}
	return sb.String()
}

// countClosingOverhead — Подсчет длины закрывающих тегов
func countClosingOverhead(stack []tagStackItem) int {
	overhead := 0
	for _, tag := range stack {
		overhead += 2 + len(tag.Name) + 1 // </tag>
	}
	return overhead
}

// updateStack — Обновление стека при встрече тегов
func updateStack(tt html.TokenType, token html.Token, tokenString string, stack *[]tagStackItem) {
	switch tt {
	case html.StartTagToken:
		if !isVoidElement(token.Data) {
			*stack = append(*stack, tagStackItem{
				Name:  token.Data,
				Token: tokenString,
			})
		}
	case html.EndTagToken:
		s := *stack
		for i := len(s) - 1; i >= 0; i-- {
			if s[i].Name == token.Data {
				*stack = append(s[:i], s[i+1:]...)
				break
			}
		}
	}
}

func validateChunkHTML(chunk string) error {
	z := html.NewTokenizer(strings.NewReader(chunk))
	var stack []string

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		token := z.Token()
		switch tt {
		case html.StartTagToken:
			// Добавляем в стек только если это НЕ void элемент (как br или img)
			if !isVoidElement(token.Data) {
				stack = append(stack, token.Data)
			}
		case html.EndTagToken:
			// Закрывающий тег. Проверяем стек.
			if len(stack) == 0 {
				return fmt.Errorf("found closing tag </%s> without opening tag", token.Data)
			}
			last := stack[len(stack)-1]
			// Если теги не совпадают
			if last != token.Data {
				return fmt.Errorf("tag nesting error: expected closing </%s>, found </%s>", last, token.Data)
			}
			// Убираем из стека
			stack = stack[:len(stack)-1]
		}
	}
	if len(stack) > 0 {
		return fmt.Errorf("unclosed tags at the end of chunk: %v", stack)
	}
	return nil
}

func isVoidElement(tagName string) bool {
	// Список тегов, которые не требуют закрытия согласно спецификации HTML
	// Для Telegram чаще всего актуальны br, img, hr.
	switch strings.ToLower(tagName) {
	case "br", "img", "hr", "input", "meta", "link", "col", "base", "area", "param":
		return true
	}
	return false
}

// Вспомогательная функция для взятия N символов (не байт!)
func splitByRune(text string, n int) (head, tail string) {
	count := 0
	for i := range text {
		if count == n {
			return text[:i], text[i:]
		}
		count++
	}
	return text, ""
}

func SanitizeForTelegram(input string) string {
	// 0. Предотвращаем слипание текста:
	// Заменяем <br> на пробелы, и после закрывающих блочных тегов добавляем пробел.
	// Это превратит "</p><p>" в "</p> <p>", и после удаления тегов текст не слипнется.
	input = reBr.ReplaceAllString(input, " ")
	// 1. Защищаем кастомные теги спойлера Telegram от вырезания парсером bluemonday.
	// Мы временно заменяем их на уникальные текстовые маркеры.
	replacer := strings.NewReplacer(
		"<tg-spoiler>", "___TGS_OPEN___",
		"</tg-spoiler>", "___TGS_CLOSE___",
	)
	restoreReplacer := strings.NewReplacer(
		"___TGS_OPEN___", "<tg-spoiler>",
		"___TGS_CLOSE___", "</tg-spoiler>",
	)
	protected := replacer.Replace(input)

	p := bluemonday.NewPolicy()
	// 1. Простые текстовые теги
	p.AllowElements(
		"b", "strong", // жирный
		"i", "em", // курсив
		"u", "ins", // подчеркивание
		"s", "strike", "del", // зачеркивание
		"code", "pre", // код
		"blockquote", // цитаты
		"tg-spoiler", // спойлеры (специфичный тег)
	)
	// 2. Ссылки (тег <a>)
	// Разрешаем только атрибут href и только безопасные протоколы
	p.AllowAttrs("href").OnElements("a")
	p.AllowURLSchemes("http", "https", "tg", "mailto")
	// 4. Кастомные эмодзи
	p.AllowAttrs("emoji-id").OnElements("tg-emoji")
	sanitized := p.Sanitize(protected)

	// 5. Возвращаем кастомные теги на место
	sanitized = restoreReplacer.Replace(sanitized)
	// 6. Возвращаем кавычки и апострофы обратно.
	// Bluemonday экранирует их для безопасности атрибутов, но в тексте сообщений ТГ они не мешают.
	sanitized = strings.ReplaceAll(sanitized, "&#34;", "\"")
	sanitized = strings.ReplaceAll(sanitized, "&quot;", "\"")
	sanitized = strings.ReplaceAll(sanitized, "&#39;", "'")

	return normalizeNewLines(sanitized)
}

func SplitTextForTelegramHtmlMarkdown(input string, chunkSize int) []string {
	return SmartSplitTextIntoChunks(SanitizeForTelegram(input), chunkSize)
}

var (
	reTrailingSpaces      = regexp.MustCompile(`(?m)[ \t]+$`)
	reMoreThanTwoNewlines = regexp.MustCompile(`\n{3,}`)
	reBr                  = regexp.MustCompile(`(?i)<br\s*/?>`)
)

func normalizeNewLines(s string) string {
	// 1. Унифицируем переносы: Windows (\r\n) -> Unix (\n)
	result := strings.ReplaceAll(s, "\r\n", "\n")
	// 2. Убираем пробелы и табы в концах строк (чтобы "\n  \n" стало "\n\n")
	// Используем флаг (?m), чтобы $ срабатывал на конец каждой строки
	result = reTrailingSpaces.ReplaceAllString(result, "")
	// 3. Схлопываем 3 и более \n в ровно 2 \n
	result = reMoreThanTwoNewlines.ReplaceAllString(result, "\n\n")
	// 4. Убираем лишнее по краям всего текста
	return strings.TrimSpace(result)
}
