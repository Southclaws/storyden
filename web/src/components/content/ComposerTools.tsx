import { Portal } from "@ark-ui/react";
import { motion } from "framer-motion";
import { PropsWithChildren, useRef, useState } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { Box, HStack, styled } from "@/styled-system/jsx";

import { Spinner } from "../ui/Spinner";

type Props = {
  icon: React.ReactNode;
  expandedIcon?: React.ReactNode;
  workingCount?: number;
  onClick?: () => void;
};

export function ComposerTools({
  icon,
  expandedIcon,
  workingCount = 0,
  onClick,
  children,
}: PropsWithChildren<Props>) {
  const [isExpanded, setIsExpanded] = useState(false);
  const closeTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const isWorking = workingCount > 0;

  const handleExpand = () => {
    if (isExpanded) {
      onClick?.();
      return;
    }

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

  return (
    <Box
      position="absolute"
      height="full"
      top="-1"
      right="-1"
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
  );
}
