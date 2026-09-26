"use client";

import {
  type ChangeEventHandler,
  type ClipboardEventHandler,
  type KeyboardEventHandler,
  type RefObject,
  useRef,
} from "react";

import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { Textarea } from "@/components/ui/textarea";
import { HStack, LStack, styled } from "@/styled-system/jsx";

import { RobotChatAttachmentShelf } from "./RobotChatAttachmentShelf";
import { RobotChatImageAttachment } from "./useRobotChatImageAttachments";

type Props = {
  value: string;
  textareaRef: RefObject<HTMLTextAreaElement | null>;
  attachments: RobotChatImageAttachment[];
  isDragging: boolean;
  canSend: boolean;
  canCancel: boolean;
  isCancelling: boolean;
  onChange: ChangeEventHandler<HTMLTextAreaElement>;
  onKeyDown: KeyboardEventHandler<HTMLTextAreaElement>;
  onPaste: ClipboardEventHandler<HTMLTextAreaElement>;
  onFiles: (files: FileList) => void;
  onRemoveAttachment: (id: string) => void;
  onCancel: () => void;
};

export function RobotChatComposer({
  value,
  textareaRef,
  attachments,
  isDragging,
  canSend,
  canCancel,
  isCancelling,
  onChange,
  onKeyDown,
  onPaste,
  onFiles,
  onRemoveAttachment,
  onCancel,
}: Props) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  return (
    <HStack w="full" minW="0" gap="2" alignItems="end">
      <LStack
        w="auto"
        minW="0"
        flex="1"
        gap="0"
        background="background.control"
        borderWidth="thin"
        borderColor={isDragging ? "accent.solid" : "border.default"}
        borderRadius="md"
      >
        <RobotChatAttachmentShelf
          attachments={attachments}
          onRemove={onRemoveAttachment}
        />
        <HStack w="full" minW="0" gap="0" pl="1" alignItems="center">
          <styled.input
            ref={fileInputRef}
            type="file"
            accept="image/gif,image/jpeg,image/png,image/webp"
            multiple
            display="none"
            onChange={(event) => {
              if (event.target.files) {
                onFiles(event.target.files);
              }
              event.target.value = "";
            }}
          />
          <IconButton
            aria-label="Attach images"
            variant="ghost"
            type="button"
            onClick={() => fileInputRef.current?.click()}
          >
            <AddIcon />
          </IconButton>
          <Textarea
            ref={textareaRef}
            aria-label="Message"
            value={value}
            onChange={onChange}
            onKeyDown={onKeyDown}
            onPaste={onPaste}
            placeholder={
              isDragging ? "Drop images to attach" : "Type a message..."
            }
            rows={1}
            variant="ghost"
            w="auto"
            minW="0"
            flex="1"
            pl="1"
            pr="3"
            py="2"
            borderWidth="none"
            borderRadius="md"
            resize="none"
            background="transparent"
          />
        </HStack>
      </LStack>
      <IconButton
        aria-label="Send message"
        variant="subtle"
        type="submit"
        disabled={!canSend}
      >
        <DiscussionIcon />
      </IconButton>
      {canCancel && (
        <IconButton
          aria-label="Cancel Robot response"
          variant="subtle"
          type="button"
          loading={isCancelling}
          onClick={onCancel}
        >
          <CancelIcon />
        </IconButton>
      )}
    </HStack>
  );
}
