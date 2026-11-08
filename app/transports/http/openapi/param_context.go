package openapi

import (
	"github.com/Southclaws/fault/fctx"
	"github.com/labstack/echo/v5"
)

// ParameterContext is a simple middleware for injecting request metadata into a
// context object for use with the fctx library. This makes diagnostics easy.
func ParameterContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		r := c.Request()
		ctx := r.Context()

		meta := []string{}
		for _, param := range c.PathValues() {
			meta = append(meta, param.Name, param.Value)
		}

		c.SetRequest(r.WithContext(fctx.WithMeta(ctx, meta...)))

		return next(c)
	}
}
