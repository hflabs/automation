package confluence

type ApiConfluence interface {
	GetContentById(id string) (string, error)
	GetVersionInfoById(id string) (VersionResponse, error)
	GetPagesByName(name, spaceKey string) ([]PageInfo, error)
	GetPagesByIncludedName(name, spaceKey string) ([]PageInfo, error)

	CreatePage(name, spaceKey, content string, parentPageId string) (string, error)

	AddLabelToPage(pageId, label string) error
	GetLabels(pageId string) ([]string, error)

	UpdatePageById(id string, content string, reCreate bool) error
	UpdatePageByIdWithCheck(id string, content string, reCreate bool) error
	UpdatePageParentById(id, parentId string) error
}
