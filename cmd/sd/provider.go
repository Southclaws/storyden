package main

import (
	"github.com/Southclaws/storyden/cmd/sd/internal/apiexec"
	"github.com/Southclaws/storyden/cmd/sd/internal/cli"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/credentials"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/login"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/register"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/remove"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/auth/switcher"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/config/path"
	infocmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/info"
	nodeassets "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/assets"
	nodechildren "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/children"
	nodecreate "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/create"
	nodedelete "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/delete"
	nodeget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/get"
	nodelist "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/list"
	nodemeta "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/meta"
	nodemove "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/move"
	nodeopen "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/open"
	propertiesget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/get"
	schemachildren "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/children"
	schemaget "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/get"
	schemaset "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/schema/set"
	propertiesset "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/properties/set"
	nodesearch "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/search"
	nodetree "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/tree"
	nodeupdate "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/update"
	nodevisibility "github.com/Southclaws/storyden/cmd/sd/internal/commands/node/visibility"
	pluginactivate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/activate"
	plugindeactivate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/deactivate"
	plugindelete "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/delete"
	plugindevdownload "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/download"
	plugindevinstall "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/install"
	plugindevnew "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/new"
	plugindevpackage "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/package"
	plugindevrun "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/run"
	plugindevsymbols "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/symbols"
	plugindevvalidate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/dev/validate"
	pluginget "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/get"
	pluginlist "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/list"
	pluginlogs "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/logs"
	plugintokenrotate "github.com/Southclaws/storyden/cmd/sd/internal/commands/plugin/token/rotate"
	searchcmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/search"
	threadcreate "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/create"
	threadget "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/get"
	threadlist "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/list"
	threadreply "github.com/Southclaws/storyden/cmd/sd/internal/commands/thread/reply"
	tuicmd "github.com/Southclaws/storyden/cmd/sd/internal/commands/tui"
	"github.com/Southclaws/storyden/cmd/sd/internal/completion"
	storeconfig "github.com/Southclaws/storyden/cmd/sd/internal/config"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(apiexec.Build(), fx.Provide(
		apiexec.New,
		completion.New,
		cligen.NewCompletionCommand,
		apiexec.NewOperations,
		apiexec.NewSchema,
		apiexec.NewRequest,
		apiexec.NewChat,
		storeconfig.NewStore,
		cli.NewStreams,
		newLogger,

		// auth
		login.New,
		register.New,
		register.NewCheck,
		register.NewWait,
		register.NewCancel,
		credentials.NewToken,
		credentials.NewHeaders,
		credentials.NewStatus,
		remove.New,
		switcher.New,
		cligen.NewAuthCommand,

		// config
		path.New,
		cligen.NewConfigCommand,

		// info
		infocmd.New,
		infocmd.NewMetadata,
		cligen.NewInfoCommand,

		// search
		searchcmd.New,
		cligen.NewSearchCommand,

		// node
		nodelist.New,
		nodetree.New,
		nodeget.New,
		nodecreate.New,
		nodeupdate.New,
		nodedelete.New,
		nodemove.New,
		nodeopen.New,
		nodesearch.New,
		nodemeta.NewGet,
		nodemeta.NewSet,
		nodeassets.NewUpload,
		nodeassets.NewAdd,
		nodeassets.NewRemove,
		nodeassets.NewDownload,
		nodeassets.NewPrimarySet,
		nodeassets.NewPrimaryClear,
		nodeassets.NewPrimaryDownload,
		nodevisibility.New,
		nodechildren.New,
		propertiesget.New,
		propertiesset.New,
		schemaget.New,
		schemaset.New,
		schemachildren.New,
		cligen.NewPageCommand,

		// plugin
		plugindevnew.New,
		plugindevrun.New,
		plugindevpackage.New,
		plugindevvalidate.New,
		plugindevinstall.New,
		plugindevdownload.New,
		plugindevsymbols.NewPackages,
		plugindevsymbols.NewPackage,
		plugindevsymbols.NewDetail,
		plugindevsymbols.NewSearch,
		pluginlist.New,
		pluginget.New,
		plugindelete.New,
		pluginactivate.New,
		plugindeactivate.New,
		pluginlogs.New,
		plugintokenrotate.New,
		cligen.NewPluginCommand,

		// thread
		threadlist.New,
		threadcreate.New,
		threadreply.New,
		threadget.New,
		cligen.NewThreadCommand,

		// tui
		tuicmd.New,
		cligen.NewTuiCommand,

		cligen.NewAdminCommand,
		cligen.NewAccountCommand,
		cligen.NewRoleCommand,
		cligen.NewInvitationCommand,
		cligen.NewNotificationCommand,
		cligen.NewReportCommand,
		cligen.NewProfileCommand,
		cligen.NewCategoryCommand,
		cligen.NewTagCommand,
		cligen.NewPostCommand,
		cligen.NewCollectionCommand,
		cligen.NewLinkCommand,
		cligen.NewEventCommand,
		cligen.NewRobotCommand,
		cligen.NewTrailCommand,
		cligen.NewAssetCommand,
		cligen.NewApiCommand,
		newRootCommand,
	))
}
