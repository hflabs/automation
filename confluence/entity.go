package confluence

type confluence struct {
	user     string
	password string
	baseUrl  string
}

type pageRequestOrResponse struct {
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Id      string      `json:"id"`
	Version pageVersion `json:"version"`
	Body    pageBody    `json:"body"`
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
