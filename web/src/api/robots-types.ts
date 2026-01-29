import {
  UIDataTypes,
  UIMessage,
  UIMessagePart,
  isToolOrDynamicToolUIPart,
} from "ai";

import { StorydenTools } from "./robots";

export type StorydenUIDataTypes = {
  session_name: string;
};

export type StorydenUIMessage = UIMessage<
  unknown,
  StorydenUIDataTypes,
  StorydenTools
>;

type Part = StorydenUIMessage["parts"][number];

type ToolType = Extract<Part["type"], `tool-${string}`>; // "tool-search" | "tool-robot_switch" | ...
export type ToolName = ToolType extends `tool-${infer N}` ? N : never; // "search" | "robot_switch" | ...

export function isToolType(t: Part["type"]): t is ToolType {
  return t.startsWith("tool-");
}

export function getToolName(
  part: UIMessagePart<UIDataTypes, StorydenTools>,
): string {
  if (!isToolOrDynamicToolUIPart(part)) {
    return "Unknown";
  }

  // Strip off 'tool-' and convert to Title Case
  const rawName = part.type.replace(/^tool-/, "");
  return rawName
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}
