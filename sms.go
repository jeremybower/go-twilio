package twilio

// SMS is a group of APIs related to Twilio's SMS service.
type SMS interface {
	SendMessage(from, to, body string) *SMSSendMessageBuilder
}

type smsImpl struct {
	opts *Options
}

// Lookup is a group of APIs related to Twilio's phone number lookup service.
func (impl *clientImpl) SMS() SMS {
	return &smsImpl{
		opts: impl.opts,
	}
}

// PhoneNumber creates a builder to lookup a phone number.
func (impl *smsImpl) SendMessage(from, to, body string) *SMSSendMessageBuilder {
	return &SMSSendMessageBuilder{
		opts: impl.opts,
		from: from,
		to:   to,
		body: body,
	}
}
