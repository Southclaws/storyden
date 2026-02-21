package bindings

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_manager"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Plugins struct {
	pm *plugin_manager.Manager
}

func NewPlugins(
	pm *plugin_manager.Manager,
	pl *plugin_logger.Reader,
	router *echo.Echo,
) Plugins {
	// The generated OpenAPI code does not expose the underlying ResponseWriter
	// which we need for streaming Q&A responses for that ✨chatgpt✨ effect.
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()

			if !strings.HasPrefix(path, "/api/plugins/") {
				return next(c)
			}

			if !strings.HasSuffix(path, "/logs") {
				return next(c)
			}

			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			ctx := c.Request().Context()

			pluginID := plugin.InstallationID(deserialiseID(c.Param("plugin_instance_id")))

			return streamPluginLogs(ctx, c, pm, pl, pluginID)
		}
	})

	return Plugins{
		pm: pm,
	}
}

func (p *Plugins) PluginList(ctx context.Context, request openapi.PluginListRequestObject) (openapi.PluginListResponseObject, error) {
	rs, err := p.pm.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginList200JSONResponse{
		PluginListOKJSONResponse: openapi.PluginListOKJSONResponse{
			Plugins: serialisePlugins(rs),
		},
	}, nil
}

func (p *Plugins) PluginAdd(ctx context.Context, request openapi.PluginAddRequestObject) (openapi.PluginAddResponseObject, error) {
	if request.JSONBody == nil {
		pl, err := p.pm.AddFromFile(ctx, request.Body)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return openapi.PluginAdd200JSONResponse{
			PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(&pl.Record)),
		}, nil
	}

	typ, err := request.JSONBody.Discriminator()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	switch typ {
	case string(openapi.Supervised):
		supervised, err := request.JSONBody.AsPluginInitialSupervised()
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		if strings.TrimSpace(supervised.Url) == "" {
			return nil, fault.Wrap(
				fault.New("url is required"),
				fctx.With(ctx),
				ftag.With(ftag.InvalidArgument),
			)
		}

		u, err := url.Parse(supervised.Url)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		pl, err := p.pm.AddFromURL(ctx, *u)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return openapi.PluginAdd200JSONResponse{
			PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(&pl.Record)),
		}, nil

	case string(openapi.External):
		external, err := request.JSONBody.AsPluginInitialExternal()
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		manifest, err := rpc.ManifestFromMap(external.Manifest)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		rec, _, err := p.pm.AddExternal(ctx, *manifest)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return openapi.PluginAdd200JSONResponse{
			PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(rec)),
		}, nil

	default:
		return nil, fault.Wrap(
			fault.Newf("unknown plugin mode: %s", typ),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}
}

func (p *Plugins) PluginDelete(ctx context.Context, request openapi.PluginDeleteRequestObject) (openapi.PluginDeleteResponseObject, error) {
	err := p.pm.Delete(ctx, plugin.InstallationID(deserialiseID(request.PluginInstanceId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginDelete204Response{}, nil
}

func (p *Plugins) PluginGet(ctx context.Context, request openapi.PluginGetRequestObject) (openapi.PluginGetResponseObject, error) {
	record, err := p.pm.Get(ctx, plugin.InstallationID(deserialiseID(request.PluginInstanceId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginGet200JSONResponse{
		PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(record)),
	}, nil
}

func (p *Plugins) PluginSetActiveState(ctx context.Context, request openapi.PluginSetActiveStateRequestObject) (openapi.PluginSetActiveStateResponseObject, error) {
	id := plugin.InstallationID(deserialiseID(request.PluginInstanceId))
	status, err := plugin.NewActiveState(string(request.Body.Active))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	err = p.pm.SetActiveState(ctx, id, status)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	record, err := p.pm.Get(ctx, plugin.InstallationID(deserialiseID(request.PluginInstanceId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginSetActiveState200JSONResponse{
		PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(record)),
	}, nil
}

func (p *Plugins) PluginGetLogs(ctx context.Context, request openapi.PluginGetLogsRequestObject) (openapi.PluginGetLogsResponseObject, error) {
	return nil, nil
}

func (p *Plugins) PluginUpdateManifest(ctx context.Context, request openapi.PluginUpdateManifestRequestObject) (openapi.PluginUpdateManifestResponseObject, error) {
	id := plugin.InstallationID(deserialiseID(request.PluginInstanceId))

	manifest, err := rpc.ManifestFromMap(*request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	rec, err := p.pm.UpdateManifest(ctx, id, *manifest)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginUpdateManifest200JSONResponse{
		PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(rec)),
	}, nil
}

func (p *Plugins) PluginUpdatePackage(ctx context.Context, request openapi.PluginUpdatePackageRequestObject) (openapi.PluginUpdatePackageResponseObject, error) {
	id := plugin.InstallationID(deserialiseID(request.PluginInstanceId))

	rec, err := p.pm.UpdatePackage(ctx, id, request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginUpdatePackage200JSONResponse{
		PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(rec)),
	}, nil
}

func (p *Plugins) PluginCycleToken(ctx context.Context, request openapi.PluginCycleTokenRequestObject) (openapi.PluginCycleTokenResponseObject, error) {
	id := plugin.InstallationID(deserialiseID(request.PluginInstanceId))

	token, err := p.pm.CycleExternalToken(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginCycleToken200JSONResponse{
		PluginCycleTokenOKJSONResponse: openapi.PluginCycleTokenOKJSONResponse{
			Token: token,
		},
	}, nil
}

func (p *Plugins) PluginGetConfigurationSchema(ctx context.Context, request openapi.PluginGetConfigurationSchemaRequestObject) (openapi.PluginGetConfigurationSchemaResponseObject, error) {
	id := plugin.InstallationID(deserialiseID(request.PluginInstanceId))

	schema, err := p.pm.GetConfigurationSchema(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginGetConfigurationSchema200JSONResponse{
		PluginGetConfigurationSchemaOKJSONResponse: openapi.PluginGetConfigurationSchemaOKJSONResponse(
			serialisePluginConfigurationSchema(schema),
		),
	}, nil
}

func (p *Plugins) PluginGetConfiguration(ctx context.Context, request openapi.PluginGetConfigurationRequestObject) (openapi.PluginGetConfigurationResponseObject, error) {
	id := plugin.InstallationID(deserialiseID(request.PluginInstanceId))

	cfg, err := p.pm.GetConfiguration(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pluginGetConfiguration200JSONResponse(cfg), nil
}

func (p *Plugins) PluginUpdateConfiguration(ctx context.Context, request openapi.PluginUpdateConfigurationRequestObject) (openapi.PluginUpdateConfigurationResponseObject, error) {
	id := plugin.InstallationID(deserialiseID(request.PluginInstanceId))

	cfg, err := p.pm.UpdateConfiguration(ctx, id, map[string]any(*request.Body))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pluginUpdateConfiguration200JSONResponse(cfg), nil
}

// NOTE: There's a bug in oapi-codegen for responses that are top-level any maps
// so we have to wrap these in a new type that implements the response manually.

type pluginGetConfiguration200JSONResponse map[string]any

func (response pluginGetConfiguration200JSONResponse) VisitPluginGetConfigurationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(map[string]any(response))
}

type pluginUpdateConfiguration200JSONResponse map[string]any

func (response pluginUpdateConfiguration200JSONResponse) VisitPluginUpdateConfigurationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(map[string]any(response))
}

func serialisePlugin(in *plugin.Record) openapi.Plugin {
	out := openapi.Plugin{
		Id:          openapi.Identifier(xid.ID(in.InstallationID).String()),
		AddedAt:     in.Created,
		Status:      serialisePluginStatus(in),
		Manifest:    in.Manifest.ToMap(),
		Name:        in.Manifest.Name,
		Description: &in.Manifest.Description,
		Version:     &in.Manifest.Version,
	}

	if in.Mode.Supervised() {
		_ = out.Connection.FromPluginSupervisedProps(openapi.PluginSupervisedProps{
			Mode: openapi.Supervised,
		})
		return out
	}

	_ = out.Connection.FromPluginExternalProps(openapi.PluginExternalProps{
		Mode:  openapi.External,
		Token: in.Token,
	})
	return out
}

func serialisePlugins(in []*plugin.Record) []openapi.Plugin {
	return dt.Map(in, serialisePlugin)
}

func serialisePluginStatus(in *plugin.Record) openapi.PluginStatus {
	as := openapi.PluginStatus{}

	switch in.ReportedState {
	case plugin.ReportedStateActive:
		as.FromPluginStatusActive(openapi.PluginStatusActive{
			ActiveState: openapi.PluginStatusActiveActiveStateActive,
			ActivatedAt: in.StartedAt,
		})

	case plugin.ReportedStateInactive:
		as.FromPluginStatusInactive(openapi.PluginStatusInactive{
			ActiveState:   openapi.Inactive,
			DeactivatedAt: in.StateChangedAt,
		})

	case plugin.ReportedStateStarting:
		as.FromPluginStatusStarting(openapi.PluginStatusStarting{
			ActiveState: openapi.Starting,
			StartingAt:  in.StateChangedAt,
		})

	case plugin.ReportedStateConnecting:
		as.FromPluginStatusConnecting(openapi.PluginStatusConnecting{
			ActiveState:  openapi.Connecting,
			ConnectingAt: in.StateChangedAt,
		})

	case plugin.ReportedStateError:
		as.FromPluginStatusError(openapi.PluginStatusError{
			ActiveState: openapi.Error,
			Message:     in.StatusMessage,
			Details:     in.Details,
		})

	case plugin.ReportedStateRestarting:
		as.FromPluginStatusRestarting(openapi.PluginStatusRestarting{
			ActiveState: openapi.Restarting,
			Message:     in.StatusMessage,
			Details:     in.Details,
		})

	default:
		as.FromPluginStatusInactive(openapi.PluginStatusInactive{
			ActiveState:   openapi.Inactive,
			DeactivatedAt: in.StateChangedAt,
		})
	}

	return as
}

func serialisePluginConfigurationSchema(in rpc.ManifestConfigurationSchema) openapi.PluginConfigurationSchema {
	fields := dt.Map(in.Fields, serialisePluginConfigurationField)
	return openapi.PluginConfigurationSchema{
		Fields: &fields,
	}
}

func serialisePluginConfigurationField(in rpc.PluginConfigurationFieldSchema) openapi.PluginConfigurationFieldUnion {
	out := openapi.PluginConfigurationFieldUnion{}

	switch v := in.PluginConfigurationFieldUnion.(type) {
	case *rpc.PluginConfigurationFieldString:
		body := openapi.PluginConfigurationFieldString{
			Type: openapi.PluginConfigurationFieldStringType("string"),
		}
		if defaultValue, ok := v.Default.Get(); ok {
			body.Default = &defaultValue
		}
		_ = out.FromPluginConfigurationFieldString(body)
		out.Description = &v.Description
		out.Id = &v.ID
		out.Label = &v.Label

	case *rpc.PluginConfigurationFieldNumber:
		body := openapi.PluginConfigurationFieldNumber{
			Type: openapi.PluginConfigurationFieldNumberType("number"),
		}
		if defaultValue, ok := v.Default.Get(); ok {
			d := float32(defaultValue)
			body.Default = &d
		}
		_ = out.FromPluginConfigurationFieldNumber(body)
		out.Description = &v.Description
		out.Id = &v.ID
		out.Label = &v.Label

	case *rpc.PluginConfigurationFieldBoolean:
		body := openapi.PluginConfigurationFieldBoolean{
			Type: openapi.PluginConfigurationFieldBooleanType("boolean"),
		}
		if defaultValue, ok := v.Default.Get(); ok {
			body.Default = &defaultValue
		}
		_ = out.FromPluginConfigurationFieldBoolean(body)
		out.Description = &v.Description
		out.Id = &v.ID
		out.Label = &v.Label

	case rpc.PluginConfigurationFieldString:
		return serialisePluginConfigurationField(rpc.PluginConfigurationFieldSchema{
			PluginConfigurationFieldUnion: &v,
		})

	case rpc.PluginConfigurationFieldNumber:
		return serialisePluginConfigurationField(rpc.PluginConfigurationFieldSchema{
			PluginConfigurationFieldUnion: &v,
		})

	case rpc.PluginConfigurationFieldBoolean:
		return serialisePluginConfigurationField(rpc.PluginConfigurationFieldSchema{
			PluginConfigurationFieldUnion: &v,
		})
	}

	return out
}

func streamPluginLogs(
	ctx context.Context,
	c echo.Context,
	pm *plugin_manager.Manager,
	pr *plugin_logger.Reader,
	pluginID plugin.InstallationID,
) error {
	rec, err := pm.Get(ctx, pluginID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if rec.Mode.External() {
		return fault.Wrap(
			fault.New("external plugins do not have host-managed logs"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	w := c.Response().Writer
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		if innerWriter := unwrapWriter(w); innerWriter != nil {
			flusher, _ = innerWriter.(http.Flusher)
		}
	}

	stream, err := pr.StreamPluginLogs(ctx, pluginID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-stream.Done:
			if _, err := w.Write([]byte("event: end\n\n")); err != nil {
				return err
			}
			if flusher != nil {
				flusher.Flush()
			}
			return nil
		case line, ok := <-stream.Lines:
			if !ok {
				if _, err := w.Write([]byte("event: end\n\n")); err != nil {
					return err
				}
				if flusher != nil {
					flusher.Flush()
				}
				return nil
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", line); err != nil {
				return err
			}
			if flusher != nil {
				flusher.Flush()
			}
		}
	}
}
