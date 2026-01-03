"use client";

import { Command, useCommandState } from "cmdk";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { datagraphMatches } from "@/api/openapi-client/datagraph";
import { DatagraphMatch } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { Spinner } from "@/components/ui/Spinner";
import { IconButton } from "@/components/ui/icon-button";
import { AdminIcon } from "@/components/ui/icons/Admin";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { HomeIcon } from "@/components/ui/icons/Home";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { SettingsIcon } from "@/components/ui/icons/Settings";
import { HStack, LStack } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

import "./styles.css";

import { CommandPaletteMode, useCommandPalette } from "./Context";
import { DatagraphSearchItem } from "./DatagraphSearch/DatagraphSearchItem";
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
          if (filteredCount === 0 || selectedValue === "/robot-chat") {
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
  const { mode, handleSelectItem, search, setOpen, loadSession } =
    useCommandPalette();
  const router = useRouter();
  const session = useSession();
  const [searchResults, setSearchResults] = useState<DatagraphMatch[]>([]);
  const [isSearching, setIsSearching] = useState(false);

  const isAdmin = session?.roles.some((role) => role.name === "admin");
  const canUseRobots = hasPermission(session, "USE_ROBOTS");

  const isCommandMode = search.startsWith("/");

  // Parse /chat {sessionId} command
  const chatCommand = search.match(/^\/chat\s+([a-z0-9]{20})$/i);
  const chatSessionId = chatCommand?.[1];

  useEffect(() => {
    if (mode !== "idle" || !search.trim() || isCommandMode) {
      setSearchResults([]);
      return;
    }

    const timeoutId = setTimeout(async () => {
      setIsSearching(true);
      try {
        const results = await datagraphMatches({ q: search });
        setSearchResults(results.items);
      } catch (error) {
        console.error("Search failed", error);
        setSearchResults([]);
      } finally {
        setIsSearching(false);
      }
    }, 250);

    return () => clearTimeout(timeoutId);
  }, [search, mode, isCommandMode]);

  const handleNavigate = (path: string) => {
    setOpen(false);
    router.push(path);
  };

  switch (mode) {
    case "chat":
      return null;

    default:
      return (
        <>
          <Command.Separator />
          <Command.List>
            <Command.Item
              value="/home"
              keywords={["home"]}
              onSelect={() => handleNavigate("/")}
            >
              <HStack gap="2">
                <HomeIcon />
                Home
              </HStack>
            </Command.Item>
            <Command.Item
              value="/settings"
              keywords={["settings"]}
              onSelect={() => handleNavigate("/settings")}
            >
              <HStack gap="2">
                <SettingsIcon />
                Settings
              </HStack>
            </Command.Item>
            {isAdmin && (
              <Command.Item
                value="/admin"
                keywords={["admin"]}
                onSelect={() => handleNavigate("/admin")}
              >
                <HStack gap="2">
                  <AdminIcon />
                  Admin
                </HStack>
              </Command.Item>
            )}
            {canUseRobots && (
              <Command.Item
                value="/robot-chat"
                keywords={["robot", "chat", "ai"]}
                onSelect={handleSelectItem}
              >
                <HStack gap="2">
                  <RobotIcon />
                  Robot Chat
                </HStack>
              </Command.Item>
            )}
            {canUseRobots && chatSessionId && (
              <Command.Item
                value={search}
                onSelect={() => loadSession?.(chatSessionId)}
              >
                <HStack gap="2">
                  <RobotIcon />
                  Resume chat session: {chatSessionId}
                </HStack>
              </Command.Item>
            )}

            {!isCommandMode &&
              searchResults.map((result) => {
                return (
                  <DatagraphSearchItem
                    key={result.id}
                    result={result}
                    handleNavigate={handleNavigate}
                  />
                );
              })}

            {!isCommandMode && isSearching && (
              <Command.Loading>
                <Spinner />
              </Command.Loading>
            )}
          </Command.List>
        </>
      );
  }
}
