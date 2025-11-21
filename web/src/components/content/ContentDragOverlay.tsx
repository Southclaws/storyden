import { Box, styled } from "@/styled-system/jsx";

export function ContentDragOverlay({ isError, message }) {
  return (
    <Box
      position="absolute"
      top="0"
      left="0"
      right="0"
      bottom="0"
      pointerEvents="none"
      display="flex"
      alignItems="center"
      justifyContent="center"
      backgroundColor="bg.emphasized"
      borderWidth="medium"
      borderStyle="dashed"
      borderColor={isError ? "border.error" : "accent.default"}
      borderRadius="md"
      style={{ opacity: 0.95 }}
      role="status"
      aria-live="polite"
      aria-label={message}
    >
      <styled.div
        fontSize="sm"
        fontWeight="medium"
        color={isError ? "fg.error" : "accent.default"}
        display="flex"
        flexDirection="column"
        alignItems="center"
        gap="2"
      >
        <span>{message}</span>
      </styled.div>
    </Box>
  );
}
