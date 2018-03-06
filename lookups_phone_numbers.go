package twilio

import (
	"net/http"
	"net/url"
)

// LookupsPhoneNumbersBuilder is a helper to build requests to lookup phone
// numbers.
type LookupsPhoneNumbersBuilder struct {
	lookups     *Lookups
	url         string
	phoneNumber string
	countryCode string
	carrier     bool
	callerName  bool
}

// LookupsPhoneNumbersCallerNameResponse is the optional caller name part of the
// response.
type LookupsPhoneNumbersCallerNameResponse struct {
	CallerName string `json:"caller_name"`
	CallerType string `json:"caller_type"`
	ErrorCode  string `json:"error_code"`
}

// LookupsPhoneNumbersCarrierResponse is the optional carrier information part
// of the response.
type LookupsPhoneNumbersCarrierResponse struct {
	MobileCountryCode string `json:"mobile_country_code"`
	MobileNetworkCode string `json:"mobile_network_code"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	ErrorCode         string `json:"error_code"`
}

// LookupsPhoneNumbersResponse is the response for requests to lookup phone
// numbers.
type LookupsPhoneNumbersResponse struct {
	CountryCode    string                                `json:"country_code"`
	PhoneNumber    string                                `json:"phone_number"`
	NationalFormat string                                `json:"national_format"`
	CallerName     LookupsPhoneNumbersCallerNameResponse `json:"caller_name"`
	Carrier        LookupsPhoneNumbersCarrierResponse    `json:"carrier"`
	URL            string                                `json:"url"`
}

// WithCountryCode is optional. It adds the ISO country code of the phone
// number. This is used to specify the country when the number is provided in a
// national format.
func (b *LookupsPhoneNumbersBuilder) WithCountryCode(
	countryCode string,
) *LookupsPhoneNumbersBuilder {
	b.countryCode = countryCode
	return b
}

// WithCarrier indicates that carrier information should be returned with the
// request. Extra charges may apply.
func (b *LookupsPhoneNumbersBuilder) WithCarrier() *LookupsPhoneNumbersBuilder {
	b.carrier = true
	return b
}

// WithCallerName indicates that caller name information should be returned with
// the request. Extra charges may apply.
func (b *LookupsPhoneNumbersBuilder) WithCallerName() *LookupsPhoneNumbersBuilder {
	b.callerName = true
	return b
}

// Build will build the request.
func (b *LookupsPhoneNumbersBuilder) Build() (*http.Request, error) {
	requestURL, err := url.Parse(
		b.lookups.host + "/v1/PhoneNumbers/" + url.PathEscape(b.phoneNumber))
	if err != nil {
		return nil, err
	}

	q := requestURL.Query()

	if b.countryCode != "" {
		q.Set("CountryCode", b.countryCode)
	}

	if b.callerName {
		q.Add("Type", "caller-name")
	}

	if b.carrier {
		q.Add("Type", "carrier")
	}

	requestURL.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Do will build and perform the request, and return the response.
func (b *LookupsPhoneNumbersBuilder) Do() (*LookupsPhoneNumbersResponse, error) {
	req, err := b.Build()
	if err != nil {
		return nil, err
	}

	var responseObject LookupsPhoneNumbersResponse
	err = b.lookups.client.do(req, true, http.StatusOK, &responseObject)
	if err != nil {
		return nil, err
	}

	return &responseObject, nil
}
