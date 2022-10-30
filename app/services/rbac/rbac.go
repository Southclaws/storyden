package rbac

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/thread"
)

func NewPolicy() *restrict.PolicyDefinition {
	return &restrict.PolicyDefinition{
		Roles: restrict.Roles{
			account.Name: {
				Grants: restrict.GrantsMap{
					// account.Name: {
					// 	&restrict.Permission{Action: "create"},
					// 	&restrict.Permission{Action: "read"},
					// 	&restrict.Permission{Action: "update"},
					// 	&restrict.Permission{Action: "delete"},
					// },
					thread.Name: {
						&restrict.Permission{Action: "create", Conditions: []restrict.Condition{}},
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
	pm, err := restrict.NewPolicyManager(storage, true)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new policy manager"))
	}

	return pm, nil
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
