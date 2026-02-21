package rpc_handler

import (
	"log/slog"
	"net/url"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/internal/config"
)

type Factory struct {
	apiBaseURL     url.URL
	accountQuerier *account_querier.Querier
	accountWriter  *account_writer.Writer
	accessKeys     *access_key.Repository
	pluginReader   *plugin_reader.Reader
}

func NewFactory(
	cfg config.Config,
	accountQuerier *account_querier.Querier,
	accountWriter *account_writer.Writer,
	accessKeys *access_key.Repository,
	pluginReader *plugin_reader.Reader,
) *Factory {
	return &Factory{
		apiBaseURL:     cfg.PublicAPIAddress,
		accountQuerier: accountQuerier,
		accountWriter:  accountWriter,
		accessKeys:     accessKeys,
		pluginReader:   pluginReader,
	}
}

func (f *Factory) New(
	logger *slog.Logger,
	installationID plugin.InstallationID,
	manifest *plugin.Validated,
) *Handler {
	return New(
		logger,
		installationID,
		manifest,
		f.apiBaseURL,
		f.accountQuerier,
		f.accountWriter,
		f.accessKeys,
		f.pluginReader,
	)
}
