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
	var resp Result
	pathFile, nameFile := path.Split(filePath)
	body := buildSearchRequestFile(a.repo, pathFile, nameFile)
	err := requests.
		URL(a.aqlBaseUrl).
		BasicAuth(a.username, a.password).
		Method(http.MethodPost).
		BodyBytes(body).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return Item{}, err
	}
	if resp.Count.Total == 0 {
		return Item{}, fmt.Errorf("file %s not found", filePath)
	}
	return resp.Items[0], err
}

func (a *artifactory) ListFolders(folderPath, sortType, sortBy string) ([]Item, error) {
	var resp Result
	body := buildSearchRequestFolder(a.repo, folderPath, TypeFolder, sortType, sortBy)
	err := requests.
		URL(a.aqlBaseUrl).
		BasicAuth(a.username, a.password).
		Method(http.MethodPost).
		BodyBytes(body).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Items, err
}

func (a *artifactory) ListFiles(folderPath, sortType, sortBy string) ([]Item, error) {
	var resp Result
	body := buildSearchRequestFolder(a.repo, folderPath, TypeFile, sortType, sortBy)

	err := requests.
		URL(a.aqlBaseUrl).
		BasicAuth(a.username, a.password).
		Method(http.MethodPost).
		BodyBytes(body).
		ToJSON(&resp).
		AddValidator(validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.Items, err
}

func (a *artifactory) FindLastCreatedFileVersion(folderPath string) (string, error) {
	folders, err := a.ListFolders(path.Clean(folderPath), SortTypeDesc, SortByCreated)
	if err != nil {
		return "", fmt.Errorf("a.ListFolders: %w", err)
	}
	if len(folders) == 0 {
		return "", fmt.Errorf("folder %s is empty", folderPath)
	}
	return folders[0].Name, nil
}

func buildSearchRequestFolder(repo, pathFolder, itemType, sortType, sortBy string) []byte {
	pathFolder = path.Clean(pathFolder)
	if itemType == "" {
		itemType = TypeAny
	}
	itemType = strings.ToLower(itemType)
	if sortBy == "" {
		request := fmt.Sprintf(searchTemplate, repo, pathFolder, itemType)
		return []byte(request)
	}
	sortType = strings.ToLower(sortType)
	switch sortType {
	case SortTypeAsc, SortTypeDesc:
		break
	default:
		sortType = SortTypeAsc
	}
	request := fmt.Sprintf(searchTemplateWithSort, repo, pathFolder, itemType, sortType, sortBy)
	return []byte(request)
}

func buildSearchRequestFile(repo, pathFolder, filename string) []byte {
	request := fmt.Sprintf(searchTemplateFile, repo, path.Clean(pathFolder), filename)
	return []byte(request)
}
