package oauth

import (
	"fmt"
	"net/url"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
)

type Configuration struct {
	Enabled      bool
	ClientID     string
	ClientSecret string
}

func Redirect(publicWebURL url.URL, svc authentication.Service) url.URL {
	name := svc.String()

	// TODO: Let the client/caller control this callback URL path.
	ref, _ := url.Parse(fmt.Sprintf("/auth/%s/callback", name))

	return *publicWebURL.ResolveReference(ref)
}
