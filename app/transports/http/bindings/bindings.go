// Package bindings is responsible for providing a bridge between the code that
// is generated from the OpenAPI specification that comprises the "Transport"
// layer and the code that implements the product, or the "Service" layer.
//
// This package is structured as a set of structs which are compose together in
// the struct called `Bindings` which implements an interface generated by the
// oapi-codegen tool which describes all the endpoints.
//
// Aside from `bindings.go`, pretty much every Go file is named after a REST
// collection from the OpenAPI specification.
//
// ## Adding Routes
//
// To add a new route, you first modify the OpenAPI specification YAML document
// and then run the `generate` task which generates all the handler declarations
// and types necessary to implement the binding.
//
// Next, you create a file in this package named after the collection. So if you
// added `/things/{id}` you'd create `things.go` and inside that file, a struct
// named `Things` and a constructor named `NewThings`. This pattern may not need
// to apply for certain cases but it's generally best to try to follow for most.
//
// You then add your struct to the `Bindings` composed struct and provide the
// implementation of your struct to the DI system using `bindingsProviders`.
//
// ## Changing Routes
//
// Updating a route is as simple as just modifying the OpenAPI specification
// and making the necessary changes to the bindings to get the code compiling.
package bindings

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	oapi_middleware "github.com/oapi-codegen/echo-middleware"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

// Everything in this package is mounted under this path. Any handlers outside
// of this path are not covered by the Storyden OpenAPI specification.
const apiPathPrefix = "/api"

// Bindings is a DI parameter struct that is used to compose together all of the
// individual service bindings in this package. When the provider below depends
// on this type, it provides all these composed bindings to the DI system so the
// invoke call can mount them onto the router using the `StrictServerInterface`.
//
// The reason this is done this way is so we split code up based on OpenAPI
// REST collections instead of bundling everything into one huge struct with
// loads of dependencies. This is just how the oapi-codegen tool works, by
// generating one big interface which the bindings layer must satisfy.
type Bindings struct {
	fx.In
	Version
	Spec
	Info
	Admin
	Roles
	Authentication
	WebAuthn
	PhoneAuth
	Accounts
	Invitations
	Notifications
	Profiles
	Categories
	Tags
	Posts
	Threads
	Replies
	Reacts
	Assets
	Likes
	Collections
	Nodes
	Links
	Datagraph
	Events
}

// bindingsProviders provides to the application the necessary implementations
// that compose the `Bindings` parameter struct which implements the OpenAPI
// server interface. When you add a new collection, add it to Bindings and here.
func bindingsProviders() fx.Option {
	return fx.Provide(
		NewVersion,
		NewSpec,
		NewInfo,
		NewAdmin,
		NewRoles,
		NewAuthentication,
		NewWebAuthn,
		NewPhoneAuth,
		NewAccounts,
		NewInvitations,
		NewNotifications,
		NewProfiles,
		NewCategories,
		NewTags,
		NewPosts,
		NewThreads,
		NewReplies,
		NewReacts,
		NewAssets,
		NewLikes,
		NewCollections,
		NewNodes,
		NewLinks,
		NewDatagraph,
		NewEvents,
	)
}

// bindings provides to the application the above struct which binds the service
// layer to the transport layer. This uses `Bindings` as an fx parameter struct.
//
// ## WHY AM I GETTING AN ERROR HERE?
//
// When you edit `openapi.yaml` and re-run the code generation task, this will
// most likely change the declaration of `StrictServerInterface` inside the
// generated package `openapi`.
//
// The error you will see is most likely something along the lines of:
//
//	*Bindings does not implement openapi.StrictServerInterface
//
// and the underlying problem is either missing methods or methods that have
// changed signature due to changes to the parameters or request or response.
//
// This API follows RESTful design so a collection in the API specification
// (such as `/accounts`) will map to a file, struct and constructor here (such
// as `accounts.go`, `Accounts` and `NewAccounts`) and everything is glued
// together in this file.
func bindings(s Bindings) openapi.StrictServerInterface {
	return &s
}

// mounts the OpenAPI routes and middleware onto the /api path. Everything that
// is outside of the `/api` path is considered separate from the OpenAPI spec.
func mount(
	logger *slog.Logger,
	router *echo.Echo,
	auth *Authorisation,
	si openapi.StrictServerInterface,
) error {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to get openapi specification"))
	}

	// Skips validation for any paths not prefixed with "/api".
	skipper := func(c echo.Context) bool {
		if !strings.HasPrefix(c.Path(), apiPathPrefix) {
			return true
		}

		// Skip validation for asset upload. This is due to the fact that the
		// OpenAPI validator performs an io.ReadAll on the request body which
		// will cause memory usage issues for large file uploads since it will
		// completely remove the ability to stream request body to the uploader.
		if c.Path() == "/api/assets" && c.Request().Method == http.MethodPost {
			return true
		}

		return false
	}

	requestValidatorMiddleware := oapi_middleware.OapiRequestValidatorWithOptions(spec, &oapi_middleware.Options{
		Skipper: skipper,
		Options: openapi3filter.Options{
			IncludeResponseStatus: true,
			AuthenticationFunc:    auth.validator,
		},
		SilenceServersWarning: true,
		// Handles validation errors that occur BEFORE the handler is called.
		ErrorHandler: openapi.ValidatorErrorHandler(),
	})

	router.GET("/", func(c echo.Context) error {
		return c.Stream(http.StatusOK, "text/html", easteregg())
	})

	openapi.RegisterHandlersWithBaseURL(router, openapi.NewStrictHandler(si, nil), apiPathPrefix)

	router.Use(
		requestValidatorMiddleware,
		openapi.ParameterContext,
	)

	logger.Debug("mounted OpenAPI to service bindings",
		slog.Any("routes", lo.Map(router.Routes(), func(r *echo.Route, _ int) string {
			return r.Path
		})),
	)

	return nil
}

func newRouter(logger *slog.Logger) *echo.Echo {
	router := echo.New()
	router.HTTPErrorHandler = openapi.HTTPErrorHandler(logger)

	return router
}

func Build() fx.Option {
	return fx.Options(
		// Provide the bindings struct which implements the generated OpenAPI
		// interface by composing together all of the service bindings into a
		// single struct.
		fx.Provide(bindings),

		fx.Provide(newAuthorisation),

		// Provide the Echo router.
		fx.Provide(newRouter),

		// Mount the bound OpenAPI routes onto the router.
		fx.Invoke(mount),

		// Provide all service layer bindings to the DI system.
		bindingsProviders(),
	)
}
