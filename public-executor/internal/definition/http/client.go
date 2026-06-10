package http

import (
	"net/http"
)

func ProvideDefaultHTTPClient() *http.Client {
	return &http.Client{}
}
