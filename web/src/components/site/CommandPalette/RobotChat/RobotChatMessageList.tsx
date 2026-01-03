import { useEffect, useRef } from "react";

import { Admonition } from "@/components/ui/admonition";
import { VStack } from "@/styled-system/jsx";

import { useRobotChat } from "./RobotChatContext";
import { RobotMessage } from "./RobotMessage";

export function RobotChatMessageList() {
  const { messages, errorState, handleDismissError } = useRobotChat();
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // This auto-scrolls the chat, it's super naïve because if you're scrolled
    // up intentionally and reading something while the agent responds, you get
    // rudely interrupted, but we'll revisit in future.
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  if (messages.length === 0) {
    return null;
  }

  return (
    <VStack
      w="full"
      gap="3"
      p="4"
      maxH="96"
      overflowY="auto"
      alignItems="stretch"
    >
      {messages.map((message: any) => (
        <RobotMessage
          key={message["id"]}
          role={message["role"]}
          parts={message["parts"] || []}
        />
      ))}
      <Admonition value={Boolean(errorState)} onChange={handleDismissError}>
        <p>{errorState}</p>
      </Admonition>
      <div ref={bottomRef} />
    </VStack>
  );
}
