"use client";

import { styled } from "@/styled-system/jsx";

type Props = {
  signature?: string;
  maxHeight?: number;
};

export function Signature({ signature, maxHeight = 160 }: Props) {
  if (!signature || signature === "<body></body>") {
    return null;
  }

  return (
    <styled.div
      className="member__signature typography"
      my="2"
      pt="1"
      width="full"
      borderTopWidth="thin"
      borderTopColor="border.subtle"
      // Force all typography to be subtle, to avoid distracting from content.
      color="fg.muted"
      fontSize="sm"
      overflow="hidden"
      style={{ maxHeight }}
      // NOTE: the signature content is typed as datagraph.Content on the API
      // forcing sanitisation on the server side, we can trust it's safe.
      dangerouslySetInnerHTML={{ __html: signature }}
    />
  );
}
