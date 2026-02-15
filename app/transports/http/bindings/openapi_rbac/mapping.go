package openapi_rbac

import "github.com/Southclaws/storyden/app/resources/rbac"

var _ OperationPermissions = &Mapping{}

type Mapping struct{}

func (m *Mapping) GetVersion() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) GetSpec() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) GetDocs() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) GetInfo() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) IconGet() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) IconUpload() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageSettings
}

func (m *Mapping) BannerGet() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) BannerUpload() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageSettings
}

func (m *Mapping) SendBeacon() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AdminSettingsGet() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageSettings
}

func (m *Mapping) AdminSettingsUpdate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageSettings
}

func (m *Mapping) AuditEventList() (bool, *rbac.Permission) {
	return true, &rbac.PermissionAdministrator
}

func (m *Mapping) AuditEventGet() (bool, *rbac.Permission) {
	return true, &rbac.PermissionAdministrator
}

func (m *Mapping) ModerationActionCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionAdministrator
}

func (m *Mapping) AdminAccountBanCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageSuspensions
}

func (m *Mapping) AdminAccessKeyList() (bool, *rbac.Permission) {
	return true, &rbac.PermissionAdministrator
}

func (m *Mapping) AdminAccessKeyDelete() (bool, *rbac.Permission) {
	return true, &rbac.PermissionAdministrator
}

func (m *Mapping) AdminAccountBanRemove() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageSuspensions
}

func (m *Mapping) RoleCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageRoles
}

func (m *Mapping) RoleList() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) RoleGet() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) RoleUpdate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageRoles
}

func (m *Mapping) RoleUpdateOrder() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageRoles
}

func (m *Mapping) RoleDelete() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageRoles
}

func (m *Mapping) AuthProviderList() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AuthPasswordSignup() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AuthPasswordSignin() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AuthPasswordCreate() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AuthPasswordUpdate() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AuthPasswordReset() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AuthEmailPasswordSignup() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AuthEmailPasswordSignin() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AuthPasswordResetRequestEmail() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AuthEmailSignup() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AuthEmailSignin() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AuthEmailVerify() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) OAuthProviderCallback() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) WebAuthnRequestCredential() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) WebAuthnMakeCredential() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) WebAuthnGetAssertion() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) WebAuthnMakeAssertion() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) PhoneRequestCode() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) PhoneSubmitCode() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) AccessKeyList() (bool, *rbac.Permission) {
	return true, &rbac.PermissionUsePersonalAccessKeys
}

func (m *Mapping) AccessKeyCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionUsePersonalAccessKeys
}

func (m *Mapping) AccessKeyDelete() (bool, *rbac.Permission) {
	return true, &rbac.PermissionUsePersonalAccessKeys
}

func (m *Mapping) AuthProviderLogout() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountGet() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountView() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) AccountUpdate() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountAuthProviderList() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountAuthMethodDelete() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountEmailAdd() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountEmailRemove() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountSetAvatar() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountGetAvatar() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountAddRole() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageRoles
}

func (m *Mapping) AccountRemoveRole() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageRoles
}

func (m *Mapping) AccountRoleSetBadge() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) AccountRoleRemoveBadge() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) InvitationList() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) InvitationCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionCreateInvitation
}

func (m *Mapping) InvitationGet() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) InvitationDelete() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) NotificationList() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) NotificationUpdate() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) NotificationUpdateMany() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) ReportCreate() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) ReportList() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) ReportUpdate() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) ProfileList() (bool, *rbac.Permission) {
	return false, &rbac.PermissionListProfiles
}

func (m *Mapping) ProfileGet() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadProfile
}

func (m *Mapping) ProfileFollowersGet() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) ProfileFollowersAdd() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) ProfileFollowersRemove() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) ProfileFollowingGet() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) CategoryCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageCategories
}

func (m *Mapping) CategoryList() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) CategoryGet() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) CategoryUpdatePosition() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageCategories
}

func (m *Mapping) CategoryUpdate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageCategories
}

func (m *Mapping) CategoryDelete() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageCategories
}

func (m *Mapping) TagList() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) TagGet() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) ThreadCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionCreatePost
}

func (m *Mapping) ThreadList() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadPublishedThreads
}

func (m *Mapping) ThreadGet() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadPublishedThreads
}

func (m *Mapping) ThreadUpdate() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) ThreadDelete() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) ReplyCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionCreatePost
}

func (m *Mapping) PostUpdate() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) PostDelete() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) PostLocationGet() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadPublishedThreads
}

func (m *Mapping) PostSearch() (bool, *rbac.Permission) {
	return true, &rbac.PermissionReadPublishedThreads
}

func (m *Mapping) PostReactAdd() (bool, *rbac.Permission) {
	return true, &rbac.PermissionCreateReaction
}

func (m *Mapping) PostReactRemove() (bool, *rbac.Permission) {
	return true, &rbac.PermissionCreateReaction
}

func (m *Mapping) AssetUpload() (bool, *rbac.Permission) {
	return true, &rbac.PermissionUploadAsset
}

func (m *Mapping) AssetGet() (bool, *rbac.Permission) {
	return false, nil // Public
}

func (m *Mapping) LikePostGet() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) LikePostAdd() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) LikePostRemove() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) LikeProfileGet() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) CollectionCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionCreateCollection
}

func (m *Mapping) CollectionList() (bool, *rbac.Permission) {
	return false, &rbac.PermissionListCollections
}

func (m *Mapping) CollectionGet() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadCollection
}

func (m *Mapping) CollectionUpdate() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) CollectionDelete() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) CollectionAddPost() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) CollectionRemovePost() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) CollectionAddNode() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) CollectionRemoveNode() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeCreate() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeList() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadPublishedLibrary
}

func (m *Mapping) NodeGet() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadPublishedLibrary
}

func (m *Mapping) NodeListChildren() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadPublishedLibrary
}

func (m *Mapping) NodeUpdate() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeGenerateContent() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeGenerateTags() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeGenerateTitle() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeDelete() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeUpdateChildrenPropertySchema() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeUpdatePropertySchema() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeUpdateProperties() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeUpdateVisibility() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeAddAsset() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeRemoveAsset() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeAddNode() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeRemoveNode() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) NodeUpdatePosition() (bool, *rbac.Permission) {
	return true, nil // See NOTE.
}

func (m *Mapping) LinkCreate() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) LinkList() (bool, *rbac.Permission) {
	return false, &rbac.PermissionReadPublishedLibrary
}

func (m *Mapping) LinkGet() (bool, *rbac.Permission) {
	return true, &rbac.PermissionReadPublishedLibrary
}

func (m *Mapping) DatagraphSearch() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) DatagraphMatches() (bool, *rbac.Permission) {
	return true, nil
}

func (m *Mapping) DatagraphAsk() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) EventList() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) EventCreate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageEvents
}

func (m *Mapping) EventGet() (bool, *rbac.Permission) {
	return false, nil
}

func (m *Mapping) EventUpdate() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageEvents
}

func (m *Mapping) EventDelete() (bool, *rbac.Permission) {
	return true, &rbac.PermissionManageEvents
}

func (m *Mapping) EventParticipantUpdate() (bool, *rbac.Permission) {
	// Requires PermissionManageEvents unless updating self
	return true, nil
}

func (m *Mapping) EventParticipantRemove() (bool, *rbac.Permission) {
	// Requires PermissionManageEvents unless deleting self
	return true, nil
}
