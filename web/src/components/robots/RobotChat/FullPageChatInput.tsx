"use client";

import { useRef, useState } from "react";

import { useRobotChat } from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { RobotChatLoadingStatus } from "@/components/site/CommandPalette/RobotChat/RobotChatLoadingStatus";
import { IconButton } from "@/components/ui/icon-button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { Textarea } from "@/components/ui/textarea";
import { HStack, LStack, styled } from "@/styled-system/jsx";
import { pluralise } from "@/utils/text";

export function FullPageChatInput() {
  const {
    activeRobotName,
    sendMessage,
    cancelActiveTurn,
    canCancelActiveTurn,
    isCancelling,
    status,
    queuedMessageCount,
  } = useRobotChat();
  const [input, setInput] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const isBusy = status === "submitted" || status === "streaming";

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!input.trim()) return;

    const text = input.trim();
    setInput("");

    try {
      await sendMessage({ text });
      // Refocus the textarea after sending
      textareaRef.current?.focus();
    } catch (err) {
      console.error("sendMessage failed", err);
      setInput((current) => current || text);
      // Also refocus on error
      textareaRef.current?.focus();
    }
  }

  async function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      await handleSubmit(e);
    }
  }

  return (
    <styled.form
      onSubmit={handleSubmit}
      w="full"
      flexShrink="0"
      aria-label="Send message to Robot"
    >
      <LStack w="full" gap="1.5">
        <RobotChatLoadingStatus active={isBusy} robotName={activeRobotName} />
        {queuedMessageCount > 0 && (
          <styled.span color="text.muted" fontSize="xs">
            {queuedMessageCount} {pluralise(queuedMessageCount, "message")}{" "}
            queued
          </styled.span>
        )}
        <HStack w="full" gap="2">
          <Textarea
            ref={textareaRef}
            aria-label="Message"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Type a message..."
            rows={1}
            px="3"
            py="2"
            borderRadius="md"
            resize="none"
          />
          <IconButton
            aria-label="Send message"
            variant="subtle"
            type="submit"
            disabled={!input.trim()}
          >
            <DiscussionIcon />
          </IconButton>
          {canCancelActiveTurn && (
            <IconButton
              aria-label="Cancel Robot response"
              variant="subtle"
              type="button"
              loading={isCancelling}
              onClick={() => void cancelActiveTurn()}
            >
              <CancelIcon />
            </IconButton>
          )}
        </HStack>
      </LStack>
    </styled.form>
  );
}
