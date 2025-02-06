package confluence

type confluence struct {
	user     string
	password string
	baseUrl  string
}

type pageRequestOrResponse struct {
	Type    string                  `json:"type,omitempty"`
	Title   string                  `json:"title,omitempty"`
	Id      string                  `json:"id,omitempty"`
	Space   *space                  `json:"space,omitempty"`
	Version *pageVersion            `json:"version,omitempty"`
	Body    *pageBody               `json:"body,omitempty"`
	Parents []pageRequestOrResponse `json:"ancestors,omitempty"`
}

type space struct {
	Id   int    `json:"id,omitempty"`
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}
type pageVersion struct {
	Number int `json:"number,omitempty"`
}

type pageBody struct {
	Storage pageStorage `json:"storage"`
}

type pageStorage struct {
	Value          string `json:"value,omitempty"`
	Representation string `json:"representation,omitempty"`
}

type Label struct {
	Prefix string `json:"prefix,omitempty"`
	Name   string `json:"name,omitempty"`
}

type labelResponse struct {
	Results    []Label `json:"results"`
	TotalCount int     `json:"totalCount"`
	StartIndex int     `json:"start"`
	Limit      int     `json:"limit"`
	Size       int     `json:"size"`
}
type VersionResponse struct {
	Title   string `json:"title,omitempty"`
	Version struct {
		Number int `json:"number,omitempty"`
		Editor struct {
			Username string `json:"username,omitempty"`
		} `json:"by"`
	} `json:"version"`
}

type searchPagesResponse struct {
	Results []PageInfo `json:"results"`
	Start   int        `json:"start,omitempty"`
	Limit   int        `json:"limit,omitempty"`
	Size    int        `json:"size,omitempty"`
}

type PageInfo struct {
	Id     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
}
