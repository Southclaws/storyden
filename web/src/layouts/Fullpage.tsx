import { PropsWithChildren } from "react";

import { Flex, styled } from "@/styled-system/jsx";

export function Fullpage(props: PropsWithChildren) {
  return (
    <Flex
      className="sd-layout sd-layout--fullpage"
      width="full"
      height="full"
      minHeight="dvh"
      justifyContent="start"
      alignItems="center"
      flexDirection="column"
    >
      <styled.main
        className="sd-screen sd-screen--fullpage"
        flexGrow={1}
        width="full"
        height="full"
      >
        {props.children}
      </styled.main>
    </Flex>
  );
}
