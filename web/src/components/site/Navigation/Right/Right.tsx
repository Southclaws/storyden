"use client";

import { VStack, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import { useRightPortalRef } from "./context";

export function Right() {
  const { targetRef } = useRightPortalRef();

  const children = targetRef.current?.children;

  // hide the entire panel if no children in the portal.
  if (children?.length === 0) return null;

  return (
    <VStack className={Floating()} justify="space-between" w="full" p="4">
      <styled.aside w="full" ref={targetRef} />
    </VStack>
  );
}
