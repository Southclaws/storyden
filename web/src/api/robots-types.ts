import { UIDataTypes, UIMessage, UIMessagePart, isToolUIPart } from "ai";

import { RobotSessionMessageList } from "./openapi-schema/robotSessionMessageList";
import { StorydenTools } from "./robots";

export type RobotRenderCardData = {
  ref: string;
  kind: string;
  id: string;
};

export type StorydenUIDataTypes = {
  session_name: string;
  render_card: RobotRenderCardData;
};

export type StorydenUIMessage = UIMessage<
  unknown,
  StorydenUIDataTypes,
  StorydenTools
>;

export function toStorydenUIMessages(
  messages: RobotSessionMessageList,
): StorydenUIMessage[] {
  return messages as unknown as StorydenUIMessage[];
}

type Part = StorydenUIMessage["parts"][number];

type ToolType = Extract<Part["type"], `tool-${string}`>; // "tool-search" | "tool-robot_switch" | ...
export type ToolName = ToolType extends `tool-${infer N}` ? N : never; // "search" | "robot_switch" | ...

export function isToolType(t: Part["type"]): t is ToolType {
  return t.startsWith("tool-");
}

export function getToolName(
  part: UIMessagePart<UIDataTypes, StorydenTools>,
): string {
  if (!isToolUIPart(part)) {
    return "Unknown";
  }

  const rawName = getRawToolName(part);

  if (rawName === "adk_request_confirmation") {
    return "Confirmation";
  }

  return rawName
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

export function getRawToolName(part: { type: string }): string {
  if (!part.type.startsWith("tool-")) {
    return "";
  }

  return String(part.type).replace(/^tool-/, "");
}
