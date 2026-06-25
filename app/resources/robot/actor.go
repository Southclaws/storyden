package robot

import (
	"fmt"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
)

type BuiltinRobotID string

func (id BuiltinRobotID) String() string {
	return string(id)
}

type Actor struct {
	DatabaseRobotID opt.Optional[xid.ID]
	BuiltinRobotID  opt.Optional[BuiltinRobotID]
}

func NewDatabaseActor(id xid.ID) Actor {
	return Actor{DatabaseRobotID: opt.New(id)}
}

func NewBuiltinActor(id string) Actor {
	return Actor{BuiltinRobotID: opt.New(BuiltinRobotID(id))}
}

func NewActorFromRobotRef(ref string) (Actor, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return Actor{}, fmt.Errorf("robot actor reference is required")
	}

	if id, err := xid.FromString(ref); err == nil {
		return NewDatabaseActor(id), nil
	}

	return NewBuiltinActor(ref), nil
}

func (a Actor) Validate() error {
	_, hasDatabase := a.DatabaseRobotID.Get()
	builtinID, hasBuiltin := a.BuiltinRobotID.Get()

	switch {
	case hasDatabase && hasBuiltin:
		return fmt.Errorf("robot actor cannot reference both database and built-in robots")
	case hasBuiltin && strings.TrimSpace(builtinID.String()) == "":
		return fmt.Errorf("built-in robot actor ID is required")
	default:
		return nil
	}
}

func (a Actor) Ref() opt.Optional[string] {
	if id, ok := a.DatabaseRobotID.Get(); ok {
		return opt.New(id.String())
	}
	if id, ok := a.BuiltinRobotID.Get(); ok {
		return opt.New(id.String())
	}
	return opt.NewEmpty[string]()
}
