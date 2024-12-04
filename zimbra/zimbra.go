package zimbra

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"strconv"
)

type zimbra struct {
	url      string
	account  string
	password string
}

func NewApiZimbra(url, account, password string) ApiZimbra {
	return &zimbra{url: url, account: account, password: password}
}

func (z *zimbra) GetMail(messageId string) (Message, error) {
	var messages Messages
	err := requests.
		URL(fmt.Sprintf("%s/%s/", z.url, z.account)).
		Param("id", messageId).
		Param("fmt", "json").
		BasicAuth(z.account, z.password).
		ToJSON(&messages).
		Fetch(context.Background())
	if err != nil {
		return Message{}, err
	}
	if len(messages.Messages) == 0 {
		return Message{}, fmt.Errorf("no messages found by id %s", messageId)
	}
	return messages.Messages[0], nil
}

func (z *zimbra) GetMails(folder string, sorting SortType, limit, offset int) (Messages, error) {
	var messages Messages
	err := requests.
		URL(fmt.Sprintf("%s/%s/%s", z.url, z.account, folder)).
		Param("limit", strconv.Itoa(limit)).
		Param("offset", strconv.Itoa(offset)).
		Param("sortBy", string(sorting)).
		Param("fmt", "json").
		BasicAuth(z.account, z.password).
		ToJSON(&messages).
		Fetch(context.Background())
	if err != nil {
		return Messages{}, err
	}
	return messages, nil
}
