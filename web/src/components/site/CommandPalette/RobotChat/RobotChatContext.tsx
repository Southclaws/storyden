"use client";

import { useChat } from "@ai-sdk/react";
import { ChatStatus, UIMessageChunk, readUIMessageStream } from "ai";
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
  createDurableChatTransport,
  observeRobotSession,
} from "@/api/durable-chat-transport";
import type { CommandAccepted } from "@/api/durable-chat-transport";
import {
  getRobotToolsetsListKey,
  getRobotsListKey,
  robotSessionGet,
  robotSessionTurnCancel,
  useRobotGet,
  useRobotSessionsList,
  useRobotWorkspacesList,
} from "@/api/openapi-client/robots";
import { getThreadListKey } from "@/api/openapi-client/threads";
import { getTrailListKey } from "@/api/openapi-client/trails";
import {
  RobotSessionList,
  RobotSessionStreamEvent,
  RobotWorkspaceList,
} from "@/api/openapi-schema";
import {
  TOOL_NAMES,
  ToolInputMap,
  ToolLibraryRequestPageOutput,
  ToolName,
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
  "toolset_create",
  "toolset_update",
  "toolset_delete",
];

const MUTATIVE_THREAD_TOOLS: ToolName[] = [
  "thread_create",
  "thread_update",
  "thread_reply",
];

const MUTATIVE_TRAIL_TOOLS: ToolName[] = [
  "trail_create",
  "trail_update",
  "trail_run_create",
  "trail_action_run_cancel",
];

export const DENBOT_NAME = "Denbot";
export const DENBOT_ID = "denbot";

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

type RobotChatContextValue = {
  sessionId: string;
  activeRobotName: string;
  selectedWorkspaceID?: string;
  setSelectedWorkspaceID: (workspaceID: string | undefined) => void;
  workspaces: RobotWorkspaceList;
  workspacesReady: boolean;
  sessions: RobotSessionList;
  sendMessage: (input: { text: string }) => Promise<void>;
  cancelActiveTurn: () => Promise<void>;
  canCancelActiveTurn: boolean;
  isCancelling: boolean;
  messages: StorydenUIMessage[];
  hasOlderMessages: boolean;
  isLoadingOlderMessages: boolean;
  loadOlderMessages: () => Promise<boolean>;
  status: ReturnType<typeof useChat>["status"];
  queuedMessageCount: number;
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
  initialStreamOffset?: string;
  initialActiveTurnID?: string;
  initialRootRobotID?: string;
  initialSelectedWorkspaceID?: string;
}>;

export function RobotChatContext({
  children,
  initialSessionID,
  initialMessages,
  initialNextBefore,
  initialStreamOffset,
  initialActiveTurnID,
  initialRootRobotID,
  initialSelectedWorkspaceID,
}: RobotChatContextProps) {
  const rootRobotID = initialRootRobotID ?? DENBOT_ID;
  const { data: rootRobot } = useRobotGet(rootRobotID, {
    swr: { enabled: rootRobotID !== DENBOT_ID },
  });
  const { data: workspacesData } = useRobotWorkspacesList();
  const { data: sessionsData, mutate: mutateSessionList } =
    useRobotSessionsList();

  const { mutate } = useSWRConfig();
  const [selectedWorkspaceID, setSelectedWorkspaceID] = useState<
    string | undefined
  >(initialSelectedWorkspaceID);
  const selectedWorkspaceIDRef = useRef<string | undefined>(
    initialSelectedWorkspaceID,
  );
  const autoSubmittedToolOutputIDsRef = useRef<Set<string>>(new Set());
  const [sessionId] = useState(() => initialSessionID ?? generateXid());
  const [streamStartOffset] = useState(initialStreamOffset ?? "-1");
  const [isSessionConfirmed, setIsSessionConfirmed] =
    useState(!!initialSessionID);
  const [nextBefore, setNextBefore] = useState<string | undefined>(
    initialNextBefore,
  );
  const [isLoadingOlderMessages, setIsLoadingOlderMessages] = useState(false);
  const [observerStatus, setObserverStatus] = useState<ChatStatus>("ready");
  const [activeTurnID, setActiveTurnID] = useState<string | undefined>(
    initialActiveTurnID,
  );
  const [isCancelling, setIsCancelling] = useState(false);

  const [errorState, setErrorState] = useState<string | undefined>(undefined);
  const getPageContext = useRobotPageContext();

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

  const fetchClient = useMemo<typeof fetch>(() => {
    return async (input, init) => {
      let request: RequestInit = {
        ...init,
        credentials: "include",
      };

      if (init?.method?.toUpperCase() === "POST") {
        const body =
          typeof init.body === "string"
            ? (JSON.parse(init.body) as Record<string, unknown>)
            : {};
        const pageContext = await getPageContext();

        request = {
          ...request,
          body: JSON.stringify({
            ...body,
            robotId: rootRobotID === DENBOT_ID ? undefined : rootRobotID,
            context: pageContext ?? body["context"],
            workspace: selectedWorkspaceID
              ? { workspace_id: selectedWorkspaceID }
              : body["workspace"],
          }),
        };
      }

      return fetch(input, request);
    };
  }, [getPageContext, rootRobotID, selectedWorkspaceID]);

  const handleCommandAccepted = useCallback(
    (command: CommandAccepted) => {
      if (command.sessionId !== sessionId) {
        return;
      }
      setIsSessionConfirmed(true);
    },
    [sessionId],
  );

  const transport = useMemo(() => {
    return createDurableChatTransport<StorydenUIMessage>({
      api: `${API_ADDRESS}/api/robots/sessions`,
      fetchClient,
      onCommandAccepted: handleCommandAccepted,
    });
  }, [fetchClient, handleCommandAccepted]);

  const handleToolCall = useCallback(
    async ({ toolCall }: HandleToolCallOptions) => {
      console.debug("[RobotChat] onToolCall", toolCall);

      if (!isStorydenToolCall(toolCall)) {
        const toolName = toolCall.toolName;
        console.warn(`Unknown tool name: ${toolName} list: ${TOOL_NAMES}`);
        return;
      }

      const toolName = toolCall.toolName;

      if (
        toolName === "robot_delete" ||
        toolName === "toolset_delete" ||
        toolName === "library_request_page"
      ) {
        return;
      }

      // NOTE: When a tool is called that internally mutates the robot list
      // (create, update, delete), we need to tell SWR to re-validate the list.
      if (MUTATIVE_ROBOT_TOOLS.includes(toolName)) {
        await Promise.all([
          mutate(getRobotsListKey()),
          mutate(getRobotToolsetsListKey()),
        ]);
      }

      // NOTE: When a tool is called that internally mutates threads
      // (create, update, reply), we need to tell SWR to re-validate the feed.
      if (MUTATIVE_THREAD_TOOLS.includes(toolName)) {
        await mutate(threadListKeyFilterFn);
      }

      if (MUTATIVE_TRAIL_TOOLS.includes(toolName)) {
        await mutate(getTrailListKey());
      }
    },
    [mutate],
  );

  const handleStreamData = useCallback(
    (message: { type: string; data: unknown }) => {
      switch (message.type) {
        case "data-session_id": {
          if (message.data === sessionId) {
            setIsSessionConfirmed(true);
          }
          break;
        }
        case "data-session_name": {
          setIsSessionConfirmed(true);
          if (!sessionsData || typeof message.data !== "string") return;

          const sessionName = message.data;
          const newData = {
            ...sessionsData,
            sessions: sessionsData.sessions.map((session) =>
              session.id === sessionId
                ? { ...session, name: sessionName }
                : session,
            ),
          };
          mutateSessionList(newData, { revalidate: true });
          break;
        }
      }
    },
    [mutateSessionList, sessionId, sessionsData, setIsSessionConfirmed],
  );

  const chat = useChat<StorydenUIMessage>({
    id: sessionId,
    messages: initialMessages,
    transport,
    resume: false,
    onError: async (e) => {
      console.error("[RobotChat] Chat error:", e);
      setErrorState(deriveError(e));
    },
    onData: async (message) => {
      console.debug(`[RobotChat] Session data`, message);
      handleStreamData(message);
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

  const observerCallbacksRef = useRef({ handleStreamData, handleToolCall });
  const setMessagesRef = useRef(chat.setMessages);
  const observerFetchClient = useMemo<typeof fetch>(() => {
    return (input, init) => fetch(input, { ...init, credentials: "include" });
  }, []);

  useEffect(() => {
    observerCallbacksRef.current = { handleStreamData, handleToolCall };
    setMessagesRef.current = chat.setMessages;
  }, [chat.setMessages, handleStreamData, handleToolCall]);

  useEffect(() => {
    if (!isSessionConfirmed) {
      return;
    }

    const abortController = new AbortController();
    const consumers = new Map<string, TurnMessageConsumer>();
    const queuedParts = new Map<string, UIMessageChunk[]>();
    const handledToolCallIDs = new Set<string>();
    const activeTurnIDs = new Set<string>();

    const updateTurnMessage = (message: StorydenUIMessage) => {
      setMessagesRef.current((messages) => upsertMessage(messages, message));
      for (const part of message.parts) {
        if (
          !part.type.startsWith("tool-") ||
          !("toolCallId" in part) ||
          !part.toolCallId ||
          !("state" in part) ||
          part.state === "output-available" ||
          part.state === "output-error" ||
          part.state === "approval-responded" ||
          handledToolCallIDs.has(part.toolCallId)
        ) {
          continue;
        }
        handledToolCallIDs.add(part.toolCallId);
        void observerCallbacksRef.current
          .handleToolCall({
            toolCall: {
              toolCallId: part.toolCallId,
              toolName: part.type.replace(/^tool-/, ""),
              input: "input" in part ? part.input : undefined,
            },
          })
          .catch((error) => setErrorState(deriveError(error)));
      }
    };

    const consumerFor = (turnID: string) => {
      let consumer = consumers.get(turnID);
      if (consumer) {
        return consumer;
      }
      consumer = createTurnMessageConsumer(updateTurnMessage, (error) =>
        setErrorState(deriveError(error)),
      );
      consumers.set(turnID, consumer);
      for (const part of queuedParts.get(turnID) ?? []) {
        consumer.push(part);
      }
      queuedParts.delete(turnID);
      return consumer;
    };

    const handleEvent = (event: RobotSessionStreamEvent) => {
      const turnID = event.turn_id;
      const parts = event.parts as UIMessageChunk[];
      for (const part of parts) {
        if (part.type.startsWith("data-") && "data" in part) {
          observerCallbacksRef.current.handleStreamData({
            type: part.type,
            data: part.data,
          });
        }
        if (part.type === "error") {
          setErrorState(part.errorText);
        }
      }

      if (event.event_kind === "turn_queued") {
        if (!turnID) return;
        activeTurnIDs.add(turnID);
        setActiveTurnID(turnID);
        queuedParts.set(turnID, parts);
        const claimedInputIDs = new Set(event.input_ids ?? []);
        if (claimedInputIDs.size > 0) {
          setMessagesRef.current((messages) =>
            messages.map((message) =>
              claimedInputIDs.has(message.id)
                ? { ...message, queued: false }
                : message,
            ),
          );
        }
        setObserverStatus("submitted");
        return;
      }

      if (event.message) {
        setMessagesRef.current((messages) =>
          upsertMessage(messages, event.message as StorydenUIMessage),
        );
      } else if (parts.length > 0 && turnID) {
        setObserverStatus("streaming");
        const consumer = consumerFor(turnID);
        for (const part of parts) {
          consumer.push(part);
        }
      }

      if (isTerminalSessionEvent(event)) {
        if (!turnID) return;
        activeTurnIDs.delete(turnID);
        setActiveTurnID((active) => (active === turnID ? undefined : active));
        setIsCancelling(false);
        consumers.get(turnID)?.close();
        consumers.delete(turnID);
        queuedParts.delete(turnID);
        setObserverStatus(activeTurnIDs.size === 0 ? "ready" : "streaming");
      }
    };

    void consumeRobotSession(
      {
        url: `${API_ADDRESS}/api/robots/sessions/${sessionId}/stream`,
        offset: streamStartOffset,
        signal: abortController.signal,
        fetchClient: observerFetchClient,
      },
      handleEvent,
    ).catch((error) => {
      if (!abortController.signal.aborted) {
        setObserverStatus("error");
        setErrorState(deriveError(error));
      }
    });

    return () => {
      abortController.abort();
      for (const consumer of consumers.values()) {
        consumer.close();
      }
    };
  }, [isSessionConfirmed, observerFetchClient, sessionId, streamStartOffset]);

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

  // Wrapper around chat.sendMessage that includes page and workspace context.
  const sendMessage = useCallback(
    async (input: { text: string }) => {
      const pageContext = await getPageContext();
      const currentWorkspaceID = selectedWorkspaceIDRef.current;
      await chat.sendMessage(
        {
          id: generateXid(),
          role: "user",
          parts: [{ type: "text", text: input.text }],
          queued: true,
        },
        {
          body: {
            context: pageContext,
            workspace: currentWorkspaceID
              ? { workspace_id: currentWorkspaceID }
              : undefined,
          },
        },
      );
    },
    [chat.sendMessage, getPageContext],
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

  const cancelActiveTurn = useCallback(async () => {
    if (!activeTurnID || isCancelling) {
      return;
    }

    setIsCancelling(true);
    try {
      await robotSessionTurnCancel(sessionId, activeTurnID);
    } catch (error) {
      setIsCancelling(false);
      setErrorState(deriveError(error));
    }
  }, [activeTurnID, isCancelling, sessionId]);

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
        await Promise.all([
          mutate(getRobotsListKey()),
          mutate(getRobotToolsetsListKey()),
        ]);
      }

      if (MUTATIVE_THREAD_TOOLS.includes(input.toolName)) {
        await mutate(threadListKeyFilterFn);
      }

      if (MUTATIVE_TRAIL_TOOLS.includes(input.toolName)) {
        await mutate(getTrailListKey());
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

  function handleDismissError() {
    setErrorState(undefined);
  }

  const value: RobotChatContextValue = {
    sessionId,
    activeRobotName:
      rootRobot?.name ??
      (rootRobotID === DENBOT_ID ? DENBOT_NAME : "Custom Robot"),
    selectedWorkspaceID,
    setSelectedWorkspaceID: handleSetSelectedWorkspaceID,
    workspaces: workspacesData?.workspaces ?? [],
    workspacesReady: !!workspacesData,
    sessions: sessionsData?.sessions ?? [],
    sendMessage,
    cancelActiveTurn,
    canCancelActiveTurn: !!activeTurnID,
    isCancelling,
    messages: chat.messages,
    hasOlderMessages: Boolean(nextBefore),
    isLoadingOlderMessages,
    loadOlderMessages,
    status: observerStatus === "ready" ? chat.status : observerStatus,
    queuedMessageCount: chat.messages.filter(
      (message) => message.role === "user" && message.queued,
    ).length,
    errorState,
    handleDismissError,
    isSessionConfirmed,
    resolveToolConfirmation,
    resolveLibraryPageRequest,
  };

  return <context.Provider value={value}>{children}</context.Provider>;
}

type TurnMessageConsumer = {
  push: (part: UIMessageChunk) => void;
  close: () => void;
};

async function consumeRobotSession(
  options: Parameters<typeof observeRobotSession>[0],
  onEvent: (event: RobotSessionStreamEvent) => void,
) {
  const events = await observeRobotSession(options);
  for await (const event of events) {
    onEvent(event);
  }
}

function createTurnMessageConsumer(
  onMessage: (message: StorydenUIMessage) => void,
  onError: (error: unknown) => void,
): TurnMessageConsumer {
  let controller: ReadableStreamDefaultController<UIMessageChunk> | undefined;
  let closed = false;
  const input = new ReadableStream<UIMessageChunk>({
    start(streamController) {
      controller = streamController;
    },
  });
  const messages = readUIMessageStream<StorydenUIMessage>({
    stream: input,
    onError,
    terminateOnError: true,
  });
  void (async () => {
    try {
      for await (const message of messages) {
        onMessage(message);
      }
    } catch (error) {
      onError(error);
    }
  })();

  return {
    push(part) {
      if (!closed) {
        controller?.enqueue(part);
      }
    },
    close() {
      if (!closed) {
        closed = true;
        controller?.close();
      }
    },
  };
}

function upsertMessage(
  messages: readonly StorydenUIMessage[],
  incoming: StorydenUIMessage,
) {
  const existingIndex = messages.findIndex(
    (message) => message.id === incoming.id,
  );
  if (existingIndex === -1) {
    return [...messages, incoming];
  }
  return messages.map((message, index) =>
    index === existingIndex ? incoming : message,
  );
}

function isTerminalSessionEvent(event: RobotSessionStreamEvent) {
  return (
    event.event_kind === "turn_completed" ||
    event.event_kind === "turn_blocked" ||
    event.event_kind === "turn_failed" ||
    event.event_kind === "turn_cancelled"
  );
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
