package rpc_handler

import (
	"log/slog"
	"net/url"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/account/role/role_writer"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/services/account/account_role_assign"
	"github.com/Southclaws/storyden/app/services/semdex/robot/session_coordinator"
	"github.com/Southclaws/storyden/internal/config"
)

type Factory struct {
	apiBaseURL     url.URL
	accountQuerier *account_querier.Querier
	accountWriter  *account_writer.Writer
	roleQuerier    *role_querier.Querier
	roleWriter     *role_writer.Writer
	roleAssigner   *account_role_assign.Manager
	accessKeys     *access_key.Repository
	pluginReader   *plugin_reader.Reader
	robotSessions  *robot_session.Repository
	robotAgent     *session_coordinator.Coordinator
}

func NewFactory(
	cfg config.Config,
	accountQuerier *account_querier.Querier,
	accountWriter *account_writer.Writer,
	roleQuerier *role_querier.Querier,
	roleWriter *role_writer.Writer,
	roleAssigner *account_role_assign.Manager,
	accessKeys *access_key.Repository,
	pluginReader *plugin_reader.Reader,
	robotSessions *robot_session.Repository,
	robotAgent *session_coordinator.Coordinator,
) *Factory {
	return &Factory{
		apiBaseURL:     cfg.PublicAPIAddress,
		accountQuerier: accountQuerier,
		accountWriter:  accountWriter,
		roleQuerier:    roleQuerier,
		roleWriter:     roleWriter,
		roleAssigner:   roleAssigner,
		accessKeys:     accessKeys,
		pluginReader:   pluginReader,
		robotSessions:  robotSessions,
		robotAgent:     robotAgent,
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
		f.roleQuerier,
		f.roleWriter,
		f.roleAssigner,
		f.accessKeys,
		f.pluginReader,
		f.robotSessions,
		f.robotAgent,
	)
}
