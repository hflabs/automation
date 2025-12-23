package confluence

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/carlmjohnson/requests"
)

func NewConfluence(baseUrl, user, password string) ApiConfluence {
	return &confluence{user: user, password: password, baseUrl: baseUrl}
}

func (c *confluence) GetPagesByName(ctx context.Context, name, spaceKey string) ([]PageInfo, error) {
	var resp searchPagesResponse
	err := requests.
		URL(c.baseUrl).
		Method(http.MethodGet).
		Param("title", name).
		Param("spaceKey", spaceKey).
		Param("type", "page").
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return resp.Results, fmt.Errorf("GetPagesByName — get confluence page by name %s in space %s err: %w", name, spaceKey, err)
	}
	return resp.Results, nil
}

func (c *confluence) GetPagesByIncludedName(ctx context.Context, name, spaceKey string) ([]PageInfo, error) {
	cqlQuery := fmt.Sprintf("space=\"%s\" AND type=\"page\" AND title~\"%s\"", spaceKey, name)
	var resp searchPagesResponse
	err := requests.
		URL(fmt.Sprintf("%s/search", c.baseUrl)).
		Method(http.MethodGet).
		Param("cql", cqlQuery).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return resp.Results, fmt.Errorf("GetPagesByIncludedName — get confluence page by include name %s in space %s err: %w", name, spaceKey, err)
	}
	return resp.Results, nil
}

func (c *confluence) GetContentById(ctx context.Context, id string) (string, error) {
	var resp PageInfo
	err := requests.
		URL(fmt.Sprintf("%s/%s", c.baseUrl, id)).
		Method(http.MethodGet).
		Param("expand", "body.storage").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return "", fmt.Errorf("GetContentById — get confluence pageId %s err: %w", id, err)
	}
	return resp.Body.Storage.Value, err
}

func (c *confluence) GetVersionById(ctx context.Context, id string) (VersionResponse, error) {
	var resp VersionResponse
	err := requests.
		URL(fmt.Sprintf("%s/%s", c.baseUrl, id)).
		Method(http.MethodGet).
		Param("expand", "version").
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)

	if err != nil {
		return resp, fmt.Errorf("GetLastVersionInfoById — get confluence pageId %s err: %w", id, err)
	}
	return resp, nil
}

func (c *confluence) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	return c.getUser(ctx, "username", username)
}

func (c *confluence) GetUserByKey(ctx context.Context, key string) (*User, error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}
	return c.getUser(ctx, "key", key)
}

func (c *confluence) getUser(ctx context.Context, paramName, paramValue string) (*User, error) {
	var u User
	// baseUrl обычно используем https://confluence.ru/rest/api/content
	baseUrl := strings.Replace(c.baseUrl, "/content", "/user", 1)
	err := requests.
		URL(baseUrl).
		Method(http.MethodGet).
		Param(paramName, paramValue).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&u).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("getUser(%s=%s) err: %w", paramName, paramValue, err)
	}
	return &u, nil
}

func (c *confluence) CreatePage(ctx context.Context, name, spaceKey, content, parentPageId string) (string, error) {
	req := PageInfo{
		Type:    "page",
		Title:   name,
		Space:   &Space{Key: spaceKey},
		Parents: []PageInfo{{Type: "page", Id: parentPageId}},
		Body:    &PageBody{Storage: PageStorage{Value: content, Representation: "storage"}},
	}
	err := requests.
		URL(c.baseUrl).
		Method(http.MethodPost).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		BodyJSON(req).
		ToJSON(&req).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return "", fmt.Errorf("CreatePage — create confluence page `%s` in space `%s` with content `%s` err: %w", name, spaceKey, content, err)
	}
	return req.Id, nil
}

func (c *confluence) CreatePageWithHash(ctx context.Context, name, spaceKey, content, parentPageId string) (string, error) {
	hash := md5.Sum([]byte(content))
	hashcode := hex.EncodeToString(hash[:])
	content = fmt.Sprintf(HideHash, hashcode) + "\n" + content
	return c.CreatePage(ctx, name, spaceKey, content, parentPageId)
}

func (c *confluence) UpdatePageById(ctx context.Context, id string, content string, reCreate bool) error {
	versionInfo, err := c.GetVersionById(ctx, id)
	if err != nil {
		return err
	}

	req := PageInfo{
		Type:    "page",
		Title:   versionInfo.Title,
		Id:      id,
		Version: &PageVersion{Number: versionInfo.Version.Number + 1},
		Body:    &PageBody{PageStorage{Value: content, Representation: "storage"}},
	}
	oldContent, err := c.GetContentById(ctx, id)
	if err != nil {
		return err
	}
	// Проверяем есть ли "линия" на странице, чтобы сохранить контент до или после неё
	if !reCreate && strings.Contains(oldContent, CheckLine) {
		partsContent := strings.Split(oldContent, CheckLine)
		switch len(partsContent) {
		case 2:
			contentToSave := partsContent[1]
			req.Body.Storage.Value = content + CheckLine + contentToSave
		case 3:
			contentToSaveStart, contentToSaveEnd := partsContent[0], partsContent[2]
			req.Body.Storage.Value = contentToSaveStart + CheckLine + content + CheckLine + contentToSaveEnd
		}
	}
	err = c.updatePage(ctx, id, req)
	if err != nil {
		return fmt.Errorf("UpdatePageById — update confluence pageId %s, content %s err: %w", id, content, err)
	}
	return nil
}

func (c *confluence) UpdatePageByIdWithCheck(ctx context.Context, id string, content string, reCreate bool) error {
	hash := md5.Sum([]byte(content))
	hashcode := hex.EncodeToString(hash[:])
	versionInfo, err := c.GetVersionById(ctx, id)
	if err != nil {
		return err
	}
	currentContent, err := c.GetContentById(ctx, id)
	if err != nil {
		return err
	}
	currentHash := extractHashcodeFromContent(currentContent)
	if hashcode == currentHash && versionInfo.Version.Editor.Username == "automation" {
		return nil
	}
	content = fmt.Sprintf(HideHash, hashcode) + "\n" + content
	return c.UpdatePageById(ctx, id, content, reCreate)
}

func (c *confluence) updatePage(ctx context.Context, id string, req PageInfo) error {
	return requests.
		URL(fmt.Sprintf("%s/%s", c.baseUrl, id)).
		Method(http.MethodPut).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		BodyJSON(req).
		AddValidator(validateStatus).
		Fetch(ctx)
}

func (c *confluence) UpdatePageParentById(ctx context.Context, id, parentPageId string) error {
	versionInfo, err := c.GetVersionById(ctx, id)
	if err != nil {
		return err
	}
	req := PageInfo{
		Id:      id,
		Type:    "page",
		Title:   versionInfo.Title,
		Version: &PageVersion{Number: versionInfo.Version.Number + 1},
		Parents: []PageInfo{{Type: "page", Id: parentPageId}},
	}
	err = c.updatePage(ctx, id, req)
	if err != nil {
		return fmt.Errorf("UpdatePageParentById — update confluence pageId %s, newParentId %s err: %w", id, parentPageId, err)
	}
	return nil
}

func (c *confluence) AddLabelById(ctx context.Context, id, label string) error {
	req := Label{Prefix: "global", Name: label}
	err := requests.
		URL(fmt.Sprintf("%s/%s/label", c.baseUrl, id)).
		Method(http.MethodPost).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		BodyJSON(req).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return fmt.Errorf("AddLabelToPage — update confluence pageId %s, label %s err: %w", id, label, err)
	}
	return nil
}

func (c *confluence) GetLabelsById(ctx context.Context, id string) ([]string, error) {
	var resp labelResponse
	err := requests.
		URL(fmt.Sprintf("%s/%s/label", c.baseUrl, id)).
		Method(http.MethodGet).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetLabels — get confluence pageId %s labels err: %w", id, err)
	}
	var labels []string
	for _, label := range resp.Results {
		labels = append(labels, label.Name)
	}
	return labels, nil
}

func (c *confluence) GetChildrenById(ctx context.Context, id string, limit int) ([]PageInfo, error) {
	var resp searchPagesResponse
	err := requests.
		URL(fmt.Sprintf("%s/%s/child/page", c.baseUrl, id)).
		Method(http.MethodGet).
		Param("limit", strconv.Itoa(limit)).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetChildrenById — get confluence pageId %s children err: %w", id, err)
	}
	return resp.Results, nil
}

func (c *confluence) GetChildrenByIdRecursive(ctx context.Context, id string, limit int) ([]PageInfo, error) {
	var children []PageInfo
	childrenFirstLevel, err := c.GetChildrenById(ctx, id, limit)
	if err != nil {
		return nil, fmt.Errorf("GetChildrenByIdRecursive — get confluence pageId %s children first level err: %w", id, err)
	}
	for _, child := range childrenFirstLevel {
		childrenNextLevel, err := c.GetChildrenById(ctx, child.Id, limit)
		if err != nil {
			return nil, fmt.Errorf("GetChildrenByIdRecursive — get confluence pageId %s children next level err: %w", id, err)
		}
		children = append(children, child)
		if len(childrenNextLevel) == 0 {
			continue
		}
		recursiveChildren, err := c.GetChildrenByIdRecursive(ctx, child.Id, limit)
		if err != nil {
			return nil, fmt.Errorf("GetChildrenByIdRecursive — get confluence pageId %s children recursive level err: %w", id, err)
		}
		children = append(children, recursiveChildren...)
	}
	return children, nil
}

func (c *confluence) SetRestrictionUser(ctx context.Context, id, username, action string) error {
	// метод работает только в экспериментальном апи, поэтому делаем подмену
	baseUrl := strings.Replace(c.baseUrl, "/api/", "/experimental/", 1)
	err := requests.
		URL(fmt.Sprintf("%s/%s/restriction/byOperation/%s/user", baseUrl, id, action)).
		Method(http.MethodPut).
		Param("userName", username).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return fmt.Errorf("SetRestrictionUser — set restriction '%s' with user '%s' on pageId %s err: %w", action, username, id, err)
	}
	return nil
}

func (c *confluence) SetRestrictionGroup(ctx context.Context, id, groupName, action string) error {
	// метод работает только в экспериментальном апи, поэтому делаем подмену
	baseUrl := strings.Replace(c.baseUrl, "/api/", "/experimental/", 1)
	err := requests.
		URL(fmt.Sprintf("%s/%s/restriction/byOperation/%s/group/%s", baseUrl, id, action, groupName)).
		Method(http.MethodPut).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		AddValidator(validateStatus).
		Fetch(ctx)
	if err != nil {
		return fmt.Errorf("SetRestrictionGroup — set restriction '%s' with group '%s' on pageId %s err: %w", action, groupName, id, err)
	}
	return nil
}

func (c *confluence) SetRestrictionsForHFLabsOnly(ctx context.Context, id string) error {
	err := c.SetRestrictionUser(ctx, id, c.user, "update")
	if err != nil {
		return err
	}
	err = c.SetRestrictionUser(ctx, id, c.user, "read")
	if err != nil {
		return err
	}
	err = c.SetRestrictionGroup(ctx, id, "hfl-conf-worker", "update")
	if err != nil {
		return err
	}
	err = c.SetRestrictionGroup(ctx, id, "hfl-conf-worker", "read")
	if err != nil {
		return err
	}
	return nil
}
