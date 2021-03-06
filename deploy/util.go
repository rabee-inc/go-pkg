package deploy

import (
	"os"
	"strings"
)

// IsLocal ... 現在の環境がローカルか判定する
func IsLocal() bool {
	d := os.Getenv("DEPLOY")
	return d == "" || d == "local"
}

// IsStaging ... 現在の環境がステージングか判定する
func IsStaging() bool {
	return strings.HasPrefix(os.Getenv("DEPLOY"), "staging")
}

// IsProduction ... 現在の環境が本番か判定する
func IsProduction() bool {
	return strings.HasPrefix(os.Getenv("DEPLOY"), "production")
}
