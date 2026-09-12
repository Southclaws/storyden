import type { UseChatHelpers } from "@ai-sdk/react";
import { useCallback, useEffect, useRef } from "react";

import type { RobotChatContext } from "@/api/openapi-schema";
import type { ToolName } from "@/api/robots";
import type { StorydenUIMessage } from "@/api/robots-types";
import { deriveError } from "@/utils/error";

import {
  executeWebMCPTool,
  getWebMCPClientID,
  readWebMCPToolCallMetadata,
  subscribeToWebMCPToolChanges,
  webMCPToolScopeMatches,
} from "./webMCPClient";

export type ClientToolCall = {
  toolName: string;
  toolCallId: string;
  input: unknown;
  dynamic?: boolean;
  providerMetadata?: unknown;
};

type ClientToolCallPart = {
  type: string;
  toolCallId: string;
  toolName?: string;
  input?: unknown;
  callProviderMetadata?: unknown;
};

export type WebMCPToolOutputSubmission = {
  toolName: string;
  toolCallId: string;
  output?: unknown;
  errorText?: string;
};

export type ToolOutputSubmitter = (
  input: WebMCPToolOutputSubmission,
) => void | PromiseLike<void>;

export function useWebMCPClientTools(
  getPageContext: () => Promise<RobotChatContext>,
) {
  const pendingToolCallsRef = useRef<Map<string, ClientToolCall>>(new Map());
  const executingToolCallIDsRef = useRef<Set<string>>(new Set());
  const submittedToolCallIDsRef = useRef<Set<string>>(new Set());
  const submitToolOutputRef = useRef<ToolOutputSubmitter>(() => undefined);

  const executePendingToolCall = useCallback(
    async (toolCall: ClientToolCall) => {
      const metadata = readWebMCPToolCallMetadata(toolCall.providerMetadata);
      if (!metadata || metadata.clientID !== getWebMCPClientID()) {
        return;
      }
      if (
        executingToolCallIDsRef.current.has(toolCall.toolCallId) ||
        submittedToolCallIDsRef.current.has(toolCall.toolCallId)
      ) {
        return;
      }

      executingToolCallIDsRef.current.add(toolCall.toolCallId);
      let submission: WebMCPToolOutputSubmission | undefined;
      try {
        const currentPageContext = await getPageContext();
        if (webMCPToolScopeMatches(metadata.scope, currentPageContext)) {
          const result = await executeWebMCPTool(
            toolCall.toolName,
            toolCall.input,
          );
          if (result.status === "completed") {
            submission = {
              toolName: toolCall.toolName,
              toolCallId: toolCall.toolCallId,
              output: result.output,
            };
          }
        }
      } catch (error) {
        submission = {
          toolName: toolCall.toolName,
          toolCallId: toolCall.toolCallId,
          errorText: deriveError(error),
        };
      }
      executingToolCallIDsRef.current.delete(toolCall.toolCallId);

      if (submission) {
        pendingToolCallsRef.current.delete(toolCall.toolCallId);
        submittedToolCallIDsRef.current.add(toolCall.toolCallId);
        void submitToolOutputRef.current(submission);
      }
    },
    [getPageContext],
  );

  const handleWebMCPToolCall = useCallback(
    async (toolCall: ClientToolCall): Promise<boolean> => {
      const metadata = readWebMCPToolCallMetadata(toolCall.providerMetadata);
      if (!metadata) {
        return false;
      }
      if (metadata.clientID === getWebMCPClientID()) {
        pendingToolCallsRef.current.set(toolCall.toolCallId, toolCall);
        await executePendingToolCall(toolCall);
      }
      return true;
    },
    [executePendingToolCall],
  );

  const setToolOutputSubmitter = useCallback(
    (submitter: ToolOutputSubmitter) => {
      submitToolOutputRef.current = submitter;
      return () => {
        submitToolOutputRef.current = () => undefined;
      };
    },
    [],
  );

  useEffect(() => {
    return subscribeToWebMCPToolChanges(() => {
      for (const toolCall of pendingToolCallsRef.current.values()) {
        void executePendingToolCall(toolCall);
      }
    });
  }, [executePendingToolCall]);

  return { handleWebMCPToolCall, setToolOutputSubmitter };
}

function submitWebMCPToolOutput(
  addToolOutput: UseChatHelpers<StorydenUIMessage>["addToolOutput"],
  input: WebMCPToolOutputSubmission,
) {
  if (input.errorText) {
    return addToolOutput({
      tool: input.toolName as ToolName,
      toolCallId: input.toolCallId,
      state: "output-error",
      errorText: input.errorText,
    });
  }
  return addToolOutput({
    tool: input.toolName as ToolName,
    toolCallId: input.toolCallId,
    state: "output-available",
    output: input.output as never,
  });
}

export function useWebMCPToolOutputSubmitter(
  setSubmitter: (submitter: ToolOutputSubmitter) => () => void,
  addToolOutput: UseChatHelpers<StorydenUIMessage>["addToolOutput"],
) {
  useEffect(
    () => setSubmitter((input) => submitWebMCPToolOutput(addToolOutput, input)),
    [addToolOutput, setSubmitter],
  );
}

export function clientToolCallFromPart(
  part: ClientToolCallPart,
): ClientToolCall {
  const dynamic = part.type === "dynamic-tool";
  return {
    toolCallId: part.toolCallId,
    toolName: dynamic ? (part.toolName ?? "") : part.type.replace(/^tool-/, ""),
    input: part.input,
    dynamic,
    providerMetadata: part.callProviderMetadata,
  };
}

export function usePendingWebMCPToolCalls(
  messages: StorydenUIMessage[],
  handleToolCall: (toolCall: ClientToolCall) => Promise<boolean>,
) {
  useEffect(() => {
    for (const toolCall of findPendingWebMCPToolCalls(messages)) {
      void handleToolCall(toolCall);
    }
  }, [handleToolCall, messages]);
}

function findPendingWebMCPToolCalls(
  messages: StorydenUIMessage[],
): ClientToolCall[] {
  const completedToolCallIDs = new Set<string>();
  const pendingParts: ClientToolCallPart[] = [];

  for (const message of messages) {
    for (const part of message.parts ?? []) {
      if (
        (part.type === "dynamic-tool" || part.type.startsWith("tool-")) &&
        "toolCallId" in part &&
        (part.state === "output-available" || part.state === "output-error")
      ) {
        completedToolCallIDs.add(part.toolCallId);
      } else if (
        part.type === "dynamic-tool" &&
        part.state === "input-available"
      ) {
        pendingParts.push(part);
      }
    }
  }

  const pendingToolCalls: ClientToolCall[] = [];
  for (const part of pendingParts) {
    if (!completedToolCallIDs.has(part.toolCallId)) {
      pendingToolCalls.push(clientToolCallFromPart(part));
    }
  }
  return pendingToolCalls;
}
