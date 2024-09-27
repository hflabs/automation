package confluence

type ApiConfluence interface {
	GetContentById(id string) (string, error)
	GetVersionInfoById(id string) (VersionResponse, error)

	UpdatePageById(id string, content string, reCreate bool) error
	UpdatePageByIdWithCheck(id string, content string, reCreate bool) error
}
