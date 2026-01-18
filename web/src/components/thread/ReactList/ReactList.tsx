import { Portal } from "@ark-ui/react";
import { EmojiClickData, EmojiStyle } from "emoji-picker-react";
import { motion } from "framer-motion";
import { throttle } from "lodash";
import dynamic from "next/dynamic";
import { useCallback, useEffect, useRef, useState } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { ReactionAddIcon } from "@/components/ui/icons/Reaction";
import * as Popover from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { useSettings } from "@/lib/settings/settings-client";
import { css } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { useDisclosure } from "@/utils/useDisclosure";

import {
  Props,
  REACTION_THROTTLE,
  ReactCount,
  useReactionList,
} from "./useReactList";

const EmojiPicker = dynamic(
  () => {
    return import("emoji-picker-react");
  },
  { ssr: false },
);

const reactButtonStyles = css({
  fontSize: "sm",
  height: "7",
  paddingX: "2",
});

const DEFAULT_QUICK_REACTIONS = ["❤️", "😂", "😮", "😢", "😠", "👍"];

/** Displays a list of reactions with counts and allows adding new reactions. */
export function ReactList(props: Props) {
  const { data, handlers } = useReactionList(props);
  const { settings } = useSettings();

  const { isLoggedIn, reacts } = data;
  const { handleReactExisting, handleReactPicker } = handlers;

  const quickReactions = settings?.quick_reactions ?? DEFAULT_QUICK_REACTIONS;

  return (
    <HStack flexWrap="wrap" gap="1">
      {reacts.map((react) => (
        <ReactTrigger
          key={react.emoji}
          react={react}
          disabled={!isLoggedIn}
          onClick={handleReactExisting}
        />
      ))}

      {isLoggedIn && (
        <ReactionPickerTrigger
          quickReactions={quickReactions}
          onSelect={handleReactPicker}
        />
      )}
    </HStack>
  );
}

type ReactionProps = {
  react: ReactCount;
  disabled: boolean;
  onClick: (emoji: string) => void;
};

/** Renders a single reaction button with animated count display. */
function ReactTrigger({ react, disabled, onClick }: ReactionProps) {
  const [count, setCount] = useState(react.count);
  const [direction, setDirection] = useState(1);
  const [hasMounted, setHasMounted] = useState(false);

  const hasReacted = useRef(false);
  useEffect(() => {
    hasReacted.current = react.hasReacted;
    setCount(react.count);
  }, [react]);

  useEffect(() => {
    setHasMounted(true);
  }, []);

  const handleAdd = () => {
    setDirection(1);
    setCount((prevCount) => prevCount + 1);
  };

  const handleRemove = () => {
    setDirection(-1);
    setCount((prevCount) => (prevCount > 0 ? prevCount - 1 : 0));
  };

  const handleClick = useCallback(
    throttle(() => {
      if (disabled) {
        return;
      }

      if (hasReacted.current) {
        handleRemove();
      } else {
        handleAdd();
      }
      onClick(react.emoji);
    }, REACTION_THROTTLE),
    [react, hasReacted.current],
  );

  if (count === 0) {
    return null;
  }

  const reacters = react.reactions.map((r) => r.author?.handle);
  const buttonLabel = `${react.emoji}: ${reacters.join(", ")}`;

  return (
    <Button
      className={reactButtonStyles}
      size="xs"
      variant="subtle"
      gap="0"
      borderRadius="md"
      color="fg.subtle"
      onClick={handleClick}
      fontVariantNumeric="tabular-nums"
      title={buttonLabel}
    >
      {react.emoji}
      <motion.span
        key={count}
        initial={
          hasMounted ? { y: direction === 1 ? 10 : -10, opacity: 0 } : false
        }
        animate={{ y: 0, opacity: 1 }}
        transition={{ type: "spring", stiffness: 300, damping: 20 }}
        style={{ marginLeft: "8px" }}
      >
        {count}
      </motion.span>
    </Button>
  );
}

type ReactionPickerTriggerProps = {
  quickReactions: string[];
  onSelect: (emoji: string) => void;
};

/** Opens a reaction picker using emoji-picker-react's built-in reactions mode. */
function ReactionPickerTrigger({
  quickReactions,
  onSelect,
}: ReactionPickerTriggerProps) {
  const { isOpen, onToggle, onClose } = useDisclosure();

  function handleSelect(e: EmojiClickData) {
    onSelect(e.emoji);
    onClose();
  }

  return (
    <Popover.Root
      lazyMount
      open={isOpen}
      positioning={{
        gutter: 12,
        overflowPadding: 12,
        fitViewport: true,
        placement: "bottom",
        flip: true,
      }}
      onInteractOutside={onClose}
      onEscapeKeyDown={onClose}
    >
      <Popover.Trigger
        type="button"
        cursor="pointer"
        onClick={onToggle}
        asChild
      >
        <IconButton
          size="xs"
          className={reactButtonStyles}
          variant="subtle"
          borderRadius="md"
          color="fg.muted"
          aria-label="Add reaction"
        >
          <ReactionAddIcon />
        </IconButton>
      </Popover.Trigger>

      <Portal>
        <Popover.Positioner minW="0">
          <Popover.Content
            padding="0"
            bgColor="transparent"
            border="none"
            boxShadow={"none" as any}
          >
            <EmojiPicker
              onEmojiClick={handleSelect}
              onReactionClick={handleSelect}
              emojiStyle={EmojiStyle.NATIVE}
              reactionsDefaultOpen
              reactions={quickReactions}
              allowExpandReactions
              previewConfig={{
                showPreview: false,
              }}
              lazyLoadEmojis
            />
          </Popover.Content>
        </Popover.Positioner>
      </Portal>
    </Popover.Root>
  );
}
