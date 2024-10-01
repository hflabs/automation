package jira

func GetJiraSuggestions() Fields {
	return Fields{
		Issue: Issue{
			Status: Status{
				Open:                        "1",
				ToDo:                        "2",
				InProgress:                  "3",
				Reopened:                    "4",
				Resolved:                    "5",
				Closed:                      "6",
				TestingInProgress:           "1006",
				CRInProgress:                "10008",
				Delay:                       "10009",
				New:                         "10011",
				Assigned:                    "10012",
				InProgressHRP:               "10013",
				ResolveSup:                  "10014",
				Confirmed:                   "10015",
				Saved:                       "10016",
				RequestInfo:                 "10018",
				Escalation:                  "10019",
				Done:                        "10020",
				Out:                         "10022",
				AgreedHRP:                   "10120",
				Cancelled:                   "10121",
				OnApproval:                  "10220",
				CRCompleted:                 "10320",
				Backlog:                     "10424",
				Selected:                    "10520",
				OfferDone:                   "10425",
				Rated:                       "10620",
				CodeReview:                  "10720",
				Deployed:                    "10820",
				SelectedForDevelopment:      "10920",
				TestReview:                  "11020",
				DesignReview:                "11021",
				AllOrganizedHRP:             "11320",
				Design:                      "11420",
				ReadyForDevelopment:         "11421",
				Development:                 "11422",
				Review:                      "11423",
				Analytics:                   "11520",
				Requirements:                "11521",
				ReadyForTesting:             "11522",
				Testing:                     "11523",
				Documentation:               "11620",
				FirstAnswer:                 "11720",
				NoNeedReaction:              "11820",
				AssignedInQueue:             "11821",
				AnsweredButNeedImprovements: "11822",
				AwaitingCustomerResponse:    "11823",
				AwaitingDecisionColleagues:  "11824",
			},
			Type: Type{
				Bug:         "1",
				NewFeature:  "2",
				Task:        "3",
				Improvement: "4",
				Epic:        "5",
			},
			Priority: Priority{
				Blocker: "2",
				Minor:   "4",
			},
			Resolution: Resolution{
				Fixed:   "1",
				WontFix: "2",
				Unknown: "100300",
			},
			Fields: IssueFields{
				StoryPoints:   "customfield_10083",
				BusinessValue: "customfield_10084",
				WeightedJob:   "customfield_12580",
			},
			Transitions: Transitions{
				INNA: ProjectTransition{
					FromBacklogToRate:        "51",
					FromBacklogToDone:        "171",
					FromRateToSelected:       "61",
					FromRateToBacklog:        "161",
					FromRateToDone:           "141",
					FromSelectedToInProgress: "71",
					FromSelectedToBacklog:    "131",
					FromSelectedToDone:       "141",
					FromInProgressToResolved: "81",
					FromInProgressToBacklog:  "131",
					FromInProgressToDelayed:  "111",
					FromResolvedToDone:       "91",
					FromResolvedToSelected:   "101",
					FromResolvedToDelayed:    "111",
					FromDelayedToInProgress:  "71",
				},
			},
		},
		Changelog: Changelog{
			SingleItem: SingleItem{
				Field: Field{
					StoryPoints:   "Story Points",
					BusinessValue: "Business Value",
					WeightedJob:   "Weighted Job",
				},
			},
		},
		EventType: EventType{
			Updated:  "issue_updated",
			Generic:  "issue_generic",
			Assigned: "issue_assigned",
		},
	}
}

type Fields struct {
	Issue     Issue     // Задача
	Changelog Changelog // Журнал изменения
	EventType EventType // Тип события
}
type Issue struct {
	Status      Status      // Статус
	Type        Type        // Тип
	Priority    Priority    // Приоритет
	Resolution  Resolution  // Результат
	Transitions Transitions // Возможные переходы статусов
	Fields      IssueFields // Идентификаторы полей
}
type Status struct {
	Open                        string // Открытый
	ToDo                        string // К выполнению
	InProgress                  string // В работе
	Reopened                    string // Переоткрыт
	Resolved                    string // Решенные
	Closed                      string // Закрыт
	TestingInProgress           string // Testing In Progress
	CRInProgress                string // CR In Progress
	Delay                       string // Отложен
	New                         string // Новый
	Assigned                    string // Назначен
	InProgressHRP               string // В работе
	ResolveSup                  string // Решен
	Confirmed                   string // Подтвержден
	Saved                       string // Сохранен
	RequestInfo                 string // Запрос информации
	Escalation                  string // Эскалация
	Done                        string // Готово
	Out                         string // Выход
	AgreedHRP                   string // Согласовано
	Cancelled                   string // Отменен
	OnApproval                  string // На согласовании
	CRCompleted                 string // CR завершен
	Backlog                     string // Очередь
	Selected                    string // Выбрано
	OfferDone                   string // Сделан оффер
	Rated                       string // Оценено
	CodeReview                  string // Ревью кода
	Deployed                    string // Deployed
	SelectedForDevelopment      string // Selected for Development
	TestReview                  string // Ревью тестирования
	DesignReview                string // Ревью дизайна
	AllOrganizedHRP             string // Всё организовано
	Design                      string // Дизайн
	ReadyForDevelopment         string // Готово к разработке
	Development                 string // Разработка
	Review                      string // Ревью
	Analytics                   string // Аналитика
	Requirements                string // Ревью требований
	ReadyForTesting             string // Готово к тестированию
	Testing                     string // Тестирование
	Documentation               string // Документация
	FirstAnswer                 string // Дан первичный ответ
	NoNeedReaction              string // Не требует реакции
	AssignedInQueue             string // Назначен (в очереди)
	AnsweredButNeedImprovements string // Отвечено, но остались доделки
	AwaitingCustomerResponse    string // Ожидает ответа заказчика
	AwaitingDecisionColleagues  string // Ожидает решения от коллег
}

type EventType struct {
	Updated  string // Обновлено
	Generic  string // Создано
	Assigned string // Назначено
}

type Type struct {
	Bug         string // Ошибка
	NewFeature  string // Новая функциональность
	Task        string // Задача
	Improvement string // Улучшение
	Epic        string // Эпик
}

type Priority struct {
	Blocker string // Блокирующий
	Minor   string // Незначительный
}

type Resolution struct {
	Fixed   string // Исправленный
	WontFix string // Не будет исправлено
	Unknown string // Неизвестно (нужно найти и заполнить)
}

type Changelog struct {
	SingleItem SingleItem // Что изменилось
}

type SingleItem struct {
	Field Field // Поле
}

type Field struct {
	StoryPoints   string // Story Points
	BusinessValue string // Business Value
	WeightedJob   string // Weighted Job
}

type Transitions struct {
	INNA ProjectTransition
}

type IssueFields struct {
	StoryPoints   string // customfield_10083
	BusinessValue string // customfield_10084
	WeightedJob   string // customfield_12580
}

type ProjectTransition struct {
	FromBacklogToRate string // Rate - Анализ задачи на исполняемость, атомарность и нужность. Оценка для получения индекса, по которому сортируется очередь задач
	FromBacklogToDone string // Trash - Закрытие задачи без разработки

	FromRateToSelected string // Plan - Планирование задач в спринт. Будет использоваться, чтобы понимать, успели запланированное в спринт или нет.
	FromRateToBacklog  string // Revert
	FromRateToDone     string // Close - Закрытие задачи без разработки

	FromSelectedToInProgress string // Start - Начало разработки
	FromSelectedToBacklog    string // Revert - Возврат задачи в бэклог
	FromSelectedToDone       string // Close - Закрытие задачи без разработки

	FromInProgressToResolved string // Finish - Окончание разработки
	FromInProgressToBacklog  string // Revert - Возврат задачи в бэклог
	FromInProgressToDelayed  string // Delay - Работа над задачей заблокирована или отложена из-за внешних факторов

	FromResolvedToDone     string // Verify and close - Успешное тестирование
	FromResolvedToSelected string // Reopen - Возврат задачи в разработку, т.к. есть ошибки
	FromResolvedToDelayed  string // Delay - Работа над задачей заблокирована или отложена из-за внешних факторов

	FromDelayedToInProgress string // Start - Начало разработки
}
