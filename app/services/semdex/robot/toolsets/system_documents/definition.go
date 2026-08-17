package system_documents

import (
	"google.golang.org/adk/v2/agent"

	"github.com/Southclaws/storyden/app/services/semdex/robot/documents"
)

const (
	ID          = "system.documents"
	Name        = "Document exploration"
	Description = "Navigate, inspect, search, and close bounded document snapshots opened by other Storyden tools."
	Instruction = `Treat open documents as immutable snapshots used for focused exploration. Start with the bounded projection returned by an open tool. Prefer the smallest relevant structural location: if a search preview already answers the question, use it directly; otherwise use document_get on the most specific matching leaf, and expand a broad heading or table only when narrower content is insufficient. Follow previous_page or next_page when a structural location spans multiple pages. Each document_get call updates the current document cursor, which is navigation state rather than evidence. Use document_search with node scopes when keywords can narrow the task, document_list when the intended snapshot is unclear, and document_close when a snapshot is no longer useful. Do not revisit unrelated locations after obtaining enough evidence. Document IDs and node IDs are conversation-local navigation handles. Base claims on content returned by tool calls, and do not imply that omitted, unvisited, or truncated sections were inspected.`
)

func InstructionProvider(ctx agent.ReadonlyContext) (string, error) {
	register, err := documents.RegisterInstruction(ctx.ReadonlyState())
	if err != nil {
		return "", err
	}
	if register == "" {
		return Instruction, nil
	}
	return Instruction + "\n\n" + register, nil
}
