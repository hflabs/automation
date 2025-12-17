package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/mymmrac/telego"
)

const (
	FileEndpoint      = "https://api.telegram.org/file/bot%s/%s"
	MsgEntityLinkType = "text_link"
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

func IsMessage(update telego.Update) bool {
	if update.Message == nil {
		return false
	}
	if IsTextOrCaption(update) {
		return true
	}
	return false
}

func IsPhoto(update telego.Update) bool {
	if update.Message.Photo != nil {
		return true
	}
	return false
}

func IsDocument(update telego.Update) bool {
	if update.Message.Document != nil {
		return true
	}
	return false
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

func SmartSplitTextIntoChunks(s string, chunkSize int) []string {
	var chunks []string
	if len(s) < chunkSize {
		return []string{s}
	}
	rows := strings.Split(s, "\n")
	var extraRows string
	var chunk string
	for _, row := range rows {
		originalRow := row + "\n"
		// Добавляем доп строки от прошлой части
		if extraRows != "" {
			originalRow = extraRows + originalRow
			extraRows = ""
		}
		if len(chunk)+len(originalRow) < chunkSize {
			chunk = chunk + originalRow
			continue
		}
		// Проверяем у текущей части базовую разметку HTML, если есть проблемы пытаемся исправить их переносом строки в следующую часть
		if !CheckBasicHTML(chunk) {
			chunk, extraRows = fixHtmlChunk(chunk)
		}
		chunks = append(chunks, chunk)
		chunk = originalRow
	}
	if chunk != "" {
		chunks = append(chunks, chunk)
	}
	return chunks
}

// Пытаемся исправить часть, отрезая от её конца по 1 строке, каждый раз проверяя стала ли часть проходить проверку по HTML-тэгам
func fixHtmlChunk(sourceChunk string) (chunk string, extraRows string) {
	rows := strings.Split(sourceChunk, "\n")
	if len(rows) > 0 {
		extraRows += rows[len(rows)-1]
		rows = rows[:len(rows)-1]
		chunk = strings.Join(rows, "\n")
		if CheckBasicHTML(chunk) {
			return chunk, extraRows
		}
	}
	return chunk, extraRows
}

func CheckBasicHTML(html string) bool {
	openTags := []string{}
	for _, tag := range strings.Split(html, "<") {
		if strings.HasPrefix(tag, "/") {
			index := strings.Index(tag, ">")
			if index == -1 {
				continue
			}
			tagName := strings.TrimSpace(tag[1:index])
			if len(openTags) > 0 && openTags[len(openTags)-1] == tagName {
				openTags = openTags[:len(openTags)-1]
			} else {
				return false
			}
		} else if strings.Contains(tag, " href='http") {
			openTags = append(openTags, "a")
		} else if strings.Contains(tag, ">") {
			index := strings.Index(tag, ">")
			if index == -1 {
				continue
			}
			tagName := strings.TrimSpace(tag[:index])
			openTags = append(openTags, tagName)
		}
	}
	return len(openTags) == 0
}
