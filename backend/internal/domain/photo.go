package domain

type Photo struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Author        string `json:"author"`
	Category      string `json:"category"`
	Tags          string `json:"tags"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	URL           string `json:"url"`
	ThumbnailURL  string `json:"thumbnail_url"`
	Source        string `json:"source"`
	DownloadCount int64  `json:"download_count"`
	Likes         int64  `json:"likes"`
	Version       int32  `json:"version"`
}
