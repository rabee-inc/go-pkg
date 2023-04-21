package images

// 画像オブジェクト
type Object struct {
	ID            string           `firestore:"id"             json:"id"`
	OriginalURL   string           `firestore:"original_url"   json:"original_url"`
	URL           string           `firestore:"url"            json:"url"`
	Filename      string           `firestore:"filename"       json:"filename"`
	ContentType   string           `firestore:"content_type"   json:"content_type"`
	DominantColor string           `firestore:"dominant_color" json:"dominant_color"`
	Sizes         map[string]*Size `firestore:"sizes"          json:"sizes"`
	IsDefault     bool             `firestore:"is_default"     json:"is_default"`
}

// サイズ
type Size struct {
	URL    string `firestore:"url"    json:"url"`
	Width  int    `firestore:"width"  json:"width"`
	Height int    `firestore:"height" json:"height"`
}

// 画像変換リクエスト
type ConvertRequest struct {
	Key                  string   `json:"key"`
	SourceID             string   `json:"source_id"`
	SourceURLs           []string `json:"source_urls"`
	DstFilePath          string   `json:"dst_file_path"`
	DstEndpoint          string   `json:"dst_endpoint"`
	DstEndpointAuthToken string   `json:"dst_endpoint_auth_token"`
}

// 画像作成リクエスト
type GenerateRequest struct {
	Key         string `json:"key"`
	SourceID    string `json:"source_id"`
	SourceURL   string `json:"source_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	DstFilePath string `json:"dst_file_path"`
}

type ClientOption struct {
	ConverterTopicID string
	GeneratorTopicID string
}
