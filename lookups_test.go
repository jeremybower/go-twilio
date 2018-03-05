package twilio

// WithHost changes the default host. This is useful for testing.
func (lookups *Lookups) WithHost(
	host string,
) *Lookups {
	lookups.host = host
	return lookups
}
