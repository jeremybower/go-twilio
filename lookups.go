package twilio

// Lookups is a group of APIs related to the lookups.twilio.com service
type Lookups struct {
	client *Client
	host   string
}

// Lookups is a group of APIs related to the lookups.twilio.com service.
func (client *Client) Lookups() *Lookups {
	return &Lookups{
		client: client,
		host:   "https://lookups.twilio.com",
	}
}

// PhoneNumbers creates a builder to lookup a phone number.
func (lookups *Lookups) PhoneNumbers(phoneNumber string) *LookupsPhoneNumbersBuilder {
	return &LookupsPhoneNumbersBuilder{
		lookups:     lookups,
		phoneNumber: phoneNumber,
	}
}
