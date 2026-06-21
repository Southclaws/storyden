import { isToolUIPart } from "ai";

import { StorydenUIMessage } from "@/api/robots-types";

import { EmptyState } from "../../EmptyState";

import { useRobotChat } from "./RobotChatContext";
import { RobotMessage } from "./RobotMessage";

export function RobotChatMessageProjectedList() {
  const { messages } = useRobotChat();

  if (messages.length === 0) {
    return (
      <EmptyState authenticatedLabel="use robots by talking to them">
        no messages yet
      </EmptyState>
    );
  }

  const projectedMessages = projectToolOutputs(messages);
  const latestUserMessageId = findLatestUserMessageId(projectedMessages);

  return projectedMessages.map((message) => (
    <RobotMessage
      key={message.id}
      id={message.id}
      role={message.role}
      parts={message.parts ?? []}
      isNewestUserMessage={message.id === latestUserMessageId}
    />
  ));
}

function findLatestUserMessageId(
  messages: readonly StorydenUIMessage[],
): string | undefined {
  for (let i = messages.length - 1; i >= 0; i -= 1) {
    const message = messages[i];
    if (!message) {
      continue;
    }

    if (message.role === "user") {
      return message.id;
    }
  }

  return undefined;
}

export function projectToolOutputs(
  messages: readonly StorydenUIMessage[],
): StorydenUIMessage[] {
  const inputCallIds = new Set<string>();
  const outputsByCallId = new Map<string, StorydenUIMessage["parts"][number]>();

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
              } as StorydenUIMessage["parts"][number],
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
