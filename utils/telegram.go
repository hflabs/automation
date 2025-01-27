package utils

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/mymmrac/telego"
	"strings"
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
