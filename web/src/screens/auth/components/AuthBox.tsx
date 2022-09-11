import { Box, Flex } from "@chakra-ui/react";
import { ReactNode } from "react";

export function AuthBox({ children }: { children: ReactNode }) {
  return (
    <Box
      p={12}
      bg="linear-gradient(141.91deg, #B7CEF1 0%, #2FD596 99.55%)"
      borderRadius={8}
    >
      <Flex flexDirection="column">{children}</Flex>
    </Box>
  );
}
