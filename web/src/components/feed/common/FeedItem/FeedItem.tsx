import { PropsWithChildren } from "react";

import { styled } from "@/styled-system/jsx";

export function FeedItem({ children }: PropsWithChildren) {
  return (
    <styled.article
      display="flex"
      flexDir="column"
      width="full"
      p={3}
      gap={2}
      boxShadow="6px 6px 16px 0px rgba(0, 0, 0, 0.03)"
      borderRadius="1rem"
      backgroundColor="white"
    >
      {children}
    </styled.article>
  );
}
