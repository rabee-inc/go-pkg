package firebaseauth

import (
	"fmt"
	"strings"
)

const (
	headerPrefix      string = "BEARER"
	debugHeaderPrefix string = "user="
)

func getTokenByAuthHeader(ah string) string {
	pLen := len(headerPrefix)
	if len(ah) > pLen && strings.ToUpper(ah[0:pLen]) == headerPrefix {
		return ah[pLen+1:]
	}
	return ""
}

func getDebugByAuthHeader(ah string) string {
	token := getTokenByAuthHeader(ah)
	fmt.Printf("token: %s\n", token)
	if strings.HasPrefix(token, debugHeaderPrefix) {
		return token[len(debugHeaderPrefix):]
	}
	return ""
}
