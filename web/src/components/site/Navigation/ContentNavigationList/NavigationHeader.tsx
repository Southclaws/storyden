import Link from "next/link";
import { PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import { WStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

type Props = {
  href: string;
  controls?: React.ReactNode;
  headingAs?: "h2" | "h3";
};

const linkStyles = cx(
  button({ variant: "ghost" }),
  css({
    p: "1",
  }),
);

export function NavigationHeader({
  children,
  href,
  controls,
  headingAs = "h2",
}: PropsWithChildren<Props>) {
  return (
    <WStack>
      <Link className={linkStyles} href={href} prefetch={false}>
        <styled.h2 as={headingAs} fontSize="xs" fontWeight="medium">
          {children}
        </styled.h2>
      </Link>

      {controls}
    </WStack>
  );
}
