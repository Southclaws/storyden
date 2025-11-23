import { autoUpdate } from "@floating-ui/react";
import { motion } from "framer-motion";
import {
  PropsWithChildren,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
} from "react";

import { IconButton } from "@/components/ui/icon-button";
import { Box, HStack, styled } from "@/styled-system/jsx";

import { Spinner } from "../ui/Spinner";

// NOTE: Should match the MD breakpoint in `panda.config.ts`.
const MD_BREAKPOINT = 768;

// NOTE: Desktop is more because we have to make space for the navigation bar.
const MOBILE_TOP_OFFSET = 16;
const DESKTOP_TOP_OFFSET = 80;

type Props = {
  enabled: boolean;
  icon: React.ReactNode;
  expandedIcon?: React.ReactNode;
  workingCount?: number;
  onClick?: () => void;
};

export function ComposerTools({
  enabled,
  icon,
  expandedIcon,
  workingCount = 0,
  onClick,
  children,
}: PropsWithChildren<Props>) {
  const [isExpanded, setIsExpanded] = useState(false);
  const closeTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    return () => {
      if (closeTimeoutRef.current) {
        clearTimeout(closeTimeoutRef.current);
      }
    };
  }, []);

  const isWorking = workingCount > 0;

  const handleClick = () => {
    // Toggle the menu state
    // On desktop (mouse): closes if already open from hover
    // On mobile (touch): toggles because pointer events are ignored
    setIsExpanded((prev) => !prev);

    onClick?.();
  };

  const handlePointerEnter = (e: React.PointerEvent) => {
    // Only respond to mouse hover, not touch events
    if (e.pointerType === "touch") return;

    if (closeTimeoutRef.current) {
      clearTimeout(closeTimeoutRef.current);
      closeTimeoutRef.current = null;
    }
    setIsExpanded(true);
  };

  const handlePointerLeave = (e: React.PointerEvent) => {
    // Only respond to mouse leave, not touch events
    if (e.pointerType === "touch") return;

    closeTimeoutRef.current = setTimeout(() => {
      setIsExpanded(false);
    }, 650);
  };

  // The anchor is the absolutely positioned container. This is fixed to its
  // parent and does not move relative to the viewport. It acts as the position
  // anchor as the user scrolls, it may be out of view or within the viewport.
  const anchorRef = useRef<HTMLDivElement>(null);
  // The floating element is the menu itself, this is where we essentially build
  // our own "position: sticky" behavior by calculating the anchor's position
  // relative to the viewport scroll position and bound its position to the Y
  // dimension of the anchor. The anchor is height=full, so fills the textarea.
  const floatingRef = useRef<HTMLDivElement>(null);

  useLayoutEffect(() => {
    if (!enabled) return;
    if (!anchorRef.current || !floatingRef.current) return;

    const updatePosition = () => {
      if (!anchorRef.current || !floatingRef.current) return;

      const anchorRect = anchorRef.current.getBoundingClientRect();
      const floatingRect = floatingRef.current.getBoundingClientRect();

      // TODO: Use a resize observer to handle viewport size changes, memoise.
      const isMobile = window.innerWidth < MD_BREAKPOINT;
      const navbarOffset = isMobile ? MOBILE_TOP_OFFSET : DESKTOP_TOP_OFFSET;

      // The desired Y position in the viewport for the floating menu.
      const desiredViewportY = navbarOffset;
      const anchorTop = anchorRect.top;
      const floatingHeight = floatingRect.height;

      // Calculate how far the anchor is from the desired viewport Y position.
      // Clamp to zero so the floating menu doesn't go above the anchor.
      const anchorRelativeY = Math.max(0, desiredViewportY - anchorTop);

      // maxY is the maximum Y offset the floating menu can have within the
      // anchor, so it doesn't overflow past the bottom of the anchor.
      const maxY = anchorRect.height - floatingHeight;

      // Constrain the floating menu's Y position so it stays within the
      // anchor's bounds. This ensures the menu is always visible and doesn't
      // overflow, even if the anchor is small.
      const constrainedY = Math.min(anchorRelativeY, maxY);

      Object.assign(floatingRef.current.style, {
        top: `${constrainedY}px`,
      });
    };

    updatePosition();

    return autoUpdate(anchorRef.current, floatingRef.current, updatePosition);
  }, [enabled]);

  if (!enabled) {
    return null;
  }

  return (
    <Box
      ref={anchorRef}
      className="composer-tools__anchor"
      position="absolute"
      height="full"
      width="full"
      top="-1"
      right="-1"
      pointerEvents="none"
    >
      <Box
        ref={floatingRef}
        className="composer-tools__sticky-container"
        position="absolute"
        width="full"
        display="flex"
        justifyContent="flex-end"
        zIndex="docked"
        pointerEvents="none"
      >
        <Box
          className="composer-tools__hover-trigger"
          opacity={isExpanded ? "full" : "5"}
          onPointerEnter={handlePointerEnter}
          onPointerLeave={handlePointerLeave}
          cursor="pointer"
          backgroundColor={isExpanded ? "bg.subtle/80" : "transparent"}
          backdropBlur={isExpanded ? "sm" : undefined}
          backdropFilter={isExpanded ? "auto" : undefined}
          borderRadius="md"
          pointerEvents="auto"
          p="1"
          maxWidth="full"
        >
          <HStack gap="2" className="composer-tools__expanding-stack">
            <motion.div
              className="composer-tools__animated-reveal"
              animate={
                isExpanded
                  ? { width: "auto" }
                  : { width: isWorking ? "24px" : "0" }
              }
              transition={{ duration: 0.2, ease: "easeInOut" }}
              style={{ overflow: "hidden", display: "flex", gap: "8px" }}
            >
              {isWorking && (
                <HStack gap="1" className="composer-tools__working-status">
                  <Spinner size="sm" />
                  {workingCount > 1 && (
                    <styled.span
                      className="composer-tools__working-status-count"
                      fontSize="xs"
                      color="fg.muted"
                    >
                      {workingCount}
                    </styled.span>
                  )}
                </HStack>
              )}

              <div
                className="composer-tools__children-container"
                style={{
                  visibility: isExpanded ? "visible" : "hidden",
                  width: isExpanded ? "auto" : 0,
                  overflowX: isExpanded ? "scroll" : "hidden",
                  scrollbarWidth: "none",
                }}
              >
                {children}
              </div>
            </motion.div>

            <IconButton
              className="composer-tools__expand-button"
              type="button"
              variant="ghost"
              size="xs"
              onClick={handleClick}
              aria-label="Show editor tools"
              aria-expanded={isExpanded}
            >
              {expandedIcon ? <>{isExpanded ? expandedIcon : icon}</> : icon}
            </IconButton>
          </HStack>
        </Box>
      </Box>
    </Box>
  );
}
