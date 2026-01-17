package headers

import (
	"net/http"

	"github.com/getkin/kin-openapi/routers/gorillamux"

	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Middleware struct{}

func New() *Middleware {
	return &Middleware{}
}

// WithHeaderContext stores in the request context header info.
func (m *Middleware) WithHeaderContext() func(next http.Handler) http.Handler {
	spec, err := openapi.GetSwagger()
	if err != nil {
		panic(err)
	}

	gr, err := gorillamux.NewRouter(spec)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			opid := "unknown"
			rt, _, err := gr.FindRoute(r)
			if err == nil && rt.Operation != nil {
				opid = rt.Operation.OperationID
			}

			newctx := reqinfo.WithRequestInfo(ctx, r, opid)

			r = r.WithContext(newctx)

			next.ServeHTTP(w, r)
		})
	}
}
