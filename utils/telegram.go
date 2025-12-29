package utils

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/carlmjohnson/requests"
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
	s = sanitizeIncomingHTML(s)
	z := html.NewTokenizer(strings.NewReader(s))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		token := z.Token()
		tokenString := token.String()
		// Обработка текстовых и обычных токенов
		// Мы используем цикл, потому что один длинный текстовый токен может
		// растянуться на несколько чанков.
		for len(tokenString) > 0 {
			// 1. Считаем, сколько места займут закрывающие теги, если прерваться сейчас
			closingOverhead := countClosingOverhead(stack)
			// 2. Считаем свободное место в текущем чанке
			availableSpace := chunkSize - currentChunk.Len() - closingOverhead
			// Если места совсем нет (или оно отрицательное из-за оверхеда), форсируем сброс чанка
			// Но только если чанк не пустой (чтобы избежать бесконечного цикла, если оверхед > chunkSize)
			if availableSpace <= 0 && currentChunk.Len() > 0 {
				chunks = append(chunks, flushChunk(&currentChunk, stack))
				currentChunk.WriteString(reopenTags(stack))
				continue
			}
			// 3. Проверяем, влезает ли токен целиком
			if len(tokenString) <= availableSpace {
				currentChunk.WriteString(tokenString)
				// Если это был тег, обновляем стек
				updateStack(tt, token, tokenString, &stack)
				break // Токен полностью обработан, выходим из внутреннего цикла
			}

			// 4. Если токен не влезает
			if tt == html.TextToken {
				// Пытаемся отрезать кусок текста по пробелу
				head, tail := splitTextBySpace(tokenString, availableSpace)
				// Если head пустой, значит даже одно слово не влезает в оставшееся место.
				// Сбрасываем текущий чанк и пробуем снова в новом (пустом) чанке.
				if head == "" && currentChunk.Len() > 0 {
					chunks = append(chunks, flushChunk(&currentChunk, stack))
					currentChunk.WriteString(reopenTags(stack))
					continue
				}
				// Добавляем кусок, закрываем чанк
				currentChunk.WriteString(head)
				chunks = append(chunks, flushChunk(&currentChunk, stack))
				// Начинаем новый чанк
				currentChunk.WriteString(reopenTags(stack))
				// Остаток текста обрабатываем в следующей итерации цикла
				tokenString = tail
			} else {
				// Если это НЕ текст (а тег), и он не влезает -> закрываем текущий чанк
				chunks = append(chunks, flushChunk(&currentChunk, stack))
				currentChunk.WriteString(reopenTags(stack))
				// Тег будет добавлен в новый чанк на следующей итерации (availableSpace будет большим)
			}
		}
	}
	// Сохраняем последний кусочек, если есть
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}
	return chunks
}

// splitTextBySpace делит текст на две части:
// head - часть, которая влезает в limit (стараясь не резать слова)
// tail - остаток
func splitTextBySpace(text string, limit int) (head, tail string) {
	if len(text) <= limit {
		return text, ""
	}
	// Берем подстроку, которая теоретически влезает
	candidate := text[:limit]
	// Ищем последний пробел в этой подстроке
	lastSpaceIndex := strings.LastIndexFunc(candidate, unicode.IsSpace)
	// Если пробел найден и он не в самом начале (чтобы не возвращать пустой head постоянно)
	if lastSpaceIndex > 0 {
		// Включаем пробел в первую часть, чтобы форматирование сохранялось
		return text[:lastSpaceIndex+1], text[lastSpaceIndex+1:]
	}
	// Если пробелов нет (очень длинное слово/URL) или пробел в начале — режем жестко
	return text[:limit], text[limit:]
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
		// Добавляем в стек, только если это не void элемент (br, img и т.д.)
		if !isVoidElement(token.Data) {
			*stack = append(*stack, struct{ Name, Token string }{
				Name:  token.Data,
				Token: tokenString,
			})
		}
	case html.EndTagToken:
		// Убираем из стека соответствующий открывающий тег
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

func sanitizeIncomingHTML(s string) string {
	// Исправляем проблему с залипшим тегом ссылки
	s = strings.ReplaceAll(s, "<<a", "&lt;<a")
	return s
}
