package ogp

// OpenGraph ... OGPでよく使うもの
type OpenGraph struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	FaviconURL  string `json:"favicon_url"`
}

// ConvRequest ... 画像変換リクエスト
type ConvRequest struct {
	Key         string `json:"key"`
	SourceID    string `json:"source_id"`
	SourceURL   string `json:"source_url"`
	DstFilePath string `json:"dst_file_path"`
}
