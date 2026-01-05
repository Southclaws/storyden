"use client";

import { useChat } from "@ai-sdk/react";
import { DefaultChatTransport } from "ai";
import {
  PropsWithChildren,
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
} from "react";
import { useSWRConfig } from "swr";

import { getRobotsListKey, useRobotsList } from "@/api/openapi-client/robots";
import { Identifier, Robot, RobotList } from "@/api/openapi-schema";
import { API_ADDRESS } from "@/config";
import { generateXid } from "@/utils/xid";

type RobotChatContextValue = {
  sessionId: string;
  handleReset: () => void;
  selectedRobot?: Robot;
  setSelectedRobot: (r: Robot | undefined) => void;
  robots: RobotList;
  sendMessage: (input: { text: string }) => Promise<void>;
  messages: ReturnType<typeof useChat>["messages"];
  status: ReturnType<typeof useChat>["status"];
};

const context = createContext<RobotChatContextValue | undefined>(undefined);

export function useRobotChat() {
  const value = useContext(context);
  if (value === undefined) {
    throw new Error("useRobotChat must be used within a RobotChatContext");
  }

  return value;
}

export function RobotChatContext({ children }: PropsWithChildren) {
  const { data, error } = useRobotsList();
  const { mutate } = useSWRConfig();
  const [selectedRobot, setSelectedRobot] = useState<Robot | undefined>(
    undefined,
  );
  const [sessionId, setSessionId] = useState(() => generateXid());

  const transport = useMemo(() => {
    return new DefaultChatTransport({
      api: `${API_ADDRESS}/sse/chat`,
      credentials: "include",
    });
  }, []);

  const chat = useChat({
    id: sessionId,
    transport,
    onToolCall: async ({ toolCall }) => {
      // Revalidate robots list for any mutative robot tool calls
      const mutativeRobotTools = [
        "create_robot",
        "update_robot",
        "delete_robot",
      ];
      if (mutativeRobotTools.includes(toolCall.toolName)) {
        await mutate(getRobotsListKey());
      }

      if (toolCall.toolName === "switch_agent") {
        // Execute the tool on the client side
        const input = toolCall.input as { robot_id: string };

        // Find the robot in the list and switch to it
        const robot = data?.robots.find((r) => r.id === input.robot_id);
        if (robot) {
          setSelectedRobot(robot);
        }

        // Add tool result (will auto-submit when complete)
        chat.addToolOutput({
          tool: toolCall.toolName,
          toolCallId: toolCall.toolCallId,
          state: "output-available",
          output: { success: true, robot_id: input.robot_id },
        });
      }
    },
    // WORKAROUND: Custom sendAutomaticallyWhen to fix Vercel AI SDK bug
    //
    // The built-in `lastAssistantMessageIsCompleteWithToolCalls` helper has a bug
    // where it returns true even AFTER the agent has responded with text, causing
    // infinite auto-submission loops.
    //
    // This custom implementation adds the missing check: only auto-send when tool
    // outputs are ready but the agent hasn't responded yet. Once the agent sends
    // a text response, we stop auto-sending.
    //
    // Related: https://github.com/vercel/ai/issues/7717
    sendAutomaticallyWhen: (opts) => {
      const lastMessage = opts.messages[opts.messages.length - 1];

      if (!lastMessage || lastMessage["role"] !== "assistant") {
        return false;
      }

      const parts = lastMessage["parts"];
      if (!parts || !Array.isArray(parts)) {
        return false;
      }

      const hasToolOutputs = parts.some(
        (p: any) =>
          p.type?.startsWith("tool-") && p.state === "output-available",
      );
      const hasTextResponse = parts.some((p: any) => p.type === "text");

      // Only auto-send if we have tool outputs but NO text response yet
      return hasToolOutputs && !hasTextResponse;
    },
  });

  const handleReset = () => {
    setSessionId(generateXid());
  };

  // Wrapper around chat.sendMessage that includes robot_id in body
  const sendMessage = useCallback(
    async (input: { text: string }) => {
      await chat.sendMessage(input, {
        body: {
          robotId: selectedRobot?.id,
        },
      });
    },
    [chat.sendMessage, selectedRobot?.id],
  );

  const value: RobotChatContextValue = {
    sessionId,
    handleReset,
    selectedRobot,
    setSelectedRobot,
    robots: data?.robots ?? [],
    sendMessage,
    messages: chat.messages,
    status: chat.status,
  };

  return <context.Provider value={value}>{children}</context.Provider>;
}
