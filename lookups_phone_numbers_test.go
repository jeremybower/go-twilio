package twilio

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupsPhoneNumbersBuilder(t *testing.T) {
	// handler will verify that all params are received and respond with JSON.
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)

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
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	host := ts.URL
	builder := NewClient("sid", "token").
		Lookups().WithHost(host).
		PhoneNumbers("+15108675310").
		WithCallerName().
		WithCarrier().
		WithCountryCode("US")

	// call library I want to test, using test server ts
	resp, err := builder.Do()
	assert.NoError(t, err)
	assert.Equal(t, &LookupsPhoneNumbersResponse{
		URL:            "https://lookups.twilio.com/v1/PhoneNumbers/+15108675310?Type=carrier",
		NationalFormat: "(510) 867-5310",
		PhoneNumber:    "+15108675310",
		CountryCode:    "US",
		Carrier: LookupsPhoneNumbersCarrierResponse{
			ErrorCode:         "",
			Type:              "mobile",
			Name:              "T-Mobile USA, Inc.",
			MobileNetworkCode: "160",
			MobileCountryCode: "310",
		},
		CallerName: LookupsPhoneNumbersCallerNameResponse{
			CallerName: "John Smith",
			CallerType: "consumer",
			ErrorCode:  "",
		},
	}, resp)
}
