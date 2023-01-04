import { Box, Flex } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { Footer } from "src/components/Footer";
import { Navigation } from "src/components/Navigation/Navigation";

export function Default(props: PropsWithChildren) {
  return (
    <Flex
      width="full"
      height="full"
      minHeight="100vh"
      justifyContent="start"
      alignItems="center"
      flexDirection="column"
    >
      <Navigation />

      <Box as="main" width="full" height="full" maxW="container.md" px={4}>
        {props.children}
      </Box>

      <Footer />
    </Flex>
  );
}
