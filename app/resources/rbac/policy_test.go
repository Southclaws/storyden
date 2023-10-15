package rbac_test

import (
	"testing"

	"github.com/el-mike/restrict"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestPolicy(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		am rbac.AccessManager,
	) {
		a := assert.New(t)
		err := am.Authorize(&restrict.AccessRequest{
			Subject:  &seed.Account_001_Odin,
			Resource: &seed.Post_01_Welcome,
			Actions: []string{
				rbac.ActionDelete,
			},
			Context: map[string]interface{}{
				"": nil,
			},
			SkipConditions: false,
		})

		a.NoError(err)
	}))
}
