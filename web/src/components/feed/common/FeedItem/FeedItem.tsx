import { PropsWithChildren } from "react";

import { styled } from "@/styled-system/jsx";

export function FeedItem({ children }: PropsWithChildren) {
  return (
    <styled.article
      display="flex"
      flexDir="column"
      width="full"
      p={2}
      gap={2}
      boxShadow="md"
      borderRadius="md"
      backgroundColor="white"
    >
      {children}
    </styled.article>
  );
}
