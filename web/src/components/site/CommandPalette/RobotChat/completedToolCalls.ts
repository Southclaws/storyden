import { TOOL_NAMES, ToolName } from "@/api/robots";
import { StorydenUIMessage } from "@/api/robots-types";

export function isKnownToolName(name: string): name is ToolName {
  return (TOOL_NAMES as readonly string[]).includes(name);
}

export function findCompletedStorydenToolCalls(
  messages: readonly StorydenUIMessage[],
) {
  const completed = new Map<
    string,
    { toolCallId: string; toolName: ToolName }
  >();

  for (const message of messages) {
    for (const part of message.parts ?? []) {
      if (
        part.type === "dynamic-tool" ||
        !part.type.startsWith("tool-") ||
        !("state" in part) ||
        part.state !== "output-available" ||
        !("toolCallId" in part) ||
        !part.toolCallId
      ) {
        continue;
      }

      const toolName = part.type.slice("tool-".length);
      if (!isKnownToolName(toolName)) {
        continue;
      }

      completed.set(part.toolCallId, {
        toolCallId: part.toolCallId,
        toolName,
      });
    }
  }

  return [...completed.values()];
}
