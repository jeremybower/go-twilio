// +build unit

package twilio

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSMSSendMessageWithTimeoutErrorUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	ch := make(chan int)
	defer func() { ch <- 0 }()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = <-ch
		w.WriteHeader(http.StatusOK)
	})

	opts := NewOptions("sid", "token")
	opts.HTTPClient = &http.Client{
		Timeout: 10 * time.Millisecond,
	}
	opts.APIBaseURL = server.URL
	builder := NewClient(opts).
		SMS().
		SendMessage("+14155552345", "+15108675310", "Hello!")

	_, err := builder.Do()
	assert.Error(t, err)
}

func TestSMSSendMessageWithReadErrorUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	opts := NewOptions("sid", "token")
	opts.ReaderFunc = func(io.Reader) io.Reader {
		return errReader(0)
	}
	opts.APIBaseURL = server.URL
	builder := NewClient(opts).
		SMS().
		SendMessage("+14155552345", "+15108675310", "Hello!")

	_, err := builder.Do()
	expectedError := "test error"
	assert.Equal(t, err.Error(), expectedError)
}

func TestSMSSendMessageWithUnexpectedStatusCodeUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	opts := NewOptions("sid", "token")
	opts.APIBaseURL = server.URL
	builder := NewClient(opts).
		SMS().
		SendMessage("+14155552345", "+15108675310", "Hello!")

	_, err := builder.Do()
	expectedError := "Unexpected response. Expected 200 but found 400"
	assert.Equal(t, err.Error(), expectedError)
}

func TestSMSSendMessageWithInvalidJSONUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "invalid JSON")
	})

	opts := NewOptions("sid", "token")
	opts.APIBaseURL = server.URL
	builder := NewClient(opts).
		SMS().
		SendMessage("+14155552345", "+15108675310", "Hello!")

	_, err := builder.Do()
	expectedError := "invalid character 'i' looking for beginning of value"
	assert.Equal(t, err.Error(), expectedError)
}

func TestSMSSendMessageUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)

		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "sid", username)
		assert.Equal(t, "token", password)

		to := ""
		from := ""
		body := ""
		for key, value := range r.Form {
			switch key {
			case "To":
				to = value[0]
			case "From":
				from = value[0]
			case "Body":
				body = value[0]
			}
		}

		assert.Equal(t, "+15108675310", to)
		assert.Equal(t, "+14155552345", from)
		assert.Equal(t, "Hello!", body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
	"account_sid": "sid",
	"api_version": "2010-04-01",
	"body": "Hello!",
	"date_created": "Thu, 30 Jul 2015 20:12:31 +0000",
	"date_sent": "Thu, 30 Jul 2015 20:12:33 +0000",
	"date_updated": "Thu, 30 Jul 2015 20:12:33 +0000",
	"direction": "outbound-api",
	"error_code": null,
	"error_message": null,
	"from": "+14155552345",
	"messaging_service_sid": "MGXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"num_media": "0",
	"num_segments": "1",
	"price": -0.00750,
	"price_unit": "USD",
	"sid": "MMXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"status": "sent",
	"subresource_uris": {
		"media": "/2010-04-01/Accounts/sid/Messages/SMXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Media.json"
	},
	"to": "+15108675310",
	"uri": "/2010-04-01/Accounts/sid/Messages/SMXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX.json"
}`))
	})

	opts := NewOptions("sid", "token")
	opts.APIBaseURL = server.URL
	builder := NewClient(opts).
		SMS().
		SendMessage("+14155552345", "+15108675310", "Hello!")

	resp, err := builder.Do()
	assert.NoError(t, err)
	assert.Equal(t, &SMSSendMessageResponse{
		ErrorMessage: "",
		ErrorCode:    0,
		Status:       "sent",
	}, resp)
}
