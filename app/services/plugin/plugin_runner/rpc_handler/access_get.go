package rpc_handler

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (h *Handler) handleAccessGet(ctx context.Context, req *rpc.RPCRequestAccessGet) (rpc.RPCResponseAccessGet, error) {
	h.mu.Lock()
	if h.cachedAccess != nil {
		cached := *h.cachedAccess
		h.mu.Unlock()
		return rpc.RPCResponseAccessGet{
			ID:      req.ID,
			Jsonrpc: "2.0",
			Method:  opt.New("access_get"),
			Result:  cached,
		}, nil
	}
	h.mu.Unlock()

	accessConfig, ok := h.manifest.Metadata.Access.Get()
	if !ok {
		return accessGetError(req, -32601, "manifest does not request access"), nil
	}

	pluginAccount, err := h.ensureAccessAccount(ctx, accessConfig)
	if err != nil {
		return accessGetError(req, -32603, err.Error()), nil
	}

	existingKeys, err := h.accessKeys.List(ctx, pluginAccount.ID)
	if err != nil {
		return accessGetError(req, -32603, err.Error()), nil
	}
	for _, existingKey := range existingKeys {
		if _, err := h.accessKeys.Revoke(ctx, pluginAccount.ID, existingKey.ID); err != nil {
			return accessGetError(req, -32603, err.Error()), nil
		}
	}

	createdKey, err := h.accessKeys.Create(
		ctx,
		pluginAccount.ID,
		access_key.AccessKeyKindBot,
		"Plugin API Access",
		opt.NewEmpty[time.Time](),
	)
	if err != nil {
		return accessGetError(req, -32603, err.Error()), nil
	}

	result := rpc.RPCResponseAccessGetResult{
		APIBaseURL: h.apiBaseURL,
		AccessKey:  createdKey.String(),
	}

	h.mu.Lock()
	h.cachedAccess = &result
	h.mu.Unlock()

	return rpc.RPCResponseAccessGet{
		ID:      req.ID,
		Jsonrpc: "2.0",
		Method:  opt.New("access_get"),
		Result:  result,
	}, nil
}

func (h *Handler) ensureAccessAccount(
	ctx context.Context,
	accessConfig rpc.ManifestAccess,
) (*account.AccountWithEdges, error) {
	pluginAccount, exists, err := h.accountQuerier.LookupByHandle(ctx, accessConfig.Handle)
	if err != nil {
		return nil, err
	}

	if !exists {
		pluginAccount, err = h.accountWriter.Create(
			ctx,
			accessConfig.Handle,
			account_writer.WithKind(account.AccountKindBot),
			account_writer.WithName(accessConfig.Name),
		)
		if err != nil {
			return nil, err
		}
	}

	if pluginAccount.Kind != account.AccountKindBot {
		return nil, fault.New("access account must be kind=bot")
	}

	updateMutations := []account_writer.Mutation{
		account_writer.SetName(accessConfig.Name),
	}

	if bio, ok := accessConfig.Bio.Get(); ok {
		updateMutations = append(updateMutations, account_writer.SetBio(bio))
	}
	if len(accessConfig.Links) > 0 {
		links := make([]account.ExternalLink, 0, len(accessConfig.Links))
		for _, link := range accessConfig.Links {
			links = append(links, account.ExternalLink{
				Text: link.Text,
				URL:  link.URL,
			})
		}
		updateMutations = append(updateMutations, account_writer.SetLinks(links))
	}
	if accessConfig.Metadata != nil {
		updateMutations = append(updateMutations, account_writer.SetMetadata(accessConfig.Metadata))
	}

	if len(updateMutations) > 0 {
		pluginAccount, err = h.accountWriter.Update(ctx, pluginAccount.ID, updateMutations...)
		if err != nil {
			return nil, err
		}
	}

	return pluginAccount, nil
}

func accessGetError(req *rpc.RPCRequestAccessGet, code int, message string) rpc.RPCResponseAccessGet {
	return rpc.RPCResponseAccessGet{
		ID:      req.ID,
		Jsonrpc: "2.0",
		Method:  opt.New("access_get"),
		Error: opt.New(rpc.RPCResponseAccessGetError{
			Code:    opt.New(code),
			Message: opt.New(message),
		}),
	}
}
