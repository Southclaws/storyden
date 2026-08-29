package robot

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"google.golang.org/adk/v2/agent"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type invocationPrincipal struct {
	AccountID   account.AccountID
	Handle      string
	Name        string
	Kind        account.AccountKind
	Permissions []string
}

type invocationActor struct {
	AccountID account.AccountID
	Handle    string
	Name      string
	Kind      account.AccountKind
}

type globalInstructionContext struct {
	Mode                  RunMode
	Source                RunSource
	EventInvocation       bool
	CurrentTime           time.Time
	Principal             invocationPrincipal
	Actor                 opt.Optional[invocationActor]
	InvocationContextJSON string
}

func newInvocationPrincipal(acc account.Account, permissions rbac.Permissions) invocationPrincipal {
	permissionNames := dt.Map(permissions.List(), func(permission rbac.Permission) string {
		return permission.String()
	})
	slices.Sort(permissionNames)

	return invocationPrincipal{
		AccountID:   acc.ID,
		Handle:      acc.Handle,
		Name:        acc.Name,
		Kind:        acc.Kind,
		Permissions: permissionNames,
	}
}

func resolveInvocationActor(source RunSource, acc account.Account) opt.Optional[invocationActor] {
	if source != SourceInteractiveChat {
		return opt.NewEmpty[invocationActor]()
	}

	return opt.New(invocationActor{
		AccountID: acc.ID,
		Handle:    acc.Handle,
		Name:      acc.Name,
		Kind:      acc.Kind,
	})
}

func newGlobalInstructionContext(
	now time.Time,
	acc account.Account,
	permissions rbac.Permissions,
	invocationContext InvocationContext,
	options RunOptions,
) (globalInstructionContext, error) {
	invocationContextJSON := ""
	if len(invocationContext) > 0 {
		encoded, err := json.MarshalIndent(invocationContext, "", "  ")
		if err != nil {
			return globalInstructionContext{}, err
		}
		invocationContextJSON = string(encoded)
	}

	return globalInstructionContext{
		Mode:                  options.Mode,
		Source:                options.Source,
		EventInvocation:       options.Source == SourceScheduled && invocationContext[InvocationContextKeyTrailTriggerKind] == "event",
		CurrentTime:           now.UTC(),
		Principal:             newInvocationPrincipal(acc, permissions),
		Actor:                 resolveInvocationActor(options.Source, acc),
		InvocationContextJSON: invocationContextJSON,
	}, nil
}

func renderGlobalInstruction(ctx globalInstructionContext) (string, error) {
	return renderInstructionTemplate(globalInstructionTemplate, ctx)
}

func renderGlobalRuntimeContext(ctx globalInstructionContext) (string, error) {
	return renderNamedInstructionTemplate(globalInstructionTemplate, "runtime-context", ctx)
}

func (s *Agent) globalInstructionProvider(invocationContext InvocationContext, options RunOptions) func(ctx agent.ReadonlyContext) (string, error) {
	return func(ctx agent.ReadonlyContext) (string, error) {
		acc, err := session.GetAccount(ctx)
		if err != nil {
			return "", fault.Wrap(err, fctx.With(ctx))
		}
		permissions, err := session.GetPermissions(ctx)
		if err != nil {
			return "", fault.Wrap(err, fctx.With(ctx))
		}

		instructionContext, err := newGlobalInstructionContext(time.Now(), acc, permissions, invocationContext, options)
		if err != nil {
			return "", fault.Wrap(err, fctx.With(ctx))
		}

		instruction, err := renderGlobalInstruction(instructionContext)
		if err != nil {
			return "", fault.Wrap(err, fctx.With(ctx))
		}

		return instruction, nil
	}
}
