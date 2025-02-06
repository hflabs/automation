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

const searchLimitTemplate = ".limit(%v)"

const searchTemplate = `
items.find({
    "repo": "%s",
    "path": "%s",
    "type": "%s"
}).include(
    "repo",
    "path",
    "name",
    "created",
    "type",
	"size"
)
`

const searchTemplateWithSort = `
items.find({
    "repo": "%s",
    "path": "%s",
    "type": "%s"
}).include(
    "repo",
    "path",
    "name",
    "created",
    "type",
	"size"
).sort({
    "$%s": [
        "%s"
    ]
})
`

const searchTemplateFile = `
items.find({
    "repo": "%s",
    "path": "%s",
	"name": "%s",
    "type": "file"
}).include(
    "repo",
    "path",
    "name",
    "created",
    "type",
	"size"
)
`
