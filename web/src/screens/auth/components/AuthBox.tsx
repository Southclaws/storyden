import { Box } from "@chakra-ui/react";
import { ReactNode } from "react";

export function AuthBox({ children }: { children: ReactNode }) {
  return (
    <Box
      // TODO: Figure out why this isn't scaling well.
      px={{ base: "10%", sm: "20%", md: 24 }}
      py={12}
      bg="linear-gradient(141.91deg, #B7CEF1 0%, #2FD596 99.55%)"
      borderRadius={8}
    >
      {children}
    </Box>
  );
}
