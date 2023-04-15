import { Box, Flex } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Navigation } from "src/components/Navigation/Navigation";

export function Default(props: PropsWithChildren) {
  return (
    <Flex
      width="full"
      height="full"
      minHeight="100vh"
      alignItems="stretch"
      flexDirection="row"
    >
      <Navigation />

      <Flex as="main" px={4} w="full">
        <Box
          w="full"
          maxW={{
            base: "full",
            sm: "container.sm",
            md: "container.md",
          }}
          px={1}
        >
          {props.children}
        </Box>
      </Flex>
    </Flex>
  );
}
