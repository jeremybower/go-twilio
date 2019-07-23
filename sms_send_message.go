package twilio

import (
	"net/http"
	"net/url"
	"strings"
)

// SMSSendMessageBuilder builds an SMS message to send.
type SMSSendMessageBuilder struct {
	opts *Options
	from string
	to   string
	body string
}

// SMSSendMessageResponse is Twilio's response after sending an SMS message.
type SMSSendMessageResponse struct {
	Status       string `json:"status"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// Build will build the request.
func (client *clientImpl) SendSMSMessage(
	from string,
	to string,
	body string,
) (*SMSSendMessageResponse, error) {
	requestURL, err := url.Parse(
		client.opts.APIBaseURL + "/Accounts/" + client.opts.SID + "/Messages.json")
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("From", from)
	v.Set("To", to)
	v.Set("Body", body)
	rb := *strings.NewReader(v.Encode())

	req, err := http.NewRequest(http.MethodPost, requestURL.String(), &rb)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var response SMSSendMessageResponse
	err = client.do(req, true, http.StatusOK, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
