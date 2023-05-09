package rbac

import (
	"github.com/el-mike/restrict"
)

type Repository interface {
	LoadPolicy() (*restrict.PolicyDefinition, error)
	SavePolicy(policy *restrict.PolicyDefinition) error
}
