package artifactory

type ApiArtifactory interface {
	ListFolders(folderPath, sortType, sortBy string) ([]Item, error)
	ListFoldersWithLimit(folderPath, sortType, sortBy string, limit int) ([]Item, error)
	ListFiles(folderPath, sortType, sortBy string) ([]Item, error)
	ListFilesWithLimit(folderPath, sortType, sortBy string, limit int) ([]Item, error)

	GetFileInfo(filePath string) (Item, error)
	FindLastCreatedFileVersion(folderPath string) (string, error)
}
