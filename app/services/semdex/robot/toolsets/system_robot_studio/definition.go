package system_robot_studio

const (
	ID          = "system.robot_studio"
	Name        = "Robot studio"
	Description = "Discover capabilities and create, inspect, update, or remove reusable Robots and Toolsets."
	Instruction = `Design Robots as focused delegation targets. Inspect before updating and preserve unrelated behaviour. Use tool_search and tool_get for one or two narrow capabilities; use toolset_search and toolset_get when a coherent group of tools or specialist instructions should be shared. Never assign a tool directly when an assigned Toolset already contains it. Keep playbooks concrete: job, boundaries, workflow, tool authority, clarification rules, success, and failure behaviour.`
)
