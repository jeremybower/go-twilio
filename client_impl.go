package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewClient will create a new client with the given options.
func NewClient(opts *Options) Client {
	return &clientImpl{
		opts: opts,
	}
}

type clientImpl struct {
	opts *Options
}

func (client *clientImpl) do(
	req *http.Request,
	authorize bool,
	expectedStatusCode int,
	responseObject interface{},
) error {
	if authorize {
		req.SetBasicAuth(client.opts.SID, client.opts.Token)
	}

	resp, err := client.opts.HTTPClient.Do(req)
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
		body, err := ioutil.ReadAll(client.opts.ReaderFunc(resp.Body))
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
