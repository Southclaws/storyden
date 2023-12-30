import { PropsWithChildren } from "react";

import { styled } from "@/styled-system/jsx";
import { Card } from "@/styled-system/patterns";

export function FeedItem({ children }: PropsWithChildren) {
  return <styled.article className={Card()}>{children}</styled.article>;
}
