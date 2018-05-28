package twilio

// Lookup is a group of APIs related to Twilio's phone number lookup service.
type Lookup interface {
	PhoneNumber(phoneNumber string) *LookupPhoneNumberBuilder
}

type lookupImpl struct {
	opts *Options
}

// Lookup is a group of APIs related to Twilio's phone number lookup service.
func (impl *clientImpl) Lookup() Lookup {
	return &lookupImpl{
		opts: impl.opts,
	}
}

// PhoneNumber creates a builder to lookup a phone number.
func (impl *lookupImpl) PhoneNumber(phoneNumber string) *LookupPhoneNumberBuilder {
	return &LookupPhoneNumberBuilder{
		opts:        impl.opts,
		phoneNumber: phoneNumber,
	}
}
