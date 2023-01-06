import { Box } from "@chakra-ui/react";
import { ReactNode } from "react";

export function AuthBox({ children }: { children: ReactNode }) {
  return (
    <Box
      id="AuthBox"
      p={6}
      bg="linear-gradient(141.91deg, #B7CEF1 0%, #2FD596 99.55%)"
      borderRadius={8}
      width="full"
      maxW="xs"
    >
      {children}
    </Box>
  );
}
