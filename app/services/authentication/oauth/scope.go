package oauth

import (
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"

	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

var standardScopes = map[string]struct{}{
	"openid":         {},
	"profile":        {},
	"email":          {},
	"offline_access": {},
}

func supportedScopes() []string {
	return append([]string{"openid", "profile", "email", "offline_access"}, rbacAllPermissionNames()...)
}

func validateScopeNames(scope string) error {
	for _, sc := range splitScope(scope) {
		if _, err := permissionFromScope(sc); err != nil {
			return err
		}
	}
	return nil
}

func authorizeScopeNames(scope string, allowedScopes []string) error {
	clientHasAdministrator := contains(allowedScopes, rbac.PermissionAdministrator.String())
	for _, sc := range splitScope(scope) {
		if !contains(allowedScopes, sc) && !clientHasAdministrator {
			return fault.New("scope is not allowed for client")
		}
	}
	return nil
}

func validatePermissionScopes(scope string, accountPermissions rbac.Permissions) error {
	permissionScopes, err := permissionsFromScopes(splitScope(scope))
	if err != nil {
		return err
	}
	if accountPermissions.HasAll(rbac.PermissionAdministrator) {
		return nil
	}
	if !accountPermissions.HasAll(permissionScopes.List()...) {
		return rbac.ErrPermissions
	}

	return nil
}

func validatePermissionOnlyScopes(scopes []string) error {
	for _, scope := range scopes {
		if scope == "" {
			continue
		}
		if _, ok := standardScopes[scope]; ok {
			return fault.New("standard oauth scopes are not valid for oauth keys")
		}
		if _, err := rbac.NewPermission(scope); err != nil {
			return err
		}
	}

	return nil
}

func grantScope(requestedScope string, client *oauthresource.Client, accountPermissions rbac.Permissions) (string, error) {
	requestedScopes := splitScope(requestedScope)
	if err := validateScopeNames(requestedScope); err != nil {
		return "", err
	}

	standard := standardScopesFrom(requestedScopes)
	requestedPermissions, err := permissionsFromScopes(requestedScopes)
	if err != nil {
		return "", err
	}

	var grantedPermissions rbac.Permissions
	if shouldInheritUserPermissions(client) && len(requestedPermissions.List()) == 0 {
		grantedPermissions = accountPermissions
	} else {
		allowedPermissions, err := permissionsFromScopes(client.AllowedScopes)
		if err != nil {
			return "", err
		}
		grantedPermissions = intersectPermissions(requestedPermissions, allowedPermissions, accountPermissions)
	}

	return joinScopes(standard, grantedPermissions.List()), nil
}

func grantClientCredentialsScope(requestedScope string, client *oauthresource.Client, accountPermissions rbac.Permissions) (string, error) {
	if requestedScope != "" {
		if err := validateScopeNames(requestedScope); err != nil {
			return "", err
		}
		if err := authorizeScopeNames(requestedScope, client.AllowedScopes); err != nil {
			return "", err
		}

		requestedScopes := splitScope(requestedScope)
		requestedPermissions, err := permissionsFromScopes(requestedScopes)
		if err != nil {
			return "", err
		}

		allowedPermissions, err := permissionsFromScopes(client.AllowedScopes)
		if err != nil {
			return "", err
		}

		grantedPermissions := intersectPermissions(requestedPermissions, allowedPermissions, accountPermissions)
		return joinScopes(nil, grantedPermissions.List()), nil
	}

	allowedPermissions, err := permissionsFromScopes(client.AllowedScopes)
	if err != nil {
		return "", err
	}

	grantedPermissions := intersectPermissions(allowedPermissions, accountPermissions)
	return joinScopes(nil, grantedPermissions.List()), nil
}

func refreshScope(existingScope string, client *oauthresource.Client, accountPermissions rbac.Permissions) (string, error) {
	existingScopes := splitScope(existingScope)
	standard := standardScopesFrom(existingScopes)

	if shouldInheritUserPermissions(client) {
		existingPermissions, err := permissionsFromScopes(existingScopes)
		if err != nil {
			return "", err
		}
		grantedPermissions := intersectPermissions(existingPermissions, accountPermissions)
		return joinScopes(standard, grantedPermissions.List()), nil
	}

	existingPermissions, err := permissionsFromScopes(existingScopes)
	if err != nil {
		return "", err
	}
	allowedPermissions, err := permissionsFromScopes(client.AllowedScopes)
	if err != nil {
		return "", err
	}

	grantedPermissions := intersectPermissions(existingPermissions, allowedPermissions, accountPermissions)
	return joinScopes(standard, grantedPermissions.List()), nil
}

func shouldInheritUserPermissions(client *oauthresource.Client) bool {
	return client.ScopePolicy == oauthresource.ScopePolicyInheritUserPermissions
}

func standardScopesFrom(scopes []string) []string {
	out := []string{}
	for _, scope := range scopes {
		if _, ok := standardScopes[scope]; ok {
			out = append(out, scope)
		}
	}

	return out
}

func intersectPermissions(permissionSets ...rbac.Permissions) rbac.Permissions {
	if len(permissionSets) == 0 {
		return rbac.NewList()
	}

	first := permissionSets[0].List()
	out := []rbac.Permission{}
	for _, permission := range first {
		inAll := true
		for _, set := range permissionSets[1:] {
			if set.HasAll(rbac.PermissionAdministrator) {
				continue
			}
			if !set.HasAll(permission) {
				inAll = false
				break
			}
		}
		if inAll {
			out = append(out, permission)
		}
	}

	return rbac.NewList(out...)
}

func joinScopes(standard []string, permissions rbac.PermissionList) string {
	scopes := append([]string{}, standard...)
	for _, permission := range permissions {
		scopes = append(scopes, permission.String())
	}

	return strings.Join(scopes, " ")
}

func permissionsFromScopes(scopes []string) (rbac.Permissions, error) {
	permissions := []rbac.Permission{}
	for _, scope := range scopes {
		permission, err := permissionFromScope(scope)
		if err != nil {
			return rbac.Permissions{}, err
		}
		if permission == nil {
			continue
		}

		permissions = append(permissions, *permission)
	}

	return rbac.NewList(permissions...), nil
}

func permissionFromScope(scope string) (*rbac.Permission, error) {
	if scope == "" {
		return nil, nil
	}
	if _, ok := standardScopes[scope]; ok {
		return nil, nil
	}

	permission, err := rbac.NewPermission(scope)
	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func rbacAllPermissionNames() []string {
	return dt.Map(rbac.AllPermissions, func(permission rbac.Permission) string {
		return permission.String()
	})
}
