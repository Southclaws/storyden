import { PropsWithChildren } from "react";

import { Box, Flex, FlexProps } from "@/styled-system/jsx";

import { Props, useToolpill } from "./useToolpill";

export function Toolpill({
  onClickOutside,
  ...props
}: PropsWithChildren<FlexProps & Props>) {
  const { ref } = useToolpill({ onClickOutside, ...props });
  return (
    <Box
      id="toolpill-overlay"
      position="fixed"
      left="0"
      bottom="safeBottom"
      width="screen"
      height="dvh"
      pointerEvents="none"
      zIndex="overlay"
    >
      <Flex
        id="toolpill-flex-outer-container"
        height="full"
        px="2"
        pb="10"
        justifyContent="end"
        alignItems="center"
        flexDir="column"
      >
        <Flex
          id="toolpill-content-container"
          ref={ref}
          p="1"
          gap="2"
          flexDirection="column"
          borderRadius="xl"
          backdropBlur="sm"
          backdropFilter="auto"
          transitionProperty="background-color"
          transitionDuration="fast"
          backgroundColor="whiteAlpha.900"
          borderWidth="thin"
          borderStyle="solid"
          borderColor="blackAlpha.50"
          width="full"
          maxW="96"
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
