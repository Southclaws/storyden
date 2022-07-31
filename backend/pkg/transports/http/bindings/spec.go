package bindings

import (
	"github.com/labstack/echo/v4"

	"github.com/Southclaws/storyden/backend/pkg/transports/http/openapi"
)

func spec(c echo.Context) error {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return err
	}
	c.JSON(200, spec)
	return nil
}
