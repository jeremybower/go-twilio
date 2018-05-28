package twilio

import (
	"net/http"
	"net/url"
)

// LookupPhoneNumberBuilder is a helper to build requests to lookup phone
// numbers.
type LookupPhoneNumberBuilder struct {
	opts                        *Options
	url                         string
	phoneNumber                 string
	countryCode                 string
	includeCarrierInResponse    bool
	includeCallerNameInResponse bool
}

// LookupPhoneNumberCallerNameResponse is the optional caller name part of the
// response.
type LookupPhoneNumberCallerNameResponse struct {
	CallerName string `json:"caller_name"`
	CallerType string `json:"caller_type"`
	ErrorCode  string `json:"error_code"`
}

// LookupPhoneNumberCarrierResponse is the optional carrier information part
// of the response.
type LookupPhoneNumberCarrierResponse struct {
	MobileCountryCode string `json:"mobile_country_code"`
	MobileNetworkCode string `json:"mobile_network_code"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	ErrorCode         string `json:"error_code"`
}

// LookupPhoneNumberResponse is the response for requests to lookup phone
// numbers.
type LookupPhoneNumberResponse struct {
	CountryCode    string                              `json:"country_code"`
	PhoneNumber    string                              `json:"phone_number"`
	NationalFormat string                              `json:"national_format"`
	CallerName     LookupPhoneNumberCallerNameResponse `json:"caller_name"`
	Carrier        LookupPhoneNumberCarrierResponse    `json:"carrier"`
	URL            string                              `json:"url"`
}

// WithCountryCode is optional. It adds the ISO country code of the phone
// number. This is used to specify the country when the number is provided in a
// national format.
func (b *LookupPhoneNumberBuilder) WithCountryCode(
	countryCode string,
) *LookupPhoneNumberBuilder {
	b.countryCode = countryCode
	return b
}

// IncludeCarrierInResponse indicates that carrier information should be returned with the
// request. Extra charges may apply.
func (b *LookupPhoneNumberBuilder) IncludeCarrierInResponse() *LookupPhoneNumberBuilder {
	b.includeCarrierInResponse = true
	return b
}

// IncludeCallerNameInResponse indicates that caller name information should be returned with
// the request. Extra charges may apply.
func (b *LookupPhoneNumberBuilder) IncludeCallerNameInResponse() *LookupPhoneNumberBuilder {
	b.includeCallerNameInResponse = true
	return b
}

// Build will build the request.
func (b *LookupPhoneNumberBuilder) Build() (*http.Request, error) {
	requestURL, err := url.Parse(
		b.opts.LookupBaseURL + "/v1/PhoneNumbers/" + url.PathEscape(b.phoneNumber))
	if err != nil {
		return nil, err
	}

	q := requestURL.Query()

	if b.countryCode != "" {
		q.Set("CountryCode", b.countryCode)
	}

	if b.includeCallerNameInResponse {
		q.Add("Type", "caller-name")
	}

	if b.includeCarrierInResponse {
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
func (b *LookupPhoneNumberBuilder) Do() (*LookupPhoneNumberResponse, error) {
	req, err := b.Build()
	if err != nil {
		return nil, err
	}

	var response LookupPhoneNumberResponse
	err = do(b.opts, req, true, http.StatusOK, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
