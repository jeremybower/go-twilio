package twilio

// CountryCodeNone can be used with the country code is optional.
const CountryCodeNone = ""

// Client is the Twilio client.
type Client interface {
	LookupPhoneNumber(
		phoneNumber string,
		countryCode string,
		includeCarrierInResponse bool,
		includeCallerNameInResponse bool,
	) (*LookupPhoneNumberResponse, error)

	SendSMSMessage(
		from string,
		to string,
		body string,
	) (*SMSSendMessageResponse, error)
}
