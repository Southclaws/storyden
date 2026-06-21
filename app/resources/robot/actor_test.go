package robot

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestNewActorFromRobotRef(t *testing.T) {
	databaseID := xid.New()

	databaseActor, err := NewActorFromRobotRef(databaseID.String())
	require.NoError(t, err)
	require.Equal(t, databaseID, databaseActor.DatabaseRobotID.OrZero())
	require.False(t, databaseActor.BuiltinRobotID.Ok())

	builtinActor, err := NewActorFromRobotRef("plugin_builder")
	require.NoError(t, err)
	require.False(t, builtinActor.DatabaseRobotID.Ok())
	require.Equal(t, BuiltinRobotID("plugin_builder"), builtinActor.BuiltinRobotID.OrZero())
}

func TestActorValidateRejectsAmbiguousActor(t *testing.T) {
	actor := Actor{
		DatabaseRobotID: NewDatabaseActor(xid.New()).DatabaseRobotID,
		BuiltinRobotID:  NewBuiltinActor("plugin_builder").BuiltinRobotID,
	}

	require.Error(t, actor.Validate())
}
