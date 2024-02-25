import { PropsWithChildren } from "react";

import { styled } from "@/styled-system/jsx";
import { CardBox } from "@/styled-system/patterns";

export function FeedItem({ children }: PropsWithChildren) {
  return <styled.article className={CardBox()}>{children}</styled.article>;
}
