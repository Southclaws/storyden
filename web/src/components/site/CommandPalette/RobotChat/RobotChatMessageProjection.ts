import { isDataUIPart, isToolUIPart } from "ai";

import {
  RobotDelegationData,
  StorydenUIMessage,
  getRawToolName,
} from "@/api/robots-types";

type MessagePart = StorydenUIMessage["parts"][number];
type DelegationPart = Extract<MessagePart, { type: "data-delegation" }>;

type DelegationGroup = {
  callId: string;
  firstMessageIndex: number;
  messages: StorydenUIMessage[];
  robot?: RobotDelegationData["robot"];
  request: string;
  status: RobotDelegationData["status"];
  error?: string;
  inserted: boolean;
};

export function projectRobotMessages(
  messages: readonly StorydenUIMessage[],
): StorydenUIMessage[] {
  const delegated = mergeDelegationParts(projectDelegations(messages));
  const nestedToolCallIds = collectDelegatedToolCallIds(delegated);

  const withoutDuplicatedNestedTools = delegated.map<StorydenUIMessage>(
    (message) => ({
      ...message,
      parts: (message.parts ?? []).filter(
        (part) =>
          !isToolUIPart(part) || !nestedToolCallIds.has(part.toolCallId),
      ),
    }),
  );

  return projectToolOutputs(withoutDuplicatedNestedTools).filter(
    (message) => (message.parts ?? []).length > 0,
  );
}

function mergeDelegationParts(
  messages: readonly StorydenUIMessage[],
): StorydenUIMessage[] {
  const mergedByCallID = new Map<string, RobotDelegationData>();
  const seenMessageIDs = new Map<string, Set<string>>();

  for (const message of messages) {
    for (const part of message.parts ?? []) {
      if (!isRobotDelegationPart(part)) {
        continue;
      }

      const current = mergedByCallID.get(part.data.callId);
      if (!current) {
        mergedByCallID.set(part.data.callId, {
          ...part.data,
          messages: [...part.data.messages],
        });
        seenMessageIDs.set(
          part.data.callId,
          new Set(part.data.messages.map((nested) => nested.id)),
        );
        continue;
      }

      const seen = seenMessageIDs.get(part.data.callId) ?? new Set<string>();
      for (const nested of part.data.messages) {
        if (!seen.has(nested.id)) {
          current.messages.push(nested);
          seen.add(nested.id);
        }
      }
      seenMessageIDs.set(part.data.callId, seen);
      current.robot =
        part.data.robot.name === "Delegated Robot"
          ? current.robot
          : part.data.robot;
      current.request = part.data.request || current.request;
      current.status = part.data.status;
      current.error = part.data.error;
    }
  }

  const emitted = new Set<string>();
  return messages.map<StorydenUIMessage>((message) => ({
    ...message,
    parts: (message.parts ?? []).flatMap<MessagePart>((part) => {
      if (!isRobotDelegationPart(part)) {
        return [part];
      }
      if (emitted.has(part.data.callId)) {
        return [];
      }

      emitted.add(part.data.callId);
      return [
        {
          ...part,
          data: mergedByCallID.get(part.data.callId) ?? part.data,
        },
      ];
    }),
  }));
}

export function projectToolOutputs(
  messages: readonly StorydenUIMessage[],
): StorydenUIMessage[] {
  const inputCallIds = new Set<string>();
  const outputsByCallId = new Map<string, MessagePart>();

  for (const message of messages) {
    for (const part of message.parts ?? []) {
      if (!isToolUIPart(part) || !part.toolCallId) {
        continue;
      }

      if (part.state === "input-available") {
        inputCallIds.add(part.toolCallId);
      }

      if (part.state === "output-available") {
        outputsByCallId.set(part.toolCallId, part);
      }
    }
  }

  const inputCallIdsWithOutput = new Set(
    [...inputCallIds].filter((id) => outputsByCallId.has(id)),
  );
  const seenStandaloneOutputCallIds = new Set<string>();

  return messages.map<StorydenUIMessage>((message) => ({
    ...message,
    parts: (message.parts ?? []).flatMap((part) => {
      if (!isToolUIPart(part) || !part.toolCallId) {
        return [part];
      }

      if (part.state === "input-available") {
        const output = outputsByCallId.get(part.toolCallId);

        return output
          ? [
              {
                ...part,
                ...output,
                input: part.input,
              } as MessagePart,
            ]
          : [part];
      }

      if (part.state !== "output-available") {
        return [part];
      }

      if (inputCallIdsWithOutput.has(part.toolCallId)) {
        return [];
      }

      if (seenStandaloneOutputCallIds.has(part.toolCallId)) {
        return [];
      }

      seenStandaloneOutputCallIds.add(part.toolCallId);

      return [part];
    }),
  }));
}

export function isRobotDelegationPart(
  part: MessagePart,
): part is DelegationPart {
  return isDataUIPart(part) && part.type === "data-delegation";
}

function projectDelegations(
  messages: readonly StorydenUIMessage[],
): StorydenUIMessage[] {
  const groups = new Map<string, DelegationGroup>();
  const inferredScopes = inferDelegationScopes(messages);

  messages.forEach((message, index) => {
    const isolationScope = delegationScope(message, index, inferredScopes);
    if (!isolationScope) {
      return;
    }

    const group = groups.get(isolationScope) ?? {
      callId: isolationScope,
      firstMessageIndex: index,
      messages: [],
      request: "",
      status: "running",
      inserted: false,
    };
    group.messages.push(message);
    group.robot ??= message.robot;
    groups.set(isolationScope, group);
  });

  for (const message of messages) {
    for (const part of message.parts ?? []) {
      if (!isToolUIPart(part) || !part.toolCallId) {
        continue;
      }

      const group = groups.get(part.toolCallId);
      if (!group || !isDelegationToolPart(part)) {
        continue;
      }

      if (part.state === "input-available") {
        group.request = readDelegationRequest(part.input);
        group.robot ??= robotReferenceFromToolPart(part);
      } else if (part.state === "output-available") {
        group.status = "completed";
      } else if (part.state === "output-error") {
        group.status = "failed";
        group.error = "errorText" in part ? String(part.errorText) : undefined;
      }
    }
  }

  return messages.flatMap<StorydenUIMessage>((message, index) => {
    const isolationScope = delegationScope(message, index, inferredScopes);
    if (isolationScope) {
      const group = groups.get(isolationScope);
      if (group && !group.inserted && index === group.firstMessageIndex) {
        group.inserted = true;
        return [delegationMessage(message, group)];
      }
      return [];
    }

    const parts = (message.parts ?? []).flatMap<MessagePart>((part) => {
      if (!isToolUIPart(part) || !part.toolCallId) {
        return [part];
      }

      const group = groups.get(part.toolCallId);
      if (!group || !isDelegationToolPart(part)) {
        return [part];
      }

      if (group.inserted || part.state !== "input-available") {
        return [];
      }

      group.inserted = true;
      return [delegationPart(group)];
    });

    return [{ ...message, parts }];
  });
}

function delegationScope(
  message: StorydenUIMessage,
  index: number,
  inferredScopes: ReadonlyMap<number, string>,
): string | undefined {
  if (!message.branch) {
    return undefined;
  }

  return message.isolation_scope ?? inferredScopes.get(index);
}

function inferDelegationScopes(
  messages: readonly StorydenUIMessage[],
): Map<number, string> {
  const inferred = new Map<number, string>();
  const activeByRobot = new Map<string, string[]>();

  messages.forEach((message, index) => {
    for (const part of message.parts ?? []) {
      if (
        !isToolUIPart(part) ||
        !part.toolCallId ||
        !isDelegationToolPart(part)
      ) {
        continue;
      }

      const robotAgentName = getRawToolName(part);
      if (part.state === "input-available") {
        const active = activeByRobot.get(robotAgentName) ?? [];
        active.push(part.toolCallId);
        activeByRobot.set(robotAgentName, active);
      }
    }

    if (!message.isolation_scope && message.branch && message.robot?.id) {
      const active = activeByRobot.get(`robot_${message.robot.id}`);
      const callId = active?.at(-1);
      if (callId) {
        inferred.set(index, callId);
      }
    }

    for (const part of message.parts ?? []) {
      if (
        !isToolUIPart(part) ||
        !part.toolCallId ||
        !isDelegationToolPart(part) ||
        (part.state !== "output-available" && part.state !== "output-error")
      ) {
        continue;
      }

      const robotAgentName = getRawToolName(part);
      const active = activeByRobot.get(robotAgentName) ?? [];
      activeByRobot.set(
        robotAgentName,
        active.filter((callId) => callId !== part.toolCallId),
      );
    }
  });

  return inferred;
}

function delegationMessage(
  source: StorydenUIMessage,
  group: DelegationGroup,
): StorydenUIMessage {
  return {
    id: `delegation-${group.callId}`,
    role: "assistant",
    parts: [delegationPart(group)],
    created_at: source.created_at,
  };
}

function delegationPart(group: DelegationGroup): DelegationPart {
  const robot = group.robot ?? {
    id: getRobotIDFromGroup(group),
    name: "Delegated Robot",
  };

  return {
    type: "data-delegation",
    id: group.callId,
    data: {
      callId: group.callId,
      robot,
      request: group.request,
      status: group.status,
      messages: group.messages,
      error: group.error,
    },
  };
}

function collectDelegatedToolCallIds(
  messages: readonly StorydenUIMessage[],
): Set<string> {
  const ids = new Set<string>();

  for (const message of messages) {
    for (const part of message.parts ?? []) {
      if (!isRobotDelegationPart(part)) {
        continue;
      }

      for (const nestedMessage of part.data.messages) {
        for (const nestedPart of nestedMessage.parts ?? []) {
          if (isToolUIPart(nestedPart) && nestedPart.toolCallId) {
            ids.add(nestedPart.toolCallId);
          }
        }
      }
    }
  }

  return ids;
}

function isDelegationToolPart(part: MessagePart): boolean {
  return isToolUIPart(part) && getRawToolName(part).startsWith("robot_");
}

function readDelegationRequest(input: unknown): string {
  if (!input || typeof input !== "object" || !("request" in input)) {
    return "";
  }

  return typeof input.request === "string" ? input.request.trim() : "";
}

function robotReferenceFromToolPart(
  part: MessagePart,
): RobotDelegationData["robot"] | undefined {
  const rawName = getRawToolName(part);
  const id = rawName.startsWith("robot_") ? rawName.slice("robot_".length) : "";
  return id ? { id, name: "Delegated Robot" } : undefined;
}

function getRobotIDFromGroup(group: DelegationGroup): string {
  for (const message of group.messages) {
    if (message.robot?.id) {
      return message.robot.id;
    }
  }
  return "unknown";
}
