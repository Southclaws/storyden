import { Portal } from "@ark-ui/react";
import { EmojiClickData, EmojiStyle } from "emoji-picker-react";
import { LazyMotion, domAnimation, m } from "framer-motion";
import throttle from "lodash/throttle";
import dynamic from "next/dynamic";
import { useEffect, useMemo, useState, useSyncExternalStore } from "react";

import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { ReactionAddIcon } from "@/components/ui/icons/Reaction";
import * as Popover from "@/components/ui/popover";
import { reactList } from "@/styled-system/recipes";
import { useDisclosure } from "@/utils/useDisclosure";
import { useHydrated } from "@/utils/useHydrated";

const REACTION_THROTTLE = 180;

const EmojiPicker = dynamic(
  () => {
    return import("emoji-picker-react");
  },
  { ssr: false },
);

export type ReactListItem = {
  emoji: string;
  count: number;
  hasReacted: boolean;
  reactedBy: string[];
};

export type ReactListProps = {
  isLoggedIn: boolean;
  pendingEmojis?: ReadonlySet<string>;
  reactions: ReactListItem[];
  onReactExisting: (emoji: string) => void | Promise<void>;
  onReactPicker: (emoji: string) => void | Promise<void>;
};

export function ReactList({
  isLoggedIn,
  pendingEmojis = new Set(),
  reactions,
  onReactExisting,
  onReactPicker,
}: ReactListProps) {
  const styles = reactList();
  const isHydrated = useHydrated();

  return (
    <div className={styles.root}>
      {reactions.map((reaction) => {
        const pending = pendingEmojis.has(reaction.emoji);

        return (
          <ReactTrigger
            key={reaction.emoji}
            reaction={reaction}
            disabled={!isLoggedIn || pending}
            pending={pending}
            isHydrated={isHydrated}
            onClick={onReactExisting}
          />
        );
      })}

      {isLoggedIn && (
        <ReactionPickerTrigger
          className={styles.picker}
          onSelect={onReactPicker}
        />
      )}
    </div>
  );
}

type ReactionProps = {
  reaction: ReactListItem;
  disabled: boolean;
  pending: boolean;
  isHydrated: boolean;
  onClick: (emoji: string) => void | Promise<void>;
};

function ReactTrigger({
  reaction,
  disabled,
  pending,
  isHydrated,
  onClick,
}: ReactionProps) {
  const styles = reactList();
  const [count, setCount] = useState(reaction.count);
  const [direction, setDirection] = useState(1);

  useEffect(() => {
    setCount(reaction.count);
  }, [reaction.count]);

  const handleClick = useMemo(
    () =>
      throttle(() => {
        if (disabled) {
          return;
        }

        if (reaction.hasReacted) {
          setDirection(-1);
          setCount((previous) => Math.max(0, previous - 1));
        } else {
          setDirection(1);
          setCount((previous) => previous + 1);
        }
        void onClick(reaction.emoji);
      }, REACTION_THROTTLE),
    [disabled, onClick, reaction.emoji, reaction.hasReacted],
  );

  useEffect(() => () => handleClick.cancel(), [handleClick]);

  if (count === 0) {
    return null;
  }

  const buttonLabel = `${reaction.emoji}: ${reaction.reactedBy.join(", ")}`;

  return (
    <Button
      className={styles.reaction}
      variant="subtle"
      disabled={disabled}
      aria-busy={pending || undefined}
      onClick={handleClick}
      title={buttonLabel}
    >
      {reaction.emoji}
      <LazyMotion features={domAnimation}>
        <m.span
          key={count}
          className={styles.count}
          initial={
            isHydrated ? { y: direction === 1 ? 10 : -10, opacity: 0 } : false
          }
          animate={{ y: 0, opacity: 1 }}
          transition={{ type: "spring", stiffness: 300, damping: 20 }}
        >
          {count}
        </m.span>
      </LazyMotion>
    </Button>
  );
}

type ReactionPickerTriggerProps = {
  className: string;
  onSelect: (emoji: string) => void | Promise<void>;
};

function ReactionPickerTrigger({
  className,
  onSelect,
}: ReactionPickerTriggerProps) {
  const { isOpen, onToggle, onClose } = useDisclosure();

  function handleSelect(event: EmojiClickData) {
    void onSelect(event.emoji);
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
          className={className}
          variant="subtle"
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
            <Popover.Arrow>
              <Popover.ArrowTip />
            </Popover.Arrow>

            <EmojiPicker
              onEmojiClick={handleSelect}
              emojiStyle={EmojiStyle.NATIVE}
              searchDisabled
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
