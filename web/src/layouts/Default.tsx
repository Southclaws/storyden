import { Box, Flex } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

export function Default(props: PropsWithChildren) {
  return (
    <Flex width="full" justifyContent="center">
      <Box as="main" width="full" maxW="container.lg" px={2}>
        {props.children}
      </Box>
    </Flex>
  );
}
