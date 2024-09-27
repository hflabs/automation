package jira

const (
	IssueEventTypeUpdated  = "issue_updated"
	IssueEventTypeGeneric  = "issue_generic"
	IssueEventTypeAssigned = "issue_assigned"

	IssueStatusOpen   = "Открытый"
	IssueStatusOpenId = "1"

	IssueStatusToDo   = "К выполнению"
	IssueStatusToDoId = "2"

	IssueStatusInProgress   = "В работе"
	IssueStatusInProgressId = "3"

	IssueStatusReopened   = "Переоткрыт"
	IssueStatusReopenedId = "4"

	IssueStatusResolved   = "Решенные"
	IssueStatusResolvedId = "5"

	IssueStatusClosed   = "Закрыт"
	IssueStatusClosedId = "6"

	IssueStatusTestingInProgress   = "Testing In Progress"
	IssueStatusTestingInProgressId = "1006"

	IssueStatusCRInProgress   = "CR In Progress"
	IssueStatusCRInProgressId = "10008"

	IssueStatusDelay   = "Отложен"
	IssueStatusDelayId = "10009"

	IssueStatusNew   = "Новый"
	IssueStatusNewId = "10011"

	IssueStatusAssigned   = "Назначен"
	IssueStatusAssignedId = "10012"

	IssueHRPInProgress   = "В работе"
	IssueHRPInProgressId = "10013"

	IssueStatusResolveSup   = "Решен"
	IssueStatusResolveSupId = "10014"

	IssueStatusConfirmed   = "Подтвержден"
	IssueStatusConfirmedId = "10015"

	IssueStatusSaved   = "Сохранен"
	IssueStatusSavedId = "10016"

	IssueStatusRequestInfo   = "Запрос информации"
	IssueStatusRequestInfoId = "10018"

	IssueStatusEscalation   = "Эскалация"
	IssueStatusEscalationId = "10019"

	IssueStatusDone   = "Готово"
	IssueStatusDoneId = "10020"

	IssueStatusOut   = "Выход"
	IssueStatusOutId = "10022"

	IssueStatusHRPAgreed   = "Согласовано"
	IssueStatusHRPAgreedId = "10120"

	IssueStatusCancelled   = "Отменен"
	IssueStatusCancelledId = "10121"

	IssueStatusOnApproval   = "На согласовании"
	IssueStatusOnApprovalId = "10220"

	IssueStatusCRCompleted   = "CR завершен"
	IssueStatusCRCompletedId = "10320"

	IssueStatusBacklog   = "Очередь"
	IssueStatusBacklogId = "10424"

	IssueStatusSelectedRu = "Выбрано"
	IssueStatusSelectedId = "10520"

	IssueStatusOfferDone   = "Сделан оффер"
	IssueStatusOfferDoneId = "10425"

	IssueStatusRated   = "Оценено"
	IssueStatusRatedId = "10620"

	IssueStatusCodeReview   = "Ревью кода"
	IssueStatusCodeReviewId = "10720"

	IssueStatusDeployed   = "Deployed"
	IssueStatusDeployedId = "10820"

	IssueStatusSelectedForDevelopment   = "Selected for Development"
	IssueStatusSelectedForDevelopmentId = "10920"

	IssueStatusTestReview   = "Ревью тестирования"
	IssueStatusTestReviewId = "11020"

	IssueStatusDesignReview   = "Ревью дизайна"
	IssueStatusDesignReviewId = "11021"

	IssueStatusHRPAllOrganizedRu = "Всё организовано"
	IssueStatusHRPAllOrganizedId = "11320"

	IssueStatusDesign   = "Дизайн"
	IssueStatusDesignId = "11420"

	IssueStatusReadyForDev   = "Готово к разработке"
	IssueStatusReadyForDevId = "11421"

	IssueStatusDevelopment   = "Разработка"
	IssueStatusDevelopmentId = "11422"

	IssueStatusReview   = "Ревью"
	IssueStatusReviewId = "11423"

	IssueStatusAnalytics   = "Аналитика"
	IssueStatusAnalyticsId = "11520"

	IssueStatusReviewRequirements   = "Ревью требований"
	IssueStatusReviewRequirementsId = "11521"

	IssueStatusReadyForTesting   = "Готово к тестированию"
	IssueStatusReadyForTestingId = "11522"

	IssueStatusTesting   = "Тестирование"
	IssueStatusTestingId = "11523"

	IssueStatusDocumentation   = "Документация"
	IssueStatusDocumentationId = "11620"

	IssueStatusFirstAnswer   = "Дан первичный ответ"
	IssueStatusFirstAnswerId = "11720"

	IssueStatusNoNeedReaction   = "Дан первичный ответ"
	IssueStatusNoNeedReactionId = "11820"

	IssueStatusAssignedInQueue   = "Назначен (в очереди)"
	IssueStatusAssignedInQueueId = "11821"

	IssueStatusAnsweredButNeedImprovements   = "Отвечено, но остались доделки"
	IssueStatusAnsweredButNeedImprovementsId = "11822"

	IssueStatusAwaitingCustomerResponse   = "Ожидает ответа заказчика"
	IssueStatusAwaitingCustomerResponseId = "11823"

	IssueStatusAwaitingDecisionColleagues   = "Ожидает решения от коллег"
	IssueStatusAwaitingDecisionColleaguesId = "11824"

	IssueTypeBugRu = "Ошибка"
	IssueTypeBugId = "1"

	IssueTypeNewFeature   = "Новая функциональность"
	IssueTypeNewFeatureId = "2"

	IssueTypeTask   = "Задача"
	IssueTypeTaskId = "3"

	IssueTypeImprovement   = "Улучшение"
	IssueTypeImprovementId = "4"

	IssueTypeEpic   = "Эпик"
	IssueTypeEpicId = "5"

	IssuePriorityBlocker   = "Блокирующий"
	IssuePriorityBlockerId = "2"

	IssuePriorityMinor   = "Незначительный"
	IssuePriorityMinorId = "4"

	IssueResolutionFixed   = "Исправленный"
	IssueResolutionFixedId = "1"

	IssueResolutionWontFix   = "Не будет исправлено"
	IssueResolutionWontFixId = "2"

	IssueResolutionUnknownId = "100300"

	IssueItemFieldStoryPoints   = "Story Points"
	IssueItemFieldBusinessValue = "Business Value"
)
