import Link from "next/link";
import { PropsWithChildren } from "react";

import { styled } from "@/styled-system/jsx";

export function Title({ children }: PropsWithChildren) {
  return (
    <styled.h1
      fontSize="lg"
      fontWeight="bold"
      textWrap="nowrap"
      overflow="hidden"
      textOverflow="ellipsis"
    >
      <Link href="/">{children}</Link>
    </styled.h1>
  );
}
