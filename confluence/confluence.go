package confluence

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func NewConfluence(baseUrl, user, password string) ApiConfluence {
	return &confluence{user, password, baseUrl}
}

func (c *confluence) GetPagesByName(name, spaceKey string) ([]PageInfo, error) {
	var resp searchPagesResponse
	err := requests.
		URL(fmt.Sprintf("%s?title=%s&spaceKey=%s&type=page", c.baseUrl, url.QueryEscape(name), url.QueryEscape(spaceKey))).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return resp.Results, fmt.Errorf("GetPagesByName — get confluence page by name %s in space %s err: %w", name, spaceKey, err)
	}
	return resp.Results, nil
}

func (c *confluence) GetPagesByIncludedName(name, spaceKey string) ([]PageInfo, error) {
	cqlQuery := fmt.Sprintf("space=\"%s\" AND type=\"page\" AND title~\"%s\"", spaceKey, name)
	var resp searchPagesResponse
	err := requests.
		URL(fmt.Sprintf("%s/search?cql=%s", c.baseUrl, url.QueryEscape(cqlQuery))).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return resp.Results, fmt.Errorf("GetPagesByIncludedName — get confluence page by include name %s in space %s err: %w", name, spaceKey, err)
	}
	return resp.Results, nil
}

func (c *confluence) GetContentById(id string) (string, error) {
	var resp pageRequestOrResponse
	err := requests.
		URL(fmt.Sprintf("%s/%s?expand=body.storage", c.baseUrl, id)).
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return resp.Body.Storage.Value, fmt.Errorf("GetContentById — get confluence pageId %s err: %w", id, err)
	}
	return resp.Body.Storage.Value, err
}

func (c *confluence) GetVersionInfoById(id string) (VersionResponse, error) {
	var resp VersionResponse
	err := requests.
		URL(fmt.Sprintf("%s/%s?expand=version", c.baseUrl, id)).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())

	if err != nil {
		return resp, fmt.Errorf("GetLastVersionInfoById — get confluence pageId %s err: %w", id, err)
	}
	return resp, nil
}

func (c *confluence) CreatePage(name, spaceKey, content string, parentPageId string) (string, error) {
	req := pageRequestOrResponse{
		Type:    "page",
		Title:   name,
		Space:   &space{Key: spaceKey},
		Parents: []pageRequestOrResponse{{Type: "page", Id: parentPageId}},
		Body: &pageBody{
			Storage: pageStorage{
				Value:          content,
				Representation: "storage",
			},
		},
	}
	err := requests.
		URL(c.baseUrl).
		Method(http.MethodPost).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		BodyJSON(req).
		ToJSON(&req).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return "", fmt.Errorf("CreatePage — create confluence page `%s` in space `%s` with content `%s` err: %w", name, spaceKey, content, err)
	}
	return req.Id, nil
}

func (c *confluence) UpdatePageById(id string, content string, reCreate bool) error {
	versionInfo, err := c.GetVersionInfoById(id)
	if err != nil {
		return err
	}

	req := pageRequestOrResponse{
		Type:    "page",
		Title:   versionInfo.Title,
		Id:      id,
		Version: &pageVersion{Number: versionInfo.Version.Number + 1},
		Body:    &pageBody{pageStorage{Value: content, Representation: "storage"}},
	}
	oldContent, err := c.GetContentById(id)
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
	err = c.updatePage(id, req)
	if err != nil {
		return fmt.Errorf("UpdatePageById — update confluence pageId %s, content %s err: %w", id, content, err)
	}
	return nil
}

func (c *confluence) UpdatePageByIdWithCheck(id string, content string, reCreate bool) error {
	hash := md5.Sum([]byte(content))
	hashcode := hex.EncodeToString(hash[:])
	versionInfo, err := c.GetVersionInfoById(id)
	if err != nil {
		return err
	}
	currentContent, err := c.GetContentById(id)
	if err != nil {
		return err
	}
	currentHash := extractHashcodeFromContent(currentContent)
	if hashcode == currentHash && versionInfo.Version.Editor.Username == "automation" {
		return nil
	}
	content = fmt.Sprintf(HideHash, hashcode) + "\n" + content
	return c.UpdatePageById(id, content, reCreate)
}

func (c *confluence) updatePage(id string, req pageRequestOrResponse) error {
	return requests.
		URL(fmt.Sprintf("%s/%s", c.baseUrl, id)).
		Method(http.MethodPut).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		BodyJSON(req).
		AddValidator(validateStatus).
		Fetch(context.Background())
}

func (c *confluence) UpdatePageParentById(id, parentId string) error {
	versionInfo, err := c.GetVersionInfoById(id)
	if err != nil {
		return err
	}
	req := pageRequestOrResponse{
		Id:      id,
		Type:    "page",
		Title:   versionInfo.Title,
		Version: &pageVersion{Number: versionInfo.Version.Number + 1},
		Parents: []pageRequestOrResponse{{Type: "page", Id: parentId}},
	}
	err = c.updatePage(id, req)
	if err != nil {
		return fmt.Errorf("UpdatePageParentById — update confluence pageId %s, newParentId %s err: %w", id, parentId, err)
	}
	return nil
}

func (c *confluence) AddLabelToPage(pageId, label string) error {
	req := Label{Prefix: "global", Name: label}
	err := requests.
		URL(fmt.Sprintf("%s/%s/label", c.baseUrl, pageId)).
		Method(http.MethodPost).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		BodyJSON(req).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return fmt.Errorf("AddLabelToPage — update confluence pageId %s, label %s err: %w", pageId, label, err)
	}
	return nil
}

func (c *confluence) GetLabels(pageId string) ([]string, error) {
	var resp labelResponse
	err := requests.
		URL(fmt.Sprintf("%s/%s/label", c.baseUrl, pageId)).
		Method(http.MethodGet).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, fmt.Errorf("GetLabels — get confluence pageId %s labels err: %w", pageId, err)
	}
	var labels []string
	for _, label := range resp.Results {
		labels = append(labels, label.Name)
	}
	return labels, nil
}

func (c *confluence) GetChildrenById(pageId string) ([]PageInfo, error) {
	var resp searchPagesResponse
	err := requests.
		URL(fmt.Sprintf("%s/%s/child/page?limit=250", c.baseUrl, pageId)).
		Method(http.MethodGet).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, fmt.Errorf("GetChildrenById — get confluence pageId %s children err: %w", pageId, err)
	}
	return resp.Results, nil
}

func (c *confluence) GetChildrenByIdRecursive(pageId string) ([]PageInfo, error) {
	var children []PageInfo
	childrenFirstLevel, err := c.GetChildrenById(pageId)
	if err != nil {
		return nil, fmt.Errorf("GetChildrenByIdRecursive — get confluence pageId %s children first level err: %w", pageId, err)
	}
	for _, child := range childrenFirstLevel {
		childrenNextLevel, err := c.GetChildrenById(child.Id)
		if err != nil {
			return nil, fmt.Errorf("GetChildrenByIdRecursive — get confluence pageId %s children next level err: %w", pageId, err)
		}
		children = append(children, child)
		if len(childrenNextLevel) == 0 {
			continue
		}
		recursiveChildren, err := c.GetChildrenByIdRecursive(child.Id)
		if err != nil {
			return nil, fmt.Errorf("GetChildrenByIdRecursive — get confluence pageId %s children recursive level err: %w", pageId, err)
		}
		children = append(children, recursiveChildren...)
	}
	return children, nil
}

func extractHashcodeFromContent(content string) string {
	match := regexp.MustCompile(hashcode_pattern).FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
