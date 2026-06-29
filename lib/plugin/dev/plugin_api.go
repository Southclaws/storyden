package dev

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

const RPCURLEnvName = "STORYDEN_RPC_URL"

type ExternalPlugin struct {
	ID    string
	Token string
}

type statusResponse interface {
	Status() string
	StatusCode() int
}

func EnsureExternalPlugin(ctx context.Context, client *openapi.ClientWithResponses, manifest rpc.Manifest, requestedID string, noUpdate bool) (*ExternalPlugin, error) {
	if requestedID != "" {
		var plugin *openapi.Plugin
		var err error
		if !noUpdate {
			plugin, err = UpdateManifest(ctx, client, requestedID, manifest)
			if err != nil {
				return nil, err
			}
		} else {
			plugin, err = GetPlugin(ctx, client, requestedID)
			if err != nil {
				return nil, err
			}
		}
		return ExternalPluginFromAPI(*plugin)
	}

	plugins, err := ListPlugins(ctx, client)
	if err != nil {
		return nil, err
	}
	for _, pl := range plugins {
		if id, _ := pl.Manifest["id"].(string); id == manifest.ID {
			if PluginMode(pl) != string(openapi.External) {
				return nil, fmt.Errorf("plugin manifest id %q is already installed as a supervised plugin; pass --instance-id for an external installation", manifest.ID)
			}
			if !noUpdate {
				updated, err := UpdateManifest(ctx, client, string(pl.Id), manifest)
				if err != nil {
					return nil, err
				}
				return ExternalPluginFromAPI(*updated)
			}
			return ExternalPluginFromAPI(pl)
		}
	}

	body := openapi.PluginInitialProps{}
	if err := body.FromPluginInitialExternal(openapi.PluginInitialExternal{
		Mode:     openapi.External,
		Manifest: openapi.PluginManifest(manifest.ToMap()),
	}); err != nil {
		return nil, err
	}

	response, err := client.PluginAddWithResponse(ctx, body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, requestError("plugin add request", response, response.Body)
	}
	return ExternalPluginFromAPI(*response.JSON200)
}

func InstallSupervisedPackage(ctx context.Context, client *openapi.ClientWithResponses, pkg *PackageArchive) (*openapi.Plugin, bool, error) {
	plugins, err := ListPlugins(ctx, client)
	if err != nil {
		return nil, false, err
	}

	for _, pl := range plugins {
		if id, _ := pl.Manifest["id"].(string); id == pkg.Manifest.ID {
			if PluginMode(pl) != string(openapi.Supervised) {
				return nil, false, fmt.Errorf("plugin manifest id %q is already installed as an external plugin", pkg.Manifest.ID)
			}

			response, err := client.PluginUpdatePackageWithBodyWithResponse(ctx, openapi.PluginIDParam(pl.Id), "application/zip", bytes.NewReader(pkg.Bytes))
			if err != nil {
				return nil, false, err
			}
			if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
				return nil, false, requestError("plugin package update request", response, response.Body)
			}
			return (*openapi.Plugin)(response.JSON200), true, nil
		}
	}

	response, err := client.PluginAddWithBodyWithResponse(ctx, "application/zip", bytes.NewReader(pkg.Bytes))
	if err != nil {
		return nil, false, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, false, requestError("plugin package install request", response, response.Body)
	}
	return (*openapi.Plugin)(response.JSON200), false, nil
}

func ListPlugins(ctx context.Context, client *openapi.ClientWithResponses) ([]openapi.Plugin, error) {
	response, err := client.PluginListWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, requestError("plugin list request", response, response.Body)
	}
	return response.JSON200.Plugins, nil
}

func ExternalPluginFromAPI(p openapi.Plugin) (*ExternalPlugin, error) {
	if PluginMode(p) != string(openapi.External) {
		return nil, fmt.Errorf("plugin %s is not an external plugin", p.Id)
	}
	props, err := p.Connection.AsPluginExternalProps()
	if err != nil {
		return nil, fmt.Errorf("plugin %s is not an external plugin: %w", p.Id, err)
	}
	if props.Token == "" {
		return nil, fmt.Errorf("plugin %s did not include an external RPC token", p.Id)
	}
	return &ExternalPlugin{ID: string(p.Id), Token: props.Token}, nil
}

func GetPlugin(ctx context.Context, client *openapi.ClientWithResponses, id string) (*openapi.Plugin, error) {
	response, err := client.PluginGetWithResponse(ctx, openapi.PluginIDParam(id))
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, requestError("plugin get request", response, response.Body)
	}
	return (*openapi.Plugin)(response.JSON200), nil
}

func DeletePlugin(ctx context.Context, client *openapi.ClientWithResponses, id string) error {
	response, err := client.PluginDeleteWithResponse(ctx, openapi.PluginIDParam(id))
	if err != nil {
		return err
	}
	if response.StatusCode() != http.StatusNoContent {
		return requestError("plugin delete request", response, response.Body)
	}
	return nil
}

func SetActiveState(ctx context.Context, client *openapi.ClientWithResponses, id string, state openapi.PluginActiveState) (*openapi.Plugin, error) {
	response, err := client.PluginSetActiveStateWithResponse(ctx, openapi.PluginIDParam(id), openapi.PluginSetActiveStateJSONRequestBody{Active: state})
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, requestError("plugin active-state request", response, response.Body)
	}
	return (*openapi.Plugin)(response.JSON200), nil
}

func UpdateManifest(ctx context.Context, client *openapi.ClientWithResponses, id string, manifest rpc.Manifest) (*openapi.Plugin, error) {
	response, err := client.PluginUpdateManifestWithResponse(ctx, openapi.PluginIDParam(id), openapi.PluginUpdateManifestJSONRequestBody(manifest.ToMap()))
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, requestError("plugin manifest update request", response, response.Body)
	}
	return (*openapi.Plugin)(response.JSON200), nil
}

func CycleToken(ctx context.Context, client *openapi.ClientWithResponses, id string) (string, error) {
	response, err := client.PluginCycleTokenWithResponse(ctx, openapi.PluginIDParam(id))
	if err != nil {
		return "", err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return "", requestError("plugin token request", response, response.Body)
	}
	return response.JSON200.Token, nil
}

func ExternalRPCURL(endpoint string, token string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
	default:
		return "", fmt.Errorf("unsupported Storyden endpoint scheme %q", u.Scheme)
	}
	u.Path = "/rpc"
	u.RawPath = ""
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func CommandFromManifest(manifest rpc.Manifest, override []string) (string, []string, error) {
	if len(override) > 0 {
		return override[0], override[1:], nil
	}
	if strings.TrimSpace(manifest.Command) == "" {
		return "", nil, fmt.Errorf("manifest command is empty; pass a command after --")
	}
	return manifest.Command, manifest.Args, nil
}

func PluginMode(p openapi.Plugin) string {
	mode, err := p.Connection.Discriminator()
	if err != nil {
		return ""
	}
	return mode
}

func PluginStatus(p openapi.Plugin) string {
	state, err := p.Status.Discriminator()
	if err != nil {
		return ""
	}
	return state
}

func requestError(operation string, response statusResponse, body []byte) error {
	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("%s was not authorised; run sd auth login again", operation)
	}

	text := strings.TrimSpace(string(body))
	if text != "" {
		return fmt.Errorf("%s failed: %s: %s", operation, response.Status(), text)
	}

	return fmt.Errorf("%s failed: %s", operation, response.Status())
}
