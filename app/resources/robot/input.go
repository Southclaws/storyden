package robot

import (
	"encoding/json"
	"time"

	"github.com/rs/xid"
)

type InputID xid.ID

func (id InputID) String() string { return xid.ID(id).String() }

type SessionInput struct {
	ID         InputID
	SessionID  SessionID
	AccountID  xid.ID
	Sequence   uint64
	SourceKind string
	BatchKey   string
	InputData  json.RawMessage
	CreatedAt  time.Time
}
