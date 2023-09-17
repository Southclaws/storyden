import { PropsWithChildren } from "react";

import { Navigation } from "src/components/site/Navigation/Navigation";

import { Box, Flex, styled } from "@/styled-system/jsx";

export function Default(props: PropsWithChildren) {
  return (
    <Flex
      minHeight="100vh"
      width="full"
      flexDirection="row"
      bgColor="white"
      vaul-drawer-wrapper=""
    >
      <Navigation />

      <styled.main
        width="full"
        maxW={{
          base: "full",
          lg: "3xl",
        }}
        px={4}
        py={2}
        backgroundColor="white"
      >
        {props.children}
        <Box height="6rem"></Box>
      </styled.main>
    </Flex>
  );
}
