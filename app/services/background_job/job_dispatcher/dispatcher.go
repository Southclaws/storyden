package job_dispatcher

import (
	"context"

	"github.com/rs/xid"
)

type Dispatcher interface {
	Dispatch(context.Context, xid.ID)
}
