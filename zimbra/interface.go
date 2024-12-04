package zimbra

type ApiZimbra interface {
	GetMails(folder string, sorting SortType, limit, offset int) (Messages, error)
	GetMail(messageId string) (Message, error)
}
