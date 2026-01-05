"use client";

import { Command, useCommandState } from "cmdk";

import { IconButton } from "@/components/ui/icon-button";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { HStack, LStack } from "@/styled-system/jsx";

import "./styles.css";

import { CommandPaletteMode, useCommandPalette } from "./Context";
import { useRobotChat } from "./RobotChat/RobotChatContext";
import { RobotChatMessageList } from "./RobotChat/RobotChatMessageList";

export function CommandPaletteContent() {
  const { search, setSearch, mode, setMode, focusInput } = useCommandPalette();
  const { sendMessage, status } = useRobotChat();
  const filteredCount = useCommandState((s) => s.filtered.count);
  const selectedValue = useCommandState((s) => s.value);

  const isBusy = status === "submitted" || status === "streaming";

  async function handleChatSend() {
    if (!search.trim() || isBusy) return;

    if (mode !== "chat") {
      setMode("chat");
    }

    const originalSearch = search;
    setSearch("");

    try {
      focusInput();
      await sendMessage({ text: search.trim() });
      focusInput();
    } catch (err) {
      // TODO: Handle this by presenting the user-friendly error in the UI.
      console.error("sendMessage failed", err);
      setSearch(originalSearch);
    }
  }

  async function handleSubmissionIntent(e: React.KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) {
      switch (mode) {
        case "idle":
          // Smart Enter behavior:
          // - No results? Start chat with the typed text
          // - "Robot Chat" selected? Start chat with the typed text
          // - Other command selected? Let CMDK handle it (do nothing)
          if (filteredCount === 0 || selectedValue === "robot-chat") {
            e.preventDefault();
            await handleChatSend();
          }
          break;
        case "chat":
          await handleChatSend();
          break;
      }
    }
  }

  return (
    <LStack gap="0">
      <CommandPaletteContentMode mode={mode} />

      <HStack w="full">
        <Command.Input
          value={search}
          onValueChange={setSearch}
          placeholder="Ask, search or command..."
          disabled={isBusy}
          onKeyDown={handleSubmissionIntent}
        />
        <IconButton
          variant="subtle"
          type="submit"
          disabled={isBusy}
          loading={isBusy}
        >
          <DiscussionIcon />
        </IconButton>
      </HStack>

      <CommandPaletteContentCommands />
    </LStack>
  );
}

function CommandPaletteContentMode({ mode }: { mode: CommandPaletteMode }) {
  switch (mode) {
    case "chat":
      return <RobotChatMessageList />;

    default:
      return null;
  }
}

function CommandPaletteContentCommands() {
  const { mode, handleSelectItem } = useCommandPalette();

  switch (mode) {
    case "chat":
      return null;

    default:
      return (
        <>
          <Command.Separator />
          <Command.List>
            <Command.Item value="robot-chat" onSelect={handleSelectItem}>
              Robot Chat
            </Command.Item>
            <Command.Item value="another" onSelect={handleSelectItem}>
              Another Item
            </Command.Item>
          </Command.List>
        </>
      );
  }
}
