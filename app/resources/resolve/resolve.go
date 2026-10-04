package resolve

import "net/url"

const PathPrefix = "/_/resolve"

func Path(resource, identifier string) string {
	base := url.URL{Path: PathPrefix}
	return base.JoinPath(url.PathEscape(resource), url.PathEscape(identifier)).Path
}

func URL(webAddress url.URL, resource, identifier string) *url.URL {
	return webAddress.JoinPath(PathPrefix, url.PathEscape(resource), url.PathEscape(identifier))
}
