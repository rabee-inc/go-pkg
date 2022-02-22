package cloudstorage

import "strings"

// GenerateFileURL ... GCSのファイルURLを作成する
func GenerateFileURL(bucket string, path string, name string) string {
	return strings.Join([]string{BaseURL, bucket, path, name}, "/")
}
