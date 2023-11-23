package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(header http.Header) (string, error) {
	val := header.Get("Authorization")
	if val == "" {
		return "", errors.New("no Authorization information available")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("maformed authorization header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("maformed first part of authorization header")
	}
	return vals[1], nil
}
