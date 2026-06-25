package bindings

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/session"

	"github.com/Southclaws/storyden/app/resources/account"
	oauth_remote "github.com/Southclaws/storyden/app/resources/oauth/remote"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	robot_mcp "github.com/Southclaws/storyden/app/resources/robot/mcp"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/robot/robot_workspace"
	"github.com/Southclaws/storyden/app/resources/robot/robot_writer"
	"github.com/Southclaws/storyden/app/resources/robot/session_ref"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/admin/settings_manager"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/mcpclient"
	robot_tools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/http/robotprojection"
)

type Robots struct {
	robotQuerier       *robot_querier.Querier
	robotWriter        *robot_writer.Writer
	workspaceRepo      *robot_workspace.Repository
	workspaceManager   *robotservice.WorkspaceManager
	workspaceProviders *workspaceprovider.Registry
	sessionRepo        *robot_session.Repository
	modelFactory       *llm_provider.Factory
	settings           *settings_manager.Manager
	mcp                *mcpclient.Manager
	tools              *robot_tools.Registry
}

func NewRobots(
	robotQuerier *robot_querier.Querier,
	robotWriter *robot_writer.Writer,
	workspaceRepo *robot_workspace.Repository,
	workspaceManager *robotservice.WorkspaceManager,
	workspaceProviders *workspaceprovider.Registry,
	sessionRepo *robot_session.Repository,
	modelFactory *llm_provider.Factory,
	settingsManager *settings_manager.Manager,
	mcpManager *mcpclient.Manager,
	toolRegistry *robot_tools.Registry,
) Robots {
	return Robots{
		robotQuerier:       robotQuerier,
		robotWriter:        robotWriter,
		workspaceRepo:      workspaceRepo,
		workspaceManager:   workspaceManager,
		workspaceProviders: workspaceProviders,
		sessionRepo:        sessionRepo,
		modelFactory:       modelFactory,
		settings:           settingsManager,
		mcp:                mcpManager,
		tools:              toolRegistry,
	}
}

func (r *Robots) RobotsList(ctx context.Context, request openapi.RobotsListRequestObject) (openapi.RobotsListResponseObject, error) {
	pageParams := deserialisePageParams(request.Params.Page, 50)

	result, err := r.robotQuerier.List(ctx, pageParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotsList200JSONResponse{
		RobotsListOKJSONResponse: openapi.RobotsListOKJSONResponse(openapi.RobotsListResult{
			CurrentPage: result.CurrentPage,
			NextPage:    result.NextPage.Ptr(),
			PageSize:    result.Size,
			Results:     result.Results,
			Robots:      serialiseRobots(result.Items),
			TotalPages:  result.TotalPages,
		}),
	}, nil
}

func (r *Robots) RobotCreate(ctx context.Context, request openapi.RobotCreateRequestObject) (openapi.RobotCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []robot_writer.Option{}
	if meta := request.Body.Meta; meta != nil {
		opts = append(opts, robot_writer.WithMeta(*meta))
	}
	if tools := request.Body.Tools; tools != nil {
		if err := r.validateRobotToolsForCreate(ctx, *tools); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		opts = append(opts, robot_writer.WithTools(dt.Map(*tools, func(t string) robot_ref.ToolName { return robot_ref.ToolName(t) })))
	}
	if workspaceID := request.Body.WorkspaceId; workspaceID != nil {
		id, err := xid.FromString(string(*workspaceID))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if _, err := r.workspaceRepo.Get(ctx, robot.WorkspaceID(id)); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		opts = append(opts, robot_writer.WithWorkspaceID(id))
	}

	var model model_ref.ModelRef
	if request.Body.Model != nil {
		var err error
		model, err = model_ref.ParseID(string(*request.Body.Model))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
	} else {
		defaultModel, err := r.modelFactory.DefaultModel(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		model = defaultModel
	}
	if err := r.validateRobotModelForSave(ctx, model); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	created, err := r.robotWriter.Create(ctx,
		request.Body.Name,
		request.Body.Description,
		request.Body.Playbook,
		model,
		accountID,
		opts...,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotCreate200JSONResponse{
		RobotCreateOKJSONResponse: openapi.RobotCreateOKJSONResponse(serialiseRobot(created)),
	}, nil
}

func (r *Robots) RobotChatSSE(ctx context.Context, request openapi.RobotChatSSERequestObject) (openapi.RobotChatSSEResponseObject, error) {
	return nil, fault.New("bindings layer should not be hit, this is a bug.", fctx.With(ctx))
}

func (r *Robots) RobotProvidersList(ctx context.Context, request openapi.RobotProvidersListRequestObject) (openapi.RobotProvidersListResponseObject, error) {
	runtime, err := r.modelFactory.RuntimeSettings(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	providers, err := r.serialiseProviderStatuses(ctx, runtime)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotProvidersList200JSONResponse{
		RobotProvidersListOKJSONResponse: openapi.RobotProvidersListOKJSONResponse{
			Providers: providers,
		},
	}, nil
}

func (r *Robots) RobotModelsList(ctx context.Context, request openapi.RobotModelsListRequestObject) (openapi.RobotModelsListResponseObject, error) {
	runtime, err := r.modelFactory.RuntimeSettings(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var enabledProviders []model_ref.Provider
	for _, provider := range r.modelFactory.Providers() {
		if p, ok := runtime.Providers[provider]; ok && p.Enabled && (p.APIKey != "" || !r.modelFactory.RequiresAPIKey(provider)) {
			enabledProviders = append(enabledProviders, provider)
			if _, _, err := r.modelFactory.ProviderStatus(ctx, provider); err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	models, err := r.modelFactory.ListCachedModels(ctx, enabledProviders)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotModelsList200JSONResponse{
		RobotModelsListOKJSONResponse: openapi.RobotModelsListOKJSONResponse(openapi.RobotModelListResult{
			Models: serialiseRobotModelInfos(models),
		}),
	}, nil
}

func (r *Robots) RobotWorkspacesList(ctx context.Context, request openapi.RobotWorkspacesListRequestObject) (openapi.RobotWorkspacesListResponseObject, error) {
	pageParams := deserialisePageParams(request.Params.Page, 50)

	result, err := r.workspaceRepo.List(ctx, pageParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotWorkspacesList200JSONResponse{
		RobotWorkspacesListOKJSONResponse: openapi.RobotWorkspacesListOKJSONResponse(openapi.RobotWorkspacesListResult{
			CurrentPage: result.CurrentPage,
			NextPage:    result.NextPage.Ptr(),
			PageSize:    result.Size,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
			Workspaces:  dt.Map(result.Items, serialiseRobotWorkspace),
		}),
	}, nil
}

func (r *Robots) RobotWorkspaceProvidersList(ctx context.Context, request openapi.RobotWorkspaceProvidersListRequestObject) (openapi.RobotWorkspaceProvidersListResponseObject, error) {
	providers := r.workspaceProviders.List()

	return openapi.RobotWorkspaceProvidersList200JSONResponse{
		RobotWorkspaceProvidersListOKJSONResponse: openapi.RobotWorkspaceProvidersListOKJSONResponse(openapi.RobotWorkspaceProviderListResult{
			Providers: dt.Map(providers, func(provider workspaceprovider.ProviderInfo) openapi.RobotWorkspaceProviderInfo {
				return openapi.RobotWorkspaceProviderInfo{
					Provider: openapi.RobotWorkspaceProvider(provider.Provider),
					Name:     provider.Name,
				}
			}),
		}),
	}, nil
}

func (r *Robots) RobotWorkspaceCreate(ctx context.Context, request openapi.RobotWorkspaceCreateRequestObject) (openapi.RobotWorkspaceCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	provider := robot.WorkspaceProviderLocal
	if request.Body.Provider != nil {
		provider = robot.WorkspaceProvider(*request.Body.Provider)
	}
	if _, err := r.workspaceProviders.Get(ctx, provider); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	opts := []robot_workspace.WorkspaceOption{}
	if request.Body.Config != nil {
		opts = append(opts, robot_workspace.WithConfig(*request.Body.Config))
	}
	if request.Body.Meta != nil {
		opts = append(opts, robot_workspace.WithMetadata(*request.Body.Meta))
	}

	created, err := r.workspaceRepo.Create(ctx, request.Body.Name, request.Body.Description, provider, accountID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotWorkspaceCreate200JSONResponse{
		RobotWorkspaceCreateOKJSONResponse: openapi.RobotWorkspaceCreateOKJSONResponse(serialiseRobotWorkspace(created)),
	}, nil
}

func (r *Robots) RobotWorkspaceGet(ctx context.Context, request openapi.RobotWorkspaceGetRequestObject) (openapi.RobotWorkspaceGetResponseObject, error) {
	id, err := robot.NewWorkspaceID(request.WorkspaceId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	workspace, err := r.workspaceRepo.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotWorkspaceGet200JSONResponse{
		RobotWorkspaceGetOKJSONResponse: openapi.RobotWorkspaceGetOKJSONResponse(serialiseRobotWorkspace(workspace)),
	}, nil
}

func (r *Robots) RobotWorkspaceUpdate(ctx context.Context, request openapi.RobotWorkspaceUpdateRequestObject) (openapi.RobotWorkspaceUpdateResponseObject, error) {
	id, err := robot.NewWorkspaceID(request.WorkspaceId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	opts := []robot_workspace.WorkspaceOption{}
	if request.Body.Name != nil {
		opts = append(opts, robot_workspace.WithName(*request.Body.Name))
	}
	if request.Body.Description != nil {
		opts = append(opts, robot_workspace.WithDescription(*request.Body.Description))
	}
	if request.Body.Config != nil {
		opts = append(opts, robot_workspace.WithConfig(*request.Body.Config))
	}
	if request.Body.Meta != nil {
		opts = append(opts, robot_workspace.WithMetadata(*request.Body.Meta))
	}

	updated, err := r.workspaceRepo.Update(ctx, id, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotWorkspaceUpdate200JSONResponse{
		RobotWorkspaceGetOKJSONResponse: openapi.RobotWorkspaceGetOKJSONResponse(serialiseRobotWorkspace(updated)),
	}, nil
}

func (r *Robots) RobotWorkspaceDelete(ctx context.Context, request openapi.RobotWorkspaceDeleteRequestObject) (openapi.RobotWorkspaceDeleteResponseObject, error) {
	id, err := robot.NewWorkspaceID(request.WorkspaceId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if err := r.workspaceRepo.Delete(ctx, id); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotWorkspaceDelete200Response{}, nil
}

func (r *Robots) RobotWorkspaceInstancesList(ctx context.Context, request openapi.RobotWorkspaceInstancesListRequestObject) (openapi.RobotWorkspaceInstancesListResponseObject, error) {
	pageParams := deserialisePageParams(request.Params.Page, 50)

	result, err := r.workspaceRepo.ListInstances(ctx, pageParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotWorkspaceInstancesList200JSONResponse{
		RobotWorkspaceInstancesListOKJSONResponse: openapi.RobotWorkspaceInstancesListOKJSONResponse(openapi.RobotWorkspaceInstancesListResult{
			CurrentPage:        result.CurrentPage,
			NextPage:           result.NextPage.Ptr(),
			PageSize:           result.Size,
			Results:            result.Results,
			TotalPages:         result.TotalPages,
			WorkspaceInstances: dt.Map(result.Items, serialiseRobotWorkspaceInstance),
		}),
	}, nil
}

func (r *Robots) RobotWorkspaceInstanceGet(ctx context.Context, request openapi.RobotWorkspaceInstanceGetRequestObject) (openapi.RobotWorkspaceInstanceGetResponseObject, error) {
	id, err := robot.NewWorkspaceInstanceID(request.WorkspaceInstanceId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	instance, err := r.workspaceRepo.GetInstance(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotWorkspaceInstanceGet200JSONResponse{
		RobotWorkspaceInstanceGetOKJSONResponse: openapi.RobotWorkspaceInstanceGetOKJSONResponse(serialiseRobotWorkspaceInstance(instance)),
	}, nil
}

func (r *Robots) RobotWorkspaceInstanceDelete(ctx context.Context, request openapi.RobotWorkspaceInstanceDeleteRequestObject) (openapi.RobotWorkspaceInstanceDeleteResponseObject, error) {
	id, err := robot.NewWorkspaceInstanceID(request.WorkspaceInstanceId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if err := r.workspaceManager.DeleteInstance(ctx, id); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotWorkspaceInstanceDelete200Response{}, nil
}

func (r *Robots) RobotProviderUpdate(ctx context.Context, request openapi.RobotProviderUpdateRequestObject) (openapi.RobotProviderUpdateResponseObject, error) {
	if request.Body == nil {
		return nil, fault.Wrap(fault.New("missing request body"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	provider := model_ref.NewProvider(string(request.Provider))
	if !r.modelFactory.HasProvider(provider) {
		return nil, fault.Wrap(fault.New("unsupported robot provider"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	runtime, err := r.modelFactory.RuntimeSettings(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	current, err := r.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	services, robots, providers := currentRobotSettings(current)
	providerSettings := providers[provider.String()]
	effectiveProvider := runtime.Providers[provider]
	finalEnabled := effectiveProvider.Enabled
	finalKey := effectiveProvider.APIKey
	refreshRequired := false

	if request.Body.Enabled != nil {
		finalEnabled = *request.Body.Enabled
		providerSettings.Enabled = opt.New(finalEnabled)
	}
	if request.Body.ClearApiKey != nil && *request.Body.ClearApiKey {
		finalKey = ""
		providerSettings.APIKey = opt.New("")
		refreshRequired = finalEnabled
	}
	if request.Body.ApiKey != nil {
		finalKey = *request.Body.ApiKey
		providerSettings.APIKey = opt.New(finalKey)
		refreshRequired = true
	}

	if finalEnabled && finalKey == "" && r.modelFactory.RequiresAPIKey(provider) {
		return nil, fault.Wrap(fault.New("provider API key is required when enabling a robot provider"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if refreshRequired {
		if _, err := r.modelFactory.RefreshProviderModelsWithKey(ctx, provider, finalKey); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
	}

	providers[provider.String()] = providerSettings
	robots.Providers = opt.New(providers)
	robots.Enabled = opt.New(anyRobotProviderEnabled(runtime, providers, r.modelFactory.Providers(), provider, finalEnabled))
	services.Robots = opt.New(robots)

	if _, err := r.settings.Set(ctx, settings.Settings{
		Services: opt.New(services),
	}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := r.serialiseProviderStatus(ctx, provider)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotProviderUpdate200JSONResponse{
		RobotProviderGetOKJSONResponse: openapi.RobotProviderGetOKJSONResponse(status),
	}, nil
}

func (r *Robots) RobotProviderModelsRefresh(ctx context.Context, request openapi.RobotProviderModelsRefreshRequestObject) (openapi.RobotProviderModelsRefreshResponseObject, error) {
	provider := model_ref.NewProvider(string(request.Provider))
	if !r.modelFactory.HasProvider(provider) {
		return nil, fault.Wrap(fault.New("unsupported robot provider"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	runtime, err := r.modelFactory.RuntimeSettings(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	providerSettings := runtime.Providers[provider]
	if providerSettings.APIKey == "" && r.modelFactory.RequiresAPIKey(provider) {
		return nil, fault.Wrap(fault.New("provider API key is required to refresh robot provider models"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if _, err := r.modelFactory.RefreshProviderModelsWithKey(ctx, provider, providerSettings.APIKey); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	status, err := r.serialiseProviderStatus(ctx, provider)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotProviderModelsRefresh200JSONResponse{
		RobotProviderGetOKJSONResponse: openapi.RobotProviderGetOKJSONResponse(status),
	}, nil
}

func (r *Robots) RobotGet(ctx context.Context, request openapi.RobotGetRequestObject) (openapi.RobotGetResponseObject, error) {
	robotID, err := robot_ref.NewID(request.RobotId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rb, err := r.robotQuerier.Get(ctx, robotID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotGet200JSONResponse{
		RobotGetOKJSONResponse: openapi.RobotGetOKJSONResponse(serialiseRobot(rb)),
	}, nil
}

func (r *Robots) RobotUpdate(ctx context.Context, request openapi.RobotUpdateRequestObject) (openapi.RobotUpdateResponseObject, error) {
	robotID, err := robot_ref.NewID(request.RobotId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []robot_writer.Option{}
	if name := request.Body.Name; name != nil {
		opts = append(opts, robot_writer.WithName(*name))
	}
	if description := request.Body.Description; description != nil {
		opts = append(opts, robot_writer.WithDescription(*description))
	}
	if playbook := request.Body.Playbook; playbook != nil {
		opts = append(opts, robot_writer.WithPlaybook(*playbook))
	}
	if model := request.Body.Model; model != nil {
		modelID, err := model_ref.ParseID(string(*model))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if err := r.validateRobotModelForSave(ctx, modelID); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		opts = append(opts, robot_writer.WithModel(modelID))
	}
	if meta := request.Body.Meta; meta != nil {
		opts = append(opts, robot_writer.WithMeta(*meta))
	}
	if request.Body.WorkspaceId.IsSpecified() {
		if request.Body.WorkspaceId.IsNull() {
			opts = append(opts, robot_writer.WithoutWorkspaceID())
		} else {
			workspaceID, err := request.Body.WorkspaceId.Get()
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
			}
			id, err := xid.FromString(string(workspaceID))
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
			}
			if _, err := r.workspaceRepo.Get(ctx, robot.WorkspaceID(id)); err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			opts = append(opts, robot_writer.WithWorkspaceID(id))
		}
	}
	if tools := request.Body.Tools; tools != nil {
		if err := r.validateRobotToolsForUpdate(ctx, robotID, *tools); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		opts = append(opts, robot_writer.WithTools(dt.Map(*tools, func(t string) robot_ref.ToolName { return robot_ref.ToolName(t) })))
	}

	updated, err := r.robotWriter.Update(ctx, robotID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotUpdate200JSONResponse{
		RobotGetOKJSONResponse: openapi.RobotGetOKJSONResponse(serialiseRobot(updated)),
	}, nil
}

func (r *Robots) RobotDelete(ctx context.Context, request openapi.RobotDeleteRequestObject) (openapi.RobotDeleteResponseObject, error) {
	robotID, err := robot_ref.NewID(request.RobotId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.robotWriter.Delete(ctx, robotID); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotDelete200Response{}, nil
}

func (r *Robots) RobotSessionsList(ctx context.Context, request openapi.RobotSessionsListRequestObject) (openapi.RobotSessionsListResponseObject, error) {
	pageParams := deserialisePageParams(request.Params.Page, 20)

	accountID := opt.NewEmpty[account.AccountID]()
	if request.Params.AccountId != nil {
		id, err := xid.FromString(string(*request.Params.AccountId))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		accountID = opt.New(account.AccountID(id))
	}

	result, err := r.sessionRepo.List(ctx, pageParams, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sessions := dt.Map(result.Items, func(s *session_ref.Ref) openapi.RobotSessionRef {
		return openapi.RobotSessionRef{
			Id:        openapi.Identifier(s.ID.String()),
			Name:      s.Name,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			CreatedBy: serialiseProfileReferenceFromAccount(s.Human),
		}
	})

	return openapi.RobotSessionsList200JSONResponse{
		RobotSessionsListOKJSONResponse: openapi.RobotSessionsListOKJSONResponse(openapi.RobotSessionsListResult{
			CurrentPage: result.CurrentPage,
			NextPage:    result.NextPage.Ptr(),
			PageSize:    result.Size,
			Results:     result.Results,
			Sessions:    sessions,
			TotalPages:  result.TotalPages,
		}),
	}, nil
}

func (r *Robots) RobotSessionGet(ctx context.Context, request openapi.RobotSessionGetRequestObject) (openapi.RobotSessionGetResponseObject, error) {
	sessionID, err := robot.NewSessionID(request.SessionId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	messageParams, err := deserialiseRobotMessageCursorParams(request.Params.Before, request.Params.Limit)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, cursor, err := r.sessionRepo.Get(ctx, robot.SessionID(sessionID), messageParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Convert messages to Vercel AI SDK UIMessage format
	hiddenToolCallIDs := robotprojection.HiddenConfirmationToolCallIDs(cursor.Items)
	messageDTOs, err := dt.MapErr(cursor.Items, func(m *robot.Message) (openapi.RobotSessionMessage, error) {
		return serialiseRobotSessionMessage(m, hiddenToolCallIDs)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotSessionGet200JSONResponse{
		RobotSessionGetOKJSONResponse: openapi.RobotSessionGetOKJSONResponse(openapi.RobotSession{
			Id:              openapi.Identifier(sess.ID.String()),
			Name:            sess.Name,
			CreatedAt:       sess.CreatedAt,
			UpdatedAt:       sess.UpdatedAt,
			CreatedBy:       serialiseProfileReferenceFromAccount(sess.Human),
			ActiveRobotId:   robotservice.CurrentRobotIDFromState(sess.State).Ptr(),
			ActiveWorkspace: serialiseRobotWorkspaceMountPtr(robotservice.WorkspaceMountFromState(sess.State).Ptr()),
			MessageList: openapi.PaginatedRobotMessageList{
				NextBefore: opt.PtrMap(cursor.NextBefore, func(id robot.MessageID) openapi.Identifier {
					return openapi.Identifier(id.String())
				}),
				PageSize: cursor.Size,
				Results:  cursor.Results,
				Messages: messageDTOs,
			},
		}),
	}, nil
}

func (r *Robots) RobotToolsList(ctx context.Context, request openapi.RobotToolsListRequestObject) (openapi.RobotToolsListResponseObject, error) {
	catalogue, err := r.mcp.ListToolCatalogue(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotToolsList200JSONResponse{RobotToolsListOKJSONResponse: openapi.RobotToolsListOKJSONResponse{
		Tools: dt.Map(catalogue, serialiseRobotToolInfo),
	}}, nil
}

func (r *Robots) RobotMCPServersList(ctx context.Context, request openapi.RobotMCPServersListRequestObject) (openapi.RobotMCPServersListResponseObject, error) {
	servers, err := r.mcp.ListServers(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotMCPServersList200JSONResponse{RobotMCPServersListOKJSONResponse: openapi.RobotMCPServersListOKJSONResponse{
		Servers: serialiseRobotMCPServers(servers),
	}}, nil
}

func (r *Robots) RobotMCPServerCreate(ctx context.Context, request openapi.RobotMCPServerCreateRequestObject) (openapi.RobotMCPServerCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if request.Body == nil {
		return nil, fault.Wrap(fault.New("missing request body"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	enabled := true
	if request.Body.Enabled != nil {
		enabled = *request.Body.Enabled
	}
	description := ""
	if request.Body.Description != nil {
		description = *request.Body.Description
	}
	slug := ""
	if request.Body.Slug != nil {
		slug = *request.Body.Slug
	}
	bearerToken := ""
	if request.Body.BearerToken != nil {
		bearerToken = *request.Body.BearerToken
	}
	var oauthRemoteConnectionID *oauth_remote.ConnectionID
	if request.Body.OauthRemoteConnectionId != nil {
		id, err := oauth_remote.NewConnectionID(string(*request.Body.OauthRemoteConnectionId))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		oauthRemoteConnectionID = &id
	}

	server, err := r.mcp.CreateServer(ctx, robot_mcp.ServerCreate{
		Name:                    request.Body.Name,
		Slug:                    slug,
		Description:             description,
		EndpointURL:             request.Body.EndpointUrl,
		Enabled:                 enabled,
		BearerToken:             bearerToken,
		OAuthRemoteConnectionID: oauthRemoteConnectionID,
		AddedBy:                 accountID,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotMCPServerCreate200JSONResponse{
		RobotMCPServerCreateOKJSONResponse: openapi.RobotMCPServerCreateOKJSONResponse(serialiseRobotMCPServer(server)),
	}, nil
}

func (r *Robots) RobotMCPServerProbe(ctx context.Context, request openapi.RobotMCPServerProbeRequestObject) (openapi.RobotMCPServerProbeResponseObject, error) {
	bearerToken := ""
	if request.Body.BearerToken != nil {
		bearerToken = *request.Body.BearerToken
	}

	result, err := r.mcp.Probe(ctx, request.Body.Url, bearerToken)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	return openapi.RobotMCPServerProbe200JSONResponse{
		RobotMCPServerProbeOKJSONResponse: openapi.RobotMCPServerProbeOKJSONResponse(serialiseRobotMCPServerProbe(result)),
	}, nil
}

func (r *Robots) RobotMCPServerGet(ctx context.Context, request openapi.RobotMCPServerGetRequestObject) (openapi.RobotMCPServerGetResponseObject, error) {
	id, err := robot_mcp.NewServerID(string(request.McpServerId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	server, err := r.mcp.GetServer(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotMCPServerGet200JSONResponse{
		RobotMCPServerGetOKJSONResponse: openapi.RobotMCPServerGetOKJSONResponse(serialiseRobotMCPServer(server)),
	}, nil
}

func (r *Robots) RobotMCPServerUpdate(ctx context.Context, request openapi.RobotMCPServerUpdateRequestObject) (openapi.RobotMCPServerUpdateResponseObject, error) {
	id, err := robot_mcp.NewServerID(string(request.McpServerId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	if request.Body == nil {
		return nil, fault.Wrap(fault.New("missing request body"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	clearBearerToken := false
	if request.Body.ClearBearerToken != nil {
		clearBearerToken = *request.Body.ClearBearerToken
	}

	server, err := r.mcp.UpdateServer(ctx, id, robot_mcp.ServerUpdate{
		Name:             request.Body.Name,
		Description:      request.Body.Description,
		EndpointURL:      request.Body.EndpointUrl,
		Enabled:          request.Body.Enabled,
		BearerToken:      request.Body.BearerToken,
		ClearBearerToken: clearBearerToken,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotMCPServerUpdate200JSONResponse{
		RobotMCPServerUpdateOKJSONResponse: openapi.RobotMCPServerUpdateOKJSONResponse(serialiseRobotMCPServer(server)),
	}, nil
}

func (r *Robots) RobotMCPServerDelete(ctx context.Context, request openapi.RobotMCPServerDeleteRequestObject) (openapi.RobotMCPServerDeleteResponseObject, error) {
	id, err := robot_mcp.NewServerID(string(request.McpServerId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if err := r.mcp.DeleteServer(ctx, id); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotMCPServerDelete200Response{}, nil
}

func (r *Robots) RobotMCPServerRefresh(ctx context.Context, request openapi.RobotMCPServerRefreshRequestObject) (openapi.RobotMCPServerRefreshResponseObject, error) {
	id, err := robot_mcp.NewServerID(string(request.McpServerId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	server, err := r.mcp.Refresh(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	return openapi.RobotMCPServerRefresh200JSONResponse{
		RobotMCPServerRefreshOKJSONResponse: openapi.RobotMCPServerRefreshOKJSONResponse(serialiseRobotMCPServer(server)),
	}, nil
}

func deserialiseRobotMessageCursorParams(before *openapi.Identifier, limit *string) (robot.MessageCursorParams, error) {
	beforeID := opt.NewEmpty[robot.MessageID]()
	if before != nil {
		id, err := xid.FromString(string(*before))
		if err != nil {
			return robot.MessageCursorParams{}, err
		}
		beforeID = opt.New(robot.MessageID(id))
	}

	size := 50
	if limit != nil {
		parsed, err := strconv.ParseInt(*limit, 10, 32)
		if err != nil {
			return robot.MessageCursorParams{}, err
		}

		size = min(max(int(parsed), 1), 100)
	}

	return robot.NewMessageCursorParams(beforeID, size), nil
}

func (r *Robots) serialiseProviderStatuses(ctx context.Context, runtime llm_provider.RuntimeSettings) (openapi.RobotProviderStatusList, error) {
	providers := r.modelFactory.Providers()
	out := make(openapi.RobotProviderStatusList, 0, len(providers))
	for _, provider := range providers {
		status, models, err := r.modelFactory.ProviderStatus(ctx, provider)
		if err != nil {
			return nil, err
		}
		out = append(out, serialiseRobotProviderStatus(provider, runtime, status, models, r.modelFactory.RequiresAPIKey(provider)))
	}
	return out, nil
}

func (r *Robots) validateRobotModelForSave(ctx context.Context, model model_ref.ModelRef) error {
	runtime, err := r.modelFactory.RuntimeSettings(ctx)
	if err != nil {
		return err
	}

	if !runtime.Enabled {
		return nil
	}

	return r.modelFactory.EnsureModelAvailable(ctx, model)
}

func (r *Robots) validateRobotToolsForCreate(ctx context.Context, toolNames []string) error {
	var invalid []string
	for _, name := range toolNames {
		if !r.tools.HasTool(name) {
			invalid = append(invalid, name)
		}
	}
	if len(invalid) > 0 {
		return fault.New("invalid robot tool names: " + strings.Join(invalid, ", "))
	}
	return nil
}

func (r *Robots) validateRobotToolsForUpdate(ctx context.Context, robotID robot_ref.ID, toolNames []string) error {
	existing, err := r.robotQuerier.Get(ctx, robotID)
	if err != nil {
		return err
	}

	preserved := map[string]struct{}{}
	for _, name := range existing.Tools {
		preserved[string(name)] = struct{}{}
	}

	var invalid []string
	for _, name := range toolNames {
		if r.tools.HasTool(name) {
			continue
		}
		if _, ok := preserved[name]; ok {
			continue
		}
		invalid = append(invalid, name)
	}
	if len(invalid) > 0 {
		return fault.New("invalid robot tool names: " + strings.Join(invalid, ", "))
	}
	return nil
}

func (r *Robots) serialiseProviderStatus(ctx context.Context, provider model_ref.Provider) (openapi.RobotProviderStatus, error) {
	runtime, err := r.modelFactory.RuntimeSettings(ctx)
	if err != nil {
		return openapi.RobotProviderStatus{}, err
	}

	status, models, err := r.modelFactory.ProviderStatus(ctx, provider)
	if err != nil {
		return openapi.RobotProviderStatus{}, err
	}

	return serialiseRobotProviderStatus(provider, runtime, status, models, r.modelFactory.RequiresAPIKey(provider)), nil
}

func serialiseRobotProviderStatus(provider model_ref.Provider, runtime llm_provider.RuntimeSettings, status model_ref.CacheStatus, models []model_ref.Info, requiresAPIKey bool) openapi.RobotProviderStatus {
	providerSettings := runtime.Providers[provider]

	lastRefreshedAt := status.LastRefreshedAt.Ptr()
	stale := true
	if last, ok := status.LastRefreshedAt.Get(); ok {
		stale = time.Since(last) > llm_provider.ModelCacheTTL
	}

	return openapi.RobotProviderStatus{
		Provider:  openapi.RobotModelProvider(provider.String()),
		Supported: true,
		Settings: openapi.RobotProviderSettings{
			Enabled:        providerSettings.Enabled,
			HasApiKey:      providerSettings.APIKey != "",
			RequiresApiKey: requiresAPIKey,
		},
		Cache: openapi.RobotModelCacheStatus{
			LastRefreshedAt: lastRefreshedAt,
			LastError:       status.LastError.Ptr(),
			Stale:           stale,
		},
		Models: serialiseRobotModelInfos(models),
	}
}

func serialiseRobotModelInfos(models []model_ref.Info) openapi.RobotModelInfoList {
	return dt.Map(models, func(model model_ref.Info) openapi.RobotModelInfo {
		return openapi.RobotModelInfo{
			Ref:      openapi.RobotModelRef(model.String()),
			Provider: openapi.RobotModelProvider(model.Provider().String()),
			Model:    openapi.RobotModelName(model.Model().String()),
		}
	})
}

func serialiseRobotToolInfo(tool robot_tools.CatalogueTool) openapi.RobotToolInfo {
	return openapi.RobotToolInfo{
		Id:           tool.ID,
		CallableName: tool.CallableName,
		Name:         &tool.Name,
		Description:  tool.Description,
		Source:       openapi.RobotToolSource(tool.Source),
		Available:    tool.Available,
	}
}

func serialiseRobotMCPServer(server robot_mcp.Server) openapi.RobotMCPServer {
	var oauthRemoteConnectionID *openapi.Identifier
	if server.OAuthRemoteConnectionID != nil {
		id := openapi.Identifier(server.OAuthRemoteConnectionID.String())
		oauthRemoteConnectionID = &id
	}

	return openapi.RobotMCPServer{
		Id:                      openapi.Identifier(server.ID.String()),
		CreatedAt:               server.CreatedAt,
		UpdatedAt:               server.UpdatedAt,
		Name:                    server.Name,
		Slug:                    server.Slug,
		Description:             server.Description,
		EndpointUrl:             server.EndpointURL,
		OauthRemoteConnectionId: oauthRemoteConnectionID,
		Enabled:                 server.Enabled,
		HasBearerToken:          server.HasBearerToken,
		HasOauthToken:           server.HasOAuthAccessToken,
		LastRefreshedAt:         server.LastRefreshedAt,
		LastError:               server.LastError,
		Tools:                   serialiseRobotMCPTools(server.Tools),
	}
}

func serialiseRobotMCPServers(servers []robot_mcp.Server) openapi.RobotMCPServerList {
	return dt.Map(servers, serialiseRobotMCPServer)
}

func serialiseRobotMCPServerProbe(result mcpclient.ProbeResult) openapi.RobotMCPServerProbeResult {
	return openapi.RobotMCPServerProbeResult{
		InputUrl:      result.InputURL,
		EndpointUrl:   result.EndpointURL,
		ServerCardUrl: optionalString(result.ServerCardURL),
		ServerCard:    serialiseRobotMCPServerCard(result.ServerCard),
		RemoteType:    optionalString(result.RemoteType),
		Active:        result.Active,
		ProbeError:    optionalString(result.ProbeError),
	}
}

func serialiseRobotMCPServerCard(card *mcpclient.ServerCard) *openapi.RobotMCPServerCard {
	if card == nil {
		return nil
	}
	remotes := dt.Map(card.Remotes, func(remote mcpclient.ServerCardRemote) openapi.RobotMCPServerCardRemote {
		return openapi.RobotMCPServerCardRemote{
			Type:                      remote.Type,
			Url:                       remote.URL,
			SupportedProtocolVersions: optionalStringSlice(remote.SupportedProtocolVersions),
		}
	})
	return &openapi.RobotMCPServerCard{
		Name:        card.Name,
		Version:     card.Version,
		Description: card.Description,
		Title:       optionalString(card.Title),
		WebsiteUrl:  optionalString(card.WebsiteURL),
		Remotes:     &remotes,
	}
}

func serialiseRobotMCPTools(tools []robot_mcp.Tool) openapi.RobotMCPToolList {
	return dt.Map(tools, func(tool robot_mcp.Tool) openapi.RobotMCPTool {
		return openapi.RobotMCPTool{
			Id:           tool.ID,
			RemoteName:   tool.RemoteName,
			CallableName: tool.CallableName,
			Title:        &tool.Title,
			Description:  tool.Description,
			Enabled:      tool.Enabled,
			Available:    tool.Enabled,
			LastSeenAt:   tool.LastSeenAt,
		}
	})
}

func currentRobotSettings(current *settings.Settings) (settings.ServiceSettings, settings.RobotServiceSettings, map[string]settings.RobotProviderSettings) {
	services := settings.ServiceSettings{}
	if value, ok := current.Services.Get(); ok {
		services = value
	}

	robots := settings.RobotServiceSettings{}
	if value, ok := services.Robots.Get(); ok {
		robots = value
	}

	providers := map[string]settings.RobotProviderSettings{}
	if currentProviders, ok := robots.Providers.Get(); ok {
		for key, value := range currentProviders {
			providers[key] = value
		}
	}

	return services, robots, providers
}

func anyRobotProviderEnabled(runtime llm_provider.RuntimeSettings, providers map[string]settings.RobotProviderSettings, providerList []model_ref.Provider, updatedProvider model_ref.Provider, updatedEnabled bool) bool {
	for _, provider := range providerList {
		if provider == updatedProvider {
			if updatedEnabled {
				return true
			}
			continue
		}

		enabled := runtime.Providers[provider].Enabled
		if stored, ok := providers[provider.String()]; ok {
			enabled = stored.Enabled.Or(enabled)
		}

		if enabled {
			return true
		}
	}

	return false
}

func serialiseRobot(r *robot.Robot) openapi.Robot {
	return openapi.Robot{
		Id:          openapi.Identifier(r.ID.String()),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		Name:        r.Name,
		Description: r.Description,
		Playbook:    r.Playbook,
		Model:       r.Model.String(),
		WorkspaceId: serialiseNullableOpt(opt.Map(r.WorkspaceID, func(id xid.ID) openapi.NullableIdentifier {
			return openapi.NullableIdentifier(id.String())
		})),
		Author: serialiseProfileReferenceFromAccount(r.Author),
		Tools:  serialiseRobotToolNameList(r.Tools),
		Meta:   (*openapi.Metadata)(&r.Metadata),
	}
}

func serialiseRobotReference(r *robot.Robot) openapi.RobotReference {
	return openapi.RobotReference{
		Id:   openapi.Identifier(r.ID.String()),
		Name: r.Name,
	}
}

func serialiseRobotToolNameList(in []robot_ref.ToolName) openapi.RobotToolNameList {
	return dt.Map(in, func(tool robot_ref.ToolName) string {
		return string(tool)
	})
}

func serialiseRobots(robots []*robot.Robot) []openapi.Robot {
	return dt.Map(robots, func(r *robot.Robot) openapi.Robot {
		return serialiseRobot(r)
	})
}

func serialiseRobotWorkspace(w *robot.Workspace) openapi.RobotWorkspace {
	return openapi.RobotWorkspace{
		Id:          openapi.Identifier(w.ID.String()),
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
		Name:        w.Name,
		Description: w.Description,
		Provider:    openapi.RobotWorkspaceProvider(w.Provider),
		Config:      ensureMap(w.Config),
		Meta:        openapi.Metadata(ensureMap(w.Metadata)),
		CreatedBy:   serialiseProfileReferenceFromAccount(w.Creator),
	}
}

func serialiseRobotWorkspaceInstance(i *robot.WorkspaceInstance) openapi.RobotWorkspaceInstance {
	return openapi.RobotWorkspaceInstance{
		Id:            openapi.Identifier(i.ID.String()),
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
		WorkspaceId:   openapi.Identifier(i.WorkspaceID.String()),
		Provider:      openapi.RobotWorkspaceProvider(i.Provider),
		ProviderState: ensureMap(i.ProviderState),
		Meta:          openapi.Metadata(ensureMap(i.Metadata)),
		CreatedBy:     serialiseProfileReferenceFromAccount(i.Creator),
	}
}

func serialiseRobotWorkspaceMountPtr(mount *robot.WorkspaceMount) *openapi.RobotWorkspaceMount {
	if mount == nil {
		return nil
	}

	meta := openapi.Metadata(ensureMap(mount.Metadata))
	return &openapi.RobotWorkspaceMount{
		WorkspaceId:         openapi.Identifier(mount.WorkspaceID.String()),
		WorkspaceInstanceId: openapi.Identifier(mount.WorkspaceInstanceID.String()),
		Provider:            openapi.RobotWorkspaceProvider(mount.Provider),
		Meta:                &meta,
	}
}

func ensureMap(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{}
	}
	return in
}

// serialiseRobotSessionMessage converts a robot.Message (containing ADK Event)
// into the Vercel AI SDK UIMessage format for the frontend.
func serialiseRobotSessionMessage(m *robot.Message, hiddenToolCallIDs map[string]bool) (openapi.RobotSessionMessage, error) {
	parts, err := robotprojection.ADKEventToUIMessageParts(m.Event, hiddenToolCallIDs)
	if err != nil {
		return openapi.RobotSessionMessage{}, err
	}

	role := serialiseMessageRole(m.Event)

	msg := openapi.RobotSessionMessage{
		Id:        m.ID.String(),
		Role:      openapi.RobotSessionMessageRole(role),
		Parts:     parts,
		CreatedAt: m.CreatedAt,
		Robot:     serialiseRobotActorReference(m),
		Author:    opt.Map(m.Author, func(a *account.Account) openapi.ProfileReference { return serialiseProfileReferenceFromAccount(*a) }).Ptr(),
	}

	return msg, nil
}

func serialiseMessageRole(event adksession.Event) string {
	if event.Author == "user" {
		return "user"
	}
	return "assistant"
}

func serialiseRobotActorReference(m *robot.Message) *openapi.RobotReference {
	if r, ok := m.Robot.Get(); ok {
		ref := serialiseRobotReference(r)
		return &ref
	}

	actor, ok := m.Actor.Get()
	if !ok {
		return nil
	}

	if builtinID, ok := actor.BuiltinRobotID.Get(); ok {
		return &openapi.RobotReference{
			Id:   openapi.Identifier(builtinID.String()),
			Name: builtinRobotName(builtinID.String()),
		}
	}

	return nil
}

func builtinRobotName(id string) string {
	switch id {
	case agent_registry.RobotBuilderID, "storyden":
		return "Storyden Robot Builder"
	case agent_registry.PluginBuilderID:
		return "Plugin Builder"
	default:
		return id
	}
}
