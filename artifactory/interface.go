package artifactory

type ApiArtifactory interface {
	ListFolders(folderPath, sortType, sortBy string) ([]Item, error)
	ListFiles(folderPath, sortType, sortBy string) ([]Item, error)

	GetFileInfo(filePath string) (Item, error)
	FindLastCreatedFileVersion(folderPath string) (string, error)
}
