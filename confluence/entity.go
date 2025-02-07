package confluence

type confluence struct {
	user     string
	password string
	baseUrl  string
}

type Space struct {
	Id   int    `json:"id,omitempty"`
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}
type PageVersion struct {
	Number int `json:"number,omitempty"`
}

type PageBody struct {
	Storage PageStorage `json:"storage"`
}

type PageStorage struct {
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
	Status  string       `json:"status,omitempty"`
	Type    string       `json:"type,omitempty"`
	Title   string       `json:"title,omitempty"`
	Id      string       `json:"id,omitempty"`
	Space   *Space       `json:"Space,omitempty"`
	Version *PageVersion `json:"version,omitempty"`
	Body    *PageBody    `json:"body,omitempty"`
	Parents []PageInfo   `json:"ancestors,omitempty"`
}
