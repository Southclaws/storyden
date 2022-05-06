package rbac

import (
	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/post"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

func NewPolicy() *restrict.PolicyDefinition {
	return &restrict.PolicyDefinition{
		Roles: restrict.Roles{
			user.Role: {
				Grants: restrict.GrantsMap{
					post.Role: {
						&restrict.Permission{Action: "read"},
						&restrict.Permission{Action: "create"},
						&restrict.Permission{
							Action: "update",
							Conditions: restrict.Conditions{
								&restrict.EqualCondition{
									ID: "isOwner",
									Left: &restrict.ValueDescriptor{
										Source: restrict.ResourceField,
										Field:  "Author.ID",
									},
									Right: &restrict.ValueDescriptor{
										Source: restrict.SubjectField,
										Field:  "ID",
									},
								},
							},
						},
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
