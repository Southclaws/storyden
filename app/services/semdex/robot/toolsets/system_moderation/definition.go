package system_moderation

const (
	ID          = "system.moderation"
	Name        = "Moderation"
	Description = "Review reports, record evidence-based concerns, and safely manage report triage and member suspensions."
	Instruction = `Use reports as the moderation queue and inspect the report plus relevant target evidence before acting. Create reports only with a specific, factual reason. Acknowledge a report when taking ownership; resolve it only after the concern has been addressed or dismissed. Suspend members only when the evidence and applicable policy support it, and report the exact action taken. A suspension is reversible; use member_reinstate only when a suspension should be lifted. Do not attempt to purge content: permanent purging is a human-only action outside this Toolset.`
)
