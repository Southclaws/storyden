import { PropsWithChildren } from "react";

import { Navigation } from "src/components/site/Navigation/Navigation";

import { Box, Flex, styled } from "@/styled-system/jsx";

export async function Default(props: PropsWithChildren) {
  return (
    <Flex
      minHeight="dvh"
      width="full"
      flexDirection="row"
      backgroundColor="bg.site"
      vaul-drawer-wrapper=""
    >
      <Navigation>
        <styled.main
          containerType="inline-size"
          width="full"
          height="full"
          minW="0"
        >
          {props.children}
          <Box height="24"></Box>
        </styled.main>
      </Navigation>
    </Flex>
  );
}
