import { JSX, PropsWithChildren } from "react";

import { Box, Flex } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";
import { useClickAway } from "@/utils/useClickAway";

export type Props = {
  isOpen: boolean;
  render: () => JSX.Element;
  onClickOutside?: () => void;
};

export function CommandDock({
  isOpen,
  render,
  onClickOutside,
  children,
}: PropsWithChildren<Props>) {
  const ref = useClickAway<HTMLDivElement>(() => {
    onClickOutside?.();
  });

  return (
    <Box
      id="command-dock__overlay"
      position="fixed"
      left="0"
      bottom="safeBottom"
      width="screen"
      height="dvh"
      pointerEvents="none"
      zIndex="overlay"
    >
      <Flex
        id="command-dock__outer-container"
        height="full"
        px="2"
        pt="2"
        pb="10"
        justifyContent="end"
        alignItems="center"
        flexDir="column"
      >
        <Flex
          id="command-dock__content-container"
          ref={ref}
          className={Floating()}
          width="full"
          maxW="96"
          maxH="full"
          minH="0"
          flexDirection="column"
          justifyContent="space-between"
          alignItems="center"
          gap="2"
          p="2"
          borderRadius="xl"
          borderWidth="thin"
          borderStyle="solid"
          borderColor="border.muted"
          pointerEvents="auto"
        >
          {isOpen && render()}

          {children}
        </Flex>
      </Flex>
    </Box>
  );
}
