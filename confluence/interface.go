package confluence

type ApiConfluence interface {
	GetContentById(id int) (string, error)
	GetVersionInfoById(id int) (VersionResponse, error)
	GetPagesByName(name, spaceKey string) ([]PageInfo, error)

	CreatePage(name, spaceKey, content string, parentPageId int) (int, error)

	UpdatePageById(id int, content string, reCreate bool) error
	UpdatePageByIdWithCheck(id int, content string, reCreate bool) error
}
