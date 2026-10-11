import {
  type KeyboardEvent,
  type MouseEvent,
  useEffect,
  useRef,
  useState,
} from "react";

import { useOverflowGradient } from "@/utils/useOverflowGradient";

import { type ContentOutlineEntry, scrollToOutlineEntry } from "./outline";

const CLOSE_DELAY = 180;
const ENTRY_KEYS = ["Enter", " ", "ArrowDown"];

export function useContentOutlineRail(activeKey?: string) {
  const [open, setOpen] = useState(false);
  const triggerRef = useRef<HTMLButtonElement>(null);
  const panelRef = useRef<HTMLDivElement>(null);
  const ticksRef = useOverflowGradient<HTMLSpanElement>();
  const menuScrollRef = useOverflowGradient<HTMLDivElement>(open);
  const closeTimer = useRef<ReturnType<typeof setTimeout>>(undefined);
  const focusFrame = useRef<number>(undefined);
  const dismissed = useRef(false);

  useEffect(() => {
    return () => {
      clearTimeout(closeTimer.current);
      if (focusFrame.current !== undefined) {
        cancelAnimationFrame(focusFrame.current);
      }
    };
  }, []);

  useEffect(() => {
    const ticks = ticksRef.current;
    const activeTick = ticks?.querySelector<HTMLElement>(
      '[data-active="true"]',
    );

    if (ticks && activeTick) {
      ticks.scrollTop = activeTick.offsetTop - ticks.clientHeight / 2;
    }
  }, [activeKey, ticksRef]);

  function handleKeepOpen() {
    clearTimeout(closeTimer.current);
  }

  function handleCloseSoon() {
    handleKeepOpen();
    closeTimer.current = setTimeout(() => {
      const focusInOutline = panelRef.current?.contains(document.activeElement);
      const focusOnTrigger = document.activeElement === triggerRef.current;

      if (focusInOutline) {
        return;
      }

      if (focusOnTrigger) {
        return;
      }

      setOpen(false);
    }, CLOSE_DELAY);
  }

  function handleOpenChange({ open: nextOpen }: { open: boolean }) {
    handleKeepOpen();
    dismissed.current = !nextOpen;
    setOpen(nextOpen);
  }

  function handleFocusOutside() {
    handleKeepOpen();
    dismissed.current = true;
    setOpen(false);
  }

  function handleEscapeKeyDown() {
    dismissed.current = true;
    triggerRef.current?.focus({ preventScroll: true });
  }

  function handleTriggerPointerEnter() {
    handleKeepOpen();
    dismissed.current = false;
    setOpen(true);
  }

  function handleTriggerFocus() {
    handleKeepOpen();
    if (!dismissed.current) {
      setOpen(true);
    }
  }

  function handleTriggerBlur() {
    dismissed.current = false;
    handleCloseSoon();
  }

  function handleTriggerClick(event: MouseEvent<HTMLButtonElement>) {
    // Focus can open the popover before its built-in click toggle runs.
    event.preventDefault();
    handleTriggerPointerEnter();
  }

  function handleTriggerKeyDown(event: KeyboardEvent<HTMLButtonElement>) {
    if (!ENTRY_KEYS.includes(event.key)) {
      return;
    }

    event.preventDefault();
    handleTriggerPointerEnter();
    focusFrame.current = requestAnimationFrame(() => {
      panelRef.current?.querySelector<HTMLButtonElement>("button")?.focus();
    });
  }

  function handlePanelFocus() {
    dismissed.current = false;
    handleKeepOpen();
  }

  function handleNavigate(entry: ContentOutlineEntry) {
    handleKeepOpen();
    dismissed.current = true;
    setOpen(false);
    scrollToOutlineEntry(entry);
  }

  return {
    open,
    triggerRef,
    panelRef,
    ticksRef,
    menuScrollRef,
    handlers: {
      handleKeepOpen,
      handleCloseSoon,
      handleOpenChange,
      handleFocusOutside,
      handleEscapeKeyDown,
      handleTriggerPointerEnter,
      handleTriggerFocus,
      handleTriggerBlur,
      handleTriggerClick,
      handleTriggerKeyDown,
      handlePanelFocus,
      handleNavigate,
    },
  };
}
