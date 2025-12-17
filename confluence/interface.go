package confluence

import "context"

type ApiConfluence interface {
	GetContentById(ctx context.Context, id string) (string, error)
	GetVersionById(ctx context.Context, id string) (VersionResponse, error)
	GetPagesByName(ctx context.Context, name, spaceKey string) ([]PageInfo, error)
	GetPagesByIncludedName(ctx context.Context, name, spaceKey string) ([]PageInfo, error)

	GetChildrenById(ctx context.Context, id string, limit int) ([]PageInfo, error)
	GetChildrenByIdRecursive(ctx context.Context, id string, limit int) ([]PageInfo, error)

	CreatePage(ctx context.Context, name, spaceKey, content string, parentPageId string) (string, error)
	CreatePageWithHash(ctx context.Context, name, spaceKey, content, parentPageId string) (string, error)

	AddLabelById(ctx context.Context, id, label string) error
	GetLabelsById(ctx context.Context, id string) ([]string, error)

	UpdatePageById(ctx context.Context, id string, content string, reCreate bool) error
	UpdatePageByIdWithCheck(ctx context.Context, id string, content string, reCreate bool) error
	UpdatePageParentById(ctx context.Context, id, parentId string) error

	SetRestrictionUser(ctx context.Context, id, username, action string) error
	SetRestrictionGroup(ctx context.Context, id, groupName, action string) error
	SetRestrictionsForHFLabsOnly(ctx context.Context, id string) error
}
