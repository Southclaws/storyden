import { PropsWithChildren } from "react";

import { Navigation } from "src/components/site/Navigation/Navigation";

import { Box, Flex, styled } from "@/styled-system/jsx";

export function Default(props: PropsWithChildren) {
  return (
    <Flex
      minHeight="dvh"
      width="full"
      flexDirection="row"
      backgroundColor="accent.50"
      vaul-drawer-wrapper=""
    >
      <Navigation>
        <styled.main width="full" minW="0">
          {props.children}
          <Box height="24"></Box>
        </styled.main>
      </Navigation>
    </Flex>
  );
}
