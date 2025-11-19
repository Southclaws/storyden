import { AnimatePresence, motion } from "framer-motion";
import { useState } from "react";

import { IconButton } from "@/components/ui/icon-button";
import { EditIcon } from "@/components/ui/icons/Edit";
import { Box, HStack } from "@/styled-system/jsx";

import { EditorMenu, Props } from "./EditorMenu";

export function EditorTools(props: Props) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <Box position="absolute" height="full" right="0" p="1" pointerEvents="none">
      <Box
        position="sticky"
        top={{ base: "4", md: "20" }}
        zIndex="popover"
        opacity={isHovered ? "full" : "5"}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        cursor="pointer"
        backgroundColor="bg.subtle"
        backdropBlur="frosted"
        backdropFilter="auto"
        borderRadius="md"
        transition="all"
        pointerEvents="auto"
      >
        <HStack gap="2">
          <AnimatePresence>
            {isHovered && (
              <motion.div
                initial={{ width: 0, opacity: 0 }}
                animate={{ width: "auto", opacity: 1 }}
                exit={{ width: 0, opacity: 0 }}
                transition={{ duration: 0.2, ease: "easeInOut" }}
                style={{ overflow: "hidden" }}
              >
                <EditorMenu {...props} />
              </motion.div>
            )}
          </AnimatePresence>

          <IconButton type="button" variant="ghost" size="xs">
            <EditIcon w="4" />
          </IconButton>
        </HStack>
      </Box>
    </Box>
  );
}
