package authentication

//go:generate go run github.com/Southclaws/enumerator

type modeEnum string

// Handle is the default and enables simple username+password signup and login
// flows. This mode is very rudimentary and will make use-cases such as sending
// marketing emails, newsletters and even certain anti-spam measures impossible.
//
// Email enables the use of either email+password or email+verification methods
// which is the most common but requires an email provider to be configured.
//
// Phone enables the use of the phone number+verification method which is best
// suited to mobile-first use-cases and WhatsApp integration. This mode requires
// a phone SMS API provider to be configured for the instance in order to work.
const (
	modeHandle modeEnum = "handle" // Username (default)
	modeEmail  modeEnum = "email"  // Email address
	modePhone  modeEnum = "phone"  // Phone number
)
