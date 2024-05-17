import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";
import NextLink, { LinkProps as NextLinkProps } from "next/link";
import { PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { ButtonVariantProps, button } from "@/styled-system/recipes";
import { StyleProps } from "@/styled-system/types";

export type LinkProps = Partial<ButtonVariantProps> &
  NextLinkProps &
  StyleProps;

export function LinkButton({
  children,
  href,
  ...props
}: PropsWithChildren<LinkProps>) {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [vp, stripped] = button.splitVariantProps(props);

  const cn = cx(button(vp), css(stripped));

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
