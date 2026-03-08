import { useEffect, useState } from "react";

import { robotSessionGet } from "@/api/openapi-client/robots";
import { StorydenUIMessage } from "@/api/robots-types";

export function useChatSessionState(initialSessionID?: string) {
  const [sessionState, setSessionState] = useState<{
    id: string | undefined;
    messages: StorydenUIMessage[] | undefined;
  }>({
    id: undefined,
    messages: undefined,
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
      setSessionState({ id: undefined, messages: undefined });
      setLoadingState({ isLoading: false, error: undefined });
      return;
    }

    async function loadSessionData(id: string) {
      setLoadingState({ isLoading: true, error: undefined });
      try {
        const session = await robotSessionGet(id);
        const messages = session.message_list.messages;
        setSessionState({ id, messages });
        setLoadingState({ isLoading: false, error: undefined });
      } catch (error) {
        setLoadingState({ isLoading: false, error });
        setSessionState({ id: undefined, messages: undefined });
      }
    }

    loadSessionData(initialSessionID);
  }, [initialSessionID]);

  return {
    sessionState,
    loadingState,
  };
}
