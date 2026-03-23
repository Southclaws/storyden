package reqinfo

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mileusna/useragent"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/cachecontrol"
)

type Info struct {
	OperationID   string
	UserAgent     useragent.UserAgent
	CacheQuery    cachecontrol.Query
	ClientAddr    string
	ClientAddrSSR string
}

type infoKey struct{}

var requestInfoContextKey = infoKey{}

func WithRequestInfo(ctx context.Context, r *http.Request, opid string, clientAddr string, clientAddrSSR string) context.Context {
	ua := useragent.Parse(r.Header.Get("User-Agent"))

	ifNoneMatch := opt.NewIf(r.Header.Get("If-None-Match"), notEmpty)

	ifModifiedSince, err := opt.MapErr(
		opt.NewIf(r.Header.Get("If-Modified-Since"), notEmpty),
		parseConditionalRequestTime,
	)
	if err != nil {
		ifModifiedSince = opt.NewEmpty[time.Time]()
	}

	info := Info{
		OperationID:   opid,
		UserAgent:     ua,
		CacheQuery:    cachecontrol.NewQuery(ifNoneMatch, ifModifiedSince),
		ClientAddr:    clientAddr,
		ClientAddrSSR: clientAddrSSR,
	}

	return context.WithValue(ctx, requestInfoContextKey, info)
}

func GetOperationID(ctx context.Context) string {
	i := getInfo(ctx)
	return i.OperationID
}

func GetDeviceName(ctx context.Context) string {
	i := getInfo(ctx)
	ua := i.UserAgent

	return fmt.Sprintf("%s (%s)", ua.Name, ua.OS)
}

func GetCacheQuery(ctx context.Context) cachecontrol.Query {
	i := getInfo(ctx)
	return i.CacheQuery
}

func GetClientAddress(ctx context.Context) string {
	i := getInfo(ctx)
	return i.ClientAddr
}

func GetSSRClientAddress(ctx context.Context) string {
	i := getInfo(ctx)
	return i.ClientAddrSSR
}

func getInfo(ctx context.Context) Info {
	v := ctx.Value(requestInfoContextKey)
	i, ok := v.(Info)
	if !ok {
		panic("reqinfo: request info missing from context; ensure headers.WithHeaderContext is first in middleware chain")
	}

	return i
}

func notEmpty(s string) bool {
	return s != ""
}

func parseConditionalRequestTime(in string) (time.Time, error) {
	return time.Parse(time.RFC1123, in)
}
