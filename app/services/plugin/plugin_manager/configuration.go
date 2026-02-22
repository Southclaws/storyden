package plugin_manager

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

const configureRPCTimeout = 30 * time.Second

func (m *Manager) GetConfigurationSchema(
	ctx context.Context,
	id plugin.InstallationID,
) (rpc.ManifestConfigurationSchema, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return rpc.ManifestConfigurationSchema{}, fault.Wrap(err, fctx.With(ctx))
	}

	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return rpc.ManifestConfigurationSchema{}, fault.Wrap(err, fctx.With(ctx))
	}

	schema, ok := rec.Manifest.ConfigurationSchema.Get()
	if !ok {
		return rpc.ManifestConfigurationSchema{
			Fields: []rpc.PluginConfigurationFieldSchema{},
		}, nil
	}

	if schema.Fields == nil {
		schema.Fields = []rpc.PluginConfigurationFieldSchema{}
	}

	return schema, nil
}

func (m *Manager) GetConfiguration(
	ctx context.Context,
	id plugin.InstallationID,
) (map[string]any, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cfg, err := m.pluginQuerier.GetConfig(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cloneConfig(cfg), nil
}

func (m *Manager) UpdateConfiguration(
	ctx context.Context,
	id plugin.InstallationID,
	config map[string]any,
) (map[string]any, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cfg := cloneConfig(config)

	rec, err := m.pluginQuerier.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := validateManifestConfiguration(rec.Manifest, cfg); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := m.runner.GetSession(ctx, id)
	if err != nil {
		return nil, fault.Wrap(
			fault.New("plugin must be connected to update configuration"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc(
				"plugin session is not connected",
				"The plugin must be connected before its configuration can be updated.",
			),
		)
	}

	if sess.GetReportedState() != plugin.ReportedStateActive {
		return nil, fault.Wrap(
			fault.New("plugin must be connected to update configuration"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc(
				"plugin is not in active state",
				"The plugin must be connected before its configuration can be updated.",
			),
		)
	}

	rpcCtx, cancel := context.WithTimeout(ctx, configureRPCTimeout)
	defer cancel()

	requestID := xid.New()
	response, err := sess.Send(rpcCtx, requestID, rpc.RPCRequestConfigure{
		ID:      requestID,
		Jsonrpc: "2.0",
		Method:  "configure",
		Params:  cfg,
	})
	if err != nil {
		return nil, fault.Wrap(
			err,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc(
				"plugin rejected or failed to apply configuration",
				"The plugin reported an error while applying the configuration. Check the plugin error message and try again.",
			),
		)
	}

	configureResponse, err := parseConfigureResponse(response)
	if err != nil {
		return nil, fault.Wrap(
			err,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc(
				"invalid configure response from plugin",
				"The plugin returned an invalid configure response.",
			),
		)
	}

	if !configureResponse.Ok {
		return nil, fault.Wrap(
			fault.New("plugin rejected configuration"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc(
				"plugin rejected configuration",
				"The plugin rejected this configuration.",
			),
		)
	}

	updated, err := m.pluginWriter.UpdateConfig(ctx, id, cfg)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cloneConfig(updated.Config), nil
}

func parseConfigureResponse(in rpc.HostToPluginResponseUnion) (*rpc.RPCResponseConfigure, error) {
	if in.HostToPluginResponseUnionUnion == nil {
		return nil, fault.New("empty configure response")
	}

	res, ok := in.HostToPluginResponseUnionUnion.(*rpc.RPCResponseConfigure)
	if !ok {
		return nil, fault.Newf("unexpected configure response type: %T", in.HostToPluginResponseUnionUnion)
	}

	return res, nil
}

func cloneConfig(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{}
	}

	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}

	return out
}
