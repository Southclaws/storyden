import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";
import NextLink, { LinkProps as NextLinkProps } from "next/link";
import { PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { ButtonVariantProps, button } from "@/styled-system/recipes";
import { JsxStyleProps } from "@/styled-system/types";

export type LinkProps = Partial<ButtonVariantProps> &
  NextLinkProps &
  JsxStyleProps;

export type LinkButtonStyleProps = Partial<ButtonVariantProps> & JsxStyleProps;

export function LinkButton({
  children,
  href,
  ...props
}: PropsWithChildren<LinkProps>) {
  const [vp, stripped] = button.splitVariantProps(props);

  const cn = cx(button(vp), css(stripped));

  const isExternal = !(
    href.toString().startsWith("/") || href.toString().startsWith("#")
  );

  const target = isExternal ? "_blank" : undefined;

  return (
    <NextLink
      className={cn}
      href={href}
      target={target}
      onClick={props.onClick}
    >
      <styled.span
        display="flex"
        // Supports overflowing children and text ellipsis
        maxW="full"
        alignItems="center"
        gap="1"
      >
        {children}
      </styled.span>
      {isExternal && (
        <styled.span color="fg.subtle">
          <ArrowTopRightOnSquareIcon />
        </styled.span>
      )}
    </NextLink>
  );
}
