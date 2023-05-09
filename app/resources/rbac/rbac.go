package rbac

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
	"go.uber.org/fx"
)

func NewAdapter() restrict.StorageAdapter {
	return adapters.NewInMemoryAdapter(defaultPolicy)
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
		NewAdapter,
		NewPolicyManager,
		NewAccessManager,
	)
}
