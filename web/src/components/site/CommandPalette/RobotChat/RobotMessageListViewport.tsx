"use client";

import {
  type RefObject,
  useCallback,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
} from "react";

import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { Box, VStack } from "@/styled-system/jsx";

import { useRobotChat } from "./RobotChatContext";
import { RobotChatMessageProjectedList } from "./RobotChatMessageProjectedList";

const TOP_LOAD_THRESHOLD = 96;
const BOTTOM_STICK_THRESHOLD = 120;

type MessageAnchor = {
  id: string;
  offsetTop: number;
};

type Props = {
  surface: "full-page" | "palette";
};

export function RobotMessageListViewport({ surface }: Props) {
  const {
    messages,
    errorState,
    handleDismissError,
    hasOlderMessages,
    isLoadingOlderMessages,
    loadOlderMessages,
  } = useRobotChat();
  const scrollerRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);
  const pendingAnchorRef = useRef<MessageAnchor | undefined>(undefined);
  const suppressNextNewMessagePillRef = useRef(false);
  const wasNearBottomRef = useRef(true);
  const previousLastMessageIDRef = useRef<string | undefined>(undefined);
  const [showNewMessages, setShowNewMessages] = useState(false);

  const isFullPage = surface === "full-page";

  const scrollToBottom = useCallback((behavior: ScrollBehavior = "smooth") => {
    const scroller = scrollerRef.current;
    if (!scroller) return;

    scroller.scrollTo({
      top: scroller.scrollHeight,
      behavior,
    });
  }, []);

  const updateNearBottom = useCallback(() => {
    const scroller = scrollerRef.current;
    if (!scroller) return true;

    const distanceFromBottom =
      scroller.scrollHeight - scroller.scrollTop - scroller.clientHeight;
    const isNearBottom = distanceFromBottom < BOTTOM_STICK_THRESHOLD;
    wasNearBottomRef.current = isNearBottom;

    if (isNearBottom) {
      setShowNewMessages(false);
    }

    return isNearBottom;
  }, []);

  const handleLoadOlder = useCallback(async () => {
    const scroller = scrollerRef.current;
    if (!scroller || !hasOlderMessages || isLoadingOlderMessages) return;

    pendingAnchorRef.current = getFirstVisibleMessageAnchor(scroller);
    suppressNextNewMessagePillRef.current = true;
    const loaded = await loadOlderMessages();

    if (!loaded) {
      pendingAnchorRef.current = undefined;
      suppressNextNewMessagePillRef.current = false;
    }
  }, [hasOlderMessages, isLoadingOlderMessages, loadOlderMessages]);

  const handleScroll = useCallback(() => {
    const scroller = scrollerRef.current;
    if (!scroller) return;

    updateNearBottom();

    if (scroller.scrollTop < TOP_LOAD_THRESHOLD) {
      void handleLoadOlder();
    }
  }, [handleLoadOlder, updateNearBottom]);

  useLayoutEffect(() => {
    const pendingAnchor = pendingAnchorRef.current;
    const scroller = scrollerRef.current;

    if (!pendingAnchor || !scroller) return;

    const anchoredElement = scroller.querySelector<HTMLElement>(
      `#robot-message-${pendingAnchor.id}`,
    );
    if (!anchoredElement) {
      pendingAnchorRef.current = undefined;
      return;
    }

    const scrollerTop = scroller.getBoundingClientRect().top;
    const anchoredTop = anchoredElement.getBoundingClientRect().top;
    scroller.scrollTop += anchoredTop - scrollerTop - pendingAnchor.offsetTop;
    pendingAnchorRef.current = undefined;
  }, [messages]);

  useEffect(() => {
    const content = contentRef.current;
    if (!content) return;

    const observer = new ResizeObserver(() => {
      if (pendingAnchorRef.current || !wasNearBottomRef.current) {
        return;
      }

      scrollToBottom("auto");
    });

    observer.observe(content);

    return () => observer.disconnect();
  }, [scrollToBottom]);

  useEffect(() => {
    const previousLastMessageID = previousLastMessageIDRef.current;
    const currentLastMessageID = messages.at(-1)?.id;
    previousLastMessageIDRef.current = currentLastMessageID;

    if (!currentLastMessageID) {
      return;
    }

    if (!previousLastMessageID) {
      scrollToBottom("auto");
      return;
    }

    if (wasNearBottomRef.current) {
      scrollToBottom();
      return;
    }

    if (suppressNextNewMessagePillRef.current) {
      suppressNextNewMessagePillRef.current = false;
      return;
    }

    setShowNewMessages(true);
  }, [messages, scrollToBottom]);

  if (!isFullPage && messages.length === 0) {
    return null;
  }

  if (isFullPage) {
    return (
      <Box
        ref={scrollerRef}
        onScroll={handleScroll}
        role="log"
        aria-label="Robot chat messages"
        aria-live="polite"
        aria-relevant="additions text"
        w="full"
        flex="1"
        minH="0"
        minW="0"
        gap="3"
        overflowY="auto"
        position="relative"
      >
        <Box
          position="sticky"
          top="0"
          left="0"
          right="0"
          h="8"
          pointerEvents="none"
          zIndex="dropdown"
          bgColor="scroll-fade-top"
        />
        {hasOlderMessages && (
          <Box display="flex" justifyContent="center" my="1">
            <Button
              type="button"
              size="xs"
              variant="ghost"
              loading={isLoadingOlderMessages}
              onClick={() => void handleLoadOlder()}
            >
              Load older
            </Button>
          </Box>
        )}
        <MessageContent
          contentRef={contentRef}
          bottomRef={bottomRef}
          errorState={errorState}
          handleDismissError={handleDismissError}
          my="4"
        />
        <NewMessagesButton
          visible={showNewMessages}
          position="sticky"
          onClick={() => {
            setShowNewMessages(false);
            scrollToBottom();
          }}
        />
      </Box>
    );
  }

  return (
    <Box position="relative" w="full" minW="0" overflowX="hidden">
      <VStack
        ref={scrollerRef}
        onScroll={handleScroll}
        role="log"
        aria-label="Robot chat messages"
        aria-live="polite"
        aria-relevant="additions text"
        w="full"
        minW="0"
        gap="3"
        p="4"
        maxH="96"
        overflowY="auto"
        overflowX="hidden"
        alignItems="stretch"
      >
        <MessageContent
          contentRef={contentRef}
          bottomRef={bottomRef}
          errorState={errorState}
          handleDismissError={handleDismissError}
        />
      </VStack>
      <NewMessagesButton
        visible={showNewMessages}
        position="absolute"
        onClick={() => {
          setShowNewMessages(false);
          scrollToBottom();
        }}
      />
    </Box>
  );
}

type MessageContentProps = {
  contentRef: RefObject<HTMLDivElement | null>;
  bottomRef: RefObject<HTMLDivElement | null>;
  errorState?: string;
  handleDismissError: () => void;
  my?: "4";
};

function MessageContent({
  contentRef,
  bottomRef,
  errorState,
  handleDismissError,
  my,
}: MessageContentProps) {
  return (
    <VStack ref={contentRef} alignItems="stretch" w="full" minW="0" my={my}>
      <RobotChatMessageProjectedList />
      <Admonition value={Boolean(errorState)} onChange={handleDismissError}>
        <p>{errorState}</p>
      </Admonition>
      <div ref={bottomRef} />
    </VStack>
  );
}

type NewMessagesButtonProps = {
  visible: boolean;
  position: "absolute" | "sticky";
  onClick: () => void;
};

function NewMessagesButton({
  visible,
  position,
  onClick,
}: NewMessagesButtonProps) {
  if (!visible) {
    return null;
  }

  return (
    <Box
      position={position}
      bottom="3"
      left="0"
      right="0"
      display="flex"
      justifyContent="center"
      pointerEvents="none"
      zIndex="dropdown"
    >
      <Button
        type="button"
        size="xs"
        variant="subtle"
        pointerEvents="auto"
        onClick={onClick}
      >
        New messages
      </Button>
    </Box>
  );
}

function getFirstVisibleMessageAnchor(
  scroller: HTMLDivElement,
): MessageAnchor | undefined {
  const scrollerTop = scroller.getBoundingClientRect().top;
  const messages = scroller.querySelectorAll<HTMLElement>(
    '[id^="robot-message-"]',
  );

  for (const message of messages) {
    const rect = message.getBoundingClientRect();

    if (rect.bottom <= scrollerTop) {
      continue;
    }

    return {
      id: message.id.replace("robot-message-", ""),
      offsetTop: rect.top - scrollerTop,
    };
  }

  return undefined;
}
