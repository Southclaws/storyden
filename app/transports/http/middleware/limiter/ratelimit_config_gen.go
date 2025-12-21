// Package limiter contains rate limiting middleware.
// THIS FILE IS GENERATED. DO NOT EDIT MANUALLY.
// To edit rate limit configuration, edit the OpenAPI spec x-storyden extensions and run codegen.
package limiter

import "time"

// OperationRateLimitConfig defines per-operation rate limiting configuration
type OperationRateLimitConfig struct {
	// Cost is the number of requests this operation counts as
	Cost int
	// Limit is the maximum number of requests allowed in the period (0 means use global default)
	Limit int
	// Period is the time window for the limit (empty means use global default)
	Period time.Duration
}

// OperationRateLimits contains per-operation rate limit configurations extracted from OpenAPI spec
var OperationRateLimits = map[string]OperationRateLimitConfig{
	"AuthEmailSignup": {
		Cost:   10,
		Limit:  10,
		Period: time.Duration(int64(3600000000000)),
	},
	"AuthPasswordReset": {
		Cost:   5,
		Limit:  10,
		Period: time.Duration(int64(3600000000000)),
	},
	"AuthPasswordSignup": {
		Cost:   5,
		Limit:  20,
		Period: time.Duration(int64(3600000000000)),
	},
}

// RouteToOperation maps HTTP method and path to operation ID
var RouteToOperation = map[string]string{
	"delete:/api/accounts/:account_handle/roles/:role_id":       "AccountRemoveRole",
	"delete:/api/accounts/:account_handle/roles/:role_id/badge": "AccountRoleRemoveBadge",
	"delete:/api/accounts/self/auth-methods/:auth_method_id":    "AccountAuthMethodDelete",
	"delete:/api/accounts/self/emails/:email_address_id":        "AccountEmailRemove",
	"delete:/api/admin/access-keys/:access_key_id":              "AdminAccessKeyDelete",
	"delete:/api/admin/bans/:account_handle":                    "AdminAccountBanRemove",
	"delete:/api/auth/access-keys/:access_key_id":               "AccessKeyDelete",
	"delete:/api/categories/:category_slug":                     "CategoryDelete",
	"delete:/api/collections/:collection_mark":                  "CollectionDelete",
	"delete:/api/collections/:collection_mark/nodes/:node_id":   "CollectionRemoveNode",
	"delete:/api/collections/:collection_mark/posts/:post_id":   "CollectionRemovePost",
	"delete:/api/events/:event_mark":                            "EventDelete",
	"delete:/api/events/:event_mark/participants/:account_id":   "EventParticipantRemove",
	"delete:/api/invitations/:invitation_id":                    "InvitationDelete",
	"delete:/api/likes/posts/:post_id":                          "LikePostRemove",
	"delete:/api/nodes/:node_slug":                              "NodeDelete",
	"delete:/api/nodes/:node_slug/assets/:asset_id":             "NodeRemoveAsset",
	"delete:/api/nodes/:node_slug/nodes/:node_slug_child":       "NodeRemoveNode",
	"delete:/api/posts/:post_id":                                "PostDelete",
	"delete:/api/posts/:post_id/reacts/:react_id":               "PostReactRemove",
	"delete:/api/profiles/:account_handle/followers":            "ProfileFollowersRemove",
	"delete:/api/roles/:role_id":                                "RoleDelete",
	"delete:/api/threads/:thread_mark":                          "ThreadDelete",
	"get:/api/accounts":                                         "AccountGet",
	"get:/api/accounts/:account_handle/avatar":                  "AccountGetAvatar",
	"get:/api/accounts/:account_id":                             "AccountView",
	"get:/api/accounts/self/auth-methods":                       "AccountAuthProviderList",
	"get:/api/admin/access-keys":                                "AdminAccessKeyList",
	"get:/api/assets/:asset_filename":                           "AssetGet",
	"get:/api/auth":                                             "AuthProviderList",
	"get:/api/auth/access-keys":                                 "AccessKeyList",
	"get:/api/auth/logout":                                      "AuthProviderLogout",
	"get:/api/auth/webauthn/assert/:account_handle":             "WebAuthnGetAssertion",
	"get:/api/auth/webauthn/make/:account_handle":               "WebAuthnRequestCredential",
	"get:/api/categories":                                       "CategoryList",
	"get:/api/categories/:category_slug":                        "CategoryGet",
	"get:/api/collections":                                      "CollectionList",
	"get:/api/collections/:collection_mark":                     "CollectionGet",
	"get:/api/datagraph":                                        "DatagraphSearch",
	"get:/api/datagraph/ask":                                    "DatagraphAsk",
	"get:/api/datagraph/matches":                                "DatagraphMatches",
	"get:/api/docs":                                             "GetDocs",
	"get:/api/events":                                           "EventList",
	"get:/api/events/:event_mark":                               "EventGet",
	"get:/api/info":                                             "GetInfo",
	"get:/api/info/banner":                                      "BannerGet",
	"get:/api/info/icon/:icon_size":                             "IconGet",
	"get:/api/invitations":                                      "InvitationList",
	"get:/api/invitations/:invitation_id":                       "InvitationGet",
	"get:/api/likes/posts/:post_id":                             "LikePostGet",
	"get:/api/likes/profiles/:account_handle":                   "LikeProfileGet",
	"get:/api/links":                                            "LinkList",
	"get:/api/links/:link_slug":                                 "LinkGet",
	"get:/api/nodes":                                            "NodeList",
	"get:/api/nodes/:node_slug":                                 "NodeGet",
	"get:/api/nodes/:node_slug/children":                        "NodeListChildren",
	"get:/api/notifications":                                    "NotificationList",
	"get:/api/openapi.json":                                     "GetSpec",
	"get:/api/posts/location":                                   "PostLocationGet",
	"get:/api/profiles":                                         "ProfileList",
	"get:/api/profiles/:account_handle":                         "ProfileGet",
	"get:/api/profiles/:account_handle/followers":               "ProfileFollowersGet",
	"get:/api/profiles/:account_handle/following":               "ProfileFollowingGet",
	"get:/api/reports":                                          "ReportList",
	"get:/api/roles":                                            "RoleList",
	"get:/api/roles/:role_id":                                   "RoleGet",
	"get:/api/tags":                                             "TagList",
	"get:/api/tags/:tag_name":                                   "TagGet",
	"get:/api/threads":                                          "ThreadList",
	"get:/api/threads/:thread_mark":                             "ThreadGet",
	"get:/api/version":                                          "GetVersion",
	"patch:/api/accounts":                                       "AccountUpdate",
	"patch:/api/admin":                                          "AdminSettingsUpdate",
	"patch:/api/auth/password":                                  "AuthPasswordUpdate",
	"patch:/api/categories/:category_slug":                      "CategoryUpdate",
	"patch:/api/categories/:category_slug/position":             "CategoryUpdatePosition",
	"patch:/api/collections/:collection_mark":                   "CollectionUpdate",
	"patch:/api/events/:event_mark":                             "EventUpdate",
	"patch:/api/nodes/:node_slug":                               "NodeUpdate",
	"patch:/api/nodes/:node_slug/children/property-schema":      "NodeUpdateChildrenPropertySchema",
	"patch:/api/nodes/:node_slug/position":                      "NodeUpdatePosition",
	"patch:/api/nodes/:node_slug/properties":                    "NodeUpdateProperties",
	"patch:/api/nodes/:node_slug/property-schema":               "NodeUpdatePropertySchema",
	"patch:/api/nodes/:node_slug/visibility":                    "NodeUpdateVisibility",
	"patch:/api/notifications":                                  "NotificationUpdateMany",
	"patch:/api/notifications/:notification_id":                 "NotificationUpdate",
	"patch:/api/posts/:post_id":                                 "PostUpdate",
	"patch:/api/reports/:report_id":                             "ReportUpdate",
	"patch:/api/roles/:role_id":                                 "RoleUpdate",
	"patch:/api/threads/:thread_mark":                           "ThreadUpdate",
	"post:/api/accounts/self/avatar":                            "AccountSetAvatar",
	"post:/api/accounts/self/emails":                            "AccountEmailAdd",
	"post:/api/admin/bans/:account_handle":                      "AdminAccountBanCreate",
	"post:/api/assets":                                          "AssetUpload",
	"post:/api/auth/access-keys":                                "AccessKeyCreate",
	"post:/api/auth/email-password/reset":                       "AuthPasswordResetRequestEmail",
	"post:/api/auth/email-password/signin":                      "AuthEmailPasswordSignin",
	"post:/api/auth/email-password/signup":                      "AuthEmailPasswordSignup",
	"post:/api/auth/email/signin":                               "AuthEmailSignin",
	"post:/api/auth/email/signup":                               "AuthEmailSignup",
	"post:/api/auth/email/verify":                               "AuthEmailVerify",
	"post:/api/auth/oauth/:oauth_provider/callback":             "OAuthProviderCallback",
	"post:/api/auth/password":                                   "AuthPasswordCreate",
	"post:/api/auth/password/reset":                             "AuthPasswordReset",
	"post:/api/auth/password/signin":                            "AuthPasswordSignin",
	"post:/api/auth/password/signup":                            "AuthPasswordSignup",
	"post:/api/auth/phone":                                      "PhoneRequestCode",
	"post:/api/auth/webauthn/assert":                            "WebAuthnMakeAssertion",
	"post:/api/auth/webauthn/make":                              "WebAuthnMakeCredential",
	"post:/api/beacon":                                          "SendBeacon",
	"post:/api/categories":                                      "CategoryCreate",
	"post:/api/collections":                                     "CollectionCreate",
	"post:/api/events":                                          "EventCreate",
	"post:/api/info/banner":                                     "BannerUpload",
	"post:/api/info/icon":                                       "IconUpload",
	"post:/api/invitations":                                     "InvitationCreate",
	"post:/api/links":                                           "LinkCreate",
	"post:/api/nodes":                                           "NodeCreate",
	"post:/api/nodes/:node_slug/content":                        "NodeGenerateContent",
	"post:/api/nodes/:node_slug/tags":                           "NodeGenerateTags",
	"post:/api/nodes/:node_slug/title":                          "NodeGenerateTitle",
	"post:/api/reports":                                         "ReportCreate",
	"post:/api/roles":                                           "RoleCreate",
	"post:/api/threads":                                         "ThreadCreate",
	"post:/api/threads/:thread_mark/replies":                    "ReplyCreate",
	"put:/api/accounts/:account_handle/roles/:role_id":          "AccountAddRole",
	"put:/api/accounts/:account_handle/roles/:role_id/badge":    "AccountRoleSetBadge",
	"put:/api/auth/phone/:account_handle":                       "PhoneSubmitCode",
	"put:/api/collections/:collection_mark/nodes/:node_id":      "CollectionAddNode",
	"put:/api/collections/:collection_mark/posts/:post_id":      "CollectionAddPost",
	"put:/api/events/:event_mark/participants/:account_id":      "EventParticipantUpdate",
	"put:/api/likes/posts/:post_id":                             "LikePostAdd",
	"put:/api/nodes/:node_slug/assets/:asset_id":                "NodeAddAsset",
	"put:/api/nodes/:node_slug/nodes/:node_slug_child":          "NodeAddNode",
	"put:/api/posts/:post_id/reacts":                            "PostReactAdd",
	"put:/api/profiles/:account_handle/followers":               "ProfileFollowersAdd",
}

// GetOperationConfig returns the rate limit config for an operation, or nil if not configured
func GetOperationConfig(operationID string) *OperationRateLimitConfig {
	if cfg, ok := OperationRateLimits[operationID]; ok {
		return &cfg
	}
	return nil
}

// GetOperationIDFromRoute returns the operation ID for a given route
func GetOperationIDFromRoute(method string, path string) string {
	key := method + ":" + path
	if opID, ok := RouteToOperation[key]; ok {
		return opID
	}
	return ""
}
