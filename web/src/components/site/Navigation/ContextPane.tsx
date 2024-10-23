import { PropsWithChildren } from "react";

import { Box, HStack, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

export function ContextPane({ children }: PropsWithChildren) {
  return (
    <styled.nav
      display="flex"
      flexDir="column"
      alignItems="center"
      gap="2"
      width="full"
      height="full"
    >
      <Box
        id="desktop-nav-right"
        className={Floating()}
        borderRadius="md"
        w="full"
        height="min"
        p="2"
        pr="0"
        overflowY="scroll"
      >
        {children}
      </Box>

      <HStack color="fg.subtle" fontSize="xs">
        {/* TODO: Provide links to privacy/terms/etc custom pages */}
        {/* <p>copyright {settings.owner}</p> */}
        {/* <a href={PrivacyRoute}>privacy</a> */}
        <p>powered by storyden</p>
      </HStack>
    </styled.nav>
  );
}
