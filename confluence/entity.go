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
	Space   space                   `json:"space"`
	Version pageVersion             `json:"version"`
	Body    pageBody                `json:"body"`
	Parents []pageRequestOrResponse `json:"ancestors"`
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
