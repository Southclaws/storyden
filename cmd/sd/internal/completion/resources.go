package completion

import "strings"

// Sources are an allowlist of read-only discovery operations and safe display
// fields. Never infer a lookup from a mutation or expose arbitrary JSON fields.
type source struct {
	operation, collection, value, description string
	children                                  string
}

func pageSource(value string) source { return source{"NodeList", "nodes", value, "name", ""} }

func resourceSource(path, name string) (source, bool) {
	if path == "plugin dev new" {
		return source{}, false
	}
	switch name {
	case "node_slug":
		return pageSource("slug"), true
	case "slug", "parent", "target":
		if strings.HasPrefix(path, "page ") {
			return pageSource("slug"), true
		}
	case "before", "after":
		if path == "page move" {
			return pageSource("slug"), true
		}
	case "asset_id", "asset":
		if strings.HasPrefix(path, "page assets ") {
			return source{"NodeGet", "assets", "id", "filename", ""}, true
		}
	case "version_id":
		return source{"NodeVersionList", "versions", "id", "name", ""}, true
	case "run_id":
		return source{"TrailRunList", "runs", "id", "status", ""}, true
	case "action_run_id":
		return source{"TrailRunGet", "actions", "id", "status", ""}, true
	case "warning_id":
		return source{"AccountWarningList", "warnings", "id", "issued_at", ""}, true
	case "moderation_note_id":
		return source{"AccountModerationNoteList", "notes", "id", "created_at", ""}, true
	case "email_address_id":
		// AccountGet is the caller; it cannot discover another account's addresses.
		if strings.HasPrefix(path, "account ") {
			return source{"AccountGet", "email_addresses", "id", "", ""}, true
		}
	case "access_key_id":
		op := "AccessKeyList"
		if strings.HasPrefix(path, "admin ") {
			op = "AdminAccessKeyList"
		}
		return source{op, "keys", "id", "name", ""}, true
	case "oauth_client_id":
		op := "OAuthClientList"
		if strings.HasPrefix(path, "admin ") {
			op = "AdminOAuthClientList"
		}
		return source{op, "clients", "id", "name", ""}, true
	case "oauth_refresh_token_id":
		op := "OAuthRefreshTokenList"
		if strings.HasPrefix(path, "admin ") {
			op = "AdminOAuthRefreshTokenList"
		}
		return source{op, "tokens", "id", "client_name", ""}, true
	}
	sources := map[string]source{
		"node_id": pageSource("id"), "page_id": pageSource("id"),
		"thread_mark":                {"ThreadList", "threads", "id", "title", ""},
		"thread_id":                  {"ThreadList", "threads", "id", "title", ""},
		"plugin_instance_id":         {"PluginList", "plugins", "id", "name", ""},
		"account_id":                 {"ProfileList", "profiles", "id", "handle", ""},
		"account_handle":             {"ProfileList", "profiles", "handle", "name", ""},
		"author":                     {"ProfileList", "profiles", "handle", "name", ""},
		"authors":                    {"ProfileList", "profiles", "handle", "name", ""},
		"owner_handle":               {"ProfileList", "profiles", "handle", "name", ""},
		"invited_by":                 {"ProfileList", "profiles", "handle", "name", ""},
		"category_slug":              {"CategoryList", "categories", "slug", "name", "children"},
		"category":                   {"CategoryList", "categories", "id", "name", "children"},
		"categories":                 {"CategoryList", "categories", "slug", "name", "children"},
		"tags":                       {"TagList", "tags", "name", "", ""},
		"tag_name":                   {"TagList", "tags", "name", "", ""},
		"role_id":                    {"RoleList", "roles", "id", "name", ""},
		"roles":                      {"RoleList", "roles", "id", "name", ""},
		"invitation_id":              {"InvitationList", "invitations", "id", "", ""},
		"notification_id":            {"NotificationList", "notifications", "id", "event", ""},
		"report_id":                  {"ReportList", "reports", "id", "target_kind", ""},
		"collection_mark":            {"CollectionList", "collections", "id", "name", ""},
		"event_mark":                 {"EventList", "events", "slug", "name", ""},
		"link_slug":                  {"LinkList", "links", "slug", "title", ""},
		"robot_id":                   {"RobotsList", "robots", "id", "name", ""},
		"robot":                      {"RobotsList", "robots", "id", "name", ""},
		"toolset_id":                 {"RobotToolsetsList", "toolsets", "id", "name", ""},
		"provider":                   {"RobotProvidersList", "providers", "provider", "", ""},
		"workspace_id":               {"RobotWorkspacesList", "workspaces", "id", "name", ""},
		"workspace_instance_id":      {"RobotWorkspaceInstancesList", "workspace_instances", "id", "provider", ""},
		"session_id":                 {"RobotSessionsList", "sessions", "id", "name", ""},
		"session":                    {"RobotSessionsList", "sessions", "id", "name", ""},
		"mcp_server_id":              {"RobotMCPServersList", "servers", "id", "name", ""},
		"trail_id":                   {"TrailList", "trails", "id", "name", ""},
		"audit_event_id":             {"AuditEventList", "events", "id", "type", ""},
		"email_id":                   {"EmailQueueList", "emails", "id", "status", ""},
		"oauth_remote_connection_id": {"OAuthRemoteConnectionList", "connections", "id", "resource_name", ""},
		"oauth_dcr_iat_id":           {"AdminOAuthDynamicRegistrationAccessTokenList", "iats", "id", "label", ""},
		"auth_method_id":             {"AccountAuthProviderList", "active", "id", "name", ""},
	}
	src, ok := sources[name]
	if name == "categories" && path == "search" {
		src.value = "id"
	}
	if name == "account_id" && strings.HasPrefix(path, "admin ") {
		src = source{"AccountList", "accounts", "id", "handle", ""}
	}
	return src, ok
}
