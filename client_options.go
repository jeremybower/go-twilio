package twilio

import (
	"io"
	"net/http"
)

// Options are the configuration options for the client.
type Options struct {
	LookupBaseURL string
	APIBaseURL    string
	HTTPClient    *http.Client
	ReaderFunc    func(io.Reader) io.Reader
	SID           string
	Token         string
}

// NewOptions will create new options with default values.
func NewOptions(sid, token string) *Options {
	readerFunc := func(r io.Reader) io.Reader {
		return r
	}

	return &Options{
		LookupBaseURL: "https://lookups.twilio.com",
		APIBaseURL:    "https://api.twilio.com/2010-04-01",
		HTTPClient:    &http.Client{},
		ReaderFunc:    readerFunc,
		SID:           sid,
		Token:         token,
	}
}
