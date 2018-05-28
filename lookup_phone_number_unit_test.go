// +build unit

package twilio

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupMockServer() (*http.ServeMux, *httptest.Server, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	return mux, server, func() {
		server.Close()
	}
}

func TestTimeoutErrorUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	ch := make(chan int)
	defer func() { ch <- 0 }()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = <-ch
		w.WriteHeader(http.StatusOK)
	})

	opts := NewOptions("sid", "token")
	opts.LookupBaseURL = server.URL
	opts.HTTPClient = &http.Client{
		Timeout: 10 * time.Millisecond,
	}
	builder := NewClient(opts).
		Lookup().
		PhoneNumber("+15108675310").
		WithCountryCode("US").
		IncludeCallerNameInResponse().
		IncludeCarrierInResponse()

	_, err := builder.Do()
	assert.Error(t, err)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestReadErrorUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	opts := NewOptions("sid", "token")
	opts.LookupBaseURL = server.URL
	opts.ReaderFunc = func(io.Reader) io.Reader {
		return errReader(0)
	}
	builder := NewClient(opts).
		Lookup().
		PhoneNumber("+15108675310").
		WithCountryCode("US").
		IncludeCallerNameInResponse().
		IncludeCarrierInResponse()

	_, err := builder.Do()
	expectedError := "test error"
	assert.Equal(t, err.Error(), expectedError)
}

func TestUnexpectedStatusCodeUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	opts := NewOptions("sid", "token")
	opts.LookupBaseURL = server.URL
	builder := NewClient(opts).
		Lookup().
		PhoneNumber("+15108675310").
		WithCountryCode("US").
		IncludeCallerNameInResponse().
		IncludeCarrierInResponse()

	_, err := builder.Do()
	expectedError := "Unexpected response. Expected 200 but found 400"
	assert.Equal(t, err.Error(), expectedError)
}

func TestInvalidJSONUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "invalid JSON")
	})

	opts := NewOptions("sid", "token")
	opts.LookupBaseURL = server.URL
	builder := NewClient(opts).
		Lookup().
		PhoneNumber("+15108675310").
		WithCountryCode("US").
		IncludeCallerNameInResponse().
		IncludeCarrierInResponse()

	_, err := builder.Do()
	expectedError := "invalid character 'i' looking for beginning of value"
	assert.Equal(t, err.Error(), expectedError)
}

func TestLookupPhoneNumberUsingMockServer(t *testing.T) {
	mux, server, shutdown := setupMockServer()
	defer shutdown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)

		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "sid", username)
		assert.Equal(t, "token", password)

		assert.Equal(t, 2, len(r.URL.Query()["Type"]))
		assert.Equal(t, []string{"caller-name", "carrier"}, r.URL.Query()["Type"])
		assert.Equal(t, "US", r.URL.Query().Get("CountryCode"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
    "url": "https://lookups.twilio.com/v1/PhoneNumbers/+15108675310?Type=carrier",
    "carrier": {
        "error_code": null,
        "type": "mobile",
        "name": "T-Mobile USA, Inc.",
        "mobile_network_code": "160",
        "mobile_country_code": "310"
    },
    "caller_name": {
      "caller_name": "John Smith",
      "caller_type": "consumer",
      "error_code": null
    },
    "national_format": "(510) 867-5310",
    "phone_number": "+15108675310",
    "country_code": "US"
}`))
	})

	opts := NewOptions("sid", "token")
	opts.LookupBaseURL = server.URL
	builder := NewClient(opts).
		Lookup().
		PhoneNumber("+15108675310").
		WithCountryCode("US").
		IncludeCallerNameInResponse().
		IncludeCarrierInResponse()

	resp, err := builder.Do()
	assert.NoError(t, err)
	assert.Equal(t, &LookupPhoneNumberResponse{
		URL:            "https://lookups.twilio.com/v1/PhoneNumbers/+15108675310?Type=carrier",
		NationalFormat: "(510) 867-5310",
		PhoneNumber:    "+15108675310",
		CountryCode:    "US",
		Carrier: LookupPhoneNumberCarrierResponse{
			ErrorCode:         "",
			Type:              "mobile",
			Name:              "T-Mobile USA, Inc.",
			MobileNetworkCode: "160",
			MobileCountryCode: "310",
		},
		CallerName: LookupPhoneNumberCallerNameResponse{
			CallerName: "John Smith",
			CallerType: "consumer",
			ErrorCode:  "",
		},
	}, resp)
}
