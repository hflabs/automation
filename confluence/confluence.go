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
		URL(fmt.Sprintf("%s/?title=%s&spaceKey=%s", c.baseUrl, url.QueryEscape(name), url.QueryEscape(spaceKey))).
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

func (c *confluence) UpdatePageById(id string, content string, reCreate bool) error {
	versionInfo, err := c.GetVersionInfoById(id)
	if err != nil {
		return err
	}

	req := pageRequestOrResponse{
		Type:    "page",
		Title:   versionInfo.Title,
		Id:      id,
		Version: pageVersion{Number: versionInfo.Version.Number + 1},
		Body:    pageBody{pageStorage{Value: content, Representation: "storage"}},
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
	err = requests.
		URL(fmt.Sprintf("%s/%s", c.baseUrl, id)).
		Method(http.MethodPut).
		ContentType("application/json").
		BasicAuth(c.user, c.password).
		BodyJSON(req).
		AddValidator(validateStatus).
		Fetch(context.Background())
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

func extractHashcodeFromContent(content string) string {
	match := regexp.MustCompile(hashcode_pattern).FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
