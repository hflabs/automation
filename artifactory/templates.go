package artifactory

const (
	TypeFile   = "file"
	TypeFolder = "folder"
	TypeAny    = "any"

	SortTypeAsc  = "asc"
	SortTypeDesc = "desc"

	SortByName    = "name"
	SortByCreated = "created"
	SortBySize    = "size"
)

const findTemplate = "items.find(%s)"

const includeTemplate = `.include(
    "repo",
    "path",
    "name",
    "created",
    "type",
	"size"
)`

const sortTemplate = `.sort({
    "$%s": [
        "%s"
    ]
})`

const limitTemplate = ".limit(%v)"
