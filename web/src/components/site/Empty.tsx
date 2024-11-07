import { PropsWithChildren } from "react";

import { VStack, styled } from "@/styled-system/jsx";

import { EmptyIcon } from "../ui/icons/Empty";

/** @deprecated use EmptyState */
export function Empty({ children }: PropsWithChildren) {
  return (
    <VStack alignItems="center" color="fg.muted">
      <EmptyIcon />
      <styled.p fontStyle="italic" textWrap="nowrap">
        {children}
      </styled.p>
    </VStack>
  );
}
