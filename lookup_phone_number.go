package twilio

import (
	"net/http"
	"net/url"
)

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

func (client *clientImpl) LookupPhoneNumber(
	phoneNumber string,
	countryCode string,
	includeCarrierInResponse bool,
	includeCallerNameInResponse bool,
) (*LookupPhoneNumberResponse, error) {
	requestURL, err := url.Parse(
		client.opts.LookupBaseURL + "/v1/PhoneNumbers/" + url.PathEscape(phoneNumber))
	if err != nil {
		return nil, err
	}

	q := requestURL.Query()

	if countryCode != CountryCodeNone {
		q.Set("CountryCode", countryCode)
	}

	if includeCallerNameInResponse {
		q.Add("Type", "caller-name")
	}

	if includeCarrierInResponse {
		q.Add("Type", "carrier")
	}

	requestURL.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	var response LookupPhoneNumberResponse
	err = client.do(req, true, http.StatusOK, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
