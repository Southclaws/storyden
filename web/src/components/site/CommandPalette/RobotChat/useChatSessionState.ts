import { useEffect, useState } from "react";

import { robotSessionGet } from "@/api/openapi-client/robots";
import { StorydenUIMessage, toStorydenUIMessages } from "@/api/robots-types";

export function useChatSessionState(initialSessionID?: string) {
  const [sessionState, setSessionState] = useState<{
    id: string | undefined;
    activeRobotID: string | undefined;
    activeWorkspaceID: string | undefined;
    messages: StorydenUIMessage[] | undefined;
    nextBefore: string | undefined;
  }>({
    id: undefined,
    activeRobotID: undefined,
    activeWorkspaceID: undefined,
    messages: undefined,
    nextBefore: undefined,
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
        activeRobotID: undefined,
        activeWorkspaceID: undefined,
        messages: undefined,
        nextBefore: undefined,
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
          activeRobotID: session.active_robot_id,
          activeWorkspaceID: session.active_workspace?.workspace_id,
          messages,
          nextBefore: session.message_list.next_before,
        });
        setLoadingState({ isLoading: false, error: undefined });
      } catch (error) {
        setLoadingState({ isLoading: false, error });
        setSessionState({
          id: undefined,
          activeRobotID: undefined,
          activeWorkspaceID: undefined,
          messages: undefined,
          nextBefore: undefined,
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
