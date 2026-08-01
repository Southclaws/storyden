import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/outline";
import NextLink, { LinkProps as NextLinkProps } from "next/link";
import { PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { ButtonVariantProps, button } from "@/styled-system/recipes";
import { JsxStyleProps } from "@/styled-system/types";

type LinkButtonPassthroughProps = {
  className?: string;
  [key: `aria-${string}`]: string | number | boolean | undefined;
  [key: `data-${string}`]: string | number | boolean | undefined;
};

export type LinkProps = Partial<ButtonVariantProps> &
  NextLinkProps &
  JsxStyleProps &
  LinkButtonPassthroughProps;

export type LinkButtonStyleProps = Partial<ButtonVariantProps> & JsxStyleProps;

export function LinkButton({
  children,
  className,
  href,
  onClick,
  ...props
}: PropsWithChildren<LinkProps>) {
  const [vp, stripped] = button.splitVariantProps(props);
  const passthroughProps: Record<
    string,
    string | number | boolean | undefined
  > = {};
  const styleProps: Record<string, unknown> = {};

  for (const [key, value] of Object.entries(stripped)) {
    if (key.startsWith("aria-") || key.startsWith("data-")) {
      passthroughProps[key] = value as string | number | boolean | undefined;
    } else {
      styleProps[key] = value;
    }
  }

  const cn = cx(button(vp), css(styleProps as JsxStyleProps), className);

  const isExternal = !(
    href.toString().startsWith("/") || href.toString().startsWith("#")
  );

  const target = isExternal ? "_blank" : undefined;

  return (
    <NextLink
      className={cn}
      href={href}
      target={target}
      onClick={onClick}
      {...passthroughProps}
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
