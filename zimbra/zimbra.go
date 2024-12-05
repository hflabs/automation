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

func (z *zimbra) GetMail(messageId string) (Mail, error) {
	var messages Mails
	err := requests.
		URL(fmt.Sprintf("%s/home/%s/", z.url, z.account)).
		Param("id", messageId).
		Param("fmt", "json").
		BasicAuth(z.account, z.password).
		ToJSON(&messages).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return Mail{}, err
	}
	if len(messages.Mails) == 0 {
		return Mail{}, fmt.Errorf("no messages found by id %s", messageId)
	}
	return messages.Mails[0], nil
}

func (z *zimbra) GetMails(folder string, sorting SortType, limit, offset int) (Mails, error) {
	var messages Mails
	err := requests.
		URL(fmt.Sprintf("%s/home/%s/%s", z.url, z.account, folder)).
		Param("limit", strconv.Itoa(limit)).
		Param("offset", strconv.Itoa(offset)).
		Param("sortBy", string(sorting)).
		Param("fmt", "json").
		BasicAuth(z.account, z.password).
		ToJSON(&messages).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return Mails{}, err
	}
	return messages, nil
}

func (z *zimbra) GetMailsByTopicId(topicId string) (Mails, error) {
	var messages Mails
	err := requests.
		URL(fmt.Sprintf("%s/home/%s/", z.url, z.account)).
		Param("query", fmt.Sprintf("cid:%s", topicId)).
		Param("fmt", "json").
		BasicAuth(z.account, z.password).
		ToJSON(&messages).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return Mails{}, err
	}
	return messages, nil
}

func (z *zimbra) SearchMails(query string) (Mails, error) {
	var messages Mails
	err := requests.
		URL(fmt.Sprintf("%s/home/%s/", z.url, z.account)).
		Param("query", query).
		Param("fmt", "json").
		BasicAuth(z.account, z.password).
		ToJSON(&messages).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return Mails{}, err
	}
	return messages, nil
}
