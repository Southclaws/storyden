package link_url

import (
	"encoding/json"
	"net"
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
	"golang.org/x/net/idna"
)

// URL is a normalized URL accepted by the Link policy. Its zero value is invalid.
// Network access still requires DNS and address checks through httpsafe.
type URL struct {
	value url.URL
}

var ErrEmpty = invalid("A link URL is required.")

func (u URL) String() string {
	return u.value.String()
}

func (u URL) Hostname() string {
	return u.value.Hostname()
}

func (u URL) IsZero() bool {
	return u.value.Host == ""
}

// NetURL returns a copy so callers cannot mutate the validated value.
func (u URL) NetURL() url.URL {
	return u.value
}

func (u URL) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		return nil, ErrEmpty
	}

	return json.Marshal(u.String())
}

func (u *URL) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return invalid("A link URL must be a string.")
	}

	parsed, err := Parse(raw)
	if err != nil {
		return err
	}

	*u = parsed
	return nil
}

// Parse accepts portless HTTP(S) links to fully qualified DNS hostnames.
// Paths, queries, fragments, and www remain significant.
func Parse(raw string) (URL, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return URL{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("invalid link URL", "Enter a valid HTTP or HTTPS URL."))
	}

	u.Scheme = strings.ToLower(u.Scheme)
	if u.Scheme != "https" && u.Scheme != "http" {
		return URL{}, invalid("Only HTTP and HTTPS URLs are supported.")
	}

	if u.Opaque != "" || u.User != nil || u.Hostname() == "" {
		return URL{}, invalid("The URL must have a hostname and cannot contain credentials.")
	}

	if strings.ContainsAny(u.Host, ":[]") {
		return URL{}, invalid("Link URLs cannot contain ports or IP addresses.")
	}

	host, err := idna.Lookup.ToASCII(u.Host)
	if err != nil {
		return URL{}, invalid("The URL must have a valid DNS hostname.")
	}

	host, err = idna.Registration.ToASCII(strings.TrimSuffix(host, "."))
	if err != nil || host == "" {
		return URL{}, invalid("The URL must have a valid DNS hostname.")
	}

	if net.ParseIP(host) != nil || httpsafe.IsDisallowedHost(host) {
		return URL{}, invalid("The URL must have a public DNS hostname.")
	}

	labels := strings.Split(host, ".")
	if len(labels) < 2 || labels[len(labels)-1] == "" {
		return URL{}, invalid("The URL must have a fully qualified DNS hostname.")
	}

	tld := labels[len(labels)-1]
	if tld[0] < 'a' || tld[0] > 'z' {
		return URL{}, invalid("The URL must have a DNS hostname, not a numeric address.")
	}

	u.Host = host
	if u.Path == "" {
		u.Path = "/"
	}

	return URL{value: *u}, nil
}

func invalid(message string) error {
	return fault.New(message, ftag.With(ftag.InvalidArgument), fmsg.WithDesc("invalid link URL", message))
}
