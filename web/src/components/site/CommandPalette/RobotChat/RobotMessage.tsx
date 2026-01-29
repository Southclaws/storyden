import { UIMessagePart } from "ai";
import Markdown from "react-markdown";

import { Box, VStack } from "@/styled-system/jsx";

import { RobotToolCall } from "./RobotToolCall";

type Props = {
  role: "user" | "assistant";
  parts: readonly UIMessagePart<any, any>[];
};

export function RobotMessage({ role, parts }: Props) {
  const textContent = parts
    .filter((p) => p.type === "text" || p.type === "reasoning")
    .map((p) => ("text" in p ? p.text : ""))
    .filter(Boolean)
    .join("\n");

  const toolCalls = parts.filter((p) => p.type?.startsWith("tool-"));

  return (
    <VStack gap="2" alignItems={role === "user" ? "flex-end" : "flex-start"}>
      {toolCalls.map((toolCall, idx) => (
        <RobotToolCall key={idx} part={toolCall} />
      ))}

      {textContent && (
        <Box
          bg={role === "user" ? "bg.subtle" : "bg.muted"}
          borderRadius="md"
          p="3"
          w="full"
          maxW="5/6"
        >
          <Markdown className="typography">{textContent}</Markdown>
        </Box>
      )}
    </VStack>
  );
}
