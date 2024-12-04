package zimbra

type ApiZimbra interface {
	GetMails(folder string, sorting SortType, limit, offset int) (Mails, error)
	GetMailsByTopicId(topicId string) (Mails, error)
	GetMail(messageId string) (Mail, error)
	SearchMails(query string) (Mails, error)
}
