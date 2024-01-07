import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";
import NextLink, { LinkProps as NextLinkProps } from "next/link";
import { PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { LinkVariant, link } from "@/styled-system/recipes";
import { StyleProps } from "@/styled-system/types";

export type LinkProps = Partial<LinkVariant> & NextLinkProps & StyleProps;

export function Link({
  children,
  href,
  ...props
}: PropsWithChildren<LinkProps>) {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [vp, stripped] = link.splitVariantProps(props);

  const cn = cx(link(vp), css(stripped));

  const isExternal = !href.toString().startsWith("/");

  return (
    <NextLink className={cn} href={href}>
      <styled.span
        display="flex"
        alignItems="center"
        gap="1"
        textOverflow="ellipsis"
        overflowX="hidden"
        textWrap="nowrap"
      >
        {children}
      </styled.span>
      {isExternal && (
        <styled.span>
          <ArrowTopRightOnSquareIcon />
        </styled.span>
      )}
    </NextLink>
  );
}
