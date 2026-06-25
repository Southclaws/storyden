"use client";

import { useChat } from "@ai-sdk/react";
import { DefaultChatTransport } from "ai";
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
  useRobotSessionsList,
  useRobotWorkspacesList,
  useRobotsList,
} from "@/api/openapi-client/robots";
import { getThreadListKey } from "@/api/openapi-client/threads";
import {
  Robot,
  RobotList,
  RobotSessionList,
  RobotWorkspaceList,
} from "@/api/openapi-schema";
import {
  TOOL_NAMES,
  ToolInputMap,
  ToolLibraryRequestPageOutput,
  ToolName,
  ToolRobotSwitchOutput,
} from "@/api/robots";
import { StorydenUIMessage, toStorydenUIMessages } from "@/api/robots-types";
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

export const DEFAULT_ROBOT_ID = "robot_builder";

export const PLUGIN_BUILDER_ROBOT_ID = "plugin_builder";

export const DEFAULT_ROBOT_NAME = "Storyden Robot Builder";

export type RobotSelection = Pick<Robot, "id" | "name">;

export const DEFAULT_ROBOT: RobotSelection = {
  id: DEFAULT_ROBOT_ID,
  name: DEFAULT_ROBOT_NAME,
};

export const BUILT_IN_ROBOTS: RobotSelection[] = [
  DEFAULT_ROBOT,
  {
    id: PLUGIN_BUILDER_ROBOT_ID,
    name: "Plugin Builder",
  },
];

function threadListKeyFilterFn(key: unknown) {
  if (!Array.isArray(key)) return false;
  return key[0] === getThreadListKey()[0];
}

const typedSchema = mcpSchema as JSONSchema7;

type ClientToolCall = {
  toolName: string;
  toolCallId: string;
  input: unknown;
  dynamic?: boolean;
};

type HandleToolCallOptions = {
  toolCall: ClientToolCall;
};

type StorydenToolCall = {
  [K in ToolName]: {
    toolName: K;
    toolCallId: string;
    input: ToolInputMap[K];
    dynamic?: false;
  };
}[ToolName];

function isKnownToolName(name: string): name is ToolName {
  return (TOOL_NAMES as readonly string[]).includes(name);
}

function isStorydenToolCall(
  toolCall: ClientToolCall,
): toolCall is StorydenToolCall {
  return !toolCall.dynamic && isKnownToolName(toolCall.toolName);
}

function findRobotSelection(
  id: string,
  robots: RobotList,
): RobotSelection | undefined {
  return (
    BUILT_IN_ROBOTS.find((robot) => robot.id === id) ??
    robots.find((robot) => robot.id === id)
  );
}

type RobotChatContextValue = {
  sessionId: string;
  activeRobotName: string;
  selectedRobot?: RobotSelection;
  setSelectedRobot: (r: RobotSelection | undefined) => void;
  robots: RobotList;
  selectedWorkspaceID?: string;
  setSelectedWorkspaceID: (workspaceID: string | undefined) => void;
  workspaces: RobotWorkspaceList;
  workspacesReady: boolean;
  sessions: RobotSessionList;
  sendMessage: (input: { text: string }) => Promise<void>;
  stopGenerating: () => Promise<void>;
  messages: StorydenUIMessage[];
  hasOlderMessages: boolean;
  isLoadingOlderMessages: boolean;
  loadOlderMessages: () => Promise<boolean>;
  status: ReturnType<typeof useChat>["status"];
  errorState?: string;
  handleDismissError: () => void;
  isSessionConfirmed: boolean;
  resolveToolConfirmation: (input: {
    approvalId: string;
    toolName?: string;
    approved: boolean;
  }) => Promise<void>;
  resolveLibraryPageRequest: (input: {
    toolCallId: string;
    page: ToolLibraryRequestPageOutput;
  }) => Promise<void>;
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
  initialSessionID?: string;
  initialMessages?: StorydenUIMessage[];
  initialNextBefore?: string;
  initialSelectedRobotID?: string;
  initialSelectedWorkspaceID?: string;
}>;

export function RobotChatContext({
  children,
  initialSessionID,
  initialMessages,
  initialNextBefore,
  initialSelectedRobotID,
  initialSelectedWorkspaceID,
}: RobotChatContextProps) {
  // TODO: Expose this error in a user-friendly manner.
  const { data, error } = useRobotsList();
  const { data: workspacesData } = useRobotWorkspacesList();
  const { data: sessionsData, mutate: mutateSessionList } =
    useRobotSessionsList();
  // NOTE: Annoying workaround for useChat caching the onToolCall function...
  const dataRef = useRef(data);
  const errorRef = useRef(error);

  const { mutate } = useSWRConfig();
  const [selectedRobot, setSelectedRobot] = useState<
    RobotSelection | undefined
  >(DEFAULT_ROBOT);
  const selectedRobotRef = useRef<RobotSelection | undefined>(DEFAULT_ROBOT);
  const [selectedWorkspaceID, setSelectedWorkspaceID] = useState<
    string | undefined
  >(initialSelectedWorkspaceID);
  const selectedWorkspaceIDRef = useRef<string | undefined>(
    initialSelectedWorkspaceID,
  );
  const autoSubmittedToolOutputIDsRef = useRef<Set<string>>(new Set());
  const didHydrateInitialSelectedRobotRef = useRef(false);
  const [sessionId] = useState(() => initialSessionID ?? generateXid());
  const [isSessionConfirmed, setIsSessionConfirmed] =
    useState(!!initialSessionID);
  const [nextBefore, setNextBefore] = useState<string | undefined>(
    initialNextBefore,
  );
  const [isLoadingOlderMessages, setIsLoadingOlderMessages] = useState(false);

  const [errorState, setErrorState] = useState<string | undefined>(undefined);
  const getPageContext = useRobotPageContext();

  useEffect(() => {
    dataRef.current = data;
    errorRef.current = error;
  }, [data, error]);

  useEffect(() => {
    if (didHydrateInitialSelectedRobotRef.current) return;
    if (!data) return;

    didHydrateInitialSelectedRobotRef.current = true;

    if (!initialSelectedRobotID) {
      selectedRobotRef.current = DEFAULT_ROBOT;
      setSelectedRobot(DEFAULT_ROBOT);
      return;
    }

    const robot = findRobotSelection(initialSelectedRobotID, data.robots);
    if (robot) {
      selectedRobotRef.current = robot;
      setSelectedRobot(robot);
    }
  }, [data, initialSelectedRobotID]);

  const handleSetSelectedRobot = useCallback(
    (robot: RobotSelection | undefined) => {
      const nextRobot = robot ?? DEFAULT_ROBOT;
      selectedRobotRef.current = nextRobot;
      setSelectedRobot(nextRobot);
    },
    [],
  );

  const handleSetSelectedWorkspaceID = useCallback(
    (workspaceID: string | undefined) => {
      selectedWorkspaceIDRef.current = workspaceID;
      setSelectedWorkspaceID(workspaceID);
    },
    [],
  );

  useEffect(() => {
    if (!selectedWorkspaceID || !workspacesData) {
      return;
    }

    const exists = workspacesData.workspaces.some(
      (workspace) => workspace.id === selectedWorkspaceID,
    );
    if (!exists) {
      handleSetSelectedWorkspaceID(undefined);
    }
  }, [handleSetSelectedWorkspaceID, selectedWorkspaceID, workspacesData]);

  const transport = useMemo(() => {
    return new DefaultChatTransport({
      api: `${API_ADDRESS}/sse/chat`,
      credentials: "include",
      prepareSendMessagesRequest: async (request) => {
        const pageContext = await getPageContext();
        const currentSelectedRobot = selectedRobotRef.current ?? DEFAULT_ROBOT;
        const currentWorkspaceID = selectedWorkspaceIDRef.current;

        return {
          ...request,
          body: {
            ...request.body,
            id: request.id,
            messages: request.messages,
            trigger: request.trigger,
            messageId: request.messageId,
            robotId: currentSelectedRobot.id ?? request.body?.["robotId"],
            context: pageContext ?? request.body?.["context"],
            workspace: currentWorkspaceID
              ? { workspace_id: currentWorkspaceID }
              : undefined,
          },
        };
      },
    });
  }, [getPageContext]);

  const handleToolCall = useCallback(
    async ({ toolCall }: HandleToolCallOptions) => {
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

      if (!isStorydenToolCall(toolCall)) {
        const toolName = toolCall.toolName;
        console.warn(`Unknown tool name: ${toolName} list: ${TOOL_NAMES}`);
        return;
      }

      const toolName = toolCall.toolName;

      if (toolName === "robot_delete" || toolName === "library_request_page") {
        return;
      }

      switch (toolName) {
        case "robot_switch": {
          const input = toolCall.input;

          const robot = findRobotSelection(input.robot_id, currentData.robots);
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

          handleSetSelectedRobot(robot);

          const output: ToolRobotSwitchOutput = {
            success: true,
            robot_id: input.robot_id,
          };

          chat.addToolOutput({
            tool: toolName,
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
    [handleSetSelectedRobot, mutate],
  );

  const chat = useChat<StorydenUIMessage>({
    id: sessionId,
    messages: initialMessages,
    transport,
    onError: async (e) => {
      console.error("[RobotChat] Chat error:", e);
      setErrorState(deriveError(e));
    },
    onData: async (message) => {
      console.debug(`[RobotChat] Session data`, message);

      switch (message.type) {
        case "data-session_name": {
          // Mark session as confirmed when we receive session name from backend
          setIsSessionConfirmed(true);

          if (!sessionsData) return;

          if (typeof message.data !== "string") return;

          const sessionName = message.data;
          const newData = {
            ...sessionsData,
            sessions: sessionsData?.sessions.map((r) => {
              if (r.id === sessionId) {
                if (r.name === sessionName) {
                  return r;
                }

                console.debug(
                  `[RobotChat] Session name updated: ${sessionName}`,
                );

                return {
                  ...r,
                  name: sessionName,
                };
              }
              return r;
            }),
          };
          mutateSessionList(newData, { revalidate: true });
          break;
        }
      }
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
    // This custom implementation adds the missing check: auto-send when tool
    // outputs are ready, unless the assistant has already responded with text
    // after the tool outputs. Text before a tool call is allowed; confirmations
    // often look like "Deleting this now." followed by a tool call.
    //
    // Keep the completeness predicate local so Robot client tools can batch
    // multiple confirmations and skip provider-executed tools predictably.
    //
    // Related: https://github.com/vercel/ai/issues/7717
    sendAutomaticallyWhen: (opts) => {
      const lastMessage = opts.messages[opts.messages.length - 1];

      if (!lastMessage || lastMessage.role !== "assistant") {
        return false;
      }

      if (assistantHasTextAfterToolOutput(lastMessage)) {
        return false;
      }

      if (!assistantToolOutputsAreComplete(lastMessage)) {
        return false;
      }

      const completedToolOutputIDs =
        assistantCompletedToolOutputIDs(lastMessage);
      if (
        completedToolOutputIDs.every((id) =>
          autoSubmittedToolOutputIDsRef.current.has(id),
        )
      ) {
        return false;
      }

      for (const id of completedToolOutputIDs) {
        autoSubmittedToolOutputIDsRef.current.add(id);
      }

      return true;
    },
  });

  const initialMessagesSignature = useMemo(
    () => messageListSignature(initialMessages),
    [initialMessages],
  );
  const chatMessagesSignature = useMemo(
    () => messageListSignature(chat.messages),
    [chat.messages],
  );

  useEffect(() => {
    if (!initialMessages) {
      return;
    }

    if (chat.status === "submitted" || chat.status === "streaming") {
      return;
    }

    if (initialMessagesSignature === chatMessagesSignature) {
      return;
    }

    if (!shouldReplaceMessages(chat.messages, initialMessages)) {
      return;
    }

    if (hasUnhydratedToolOutput(chat.messages, initialMessages)) {
      return;
    }

    chat.setMessages(reconcileMessages(chat.messages, initialMessages));
  }, [
    chat,
    chat.messages.length,
    chat.status,
    chatMessagesSignature,
    initialMessages,
    initialMessagesSignature,
  ]);

  // Wrapper around chat.sendMessage that includes robot_id and context
  const sendMessage = useCallback(
    async (input: { text: string }) => {
      const pageContext = await getPageContext();
      const currentWorkspaceID = selectedWorkspaceIDRef.current;
      await chat.sendMessage(input, {
        body: {
          robotId: selectedRobot?.id,
          context: pageContext,
          workspace: currentWorkspaceID
            ? { workspace_id: currentWorkspaceID }
            : undefined,
        },
      });
    },
    [chat.sendMessage, selectedRobot?.id, getPageContext],
  );

  const loadOlderMessages = useCallback(async () => {
    if (!nextBefore || isLoadingOlderMessages || !isSessionConfirmed) {
      return false;
    }

    setIsLoadingOlderMessages(true);
    try {
      const session = await robotSessionGet(sessionId, {
        before: nextBefore,
        limit: "50",
      });
      const olderMessages = toStorydenUIMessages(
        session.message_list.messages ?? [],
      );

      if (olderMessages.length === 0) {
        setNextBefore(undefined);
        return false;
      }

      const existingIDs = new Set(chat.messages.map((message) => message.id));
      const uniqueOlderMessages = olderMessages.filter(
        (message) => !existingIDs.has(message.id),
      );

      setNextBefore(session.message_list.next_before);

      if (uniqueOlderMessages.length === 0) {
        return false;
      }

      chat.setMessages([...uniqueOlderMessages, ...chat.messages]);
      return true;
    } catch (e) {
      setErrorState(deriveError(e));
      return false;
    } finally {
      setIsLoadingOlderMessages(false);
    }
  }, [
    chat,
    chat.messages,
    isLoadingOlderMessages,
    isSessionConfirmed,
    nextBefore,
    sessionId,
  ]);

  const resolveToolConfirmation = useCallback(
    async (input: {
      approvalId: string;
      toolName?: string;
      approved: boolean;
    }) => {
      await chat.addToolApprovalResponse({
        id: input.approvalId,
        approved: input.approved,
      });

      if (
        !input.approved ||
        !input.toolName ||
        !isKnownToolName(input.toolName)
      ) {
        return;
      }

      if (MUTATIVE_ROBOT_TOOLS.includes(input.toolName)) {
        await mutate(getRobotsListKey());
      }

      if (MUTATIVE_THREAD_TOOLS.includes(input.toolName)) {
        await mutate(threadListKeyFilterFn);
      }
    },
    [chat, mutate],
  );

  const resolveLibraryPageRequest = useCallback(
    async (input: {
      toolCallId: string;
      page: ToolLibraryRequestPageOutput;
    }) => {
      chat.addToolOutput({
        tool: "library_request_page",
        toolCallId: input.toolCallId,
        state: "output-available",
        output: input.page,
      });
    },
    [chat],
  );

  const stopGenerating = useCallback(async () => {
    await chat.stop();
  }, [chat]);

  function handleDismissError() {
    setErrorState(undefined);
  }

  const value: RobotChatContextValue = {
    sessionId,
    activeRobotName: selectedRobot?.name ?? DEFAULT_ROBOT_NAME,
    selectedRobot,
    setSelectedRobot: handleSetSelectedRobot,
    robots: data?.robots ?? [],
    selectedWorkspaceID,
    setSelectedWorkspaceID: handleSetSelectedWorkspaceID,
    workspaces: workspacesData?.workspaces ?? [],
    workspacesReady: !!workspacesData,
    sessions: sessionsData?.sessions ?? [],
    sendMessage,
    stopGenerating,
    messages: chat.messages,
    hasOlderMessages: Boolean(nextBefore),
    isLoadingOlderMessages,
    loadOlderMessages,
    status: chat.status,
    errorState,
    handleDismissError,
    isSessionConfirmed,
    resolveToolConfirmation,
    resolveLibraryPageRequest,
  };

  return <context.Provider value={value}>{children}</context.Provider>;
}

function messageListSignature(messages?: readonly StorydenUIMessage[]) {
  return (messages ?? [])
    .map((message) => {
      const parts = (message.parts ?? [])
        .map((part) => {
          if ("toolCallId" in part && part.toolCallId) {
            return `${part.type}:${part.toolCallId}:${"state" in part ? part.state : ""}`;
          }

          if ("id" in part && part.id) {
            return `${part.type}:${part.id}`;
          }

          if (part.type === "text" && "text" in part) {
            return `${part.type}:${part.text}`;
          }

          return part.type;
        })
        .join(",");

      return `${message.id}:${message.role}:${parts}`;
    })
    .join("|");
}

export function assistantHasTextAfterToolOutput(message: StorydenUIMessage) {
  let sawToolOutput = false;

  for (const part of message.parts) {
    if (
      part.type.startsWith("tool-") &&
      "state" in part &&
      part.state === "output-available"
    ) {
      sawToolOutput = true;
      continue;
    }

    if (
      sawToolOutput &&
      part.type === "text" &&
      "text" in part &&
      part.text.trim().length > 0
    ) {
      return true;
    }
  }

  return false;
}

export function assistantToolOutputsAreComplete(message: StorydenUIMessage) {
  const toolParts = assistantToolPartsInCurrentStep(message);

  return (
    toolParts.length > 0 &&
    toolParts.every(
      (part) =>
        "state" in part &&
        (part.state === "output-available" ||
          part.state === "output-error" ||
          part.state === "approval-responded"),
    )
  );
}

export function assistantCompletedToolOutputIDs(message: StorydenUIMessage) {
  return assistantToolPartsInCurrentStep(message)
    .filter(
      (part) =>
        "state" in part &&
        (part.state === "output-available" ||
          part.state === "output-error" ||
          part.state === "approval-responded"),
    )
    .map((part) => ("toolCallId" in part ? part.toolCallId : undefined))
    .filter((id): id is string => !!id);
}

function assistantToolPartsInCurrentStep(message: StorydenUIMessage) {
  const lastStepStartIndex = message.parts.reduce((lastIndex, part, index) => {
    return part.type === "step-start" ? index : lastIndex;
  }, -1);

  return message.parts
    .slice(lastStepStartIndex + 1)
    .filter((part) => part.type.startsWith("tool-"))
    .filter((part) => !("providerExecuted" in part && part.providerExecuted));
}

export function hasUnhydratedToolOutput(
  localMessages: readonly StorydenUIMessage[],
  incomingMessages: readonly StorydenUIMessage[],
) {
  const incomingToolStates = new Map<string, string>();

  for (const message of incomingMessages) {
    for (const part of message.parts ?? []) {
      if (
        part.type.startsWith("tool-") &&
        "toolCallId" in part &&
        part.toolCallId &&
        "state" in part
      ) {
        incomingToolStates.set(part.toolCallId, part.state);
      }
    }
  }

  for (const message of localMessages) {
    for (const part of message.parts ?? []) {
      if (
        !part.type.startsWith("tool-") ||
        !("toolCallId" in part) ||
        !part.toolCallId ||
        !("state" in part) ||
        part.state !== "output-available"
      ) {
        continue;
      }

      if (incomingToolStates.get(part.toolCallId) !== "output-available") {
        return true;
      }
    }
  }

  return false;
}

export function shouldReplaceMessages(
  localMessages: readonly StorydenUIMessage[],
  incomingMessages: readonly StorydenUIMessage[],
) {
  if (incomingMessages.length === 0) {
    return localMessages.length === 0;
  }

  const localByID = new Map(
    localMessages.map((message) => [message.id, message]),
  );

  for (const incomingMessage of incomingMessages) {
    const localMessage = localByID.get(incomingMessage.id);

    if (!localMessage) {
      return true;
    }

    if (
      messageListSignature([localMessage]) !==
      messageListSignature([incomingMessage])
    ) {
      return true;
    }
  }

  return false;
}

export function reconcileMessages(
  localMessages: readonly StorydenUIMessage[],
  incomingMessages: readonly StorydenUIMessage[],
) {
  if (incomingMessages.length === 0) {
    return [];
  }

  const firstIncomingID = incomingMessages[0]?.id;
  const firstOverlapIndex = localMessages.findIndex(
    (message) => message.id === firstIncomingID,
  );

  if (firstOverlapIndex === -1) {
    return [...incomingMessages];
  }

  return [...localMessages.slice(0, firstOverlapIndex), ...incomingMessages];
}
