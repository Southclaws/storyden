package robot

import (
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adkagent "google.golang.org/adk/v2/agent"

	"github.com/Southclaws/storyden/app/resources/robot/robot_memory"
)

type robotIdentity struct {
	ID               opt.Optional[xid.ID]
	Name             string
	Description      string
	Capabilities     []string
	UnavailableTools []string
}

type robotIdentityContext struct {
	Current    robotIdentity
	Delegated  bool
	MemoryRoot opt.Optional[robotMemoryRoot]
}

type robotMemoryRoot struct {
	Items   []robotMemoryRootItem
	HasMore bool
}

type robotMemoryRootItem struct {
	ID       string
	Children int
}

type robotInstructionContext struct {
	robotIdentityContext
	PlaybookLines []string
}

func mapMemoryRoot(items []*robot_memory.Item, hasMore bool) robotMemoryRoot {
	return robotMemoryRoot{
		Items: dt.Map(items, func(item *robot_memory.Item) robotMemoryRootItem {
			return robotMemoryRootItem{
				ID:       item.Memory.ID.String(),
				Children: item.Children,
			}
		}),
		HasMore: hasMore,
	}
}

func newRobotIdentityContext(current robotIdentity, delegated bool, memoryRoot opt.Optional[robotMemoryRoot]) robotIdentityContext {
	current.Description = strings.TrimSpace(current.Description)
	return robotIdentityContext{Current: current, Delegated: delegated, MemoryRoot: memoryRoot}
}

func renderRobotInstruction(ctx robotIdentityContext, playbook string) (string, error) {
	playbook = strings.ReplaceAll(playbook, "\r", `\r`)

	return renderInstructionTemplate(identityInstructionTemplate, robotInstructionContext{
		robotIdentityContext: ctx,
		PlaybookLines:        strings.Split(playbook, "\n"),
	})
}

func robotInstructionProvider(spec *resolvedAgentSpec, identity robotIdentityContext) func(adkagent.ReadonlyContext) (string, error) {
	return func(ctx adkagent.ReadonlyContext) (string, error) {
		instruction := spec.Instruction
		if spec.InstructionProvider != nil {
			resolved, err := spec.InstructionProvider(ctx)
			if err != nil {
				return "", err
			}
			instruction = resolved
		}

		return renderRobotInstruction(identity, instruction)
	}
}
