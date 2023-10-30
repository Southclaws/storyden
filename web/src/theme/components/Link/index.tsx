import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";
import NextLink, { LinkProps as NextLinkProps } from "next/link";
import { PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import { LinkVariant, link } from "@/styled-system/recipes";
import { StyleProps } from "@/styled-system/types";

export type LinkProps = Partial<LinkVariant> & NextLinkProps & StyleProps;

export function Link({ children, ...props }: PropsWithChildren<LinkProps>) {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [_, stripped] = link.splitVariantProps(props);

  const cn = cx(link(props), css(stripped));

  const isExternal = !props.href.toString().startsWith("/");

  return (
    <NextLink className={cn} {...(stripped as NextLinkProps)}>
      {children}
      {isExternal && <ArrowTopRightOnSquareIcon />}
    </NextLink>
  );
}
