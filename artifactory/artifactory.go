package artifactory

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/http"
	"path"
)

func NewArtifactory(baseUrl, repo, username, password string) ApiArtifactory {
	aqlUrl := baseUrl + "/api/search/aql"
	return &artifactory{baseUrl, repo, aqlUrl, username, password}
}

func (a *artifactory) GetFileInfo(filePath string) (Item, error) {
	pathFile, nameFile := path.Split(filePath)
	body, err := buildFindQuery(a.repo, pathFile, ListOptions{
		Name: nameFile,
		Type: TypeFile,
	})
	if err != nil {
		return Item{}, err
	}

	items, err := a.makeRequest(body)
	if len(items) == 0 {
		return Item{}, fmt.Errorf("file %s not found", filePath)
	}
	return items[0], err
}

func (a *artifactory) ListItems(folderPath string, options ListOptions) ([]Item, error) {
	body, err := buildFindQuery(a.repo, folderPath, options)
	if err != nil {
		return nil, err
	}
	return a.makeRequest(body)
}

func (a *artifactory) FindLastCreatedFileVersion(folderPath string) (string, error) {
	folders, err := a.ListItems(path.Clean(folderPath), ListOptions{
		SortDirection: SortTypeDesc,
		SortValue:     SortByCreated,
		Limit:         1,
		Type:          TypeFolder,
	})
	if err != nil {
		return "", fmt.Errorf("a.ListItems: %w", err)
	}
	if len(folders) == 0 {
		return "", fmt.Errorf("folder %s is empty", folderPath)
	}
	return folders[0].Name, nil
}

func (a *artifactory) makeRequest(body string) ([]Item, error) {
	var resp Result
	err := requests.
		URL(a.aqlBaseUrl).
		BasicAuth(a.username, a.password).
		Method(http.MethodPost).
		BodyBytes([]byte(body)).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Items, err
}
