package twilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Client is the Twilio client.
type Client struct {
	sid   string
	token string
}

// NewClient creates a new Twilio client.
func NewClient(sid, token string) *Client {
	return &Client{
		sid:   sid,
		token: token,
	}
}

func (client *Client) do(
	req *http.Request,
	authorize bool,
	responseObject interface{},
) error {
	if authorize {
		req.SetBasicAuth(client.sid, client.token)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if responseObject != nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
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
