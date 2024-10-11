import React, { PropsWithChildren } from "react";

import { Box, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

export function ContextPane({ children }: PropsWithChildren) {
  return (
    <styled.nav
      className={Floating()}
      display="flex"
      borderRadius="md"
      flexDir="column"
      gap="4"
      width="full"
      alignItems="start"
      justifyContent="space-between"
      overflowY="scroll"
    >
      <Box id="desktop-nav-right" w="full" height="full" p="2" pr="0">
        <styled.aside w="full">{children}</styled.aside>
      </Box>
    </styled.nav>
  );
}
