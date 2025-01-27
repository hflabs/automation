package utils

import (
	"github.com/mymmrac/telego"
	"github.com/stretchr/testify/require"
	"testing"
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
