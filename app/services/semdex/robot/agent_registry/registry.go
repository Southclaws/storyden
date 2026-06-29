package agent_registry

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adkagent "google.golang.org/adk/agent"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/genai"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
)

const (
	RobotBuilderID  = "robot_builder"
	PluginBuilderID = "plugin_builder"
)

type RunMode string

const (
	ModeInteractive RunMode = "interactive"
	ModeUnattended  RunMode = "unattended"
)

type RunSource string

const (
	SourceInteractiveChat RunSource = "interactive_chat"
	SourcePluginRPC       RunSource = "plugin_rpc"
)

type WorkspaceMountSpec struct {
	WorkspaceID         opt.Optional[robotresource.WorkspaceID]
	WorkspaceInstanceID opt.Optional[robotresource.WorkspaceInstanceID]
	Metadata            map[string]any
}

type RunOptions struct {
	Mode   RunMode
	Source RunSource

	Workspace opt.Optional[WorkspaceMountSpec]
}

type ADKRunRequest struct {
	AppName              string
	Agent                adkagent.Agent
	RobotRef             string
	UserID               string
	SessionID            string
	Content              *genai.Content
	DefaultWorkspaceID   opt.Optional[xid.ID]
	WorkspaceRequirement *Definition
	Options              RunOptions
}

type ToolsetBuilder func(context.Context) (adktool.Toolset, error)

type Definition struct {
	ID                  string
	Name                string
	Description         string
	RequiresWorkspace   bool
	Hidden              bool
	AppName             string
	AgentName           string
	Instruction         string
	InstructionProvider func(adkagent.ReadonlyContext) (string, error)
	ToolNames           []string
	ToolsetBuilders     []ToolsetBuilder
	Capabilities        []string
}

type Registry struct {
	logger *slog.Logger
	mu     sync.RWMutex
	agents map[string]Definition
}

func New(logger *slog.Logger) *Registry {
	if logger == nil {
		logger = slog.Default()
	}
	return &Registry{
		logger: logger,
		agents: make(map[string]Definition),
	}
}

func (r *Registry) Register(def Definition) error {
	if def.ID == "" {
		return fmt.Errorf("agent definition ID is required")
	}
	if def.Name == "" {
		return fmt.Errorf("agent definition %q name is required", def.ID)
	}
	if def.AppName == "" {
		return fmt.Errorf("agent definition %q app name is required", def.ID)
	}
	if def.AgentName == "" {
		return fmt.Errorf("agent definition %q agent name is required", def.ID)
	}
	if def.Instruction == "" && def.InstructionProvider == nil {
		return fmt.Errorf("agent definition %q instruction is required", def.ID)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[def.ID]; exists {
		return fmt.Errorf("agent definition %q already registered", def.ID)
	}

	r.agents[def.ID] = def
	r.logger.Debug("registered robot agent", slog.String("robot_id", def.ID), slog.String("name", def.Name))

	return nil
}

func (r *Registry) Get(id string) (Definition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.agents[id]
	return def, ok
}

func (r *Registry) List(includeHidden bool) []Definition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Definition, 0, len(r.agents))
	for _, def := range r.agents {
		if def.Hidden && !includeHidden {
			continue
		}
		out = append(out, def)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})

	return out
}
