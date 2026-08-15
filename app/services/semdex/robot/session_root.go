package robot

import (
	"strings"

	"github.com/Southclaws/opt"
)

const sessionRootRobotRefStateKey = "root_robot_ref"

func SessionRootRobotRef(state map[string]any) opt.Optional[string] {
	if state == nil {
		return opt.NewEmpty[string]()
	}

	ref, ok := state[sessionRootRobotRefStateKey].(string)
	if !ok {
		return opt.NewEmpty[string]()
	}

	ref = strings.TrimSpace(ref)
	return opt.NewIf(ref, func(value string) bool { return value != "" })
}

func sessionInitialState(rootRobotRef string) map[string]any {
	return map[string]any{
		sessionRootRobotRefStateKey: strings.TrimSpace(rootRobotRef),
	}
}
