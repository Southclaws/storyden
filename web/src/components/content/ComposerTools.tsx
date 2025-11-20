import { AnimatePresence, motion } from "framer-motion";
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

  // Logic: reveal the animated presence container if working or hovered, this
  // allows the spinner to appear on its own when the editor is working. When
  // hovered, we also include the preview switch. If the user hovers while the
  // editor is working, the preview switch appears alongside the spinner.
  const reveal = isWorking || isExpanded;

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
        backgroundColor={isExpanded ? "bg.subtle" : "transparent"}
        backdropBlur="frosted"
        backdropFilter="auto"
        borderRadius="md"
        transition="all"
        pointerEvents="auto"
        p="1"
      >
        <HStack gap="2">
          <AnimatePresence>
            {reveal && (
              <motion.div
                initial={{ width: 0, opacity: 0 }}
                animate={{ width: "auto", opacity: 1 }}
                exit={{ width: 0, opacity: 0 }}
                transition={{ duration: 0.2, ease: "easeInOut" }}
                style={{ overflow: "hidden" }}
              >
                {isWorking && (
                  <HStack gap="1">
                    <Spinner />
                    {workingCount > 1 && (
                      <styled.span fontSize="xs" color="fg.muted">
                        {workingCount}
                      </styled.span>
                    )}
                  </HStack>
                )}

                {isExpanded && children}
              </motion.div>
            )}
          </AnimatePresence>

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
