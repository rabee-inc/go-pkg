package cloudstorage

import "time"

// キャッシュ設定
type CacheMode struct {
	Disabled bool
	Expire   time.Duration
}

// アップロードのレスポンス
type UploadResponse struct {
	URL string
}
