package glue

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/mileusna/useragent"
)

var uaKey = "ua_context"

// UserAgentContext stores in the request context the user agent info.
func UserAgentContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		ctx := r.Context()

		ua := useragent.Parse(r.Header.Get("User-Agent"))

		newctx := context.WithValue(ctx, uaKey, ua)

		c.SetRequest(r.WithContext(newctx))

		return next(c)
	}
}

func GetDeviceName(ctx context.Context) string {
	v := ctx.Value(uaKey)
	ua, ok := v.(useragent.UserAgent)
	if !ok {
		return "Unknown"
	}

	return fmt.Sprintf("%s (%s)", ua.Name, ua.OS)
}
