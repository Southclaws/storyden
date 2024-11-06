import { PropsWithChildren } from "react";

import { HStack, styled } from "@/styled-system/jsx";

import { EmptyIcon } from "../ui/icons/Empty";

export function Empty({ children }: PropsWithChildren) {
  return (
    <HStack alignItems="center" color="fg.muted">
      <EmptyIcon />
      <styled.p fontStyle="italic" textWrap="nowrap">
        {children}
      </styled.p>
    </HStack>
  );
}
