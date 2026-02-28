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
	OperationID string
	UserAgent   useragent.UserAgent
	CacheQuery  cachecontrol.Query
	ClientAddr  string
}

type infoKey struct{}

func WithRequestInfo(ctx context.Context, r *http.Request, opid string, clientAddr string) context.Context {
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
		OperationID: opid,
		UserAgent:   ua,
		CacheQuery:  cachecontrol.NewQuery(ifNoneMatch, ifModifiedSince),
		ClientAddr:  clientAddr,
	}

	return context.WithValue(ctx, infoKey{}, info)
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

func getInfo(ctx context.Context) Info {
	v := ctx.Value(infoKey{})
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
