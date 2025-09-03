import Link from "next/link";
import { PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import { HStack, WStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

type Props = {
  href: string;
  controls?: React.ReactNode;
};

const linkStyles = cx(
  button({ variant: "ghost", size: "xs" }),
  css({
    p: "1",
  }),
);

export function NavigationHeader({
  children,
  href,
  controls,
}: PropsWithChildren<Props>) {
  return (
    <WStack>
      <Link className={linkStyles} href={href}>
        <styled.h1 fontSize="xs" fontWeight="medium">
          {children}
        </styled.h1>
      </Link>

      {controls}
    </WStack>
  );
}
