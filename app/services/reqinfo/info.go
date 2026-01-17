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
}

type infoKey struct{}

func WithRequestInfo(ctx context.Context, r *http.Request, opid string) context.Context {
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
	}

	return context.WithValue(ctx, infoKey{}, info)
}

func GetOperationID(ctx context.Context) string {
	v := ctx.Value(infoKey{})
	i, ok := v.(Info)
	if !ok {
		return "unknown"
	}

	return i.OperationID
}

func GetDeviceName(ctx context.Context) string {
	v := ctx.Value(infoKey{})
	i, ok := v.(Info)
	if !ok {
		return "Unknown"
	}

	ua := i.UserAgent

	return fmt.Sprintf("%s (%s)", ua.Name, ua.OS)
}

func GetCacheQuery(ctx context.Context) cachecontrol.Query {
	v := ctx.Value(infoKey{})
	i, ok := v.(Info)
	if !ok {
		return cachecontrol.Query{}
	}

	return i.CacheQuery
}

func notEmpty(s string) bool {
	return s != ""
}

func parseConditionalRequestTime(in string) (time.Time, error) {
	return time.Parse(time.RFC1123, in)
}
