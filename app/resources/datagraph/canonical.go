package datagraph

import (
	"net/url"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/resolve"
)

// CanonicalResolvePathPrefix is the frontend route prefix that vendors are
// expected to implement in order to resolve a datagraph item to its canonical
// UI location. The API never knows the concrete frontend route for a given
// resource kind, so instead of emitting frontend paths directly it emits URLs
// under this prefix which the frontend redirects to the correct page.
const CanonicalResolvePathPrefix = resolve.PathPrefix

// CanonicalResolvePath builds the frontend-relative path that resolves to the
// given datagraph item. The frontend redirects this to the appropriate UI route
// for the resource kind, for example "/_/resolve/thread/foo" -> "/t/foo".
func CanonicalResolvePath(kind Kind, mark string) string {
	return resolve.Path(kind.String(), mark)
}

// CanonicalResolveURL builds an absolute URL to the frontend resolve route for
// the given datagraph item, using the provided public web address as the base.
// This is the URL safe to share in emails, MCP responses, CLI output, etc.
func CanonicalResolveURL(webAddress url.URL, kind Kind, mark string) *url.URL {
	return resolve.URL(webAddress, kind.String(), mark)
}

// CanonicalResolvePath builds the frontend-relative path that resolves to the
// given datagraph item. The frontend redirects this to the appropriate UI route
// for the resource kind, for example "/_/resolve/thread/foo" -> "/t/foo".
func CanonicalResolveMarkPath(kind Kind, mark mark.Mark) string {
	return CanonicalResolvePath(kind, mark.String())
}

// CanonicalResolveURL builds an absolute URL to the frontend resolve route for
// the given datagraph item, using the provided public web address as the base.
// This is the URL safe to share in emails, MCP responses, CLI output, etc.
func CanonicalResolveMarkURL(webAddress url.URL, kind Kind, mark mark.Mark) *url.URL {
	return CanonicalResolveURL(webAddress, kind, mark.String())
}
