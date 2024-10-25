package openapi_rbac

import "github.com/Southclaws/storyden/app/resources/rbac"

type OperationPermissions interface {
	GetVersion() (bool, *rbac.Permission)
	GetSpec() (bool, *rbac.Permission)
	GetInfo() (bool, *rbac.Permission)
	IconGet() (bool, *rbac.Permission)
	IconUpload() (bool, *rbac.Permission)
	AdminSettingsUpdate() (bool, *rbac.Permission)
	AdminAccountBanCreate() (bool, *rbac.Permission)
	AdminAccountBanRemove() (bool, *rbac.Permission)
	RoleCreate() (bool, *rbac.Permission)
	RoleList() (bool, *rbac.Permission)
	RoleGet() (bool, *rbac.Permission)
	RoleUpdate() (bool, *rbac.Permission)
	RoleDelete() (bool, *rbac.Permission)
	AuthProviderList() (bool, *rbac.Permission)
	AuthPasswordSignup() (bool, *rbac.Permission)
	AuthPasswordSignin() (bool, *rbac.Permission)
	AuthPasswordCreate() (bool, *rbac.Permission)
	AuthPasswordUpdate() (bool, *rbac.Permission)
	AuthEmailPasswordSignup() (bool, *rbac.Permission)
	AuthEmailPasswordSignin() (bool, *rbac.Permission)
	AuthEmailSignup() (bool, *rbac.Permission)
	AuthEmailSignin() (bool, *rbac.Permission)
	AuthEmailVerify() (bool, *rbac.Permission)
	OAuthProviderCallback() (bool, *rbac.Permission)
	WebAuthnRequestCredential() (bool, *rbac.Permission)
	WebAuthnMakeCredential() (bool, *rbac.Permission)
	WebAuthnGetAssertion() (bool, *rbac.Permission)
	WebAuthnMakeAssertion() (bool, *rbac.Permission)
	PhoneRequestCode() (bool, *rbac.Permission)
	PhoneSubmitCode() (bool, *rbac.Permission)
	AuthProviderLogout() (bool, *rbac.Permission)
	AccountGet() (bool, *rbac.Permission)
	AccountUpdate() (bool, *rbac.Permission)
	AccountAuthProviderList() (bool, *rbac.Permission)
	AccountAuthMethodDelete() (bool, *rbac.Permission)
	AccountSetAvatar() (bool, *rbac.Permission)
	AccountGetAvatar() (bool, *rbac.Permission)
	AccountAddRole() (bool, *rbac.Permission)
	AccountRemoveRole() (bool, *rbac.Permission)
	AccountRoleSetBadge() (bool, *rbac.Permission)
	AccountRoleRemoveBadge() (bool, *rbac.Permission)
	InvitationList() (bool, *rbac.Permission)
	InvitationCreate() (bool, *rbac.Permission)
	InvitationGet() (bool, *rbac.Permission)
	InvitationDelete() (bool, *rbac.Permission)
	NotificationList() (bool, *rbac.Permission)
	NotificationUpdate() (bool, *rbac.Permission)
	ProfileList() (bool, *rbac.Permission)
	ProfileGet() (bool, *rbac.Permission)
	ProfileFollowersGet() (bool, *rbac.Permission)
	ProfileFollowersAdd() (bool, *rbac.Permission)
	ProfileFollowersRemove() (bool, *rbac.Permission)
	ProfileFollowingGet() (bool, *rbac.Permission)
	CategoryCreate() (bool, *rbac.Permission)
	CategoryList() (bool, *rbac.Permission)
	CategoryUpdateOrder() (bool, *rbac.Permission)
	CategoryUpdate() (bool, *rbac.Permission)
	TagList() (bool, *rbac.Permission)
	ThreadCreate() (bool, *rbac.Permission)
	ThreadList() (bool, *rbac.Permission)
	ThreadGet() (bool, *rbac.Permission)
	ThreadUpdate() (bool, *rbac.Permission)
	ThreadDelete() (bool, *rbac.Permission)
	ReplyCreate() (bool, *rbac.Permission)
	PostUpdate() (bool, *rbac.Permission)
	PostDelete() (bool, *rbac.Permission)
	PostSearch() (bool, *rbac.Permission)
	PostReactAdd() (bool, *rbac.Permission)
	PostReactRemove() (bool, *rbac.Permission)
	AssetUpload() (bool, *rbac.Permission)
	AssetGet() (bool, *rbac.Permission)
	LikePostGet() (bool, *rbac.Permission)
	LikePostAdd() (bool, *rbac.Permission)
	LikePostRemove() (bool, *rbac.Permission)
	LikeProfileGet() (bool, *rbac.Permission)
	CollectionCreate() (bool, *rbac.Permission)
	CollectionList() (bool, *rbac.Permission)
	CollectionGet() (bool, *rbac.Permission)
	CollectionUpdate() (bool, *rbac.Permission)
	CollectionDelete() (bool, *rbac.Permission)
	CollectionAddPost() (bool, *rbac.Permission)
	CollectionRemovePost() (bool, *rbac.Permission)
	CollectionAddNode() (bool, *rbac.Permission)
	CollectionRemoveNode() (bool, *rbac.Permission)
	NodeCreate() (bool, *rbac.Permission)
	NodeList() (bool, *rbac.Permission)
	NodeGet() (bool, *rbac.Permission)
	NodeUpdate() (bool, *rbac.Permission)
	NodeDelete() (bool, *rbac.Permission)
	NodeUpdateVisibility() (bool, *rbac.Permission)
	NodeAddAsset() (bool, *rbac.Permission)
	NodeRemoveAsset() (bool, *rbac.Permission)
	NodeAddNode() (bool, *rbac.Permission)
	NodeRemoveNode() (bool, *rbac.Permission)
	LinkCreate() (bool, *rbac.Permission)
	LinkList() (bool, *rbac.Permission)
	LinkGet() (bool, *rbac.Permission)
	DatagraphSearch() (bool, *rbac.Permission)
	EventList() (bool, *rbac.Permission)
	EventCreate() (bool, *rbac.Permission)
	EventGet() (bool, *rbac.Permission)
	EventUpdate() (bool, *rbac.Permission)
	EventDelete() (bool, *rbac.Permission)
	EventParticipantUpdate() (bool, *rbac.Permission)
	EventParticipantRemove() (bool, *rbac.Permission)
}

func GetOperationPermission(optable OperationPermissions, op string) (bool, *rbac.Permission) {
	switch op {
	case "GetVersion":
		return optable.GetVersion()
	case "GetSpec":
		return optable.GetSpec()
	case "GetInfo":
		return optable.GetInfo()
	case "IconGet":
		return optable.IconGet()
	case "IconUpload":
		return optable.IconUpload()
	case "AdminSettingsUpdate":
		return optable.AdminSettingsUpdate()
	case "AdminAccountBanCreate":
		return optable.AdminAccountBanCreate()
	case "AdminAccountBanRemove":
		return optable.AdminAccountBanRemove()
	case "RoleCreate":
		return optable.RoleCreate()
	case "RoleList":
		return optable.RoleList()
	case "RoleGet":
		return optable.RoleGet()
	case "RoleUpdate":
		return optable.RoleUpdate()
	case "RoleDelete":
		return optable.RoleDelete()
	case "AuthProviderList":
		return optable.AuthProviderList()
	case "AuthPasswordSignup":
		return optable.AuthPasswordSignup()
	case "AuthPasswordSignin":
		return optable.AuthPasswordSignin()
	case "AuthPasswordCreate":
		return optable.AuthPasswordCreate()
	case "AuthPasswordUpdate":
		return optable.AuthPasswordUpdate()
	case "AuthEmailPasswordSignup":
		return optable.AuthEmailPasswordSignup()
	case "AuthEmailPasswordSignin":
		return optable.AuthEmailPasswordSignin()
	case "AuthEmailSignup":
		return optable.AuthEmailSignup()
	case "AuthEmailSignin":
		return optable.AuthEmailSignin()
	case "AuthEmailVerify":
		return optable.AuthEmailVerify()
	case "OAuthProviderCallback":
		return optable.OAuthProviderCallback()
	case "WebAuthnRequestCredential":
		return optable.WebAuthnRequestCredential()
	case "WebAuthnMakeCredential":
		return optable.WebAuthnMakeCredential()
	case "WebAuthnGetAssertion":
		return optable.WebAuthnGetAssertion()
	case "WebAuthnMakeAssertion":
		return optable.WebAuthnMakeAssertion()
	case "PhoneRequestCode":
		return optable.PhoneRequestCode()
	case "PhoneSubmitCode":
		return optable.PhoneSubmitCode()
	case "AuthProviderLogout":
		return optable.AuthProviderLogout()
	case "AccountGet":
		return optable.AccountGet()
	case "AccountUpdate":
		return optable.AccountUpdate()
	case "AccountAuthProviderList":
		return optable.AccountAuthProviderList()
	case "AccountAuthMethodDelete":
		return optable.AccountAuthMethodDelete()
	case "AccountSetAvatar":
		return optable.AccountSetAvatar()
	case "AccountGetAvatar":
		return optable.AccountGetAvatar()
	case "AccountAddRole":
		return optable.AccountAddRole()
	case "AccountRemoveRole":
		return optable.AccountRemoveRole()
	case "AccountRoleSetBadge":
		return optable.AccountRoleSetBadge()
	case "AccountRoleRemoveBadge":
		return optable.AccountRoleRemoveBadge()
	case "InvitationList":
		return optable.InvitationList()
	case "InvitationCreate":
		return optable.InvitationCreate()
	case "InvitationGet":
		return optable.InvitationGet()
	case "InvitationDelete":
		return optable.InvitationDelete()
	case "NotificationList":
		return optable.NotificationList()
	case "NotificationUpdate":
		return optable.NotificationUpdate()
	case "ProfileList":
		return optable.ProfileList()
	case "ProfileGet":
		return optable.ProfileGet()
	case "ProfileFollowersGet":
		return optable.ProfileFollowersGet()
	case "ProfileFollowersAdd":
		return optable.ProfileFollowersAdd()
	case "ProfileFollowersRemove":
		return optable.ProfileFollowersRemove()
	case "ProfileFollowingGet":
		return optable.ProfileFollowingGet()
	case "CategoryCreate":
		return optable.CategoryCreate()
	case "CategoryList":
		return optable.CategoryList()
	case "CategoryUpdateOrder":
		return optable.CategoryUpdateOrder()
	case "CategoryUpdate":
		return optable.CategoryUpdate()
	case "TagList":
		return optable.TagList()
	case "ThreadCreate":
		return optable.ThreadCreate()
	case "ThreadList":
		return optable.ThreadList()
	case "ThreadGet":
		return optable.ThreadGet()
	case "ThreadUpdate":
		return optable.ThreadUpdate()
	case "ThreadDelete":
		return optable.ThreadDelete()
	case "ReplyCreate":
		return optable.ReplyCreate()
	case "PostUpdate":
		return optable.PostUpdate()
	case "PostDelete":
		return optable.PostDelete()
	case "PostSearch":
		return optable.PostSearch()
	case "PostReactAdd":
		return optable.PostReactAdd()
	case "PostReactRemove":
		return optable.PostReactRemove()
	case "AssetUpload":
		return optable.AssetUpload()
	case "AssetGet":
		return optable.AssetGet()
	case "LikePostGet":
		return optable.LikePostGet()
	case "LikePostAdd":
		return optable.LikePostAdd()
	case "LikePostRemove":
		return optable.LikePostRemove()
	case "LikeProfileGet":
		return optable.LikeProfileGet()
	case "CollectionCreate":
		return optable.CollectionCreate()
	case "CollectionList":
		return optable.CollectionList()
	case "CollectionGet":
		return optable.CollectionGet()
	case "CollectionUpdate":
		return optable.CollectionUpdate()
	case "CollectionDelete":
		return optable.CollectionDelete()
	case "CollectionAddPost":
		return optable.CollectionAddPost()
	case "CollectionRemovePost":
		return optable.CollectionRemovePost()
	case "CollectionAddNode":
		return optable.CollectionAddNode()
	case "CollectionRemoveNode":
		return optable.CollectionRemoveNode()
	case "NodeCreate":
		return optable.NodeCreate()
	case "NodeList":
		return optable.NodeList()
	case "NodeGet":
		return optable.NodeGet()
	case "NodeUpdate":
		return optable.NodeUpdate()
	case "NodeDelete":
		return optable.NodeDelete()
	case "NodeUpdateVisibility":
		return optable.NodeUpdateVisibility()
	case "NodeAddAsset":
		return optable.NodeAddAsset()
	case "NodeRemoveAsset":
		return optable.NodeRemoveAsset()
	case "NodeAddNode":
		return optable.NodeAddNode()
	case "NodeRemoveNode":
		return optable.NodeRemoveNode()
	case "LinkCreate":
		return optable.LinkCreate()
	case "LinkList":
		return optable.LinkList()
	case "LinkGet":
		return optable.LinkGet()
	case "DatagraphSearch":
		return optable.DatagraphSearch()
	case "EventList":
		return optable.EventList()
	case "EventCreate":
		return optable.EventCreate()
	case "EventGet":
		return optable.EventGet()
	case "EventUpdate":
		return optable.EventUpdate()
	case "EventDelete":
		return optable.EventDelete()
	case "EventParticipantUpdate":
		return optable.EventParticipantUpdate()
	case "EventParticipantRemove":
		return optable.EventParticipantRemove()
	default:
		panic("unknown operation, must re-run rbacgen")
	}
}
