"use client";

import { Command } from "cmdk";
import { useEffect, useState } from "react";

import { robotSessionGet } from "@/api/openapi-client/robots";
import { StorydenUIMessage } from "@/api/robots-types";
import { WStack, styled } from "@/styled-system/jsx";

import { Unready, UnreadyBanner } from "../Unready";

import "./styles.css";

import { CommandPaletteContent } from "./CommandPaletteContent";
import { CommandPaletteProvider, useCommandPalette } from "./Context";
import { RobotChatContext } from "./RobotChat/RobotChatContext";
import { RobotCommandPaletteStatusBar } from "./RobotChat/RobotCommandPaletteStatusBar";
import { useChatSessionState } from "./RobotChat/useChatSessionState";

export function CommandPalette() {
  return (
    <CommandPaletteProvider>
      <CommandPaletteDialog />
    </CommandPaletteProvider>
  );
}

function CommandPaletteDialog() {
  const { open, dialogRef, initialSessionID } = useCommandPalette();
  const { loadingState, sessionState } = useChatSessionState(initialSessionID);

  return (
    <Command.Dialog
      ref={dialogRef}
      open={open}
      label="Command Menu"
      aria-description="The command palette allows you to quickly navigate and perform actions in Storyden."
    >
      <RobotChatContext
        key={sessionState.id}
        initialSessionID={sessionState.id}
        initialMessages={sessionState.messages}
      >
        {loadingState.isLoading || loadingState.error ? (
          <UnreadyBanner error={loadingState.error} />
        ) : (
          <CommandPaletteContent />
        )}

        <WStack fontSize="xs" color="fg.muted" lineHeight="tight" px="1" pt="2">
          <CommandPaletteStatusBar />
        </WStack>
      </RobotChatContext>
    </Command.Dialog>
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
          <styled.p>Storyden</styled.p>
          <styled.p>{mode}</styled.p>
        </>
      );
  }
}
