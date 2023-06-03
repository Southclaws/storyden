"use client";

import { Box, Flex } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

export function Fullpage(props: PropsWithChildren) {
  return (
    <Flex
      width="full"
      height="full"
      minHeight="100vh"
      justifyContent="start"
      alignItems="center"
      flexDirection="column"
    >
      <Box as="main" flexGrow={1} width="full" height="full">
        {props.children}
      </Box>
    </Flex>
  );
}
