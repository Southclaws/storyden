"use client";

import { useRef, useState } from "react";

import { useRobotChat } from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { RobotChatLoadingStatus } from "@/components/site/CommandPalette/RobotChat/RobotChatLoadingStatus";
import { IconButton } from "@/components/ui/icon-button";
import { CloseIcon } from "@/components/ui/icons/Close";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { HStack, LStack, styled } from "@/styled-system/jsx";

export function FullPageChatInput() {
  const { activeRobotName, sendMessage, status, stopGenerating } =
    useRobotChat();
  const [input, setInput] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const isBusy = status === "submitted" || status === "streaming";

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!input.trim() || isBusy) return;

    const text = input.trim();
    setInput("");

    try {
      await sendMessage({ text });
      // Refocus the textarea after sending
      textareaRef.current?.focus();
    } catch (err) {
      console.error("sendMessage failed", err);
      setInput(text);
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
      aria-label="Send message to Robot"
    >
      <LStack w="full" gap="1.5">
        <RobotChatLoadingStatus active={isBusy} robotName={activeRobotName} />
        <HStack w="full" gap="2">
          <styled.textarea
            ref={textareaRef}
            aria-label="Message"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Type a message..."
            rows={1}
            w="full"
            px="3"
            py="2"
            borderRadius="md"
            borderWidth="thin"
            borderColor="border.default"
            bg="bg.default"
            color="fg.default"
            fontSize="sm"
            resize="none"
            _focus={{
              borderColor: "accent.default",
              outline: "none",
            }}
            _disabled={{
              cursor: "not-allowed",
            }}
            disabled={isBusy}
            style={
              isBusy
                ? {
                    opacity: 0.5,
                  }
                : undefined
            }
          />
          <IconButton
            aria-label={isBusy ? "Cancel Robot response" : "Send message"}
            variant="subtle"
            type={isBusy ? "button" : "submit"}
            disabled={!isBusy && !input.trim()}
            onClick={isBusy ? () => void stopGenerating() : undefined}
          >
            {isBusy ? <CloseIcon /> : <DiscussionIcon />}
          </IconButton>
        </HStack>
      </LStack>
    </styled.form>
  );
}
