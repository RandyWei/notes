package repo

type Release struct {
	TagName     string  `json:"tag_name"`
	PreRelease  bool    `json:"prerelease"`
	PublishedAt string  `json:"published_at"`
	Assets      []Asset `json:"assets"`
	Body        string  `json:"body"`
	HtmlUrl     string  `json:"html_url"`
}

type Asset struct {
	DownloadUrl string `json:"browser_download_url"`
	Size        int    `json:"size"`
	ContentType string `json:"content_type"`
}
