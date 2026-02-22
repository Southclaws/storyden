package role_test

import "github.com/Southclaws/storyden/app/transports/http/openapi"

func findRole(roles []openapi.AccountRole, id string) *openapi.AccountRole {
	for _, role := range roles {
		if role.Id == id {
			return &role
		}
	}
	return nil
}

func findRoleRef(roles []openapi.AccountRoleRef, id openapi.Identifier) *openapi.AccountRoleRef {
	for _, role := range roles {
		if role.Id == id {
			return &role
		}
	}
	return nil
}
