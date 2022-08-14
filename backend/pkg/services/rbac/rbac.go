package rbac

import (
	"errors"

	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/resources/account"
	"github.com/Southclaws/storyden/backend/pkg/resources/post"
)

type admin struct{}

func (c admin) Type() string { return "ADMIN" }

func (c admin) Check(request *restrict.AccessRequest) error {
	sub := request.Subject.(*account.Account)
	res := request.Resource.(*account.Account)

	if sub.ID == res.ID {
		return nil
	}

	if sub.Admin {
		return nil
	}

	return restrict.NewConditionNotSatisfiedError(c, request, errors.New("not authorised"))
}

func NewPolicy() *restrict.PolicyDefinition {
	return &restrict.PolicyDefinition{
		Roles: restrict.Roles{
			account.Name: {
				Grants: restrict.GrantsMap{
					account.Name: {
						&restrict.Permission{Action: "create", Conditions: restrict.Conditions{admin{}}},
						&restrict.Permission{Action: "read", Conditions: restrict.Conditions{admin{}}},
						&restrict.Permission{Action: "update", Conditions: restrict.Conditions{admin{}}},
						&restrict.Permission{Action: "delete", Conditions: restrict.Conditions{admin{}}},
					},
					post.Role: {
						&restrict.Permission{Action: "create"},
						&restrict.Permission{Action: "read"},
						&restrict.Permission{Action: "update"}, // TODO: Ownership stuff
						&restrict.Permission{Action: "delete"},
					},
				},
			},
		},
	}
}

func NewAdapter(policy *restrict.PolicyDefinition) restrict.StorageAdapter {
	return adapters.NewInMemoryAdapter(policy)
}

func NewPolicyManager(storage restrict.StorageAdapter) (*restrict.PolicyManager, error) {
	return restrict.NewPolicyManager(storage, true)
}

func NewAccessManager(policyMananger *restrict.PolicyManager) *restrict.AccessManager {
	return restrict.NewAccessManager(policyMananger)
}

func Build() fx.Option {
	return fx.Provide(
		NewPolicy,
		NewAdapter,
		NewPolicyManager,
		NewAccessManager,
	)
}
