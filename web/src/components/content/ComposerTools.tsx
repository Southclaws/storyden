import { Portal } from "@ark-ui/react";
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

  // Positioning logic - this is needed because in order to show editor tools
  // within a container with overflow:hidden, we need to portal the component.
  const anchorRef = useRef<HTMLDivElement>(null); // inside overflow container
  const floatingRef = useRef<HTMLDivElement>(null); // portalled wrapper

  useLayoutEffect(() => {
    if (!anchorRef.current || !floatingRef.current) return;

    return autoUpdate(anchorRef.current, floatingRef.current, () => {
      const rect = anchorRef.current!.getBoundingClientRect();

      const top = window.scrollY + rect.top;
      const right = window.innerWidth - rect.right;

      Object.assign(floatingRef.current!.style, {
        top: `${top}px`,
        height: `${rect.height}px`,
        right: `${right}px`,
        maxWidth: `${rect.width}px`,
      });
    });
  }, []);

  const isWorking = workingCount > 0;

  const handleExpand = () => {
    if (closeTimeoutRef.current) {
      clearTimeout(closeTimeoutRef.current);
      closeTimeoutRef.current = null;
    }
    setIsExpanded(true);
  };

  const handleContract = () => {
    closeTimeoutRef.current = setTimeout(() => {
      setIsExpanded(false);
    }, 650);
  };

  if (!enabled) {
    return null;
  }

  return (
    <Box
      ref={anchorRef}
      position="absolute"
      height="full"
      width="full"
      top="-1"
      right="-1"
      pointerEvents="none"
    >
      <Portal>
        <Box
          ref={floatingRef}
          position="absolute"
          height="full"
          pointerEvents="none"
        >
          <Box
            position="sticky"
            top={{ base: "4", md: "20" }}
            zIndex="popover"
            opacity={isExpanded ? "full" : "5"}
            onMouseEnter={handleExpand}
            onMouseLeave={handleContract}
            cursor="pointer"
            backgroundColor={isExpanded ? "bg.subtle/80" : "transparent"}
            backdropBlur={isExpanded ? "sm" : undefined}
            backdropFilter={isExpanded ? "auto" : undefined}
            borderRadius="md"
            transition="all"
            pointerEvents="auto"
            p="1"
          >
            <HStack gap="2">
              <motion.div
                animate={
                  isExpanded
                    ? {
                        width: "auto",
                      }
                    : {
                        width: isWorking ? "24px" : "0",
                      }
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
                onClick={handleExpand}
                aria-label="Show editor tools"
                aria-expanded={isExpanded}
              >
                {expandedIcon ? (
                  // Switch icons if there's an expanded alternate.
                  <>{isExpanded ? expandedIcon : icon}</>
                ) : (
                  // Otherwise, just render the icon.
                  icon
                )}
              </IconButton>
            </HStack>
          </Box>
        </Box>
      </Portal>
    </Box>
  );
}
