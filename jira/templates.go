package jira

var Issue = newIssue()
var EventType = newEventTypes()
var Changelog = newChangelogs()

type Issues struct {
	Status      status      // Статус
	Type        typeIssue   // Тип
	Priority    priority    // Приоритет
	Resolution  resolution  // Результат
	Transitions transitions // Возможные переходы статусов
	Fields      fieldsIssue // Идентификаторы полей
}

func newIssue() Issues {
	return Issues{
		Status:      newStatuses(),
		Type:        newTypes(),
		Priority:    newPriorities(),
		Resolution:  newResolutions(),
		Fields:      newIssueFields(),
		Transitions: newTransitions(),
	}
}

type status struct {
	Open                        string // Открытый
	ToDo                        string // К выполнению
	InProgress                  string // В работе
	Reopened                    string // Переоткрыт
	Resolved                    string // Решенные
	ToRelease                   string // To Release
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

func newStatuses() status {
	return status{
		Open:                        "1",
		ToDo:                        "2",
		InProgress:                  "3",
		Reopened:                    "4",
		Resolved:                    "5",
		ToRelease:                   "11920",
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
	}
}

type EventTypes struct {
	Created  string // Создано
	Updated  string // Обновлено
	Generic  string // Создано
	Assigned string // Назначено
	Closed   string // Закрыто
}

func newEventTypes() EventTypes {
	return EventTypes{
		Created:  "issue_created",
		Updated:  "issue_updated",
		Generic:  "issue_generic",
		Assigned: "issue_assigned",
		Closed:   "issue_closed",
	}
}

type typeIssue struct {
	Bug         string // Ошибка
	NewFeature  string // Новая функциональность
	Task        string // Задача
	Improvement string // Улучшение
	Epic        string // Эпик
	Learning    string // Обучение
}

func newTypes() typeIssue {
	return typeIssue{
		Bug:         "1",
		NewFeature:  "2",
		Task:        "3",
		Improvement: "4",
		Epic:        "5",
		Learning:    "10600",
	}
}

type priority struct {
	Blocker  string // Блокирующий
	Critical string // Критический
	Minor    string // Незначительный
}

func newPriorities() priority {
	return priority{
		Blocker:  "1",
		Critical: "2",
		Minor:    "4",
	}
}

type resolution struct {
	Fixed      string // Исправленный
	WontFix    string // Не будет исправлено
	NeedBugFix string // Нужен багфикс
}

func newResolutions() resolution {
	return resolution{
		Fixed:      "1",
		WontFix:    "2",
		NeedBugFix: "10300",
	}
}

type Changelogs struct {
	SingleItem singleItem // Что изменилось
}

func newChangelogs() Changelogs {
	return Changelogs{SingleItem: newSingleItem()}
}

type singleItem struct {
	Field changelogField // Поле
}

func newSingleItem() singleItem {
	return singleItem{Field: newChangelogFields()}
}

type changelogField struct {
	Status        string // Статус
	StoryPoints   string // Story Points
	BusinessValue string // Business Value
	WeightedJob   string // Weighted Job
	Assignee      string // Исполнитель
	ReleaseNotes  string // Ссылка на RN
	Components    string // Компоненты
	DueDate       string // Срок исполнения
	Participants  string // Участники

	LearnTime          string // Месяц и год обучения
	LearnForWho        string // Для кого еще подходит это обучение?
	LearnWhatLike      string // Что тебе понравилось в обучении?
	LearnWhatUseful    string // Что было полезного/ценного в обучении для тебя?
	LearnWhatBad       string // Было ли что-то, что не понравилось на обучении? Расскажи
	LearnWhatLearned   string // Чему научился(ась)? Что планируешь применять в работе?
	LearnWillRecommend string // Порекомендуешь ли обучение коллегам?
	LearnPeople        string // Участники обучения
	LearnField         string // Сфера обучения
	LearnLink          string // Ссылка на курс
}

func newChangelogFields() changelogField {
	return changelogField{
		Status:        "status",
		StoryPoints:   "Story Points",
		BusinessValue: "Business Value",
		WeightedJob:   "Weighted Job",
		Assignee:      "assignee",
		ReleaseNotes:  "Release Notes",
		Components:    "components",
		DueDate:       "Срок исполнения",
		Participants:  "Участники",

		LearnTime:          "Месяц и год обучения",
		LearnForWho:        "Для кого еще подходит это обучение?",
		LearnWhatLike:      "Что тебе понравилось в обучении?",
		LearnWhatUseful:    "Что было полезного/ценного в обучении для тебя?",
		LearnWhatBad:       "Было ли что-то, что не понравилось на обучении? Расскажи",
		LearnWhatLearned:   "Чему научился(ась)? Что планируешь применять в работе?",
		LearnWillRecommend: "Порекомендуешь ли обучение коллегам?",
		LearnPeople:        "Участники обучения",
		LearnField:         "Сфера обучения",
		LearnLink:          "Ссылка на курс",
	}
}

type transitions struct {
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
	FromResolvedToRelease  string // Verify and Release — Успешное тестирование

	FromDelayedToInProgress string // Start - Начало разработки

	FromCodeReviewToResolved string // Resolved — Резолюция выбрана и ожидает подтверждения автора. Здесь проблемы либо открываются вновь, либо закрываются.
}

func newTransitions() transitions {
	return transitions{
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
		FromResolvedToRelease:    "191",
		FromResolvedToSelected:   "101",
		FromResolvedToDelayed:    "111",
		FromDelayedToInProgress:  "71",
		FromCodeReviewToResolved: "5",
	}
}

type fieldsIssue struct {
	StoryPoints        string // customfield_10083
	BusinessValue      string // customfield_10084
	WeightedJob        string // customfield_12580
	ReleaseNotes       string // customfield_13082
	ReleaseInstruction string // customfield_13081
	Components         string // components
	DueDate            string // duedate
	Participants       string // customfield_10380

	LearnTime          string // Месяц и год обучения
	LearnForWho        string // Для кого еще подходит это обучение?
	LearnWhatLike      string // Что тебе понравилось в обучении?
	LearnWhatUseful    string // Что было полезного/ценного в обучении для тебя?
	LearnWhatBad       string // Было ли что-то, что не понравилось на обучении? Расскажи
	LearnWhatLearned   string // Чему научился(ась)? Что планируешь применять в работе?
	LearnWillRecommend string // Порекомендуешь ли обучение коллегам?
	LearnPeople        string // Участники обучения
	LearnField         string // Сфера обучения
	LearnLink          string // Ссылка на курс
}

func newIssueFields() fieldsIssue {
	return fieldsIssue{
		StoryPoints:        "customfield_10083",
		BusinessValue:      "customfield_10084",
		WeightedJob:        "customfield_12580",
		ReleaseNotes:       "customfield_13082",
		ReleaseInstruction: "customfield_13081",
		Components:         "components",
		DueDate:            "duedate",
		Participants:       "customfield_10380",

		LearnTime:          "customfield_14481",
		LearnForWho:        "customfield_13881",
		LearnWhatLike:      "customfield_14483",
		LearnWhatUseful:    "customfield_13784",
		LearnWhatBad:       "customfield_14484",
		LearnWhatLearned:   "customfield_13880",
		LearnWillRecommend: "customfield_13882",
		LearnPeople:        "customfield_14480",
		LearnField:         "customfield_14380",
		LearnLink:          "customfield_13782",
	}
}
