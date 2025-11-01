package bindings

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Spec struct {
	address url.URL
	spec    *openapi3.T
}

func NewSpec(cfg config.Config, router *echo.Echo) (Spec, error) {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return Spec{}, fault.Wrap(err)
	}

	spec.Servers[0].URL = cfg.PublicAPIAddress.JoinPath("/api").String()

	route := Spec{
		address: cfg.PublicAPIAddress,
		spec:    spec,
	}

	router.Use(echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/openapi.json" {
				route.getSpecOverride(w)
				return
			}
			h.ServeHTTP(w, r)
		})
	}))

	return route, nil
}

func (v *Spec) getSpecOverride(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/vnd.oai.openapi+json;version=3.1.0")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(v.spec)
}

// NOTE: Unused, overridden by middleware.
func (v *Spec) GetSpec(context.Context, openapi.GetSpecRequestObject) (openapi.GetSpecResponseObject, error) {
	return nil, nil
}

var docsTemplate = template.Must(template.New("docs").Parse(`<!doctype html>
<html>
  <head>
    <title>Storyden API documentation</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url={{ .SchemaURL }}
	></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`))

func (v *Spec) GetDocs(ctx context.Context, request openapi.GetDocsRequestObject) (openapi.GetDocsResponseObject, error) {
	schemaURL := v.address.JoinPath("/api/openapi.json")

	b := bytes.NewBuffer(nil)

	docsTemplate.Execute(b, map[string]string{
		"SchemaURL": schemaURL.String(),
	})

	return openapi.GetDocs200TexthtmlResponse{
		Body:          b,
		ContentLength: int64(b.Len()),
	}, nil
}
