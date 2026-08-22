import { UIDataTypes, UIMessage, UIMessagePart, isToolUIPart } from "ai";

import { capitalise } from "@/utils/text";

import { RobotReference } from "./openapi-schema/robotReference";
import { RobotSessionMessage } from "./openapi-schema/robotSessionMessage";
import { RobotSessionMessageList } from "./openapi-schema/robotSessionMessageList";
import { StorydenTools } from "./robots";

export type RobotRenderCardData = {
  ref: string;
  kind: string;
  id: string;
};

export type RobotDelegationStatus = "running" | "completed" | "failed";

export interface RobotDelegationData {
  callId: string;
  robot: RobotReference;
  request: string;
  status: RobotDelegationStatus;
  messages: StorydenUIMessage[];
  error?: string;
}

export type StorydenUIDataTypes = {
  session_id: string;
  session_name: string;
  render_card: RobotRenderCardData;
  delegation: RobotDelegationData;
};

export type StorydenUIMessage = UIMessage<
  unknown,
  StorydenUIDataTypes,
  StorydenTools
> &
  Partial<
    Pick<
      RobotSessionMessage,
      | "created_at"
      | "queued"
      | "robot"
      | "author"
      | "branch"
      | "isolation_scope"
    >
  >;

export function toStorydenUIMessages(
  messages: RobotSessionMessageList,
): StorydenUIMessage[] {
  return messages as unknown as StorydenUIMessage[];
}

type Part = StorydenUIMessage["parts"][number];

type ToolType = Extract<Part["type"], `tool-${string}`>; // "tool-search" | "tool-toolset_load" | ...
export type ToolName = ToolType extends `tool-${infer N}` ? N : never; // "search" | "toolset_load" | ...

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

  return rawName.split("_").map(capitalise).join(" ");
}

export function getRawToolName(part: { type: string }): string {
  if (!part.type.startsWith("tool-")) {
    return "";
  }

  return String(part.type).replace(/^tool-/, "");
}
