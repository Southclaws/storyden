package glue

import (
	"github.com/labstack/echo/v4"

	"github.com/Southclaws/fault/errctx"
)

// ParameterContext is a simple middleware for injecting request metadata into a
// context object for use with the errctx library. This makes diagnostics easy.
func ParameterContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		ctx := r.Context()

		meta := []string{}
		for _, k := range c.ParamNames() {
			meta = append(meta, k, c.Param(k))
		}

		c.SetRequest(r.WithContext(errctx.WithMeta(ctx, meta...)))

		return next(c)
	}
}
