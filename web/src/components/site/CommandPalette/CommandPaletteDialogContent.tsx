"use client";

import { useEffect } from "react";

import { WStack } from "@/styled-system/jsx";

import { UnreadyBanner } from "../Unready";

import { CommandPaletteContent } from "./CommandPaletteContent";
import { useCommandPalette } from "./Context";
import { RobotChatContext } from "./RobotChat/RobotChatContext";
import { RobotCommandPaletteStatusBar } from "./RobotChat/RobotCommandPaletteStatusBar";
import { useChatSessionState } from "./RobotChat/useChatSessionState";

export function CommandPaletteDialogContent() {
  const { initialSessionID, dialogRef } = useCommandPalette();
  const { loadingState, sessionState } = useChatSessionState(initialSessionID);

  useEffect(() => {
    dialogRef.current?.querySelector<HTMLInputElement>("[cmdk-input]")?.focus();
  }, [dialogRef, loadingState.isLoading, loadingState.error]);

  return (
    <RobotChatContext
      key={sessionState.id}
      initialSessionID={sessionState.id}
      initialActiveTurnID={sessionState.activeTurnID}
      initialSelectedWorkspaceID={sessionState.activeWorkspaceID}
      initialMessages={sessionState.messages}
      initialNextBefore={sessionState.nextBefore}
      initialStreamOffset={sessionState.streamOffset}
    >
      {loadingState.isLoading || loadingState.error ? (
        <UnreadyBanner error={loadingState.error} />
      ) : (
        <CommandPaletteContent />
      )}

      <WStack
        fontSize="xs"
        color="text.subtle"
        lineHeight="tight"
        px="1"
        pt="2"
      >
        <CommandPaletteStatusBar />
      </WStack>
    </RobotChatContext>
  );
}

function CommandPaletteStatusBar() {
  const { mode } = useCommandPalette();

  switch (mode) {
    case "chat":
      return <RobotCommandPaletteStatusBar />;

    default:
      return (
        <>
          <p>Storyden</p>
          <p>{mode}</p>
        </>
      );
  }
}
