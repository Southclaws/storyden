import { motion } from "framer-motion";
import { PropsWithChildren, useEffect, useRef, useState } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { Box, HStack, styled } from "@/styled-system/jsx";

import { Spinner } from "../ui/Spinner";

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

  if (!enabled) {
    return null;
  }

  return (
    <Box
      position="absolute"
      height="full"
      width="full"
      top="-1"
      right="-1"
      pointerEvents="none"
    >
      <Box
        position="sticky"
        top={{ base: "4", md: "20" }}
        width="full"
        display="flex"
        justifyContent="flex-end"
        zIndex="docked"
        pointerEvents="none"
      >
        <Box
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
          <HStack gap="2">
            <motion.div
              animate={
                isExpanded
                  ? { width: "auto" }
                  : { width: isWorking ? "24px" : "0" }
              }
              transition={{ duration: 0.2, ease: "easeInOut" }}
              style={{ overflow: "hidden", display: "flex", gap: "8px" }}
            >
              {isWorking && (
                <HStack gap="1">
                  <Spinner size="sm" />
                  {workingCount > 1 && (
                    <styled.span fontSize="xs" color="fg.muted">
                      {workingCount}
                    </styled.span>
                  )}
                </HStack>
              )}

              <div
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
