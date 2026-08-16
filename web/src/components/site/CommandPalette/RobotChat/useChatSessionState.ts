import { useEffect, useState } from "react";

import { robotSessionGet } from "@/api/openapi-client/robots";
import { StorydenUIMessage, toStorydenUIMessages } from "@/api/robots-types";

export function useChatSessionState(initialSessionID?: string) {
  const [sessionState, setSessionState] = useState<{
    id: string | undefined;
    activeWorkspaceID: string | undefined;
    messages: StorydenUIMessage[] | undefined;
    nextBefore: string | undefined;
    streamOffset: string | undefined;
  }>({
    id: undefined,
    activeWorkspaceID: undefined,
    messages: undefined,
    nextBefore: undefined,
    streamOffset: undefined,
  });

  const [loadingState, setLoadingState] = useState<{
    isLoading: boolean;
    error: unknown;
  }>({
    isLoading: false,
    error: undefined,
  });

  useEffect(() => {
    if (!initialSessionID) {
      setSessionState({
        id: undefined,
        activeWorkspaceID: undefined,
        messages: undefined,
        nextBefore: undefined,
        streamOffset: undefined,
      });
      setLoadingState({ isLoading: false, error: undefined });
      return;
    }

    async function loadSessionData(id: string) {
      setLoadingState({ isLoading: true, error: undefined });
      try {
        const session = await robotSessionGet(id);
        const messages = toStorydenUIMessages(session.message_list.messages);
        setSessionState({
          id,
          activeWorkspaceID: session.active_workspace?.workspace_id,
          messages,
          nextBefore: session.message_list.next_before,
          streamOffset: session.stream_offset,
        });
        setLoadingState({ isLoading: false, error: undefined });
      } catch (error) {
        setLoadingState({ isLoading: false, error });
        setSessionState({
          id: undefined,
          activeWorkspaceID: undefined,
          messages: undefined,
          nextBefore: undefined,
          streamOffset: undefined,
        });
      }
    }

    loadSessionData(initialSessionID);
  }, [initialSessionID]);

  return {
    sessionState,
    loadingState,
  };
}
