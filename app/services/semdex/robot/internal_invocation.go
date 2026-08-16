package robot

import (
	"context"
	"crypto/sha256"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"google.golang.org/genai"
)

type InternalInvocation struct {
	InputID  robotresource.InputID
	RobotRef string
	Content  *genai.Content
	Options  RunOptions
}

type InternalInvocationEnqueuer func(context.Context, InternalInvocation) error

type internalInvocationEnqueuerKey struct{}
type delegationAttributionKey struct{}
type hiddenRuntimeInputKey struct{}

type delegationAttribution struct {
	CallID   string
	ToolName string
}

func WithInternalInvocationEnqueuer(ctx context.Context, enqueue InternalInvocationEnqueuer) context.Context {
	return context.WithValue(ctx, internalInvocationEnqueuerKey{}, enqueue)
}

func enqueueInternalInvocation(ctx context.Context, invocation InternalInvocation) error {
	enqueue, ok := ctx.Value(internalInvocationEnqueuerKey{}).(InternalInvocationEnqueuer)
	if !ok || enqueue == nil {
		return ErrInternalInvocationUnavailable
	}
	return enqueue(ctx, invocation)
}

func withDelegationAttribution(ctx context.Context, run *DelegationRun) context.Context {
	if run == nil {
		return ctx
	}
	return context.WithValue(ctx, delegationAttributionKey{}, delegationAttribution{
		CallID:   run.CallID,
		ToolName: run.ToolName,
	})
}

func delegationAttributionFromContext(ctx context.Context) (delegationAttribution, bool) {
	attribution, ok := ctx.Value(delegationAttributionKey{}).(delegationAttribution)
	return attribution, ok && attribution.CallID != "" && attribution.ToolName != ""
}

func withHiddenRuntimeInput(ctx context.Context) context.Context {
	return context.WithValue(ctx, hiddenRuntimeInputKey{}, true)
}

func hasHiddenRuntimeInput(ctx context.Context) bool {
	hidden, _ := ctx.Value(hiddenRuntimeInputKey{}).(bool)
	return hidden
}

func InternalInvocationID(namespace, sessionID, callID string) robotresource.InputID {
	sum := sha256.Sum256([]byte(namespace + "\x00" + sessionID + "\x00" + callID))
	var id [12]byte
	copy(id[:], sum[:len(id)])
	return robotresource.InputID(id)
}
