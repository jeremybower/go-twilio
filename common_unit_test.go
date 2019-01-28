package twilio

import (
	"errors"
	"net/http"
	"net/http/httptest"
)

func setupMockServer() (*http.ServeMux, *httptest.Server, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return mux, server, func() {
		server.Close()
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
