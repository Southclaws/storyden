import { Portal } from "@ark-ui/react";
import { EmojiClickData, EmojiStyle } from "emoji-picker-react";
import { motion } from "framer-motion";
import { throttle } from "lodash";
import { SmilePlusIcon } from "lucide-react";
import dynamic from "next/dynamic";
import { useCallback, useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import * as Popover from "@/components/ui/popover";
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

export function ReactList(props: Props) {
  const { data, handlers } = useReactionList(props);

  const { reacts } = data;
  const { handleReactExisting, handleReactPicker } = handlers;

  return (
    <HStack flexWrap="wrap">
      {reacts.map((react) => (
        <ReactTrigger
          key={react.emoji}
          react={react}
          onClick={handleReactExisting}
        />
      ))}

      <ReactionPickerTrigger onSelect={handleReactPicker} />
    </HStack>
  );
}

type ReactionProps = {
  react: ReactCount;
  onClick: (emoji: string) => void;
};

function ReactTrigger({ react, onClick }: ReactionProps) {
  const [count, setCount] = useState(react.count);
  const [direction, setDirection] = useState(1); // To track up or down animation
  const [hasMounted, setHasMounted] = useState(false);

  // Ref the has-reacted state, in order for the debounce to not capture value.
  const hasReacted = useRef(false);
  useEffect(() => {
    hasReacted.current = react.hasReacted;
    setCount(react.count);
  }, [react]);

  // Prevents the animation from playing on first render. Unfortunately also has
  // the unintended effect of not playing the animation on the first reaction.
  useEffect(() => {
    setHasMounted(true);
  }, []);

  const handleAdd = () => {
    setDirection(1); // Set direction upwards for increase
    setCount((prevCount) => prevCount + 1);
  };

  const handleRemove = () => {
    setDirection(-1); // Set direction downwards for decrease
    setCount((prevCount) => (prevCount > 0 ? prevCount - 1 : 0));
  };

  //
  // Actual reaction events are client-side rate limited here for two reasons:
  //
  // 1. To prevent spamming the server with requests (though the server does
  //    implement its own rate-limiting, this is an additional layer).
  //
  // 2. To allow for the revalidation to happen after a short delay. The reason
  //    for this is that revalidation triggers the entire thread to re-render,
  //    which, while mostly holding identical component state aside from the
  //    ReactList, it causes a noticeable frame rate drop for the motion.span
  //    below. To solve this, revalidation is delayed by `REACTION_THROTTLE`
  //    milliseconds and this also means no further mutations can happen until
  //    that time is up (at least not without a lot more complexity.)
  //    This throttling is only applied at the reaction level, so the user can
  //    still trigger other reactions quickly, it just causes more revalidation.
  //
  const handleClick = useCallback(
    throttle(() => {
      if (hasReacted.current) {
        handleRemove();
      } else {
        handleAdd();
      }
      onClick(react.emoji);
    }, REACTION_THROTTLE),
    [hasReacted],
  );

  // When removing an emoji completely (its count has reached zero) we need to
  // remove it, but there's a short period between the interaction and the
  // mutation where the ReactTrigger is still rendered. In this case, we don't
  // want to trigger the animation as it looks strange since on fast network
  // connections only a few frames of animation play before the revalidation
  // kicks in and re-renders the list without the component. This branch below
  // ensures that in this case, the ReactTrigger component is simply removed.
  if (count === 0) {
    return null;
  }

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
  onSelect: (emoji: string) => void;
};

function ReactionPickerTrigger(props: ReactionPickerTriggerProps) {
  const { isOpen, onToggle, onClose } = useDisclosure();

  function handleSelect(e: EmojiClickData) {
    props.onSelect(e.emoji);
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
        >
          <SmilePlusIcon />
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
