package plugin_auth

import (
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"

	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
)

const (
	QueryParamPluginID = "plugin_id"
	QueryParamToken    = "token"
)

func BuildConnectionURL(baseURL string, id resource_plugin.InstallationID, authSecret string) (*url.URL, error) {
	token, err := SealToken(id, authSecret)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to generate token"))
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to parse base URL"))
	}

	q := u.Query()
	q.Set(QueryParamPluginID, id.String())
	q.Set(QueryParamToken, token)
	u.RawQuery = q.Encode()

	return u, nil
}

type ConnectionParams struct {
	PluginID resource_plugin.InstallationID
	Token    string
}

func ParseConnectionURL(u *url.URL) (*ConnectionParams, error) {
	pluginIDStr := u.Query().Get(QueryParamPluginID)
	if pluginIDStr == "" {
		return nil, fault.Newf("missing %s query parameter", QueryParamPluginID)
	}

	token := u.Query().Get(QueryParamToken)
	if token == "" {
		return nil, fault.Newf("missing %s query parameter", QueryParamToken)
	}

	xidPluginID, err := xid.FromString(pluginIDStr)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("invalid plugin ID format"))
	}

	return &ConnectionParams{
		PluginID: resource_plugin.InstallationID(xidPluginID),
		Token:    token,
	}, nil
}
