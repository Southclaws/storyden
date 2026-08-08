import { PropsWithChildren } from "react";

import { HStack, styled } from "@/styled-system/jsx";

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
      <div id="desktop-nav-right">{children}</div>

      <HStack color="text.muted" fontSize="xs">
        {/* TODO: Provide links to privacy/terms/etc custom pages */}
        {/* <p>copyright {settings.owner}</p> */}
        {/* <a href={PrivacyRoute}>privacy</a> */}
        <p>powered by storyden</p>
      </HStack>
    </styled.nav>
  );
}
