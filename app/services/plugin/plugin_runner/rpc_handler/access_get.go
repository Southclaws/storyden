package rpc_handler

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_writer"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

const (
	pluginAccessRoleColour = "hsl(157, 65%, 44%)"

	pluginAccessRoleMetaInstallationID = "plugin.installation_id"
	pluginAccessRoleMetaAccountID      = "plugin.account_id"
	pluginAccessRoleMetaManaged        = "plugin.managed"
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
	if err := h.ensureAccessRole(ctx, pluginAccount, accessConfig); err != nil {
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
	handle := pluginAccessHandle(accessConfig.Handle, h.installationID)

	pluginAccount, exists, err := h.accountQuerier.LookupByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}

	if !exists {
		pluginAccount, err = h.accountWriter.Create(
			ctx,
			handle,
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

func (h *Handler) ensureAccessRole(
	ctx context.Context,
	pluginAccount *account.AccountWithEdges,
	accessConfig rpc.ManifestAccess,
) error {
	perms, err := rbac.NewPermissions(accessConfig.Permissions)
	if err != nil {
		return err
	}

	roleName := pluginAccessRoleName(pluginAccount.Name, pluginAccount.Handle)
	roleMeta := map[string]any{
		pluginAccessRoleMetaInstallationID: h.installationID.String(),
		pluginAccessRoleMetaAccountID:      pluginAccount.ID.String(),
		pluginAccessRoleMetaManaged:        true,
	}

	accessRole, err := h.getOrCreateAccessRole(ctx, roleName, perms.List(), roleMeta)
	if err != nil {
		return err
	}

	accessRole, err = h.roleWriter.Update(
		ctx,
		accessRole.ID,
		role_writer.WithName(roleName),
		role_writer.WithPermissions(perms.List()),
		role_writer.WithMeta(roleMeta),
	)
	if err != nil {
		return err
	}

	if accountHasRole(pluginAccount, accessRole.ID) {
		return nil
	}

	err = h.roleAssigner.UpdateRoles(ctx, pluginAccount.ID, role_assign.Add(accessRole.ID))
	if err == nil {
		return nil
	}
	if !ent.IsConstraintError(err) {
		return err
	}

	// Concurrent access_get requests can race assigning the same role. Recheck
	// the account before returning the constraint violation.
	pluginAccount, lookupErr := h.accountQuerier.GetByID(ctx, pluginAccount.ID)
	if lookupErr != nil {
		return lookupErr
	}
	if accountHasRole(pluginAccount, accessRole.ID) {
		return nil
	}

	return err
}

func (h *Handler) getOrCreateAccessRole(
	ctx context.Context,
	roleName string,
	perms rbac.PermissionList,
	meta map[string]any,
) (*role.Role, error) {
	roles, err := h.roleQuerier.List(ctx)
	if err != nil {
		return nil, err
	}

	if existing := findManagedAccessRole(roles, h.installationID); existing != nil {
		return existing, nil
	}

	created, err := h.roleWriter.Create(
		ctx,
		roleName,
		pluginAccessRoleColour,
		perms,
		role_writer.WithMeta(meta),
	)
	if err == nil {
		return created, nil
	}
	if !ent.IsConstraintError(err) {
		return nil, err
	}

	// Another request likely created this role concurrently.
	roles, listErr := h.roleQuerier.List(ctx)
	if listErr != nil {
		return nil, listErr
	}

	if existing := findManagedAccessRole(roles, h.installationID); existing != nil {
		return existing, nil
	}

	return nil, err
}

func findManagedAccessRole(
	roles role.Roles,
	installationID plugin.InstallationID,
) *role.Role {
	installationIDStr := installationID.String()

	for _, candidate := range roles {
		metadataInstallationID, ok := metadataString(candidate.Metadata, pluginAccessRoleMetaInstallationID)
		if ok && metadataInstallationID == installationIDStr {
			return candidate
		}
	}

	return nil
}

func metadataString(meta map[string]any, key string) (string, bool) {
	if meta == nil {
		return "", false
	}

	raw, ok := meta[key]
	if !ok {
		return "", false
	}

	value, ok := raw.(string)
	if !ok {
		return "", false
	}

	return value, true
}

func accountHasRole(acc *account.AccountWithEdges, roleID role.RoleID) bool {
	if acc == nil {
		return false
	}

	for _, heldRole := range acc.Roles {
		if heldRole.ID == roleID {
			return true
		}
	}

	return false
}

func pluginAccessRoleName(displayName, handle string) string {
	name := strings.TrimSpace(displayName)
	if name == "" {
		name = "Plugin Bot"
	}

	return fmt.Sprintf("%s (Bot %s)", name, pluginAccessRoleShortID(handle))
}

func pluginAccessRoleShortID(handle string) string {
	handle = strings.ToLower(strings.TrimSpace(handle))
	if handle == "" {
		return "0000"
	}

	clean := strings.Builder{}
	clean.Grow(len(handle))

	for _, r := range handle {
		if unicode.IsLower(r) || unicode.IsDigit(r) {
			clean.WriteRune(r)
		}
	}

	shortID := clean.String()
	if len(shortID) >= 4 {
		return shortID[len(shortID)-4:]
	}

	if shortID == "" {
		return "0000"
	}

	for len(shortID) < 4 {
		shortID += "0"
	}

	return shortID
}

func pluginAccessHandle(base string, installationID plugin.InstallationID) string {
	// ref: app/resources/account/validation.go
	// TODO: Make global const
	const maxHandleLen = 30
	const defaultBase = "plugin"

	base = strings.TrimSpace(strings.ToLower(base))
	if base == "" {
		base = defaultBase
	}

	clean := strings.Builder{}
	clean.Grow(len(base))
	lastWasSep := false

	for _, r := range base {
		valid := unicode.IsLower(r) || unicode.IsDigit(r) || r == '-' || r == '_'
		if !valid {
			if lastWasSep {
				continue
			}
			clean.WriteRune('-')
			lastWasSep = true
			continue
		}

		clean.WriteRune(r)
		lastWasSep = r == '-' || r == '_'
	}

	sanitised := strings.Trim(clean.String(), "-_")
	if sanitised == "" {
		sanitised = defaultBase
	}

	suffix := installationID.String()
	if len(suffix) > 4 {
		suffix = suffix[len(suffix)-4:]
	}

	maxBaseLen := maxHandleLen - 1 - len(suffix)
	if maxBaseLen < 1 {
		maxBaseLen = 1
	}
	if len(sanitised) > maxBaseLen {
		sanitised = strings.TrimRight(sanitised[:maxBaseLen], "-_")
		if sanitised == "" {
			sanitised = defaultBase
			if len(sanitised) > maxBaseLen {
				sanitised = sanitised[:maxBaseLen]
			}
		}
	}

	return sanitised + "-" + suffix
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
