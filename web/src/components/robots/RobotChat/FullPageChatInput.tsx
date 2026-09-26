"use client";

import { useRef, useState } from "react";

import { useRobotChat } from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { RobotChatLoadingStatus } from "@/components/site/CommandPalette/RobotChat/RobotChatLoadingStatus";
import { LStack, styled } from "@/styled-system/jsx";
import { pluralise } from "@/utils/text";

import { RobotChatComposer } from "./RobotChatComposer";
import { useRobotChatImageAttachments } from "./useRobotChatImageAttachments";

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
  const {
    attachments,
    assets,
    isDragging,
    isUploading,
    addFiles,
    handlePaste,
    handleDragEnter,
    handleDragLeave,
    handleDragOver,
    handleDrop,
    removeAttachment,
    clearAttachments,
    restoreAssets,
  } = useRobotChatImageAttachments();

  const isBusy = status === "submitted" || status === "streaming";
  const canSend =
    (input.trim().length > 0 || assets.length > 0) && !isUploading;

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!canSend) return;

    const text = input.trim();
    const submittedAssets = assets;
    setInput("");
    clearAttachments();

    try {
      await sendMessage({ text, assets: submittedAssets });
      textareaRef.current?.focus();
    } catch (err) {
      console.error("sendMessage failed", err);
      setInput((current) => current || text);
      restoreAssets(submittedAssets);
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
      onDragEnter={handleDragEnter}
      onDragLeave={handleDragLeave}
      onDragOver={handleDragOver}
      onDrop={handleDrop}
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
        <RobotChatComposer
          value={input}
          textareaRef={textareaRef}
          attachments={attachments}
          isDragging={isDragging}
          canSend={canSend}
          canCancel={canCancelActiveTurn}
          isCancelling={isCancelling}
          onChange={(event) => setInput(event.target.value)}
          onKeyDown={(event) => void handleKeyDown(event)}
          onPaste={handlePaste}
          onFiles={addFiles}
          onRemoveAttachment={removeAttachment}
          onCancel={() => void cancelActiveTurn()}
        />
      </LStack>
    </styled.form>
  );
}
