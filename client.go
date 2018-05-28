package twilio

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Client is the Twilio client.
type Client interface {
	Lookup() Lookup
}

// Options are the configuration options for the client.
type Options struct {
	LookupBaseURL string
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
		HTTPClient:    &http.Client{},
		ReaderFunc:    readerFunc,
		SID:           sid,
		Token:         token,
	}
}

// NewClient will create a new client with the given options.
func NewClient(opts *Options) Client {
	return &clientImpl{
		opts: opts,
	}
}

type clientImpl struct {
	opts *Options
}

func do(
	opts *Options,
	req *http.Request,
	authorize bool,
	expectedStatusCode int,
	responseObject interface{},
) error {
	if authorize {
		req.SetBasicAuth(opts.SID, opts.Token)
	}

	resp, err := opts.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != expectedStatusCode {
		return fmt.Errorf(
			"Unexpected response. Expected %d but found %d",
			expectedStatusCode,
			resp.StatusCode)
	}

	if responseObject != nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(opts.ReaderFunc(resp.Body))
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, responseObject)
		if err != nil {
			return err
		}
	}

	return nil
}
