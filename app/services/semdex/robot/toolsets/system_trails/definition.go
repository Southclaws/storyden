package system_trails

const (
	ID          = "system.trails"
	Name        = "Trail management"
	Description = "Create, inspect, operate, and troubleshoot durable scheduled or event-driven Trails."
	Instruction = `Use Trails as the source of truth for durable scheduled or event-driven work. Inspect a Trail before replacing its definition when existing configuration matters. Preview schedules before creating or changing recurring work; use a count of 1 for a one-shot wake-up. Use trail_run_list to choose a run and trail_run_get to inspect its immutable trigger context and independent action results. Create a manual run only when the user wants immediate execution; it does not move the schedule. Cancel only the selected queued or running action, because sibling actions continue independently. Never infer that one failed action stopped the others, and never retry a failed run without explicit user intent.`
)
