"use client";

import { useEffect, useRef } from "react";

import { useRobotChat } from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { RobotMessage } from "@/components/site/CommandPalette/RobotChat/RobotMessage";
import { EmptyState } from "@/components/site/EmptyState";
import { Admonition } from "@/components/ui/admonition";
import { Box, VStack } from "@/styled-system/jsx";

export function FullPageMessageList() {
  const { messages, errorState, handleDismissError } = useRobotChat();
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <Box
      w="full"
      flex="1"
      minH="0"
      gap="3"
      overflowY="auto"
      position="relative"
    >
      <Box
        position="sticky"
        top="0"
        left="0"
        right="0"
        h="8"
        pointerEvents="none"
        zIndex="dropdown"
        background="scroll-fade-top"
      />
      <VStack alignItems="stretch" my="4">
        {messages.length === 0 ? (
          <EmptyState authenticatedLabel="use robots by talking to them">
            no messages yet
          </EmptyState>
        ) : (
          messages.map((message: any) => (
            <RobotMessage
              key={message["id"]}
              role={message["role"]}
              parts={message["parts"] || []}
            />
          ))
        )}
        <Admonition value={Boolean(errorState)} onChange={handleDismissError}>
          <p>{errorState}</p>
        </Admonition>
        <div ref={bottomRef} />
      </VStack>
    </Box>
  );
}
