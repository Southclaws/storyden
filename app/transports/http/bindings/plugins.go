package bindings

import (
	"context"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_manager"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	lib_plugin "github.com/Southclaws/storyden/lib/plugin"
)

type Plugins struct {
	pm *plugin_manager.Manager
}

func NewPlugins(
	pm *plugin_manager.Manager,
) Plugins {
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
	var pl *plugin.Available

	if request.JSONBody != nil {
		u, err := url.Parse(request.JSONBody.Url)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		pl, err = p.pm.AddFromURL(ctx, *u)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	} else {
		var err error
		pl, err = p.pm.AddFromFile(ctx, request.Body)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return openapi.PluginAdd200JSONResponse{
		PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(&pl.Record)),
	}, nil
}

func (p *Plugins) PluginDelete(ctx context.Context, request openapi.PluginDeleteRequestObject) (openapi.PluginDeleteResponseObject, error) {
	err := p.pm.Delete(ctx, plugin.ID(deserialiseID(request.PluginInstanceId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginDelete204Response{}, nil
}

func (p *Plugins) PluginGet(ctx context.Context, request openapi.PluginGetRequestObject) (openapi.PluginGetResponseObject, error) {
	record, err := p.pm.Get(ctx, plugin.ID(deserialiseID(request.PluginInstanceId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PluginGet200JSONResponse{
		PluginGetOKJSONResponse: openapi.PluginGetOKJSONResponse(serialisePlugin(record)),
	}, nil
}

func (p *Plugins) PluginSetActiveState(ctx context.Context, request openapi.PluginSetActiveStateRequestObject) (openapi.PluginSetActiveStateResponseObject, error) {
	// Call plugin runner, which will handle the active state change
	// set active state in db
	// re-build runner state
	// register plugin events etc.
	// set up for run once or run in background
	return nil, nil
}

func serialisePlugin(in *plugin.Record) openapi.Plugin {
	return openapi.Plugin{
		Id:       openapi.Identifier(xid.ID(in.ID).String()),
		AddedAt:  in.Created,
		Status:   serialisePluginStatus(in),
		Manifest: serialisePluginManifest(in.Manifest),
	}
}

func serialisePlugins(in []*plugin.Record) []openapi.Plugin {
	return dt.Map(in, serialisePlugin)
}

func serialisePluginStatus(in *plugin.Record) openapi.PluginStatus {
	as := openapi.PluginStatus{}

	switch in.State {
	case plugin.ActiveStateActive:
		as.FromPluginStatusActive(openapi.PluginStatusActive{
			ActiveState: openapi.PluginStatusActiveActiveStateActive,
			ActivatedAt: in.StateChangedAt,
		})

	case plugin.ActiveStateInactive:
		as.FromPluginStatusInactive(openapi.PluginStatusInactive{
			ActiveState:   openapi.Inactive,
			DeactivatedAt: in.StateChangedAt,
		})

	case plugin.ActiveStateError:
		as.FromPluginStatusError(openapi.PluginStatusError{
			ActiveState: openapi.Error,
			Message:     in.StatusMessage,
			Details:     in.Details,
		})

	default:
		as.FromPluginStatusError(openapi.PluginStatusError{
			ActiveState: openapi.Error,
			Message:     in.StatusMessage,
			Details:     in.Details,
		})
	}

	return as
}

func serialisePluginManifest(in lib_plugin.Manifest) openapi.PluginManifest {
	return openapi.PluginManifest{
		Name:    in.Name.String(),
		Version: in.Version.String(),
	}
}
