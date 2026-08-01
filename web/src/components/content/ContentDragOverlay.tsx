import { Box, styled } from "@/styled-system/jsx";

type Props = {
  isError: boolean;
  message: string;
};

export function ContentDragOverlay({ isError, message }: Props) {
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
      backgroundColor="interactive.emphasized.surface"
      borderWidth="medium"
      borderStyle="dashed"
      borderColor={isError ? "status.danger.border" : "accent.default"}
      borderRadius="md"
      style={{ opacity: 0.95 }}
      role="status"
      aria-live="polite"
      aria-label={message}
    >
      <styled.div
        fontSize="sm"
        fontWeight="medium"
        color={isError ? "status.danger.content" : "accent.default"}
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
