import { StorydenUIMessage } from "@/api/robots-types";
import { useSession } from "@/auth";

import { EmptyState } from "../../EmptyState";

import { useRobotChat } from "./RobotChatContext";
import { projectRobotMessages } from "./RobotChatMessageProjection";
import { RobotMessage } from "./RobotMessage";

export function RobotChatMessageProjectedList() {
  const { messages } = useRobotChat();
  const session = useSession();

  if (messages.length === 0) {
    return (
      <EmptyState authenticatedLabel="use robots by talking to them">
        no messages yet
      </EmptyState>
    );
  }

  const projectedMessages = projectRobotMessages(messages);
  const latestUserMessageId = findLatestUserMessageId(projectedMessages);

  return projectedMessages.map((message) => (
    <RobotMessage
      key={message.id}
      id={message.id}
      role={message.role}
      parts={message.parts ?? []}
      author={message.author}
      isCurrentMemberMessage={message.author?.id === session?.id}
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
