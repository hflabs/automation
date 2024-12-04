package zimbra

import (
	"encoding/json"
	"time"
)

type Mails struct {
	Mails []Mail `json:"m"`
	More  bool   `json:"more"`
}

type Mail struct {
	Id          string           `json:"id"`
	TopicId     string           `json:"cid"`
	Date        time.Time        `json:"d"`
	Flags       string           `json:"f"`
	SizeBytes   int32            `json:"s"`
	Addressees  []Addressee      `json:"a"`
	Title       string           `json:"su"`
	Fragment    string           `json:"fr"`
	Folder      string           `json:"folder"`
	Tags        []string         `json:"tags"`
	Attachments []AttachmentFile `json:"attachments"`
	Content     Content          `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type Addressee struct {
	Email string        `json:"a"`
	Type  TypeAddressee `json:"t"`
}

type AttachmentFile struct {
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
}

type SortType string

const (
	// DateAsc — сортировка по дате отправки письма в порядке возрастания (старые письма первыми).
	DateAsc SortType = "dateAsc"
	// DateDesc — сортировка по дате отправки письма в порядке убывания (новые письма первыми).
	DateDesc SortType = "dateDesc"
	// SubjAsc — сортировка по теме письма в алфавитном порядке.
	SubjAsc SortType = "subjAsc"
	// SubjDesc — сортировка по теме письма в обратном алфавитном порядке.
	SubjDesc SortType = "subjDesc"
	// SizeAsc — сортировка по размеру письма в порядке возрастания (меньшие письма первыми).
	SizeAsc SortType = "sizeAsc"
	// SizeDesc — сортировка по размеру письма в порядке убывания (большие письма первыми).
	SizeDesc SortType = "sizeDesc"
	// NameAsc — сортировка по имени отправителя или получателя в алфавитном порядке.
	NameAsc SortType = "nameAsc"
	// NameDesc — сортировка по имени отправителя или получателя в обратном алфавитном порядке.
	NameDesc SortType = "nameDesc"
	// FlagAsc — сортировка по статусу флагов письма (например, непрочитанные сначала).
	FlagAsc SortType = "flagAsc"
	// FlagDesc — сортировка по статусу флагов письма в обратном порядке.
	FlagDesc SortType = "flagDesc"
	// PriorityAsc — сортировка по приоритету письма (низкий приоритет первым).
	PriorityAsc SortType = "priorityAsc"
	// PriorityDesc — сортировка по приоритету письма (высокий приоритет первым).
	PriorityDesc SortType = "priorityDesc"
)

type TypeAddressee string

const (
	// From — адрес отправителя
	From TypeAddressee = "f"
	// To — основной получатель письма
	To TypeAddressee = "t"
	// Copy — получатель письма в копии
	Copy TypeAddressee = "c"
	// HideCopy — получатель письма в скрытой копии
	HideCopy TypeAddressee = "b"
)

type Flag string

const (
	// Unread — письмо не прочитано
	Unread Flag = "u"
	// Replied — на письмо был отправлен ответ
	Replied Flag = "r"
	// Flagged — письмо помечено (например, как важное)
	Flagged Flag = "f"
	// SentByMe — письмо отправлено текущим пользователем
	SentByMe Flag = "s"
	// Draft — письмо является черновиком
	Draft Flag = "d"
	// Attachment — письмо содержит вложения
	Attachment Flag = "a"
	// Deleted — письмо удалено
	Deleted Flag = "x"
	// Notification — письмо является уведомлением
	Notification Flag = "n"
	// Invite — письмо содержит приглашение на событие
	Invite Flag = "i"
)

func (m *Mail) UnmarshalJSON(data []byte) error {
	type Alias Mail
	aux := &struct {
		D int64 `json:"d"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	m.Date = time.UnixMilli(aux.D)
	return nil
}
