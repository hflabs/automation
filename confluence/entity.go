package confluence

type confluence struct {
	user     string
	password string
	baseUrl  string
}

type pageRequestOrResponse struct {
	Type    string                  `json:"type"`
	Title   string                  `json:"title"`
	Id      int                     `json:"id,omitempty"`
	Space   space                   `json:"space"`
	Version pageVersion             `json:"version"`
	Body    pageBody                `json:"body"`
	Parents []pageRequestOrResponse `json:"ancestors"`
}

type space struct {
	Id   int    `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}
type pageVersion struct {
	Number int `json:"number"`
}

type pageBody struct {
	Storage pageStorage `json:"storage"`
}

type pageStorage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type VersionResponse struct {
	Title   string `json:"title"`
	Version struct {
		Number int `json:"number"`
		Editor struct {
			Username string `json:"username"`
		} `json:"by"`
	} `json:"version"`
}

type searchPagesResponse struct {
	Results []PageInfo `json:"results"`
	Start   int        `json:"start"`
	Limit   int        `json:"limit"`
	Size    int        `json:"size"`
}

type PageInfo struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Title  string `json:"title"`
}
