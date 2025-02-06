package artifactory

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"net/http"
	"path"
	"strings"
)

func NewArtifactory(baseUrl, repo, username, password string) ApiArtifactory {
	aqlUrl := baseUrl + "/api/search/aql"
	return &artifactory{baseUrl, repo, aqlUrl, username, password}
}

func (a *artifactory) GetFileInfo(filePath string) (Item, error) {
	pathFile, nameFile := path.Split(filePath)
	body := buildSearchRequestFile(a.repo, pathFile, nameFile)
	items, err := a.makeRequest(body)
	if len(items) == 0 {
		return Item{}, fmt.Errorf("file %s not found", filePath)
	}
	return items[0], err
}

func (a *artifactory) ListFoldersWithLimit(folderPath, sortType, sortBy string, limit int) ([]Item, error) {
	body := buildSearchRequestFolderWithLimit(a.repo, folderPath, TypeFolder, sortType, sortBy, limit)
	return a.makeRequest(body)
}

func (a *artifactory) ListFolders(folderPath, sortType, sortBy string) ([]Item, error) {
	body := buildSearchRequestFolder(a.repo, folderPath, TypeFolder, sortType, sortBy)
	return a.makeRequest(body)
}

func (a *artifactory) ListFilesWithLimit(folderPath, sortType, sortBy string, limit int) ([]Item, error) {
	body := buildSearchRequestFolderWithLimit(a.repo, folderPath, TypeFile, sortType, sortBy, limit)
	return a.makeRequest(body)
}

func (a *artifactory) ListFiles(folderPath, sortType, sortBy string) ([]Item, error) {
	body := buildSearchRequestFolder(a.repo, folderPath, TypeFile, sortType, sortBy)
	return a.makeRequest(body)
}

func (a *artifactory) FindLastCreatedFileVersion(folderPath string) (string, error) {
	folders, err := a.ListFoldersWithLimit(path.Clean(folderPath), SortTypeDesc, SortByCreated, 1)
	if err != nil {
		return "", fmt.Errorf("a.ListFolders: %w", err)
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

func buildSearchRequestFolderWithLimit(repo, pathFolder, itemType, sortType, sortBy string, limit int) string {
	request := buildSearchRequestFolder(repo, pathFolder, itemType, sortType, sortBy)
	return request + fmt.Sprintf(searchLimitTemplate, limit)
}

func buildSearchRequestFolder(repo, pathFolder, itemType, sortType, sortBy string) string {
	pathFolder = path.Clean(pathFolder)
	if itemType == "" {
		itemType = TypeAny
	}
	itemType = strings.ToLower(itemType)
	if sortBy == "" {
		return fmt.Sprintf(searchTemplate, repo, pathFolder, itemType)
	}
	sortType = strings.ToLower(sortType)
	switch sortType {
	case SortTypeAsc, SortTypeDesc:
		break
	default:
		sortType = SortTypeAsc
	}
	request := fmt.Sprintf(searchTemplateWithSort, repo, pathFolder, itemType, sortType, sortBy)
	return request
}

func buildSearchRequestFile(repo, pathFolder, filename string) string {
	request := fmt.Sprintf(searchTemplateFile, repo, path.Clean(pathFolder), filename)
	return request
}
