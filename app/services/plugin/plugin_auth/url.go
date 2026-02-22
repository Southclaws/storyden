package plugin_auth

import (
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
)

const (
	QueryParamPluginID = "plugin_id"
	QueryParamToken    = "token"
)

func BuildConnectionURL(baseURL url.URL, id resource_plugin.InstallationID, authSecret string) (*url.URL, error) {
	token, err := SealToken(id, authSecret)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to generate token"))
	}

	u := baseURL

	u.Scheme = websocketSchemeFor(u.Scheme)

	u.Path = "/rpc"

	q := u.Query()
	q.Set(QueryParamPluginID, id.String())
	q.Set(QueryParamToken, token)
	u.RawQuery = q.Encode()

	return &u, nil
}

func websocketSchemeFor(scheme string) string {
	switch strings.ToLower(strings.TrimSpace(scheme)) {
	case "https":
		return "wss"
	default:
		return "ws"
	}
}

type ConnectionParams struct {
	PluginID opt.Optional[resource_plugin.InstallationID]
	Token    string
}

func ParseConnectionURL(u *url.URL) (*ConnectionParams, error) {
	token := u.Query().Get(QueryParamToken)
	if token == "" {
		return nil, fault.Newf("missing %s query parameter", QueryParamToken)
	}

	var pluginID opt.Optional[resource_plugin.InstallationID]
	pluginIDStr := u.Query().Get(QueryParamPluginID)
	if pluginIDStr != "" {
		xidPluginID, err := xid.FromString(pluginIDStr)
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("invalid plugin ID format"))
		}

		pluginID = opt.New(resource_plugin.InstallationID(xidPluginID))
	} else if !IsExternalToken(token) {
		return nil, fault.Newf("missing %s query parameter", QueryParamPluginID)
	}

	return &ConnectionParams{
		PluginID: pluginID,
		Token:    token,
	}, nil
}
