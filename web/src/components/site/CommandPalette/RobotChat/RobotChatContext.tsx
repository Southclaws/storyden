"use client";

import { useChat } from "@ai-sdk/react";
import { ChatOnToolCallCallback, DefaultChatTransport, UIMessage } from "ai";
import type { JSONSchema7 } from "json-schema";
import {
  PropsWithChildren,
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { useSWRConfig } from "swr";

import {
  getRobotsListKey,
  robotSessionGet,
  useRobotsList,
} from "@/api/openapi-client/robots";
import { getThreadListKey } from "@/api/openapi-client/threads";
import { Robot, RobotList } from "@/api/openapi-schema";
import {
  TOOL_NAMES,
  ToolName,
  ToolRobotSwitchInput,
  ToolRobotSwitchOutput,
} from "@/api/robots";
import mcpSchema from "@/api/robots.json";
import { API_ADDRESS } from "@/config";
import { deriveError } from "@/utils/error";
import { generateXid } from "@/utils/xid";

import { useRobotPageContext } from "./useRobotChatContext";

const MUTATIVE_ROBOT_TOOLS: ToolName[] = [
  "robot_create",
  "robot_update",
  "robot_delete",
];

const MUTATIVE_THREAD_TOOLS: ToolName[] = [
  "thread_create",
  "thread_update",
  "thread_reply",
];

function threadListKeyFilterFn(key: unknown) {
  if (!Array.isArray(key)) return false;
  return key[0] === getThreadListKey()[0];
}

const typedSchema = mcpSchema as JSONSchema7;

type RobotChatContextValue = {
  sessionId: string;
  handleReset: () => void;
  selectedRobot?: Robot;
  setSelectedRobot: (r: Robot | undefined) => void;
  robots: RobotList;
  sendMessage: (input: { text: string }) => Promise<void>;
  messages: ReturnType<typeof useChat>["messages"];
  status: ReturnType<typeof useChat>["status"];
  errorState?: string;
  handleDismissError: () => void;
};

const context = createContext<RobotChatContextValue | undefined>(undefined);

export function useRobotChat() {
  const value = useContext(context);
  if (value === undefined) {
    throw new Error("useRobotChat must be used within a RobotChatContext");
  }

  return value;
}

type RobotChatContextProps = PropsWithChildren<{
  sessionIdToLoad?: string;
}>;

export function RobotChatContext({
  children,
  sessionIdToLoad,
}: RobotChatContextProps) {
  // TODO: Expose this error in a user-friendly manner.
  const { data, error } = useRobotsList();
  // NOTE: Annoying workaround for useChat caching the onToolCall function...
  const dataRef = useRef(data);
  const errorRef = useRef(error);

  const { mutate } = useSWRConfig();
  const [selectedRobot, setSelectedRobot] = useState<Robot | undefined>(
    undefined,
  );
  const [sessionId, setSessionId] = useState(
    () => sessionIdToLoad ?? generateXid(),
  );
  const [initialMessages, setInitialMessages] = useState<
    UIMessage[] | undefined
  >(undefined);
  const [errorState, setErrorState] = useState<string | undefined>(undefined);
  const getPageContext = useRobotPageContext();

  useEffect(() => {
    dataRef.current = data;
    errorRef.current = error;
  }, [data, error]);

  // Load session messages if sessionIdToLoad is provided
  useEffect(() => {
    if (!sessionIdToLoad) return;

    async function loadSessionData(id: string) {
      try {
        const session = await robotSessionGet(id);

        // Extract messages from the paginated result
        const messages = session.message_list.messages;

        setInitialMessages(messages);
      } catch (error) {
        console.error("Failed to load session:", error);
        setErrorState("Failed to load chat session");
      }
    }

    loadSessionData(sessionIdToLoad);
  }, [sessionIdToLoad]);

  const transport = useMemo(() => {
    return new DefaultChatTransport({
      api: `${API_ADDRESS}/sse/chat`,
      credentials: "include",
    });
  }, []);

  const handleToolCall = useCallback<ChatOnToolCallCallback<UIMessage>>(
    async ({ toolCall }) => {
      console.debug("[RobotChat] onToolCall", toolCall);

      const currentData = dataRef.current;
      const currentError = errorRef.current;

      if (!currentData) {
        throw new Error("Robot list not loaded yet");
      }

      if (currentError) {
        throw new Error(
          `Cannot perform tool call: ${deriveError(currentError)}`,
        );
      }

      const toolName = toolCall.toolName as ToolName;
      if (!TOOL_NAMES.includes(toolName)) {
        console.warn(`Unknown tool name: ${toolName} list: ${TOOL_NAMES}`);
        return;
      }

      switch (toolName) {
        case "robot_switch": {
          const input = toolCall.input as ToolRobotSwitchInput;

          const robot = currentData.robots.find((r) => r.id === input.robot_id);
          if (!robot) {
            console.error(
              `Robot not found: ${input.robot_id} list: ${currentData.robots}`,
            );
            return;
          }

          console.debug(
            "[RobotChat] Switching to robot:",
            robot.id,
            robot.name,
          );

          setSelectedRobot(robot);

          const output: ToolRobotSwitchOutput = {
            success: true,
            robot_id: input.robot_id,
          };

          chat.addToolOutput({
            tool: toolCall.toolName,
            toolCallId: toolCall.toolCallId,
            state: "output-available",
            output,
          });

          return;
        }
      }

      // NOTE: When a tool is called that internally mutates the robot list
      // (create, update, delete), we need to tell SWR to re-validate the list.
      if (MUTATIVE_ROBOT_TOOLS.includes(toolName)) {
        await mutate(getRobotsListKey());
      }

      // NOTE: When a tool is called that internally mutates threads
      // (create, update, reply), we need to tell SWR to re-validate the feed.
      if (MUTATIVE_THREAD_TOOLS.includes(toolName)) {
        await mutate(threadListKeyFilterFn);
      }
    },
    [mutate],
  );

  const chat = useChat({
    id: sessionId,
    messages: initialMessages,
    transport,
    onError: async (e) => {
      console.error("[RobotChat] Chat error:", e);
      setErrorState(deriveError(e));
    },
    onToolCall: async (p) => {
      try {
        await handleToolCall(p);
      } catch (e) {
        chat.sendMessage({
          role: "system",
          parts: [
            {
              type: "text",
              text: `An error occurred while executing the tool "${p.toolCall.toolName}": ${(e as Error).message}`,
            },
          ],
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

  // Wrapper around chat.sendMessage that includes robot_id and context
  const sendMessage = useCallback(
    async (input: { text: string }) => {
      const pageContext = await getPageContext();
      await chat.sendMessage(input, {
        body: {
          robotId: selectedRobot?.id,
          context: pageContext,
        },
      });
    },
    [chat.sendMessage, selectedRobot?.id, getPageContext],
  );

  function handleDismissError() {
    setErrorState(undefined);
  }

  const value: RobotChatContextValue = {
    sessionId,
    handleReset,
    selectedRobot,
    setSelectedRobot,
    robots: data?.robots ?? [],
    sendMessage,
    messages: chat.messages,
    status: chat.status,
    errorState,
    handleDismissError,
  };

  return <context.Provider value={value}>{children}</context.Provider>;
}
