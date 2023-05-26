import { Box, Flex, FlexProps } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

export function Toolpill(props: PropsWithChildren<FlexProps>) {
  return (
    <Box
      id="toolpill-overlay"
      position="fixed"
      left={0}
      bottom="env(safe-area-inset-bottom)"
      pb={2}
      width="100vw"
      height="100vh"
      pointerEvents="none"
      zIndex="overlay"
    >
      <Flex
        id="flex-container"
        height="full"
        gap={3}
        p="min(3vh, 1em)"
        justifyContent="end"
        alignItems="center"
        flexDir="column"
      >
        <Flex
          p={1}
          gap={2}
          flexDirection="column"
          borderRadius={20}
          backdropFilter="blur(4px)"
          transitionProperty="background-color"
          transitionDuration="0.5s"
          backgroundColor="hsla(210, 38.5%, 94.9%, 0.8)"
          border="1px solid hsla(209, 100%, 20%, 0.02)"
          width="full"
          maxW={{ base: "23em", md: "container.sm" }}
          justifyContent="space-between"
          alignItems="center"
          pointerEvents="auto"
          {...props}
        >
          {props.children}
        </Flex>
      </Flex>
    </Box>
  );
}
