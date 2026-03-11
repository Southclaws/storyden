package robot

import _ "embed"

const (
	defaultAgentName        = "storyden"
	defaultAgentDescription = "Storyden's default agent that helps users get started and manage their community knowledge base."
)

//go:embed default.md
var defaultInstruction string

//go:embed global.md
var globalInstruction string
