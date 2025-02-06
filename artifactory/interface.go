package artifactory

type ApiArtifactory interface {
	ListItems(folderPath string, options ListOptions) ([]Item, error)
	GetFileInfo(filePath string) (Item, error)
	FindLastCreatedFileVersion(folderPath string) (string, error)
}
