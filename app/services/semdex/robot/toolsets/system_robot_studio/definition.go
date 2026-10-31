package system_robot_studio

const (
	ID          = "system.robot_studio"
	Name        = "Robot studio"
	Description = "Discover capabilities and create, inspect, update, or remove reusable Robots and Toolsets."
	Instruction = `Create and maintain focused Robots and reusable Toolsets.

When the user names an existing Robot, find it with robot_list. Read its current configuration with robot_get before updating it, preserving unrelated instructions and capabilities. Delete a Robot only when the user clearly requests its removal.

A Robot playbook should define its job, boundaries, normal workflow, when to ask versus act, which tools establish facts or perform actions, and how to report completion or blocked work. Use concrete decision rules rather than generic advice. Keep instructions self-contained except for guidance supplied by assigned Toolsets.

Use direct tools for a narrow capability that needs no shared guidance. Prefer an existing Toolset for related capabilities and their instructions. Create a Toolset when that grouping or guidance should be reused. Do not assign a direct tool already supplied by an assigned Toolset.

Refer to Robots by their names in ordinary responses. Include raw IDs only when needed for the task or requested.`
)
