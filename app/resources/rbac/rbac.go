package rbac

import (
	"errors"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
	"go.uber.org/fx"
)

var lock sync.RWMutex

type AccessManager interface {
	Authorize(request *restrict.AccessRequest) error
}

type withFault struct {
	*restrict.AccessManager
}

func (m *withFault) Authorize(request *restrict.AccessRequest) error {
	if err := m.AccessManager.Authorize(request); err != nil {
		ae := &restrict.AccessDeniedError{}
		if errors.As(err, &ae) {
			return fault.Wrap(err,
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc("access request denied", ae.Error()),
			)
		}

		return fault.Wrap(err, ftag.With(ftag.PermissionDenied))
	}

	return nil
}

func NewAdapter() restrict.StorageAdapter {
	return adapters.NewInMemoryAdapter(defaultPolicy)
}

func NewPolicyManager(storage restrict.StorageAdapter) (*restrict.PolicyManager, error) {
	lock.Lock()
	defer lock.Unlock()

	pm, err := restrict.NewPolicyManager(storage, true)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new policy manager"))
	}

	return pm, nil
}

func NewAccessManager(policyMananger *restrict.PolicyManager) AccessManager {
	return &withFault{restrict.NewAccessManager(policyMananger)}
}

func Build() fx.Option {
	return fx.Provide(
		NewAdapter,
		NewPolicyManager,
		NewAccessManager,
	)
}
