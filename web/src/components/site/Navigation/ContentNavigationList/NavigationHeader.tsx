import { PropsWithChildren } from "react";

import { HStack, styled } from "@/styled-system/jsx";

import { AddAction } from "../../Action/Add";

type Props = {
  controls?: React.ReactNode;
};

export function NavigationHeader({
  children,
  controls,
}: PropsWithChildren<Props>) {
  return (
    <HStack w="full" justify="space-between">
      <styled.h1 pl="1" fontSize="xs" fontWeight="medium" color="fg.subtle">
        {children}
      </styled.h1>

      {controls}
    </HStack>
  );
}
