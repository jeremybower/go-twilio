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
func (b *SMSSendMessageBuilder) Build() (*http.Request, error) {
	requestURL, err := url.Parse(
		b.opts.APIBaseURL + "/Accounts/" + b.opts.SID + "/Messages.json")
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	v.Set("From", b.from)
	v.Set("To", b.to)
	v.Set("Body", b.body)
	rb := *strings.NewReader(v.Encode())

	req, err := http.NewRequest(http.MethodPost, requestURL.String(), &rb)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

// Do will build and perform the request, and return the response.
func (b *SMSSendMessageBuilder) Do() (*SMSSendMessageResponse, error) {
	req, err := b.Build()
	if err != nil {
		return nil, err
	}

	var response SMSSendMessageResponse
	err = do(b.opts, req, true, http.StatusOK, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
