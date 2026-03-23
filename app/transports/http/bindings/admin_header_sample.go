package bindings

import (
	"context"
	"net/http"
	"strings"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

const (
	adminSettingsPath       = "/api/admin"
	ssrRequestHeaderName    = "x-storyden-ssr"
	xForwardedForHeaderName = "x-forwarded-for"
)

var sampledNetworkHeaderWhitelist = map[string]struct{}{
	xForwardedForHeaderName: {},
}

var sensitiveHeaderSampleDenylist = map[string]struct{}{
	"cookie":              {},
	"set-cookie":          {},
	"authorization":       {},
	"proxy-authorization": {},
}

type networkHeadersSampleContextKey struct{}

func getNetworkHeadersSample(ctx context.Context) *openapi.NetworkHeadersSample {
	v := ctx.Value(networkHeadersSampleContextKey{})
	sample, ok := v.(*openapi.NetworkHeadersSample)
	if !ok {
		return nil
	}
	return sample
}

func buildNetworkHeadersSample(r *http.Request) *openapi.NetworkHeadersSample {
	direct := sampleWhitelistedHeaders(r.Header)
	ssr := sampleSSRHeaders(r.Header)
	rawClientAddress := strings.TrimSpace(r.RemoteAddr)

	if len(direct) == 0 && len(ssr) == 0 && rawClientAddress == "" {
		return nil
	}

	out := &openapi.NetworkHeadersSample{}
	if len(direct) > 0 {
		out.Headers = &direct
	}
	if len(ssr) > 0 {
		out.HeadersSsr = &ssr
	}
	if rawClientAddress != "" {
		out.RawClientAddress = &rawClientAddress
	}

	return out
}

func sampleWhitelistedHeaders(headers http.Header) map[string]string {
	out := make(map[string]string)
	for name, values := range headers {
		lower := strings.ToLower(strings.TrimSpace(name))
		if lower == "" {
			continue
		}
		if _, denied := sensitiveHeaderSampleDenylist[lower]; denied {
			continue
		}
		if _, ok := sampledNetworkHeaderWhitelist[lower]; !ok {
			continue
		}

		value := joinHeaderValues(values)
		if value == "" {
			continue
		}
		out[lower] = value
	}
	return out
}

func sampleSSRHeaders(headers http.Header) map[string]string {
	if strings.TrimSpace(headers.Get(ssrRequestHeaderName)) == "" {
		return map[string]string{}
	}

	return sampleWhitelistedHeaders(headers)
}

func joinHeaderValues(values []string) string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		v := strings.TrimSpace(value)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return strings.Join(out, ", ")
}
