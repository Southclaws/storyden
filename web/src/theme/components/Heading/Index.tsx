import { PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import { HeadingVariantProps, heading } from "@/styled-system/recipes";
import { StyleProps } from "@/styled-system/types";

type HeadingProps = PropsWithChildren<HeadingVariantProps & StyleProps>;

export function Heading1(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = heading.splitVariantProps(rest);

  return (
    <h1 className={cx(heading({ size: "2xl", ...recipeProps }), css(cssProps))}>
      {children}
    </h1>
  );
}

export function Heading2(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = heading.splitVariantProps(rest);

  return (
    <h2 className={cx(heading({ size: "xl", ...recipeProps }), css(cssProps))}>
      {children}
    </h2>
  );
}

export function Heading3(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = heading.splitVariantProps(rest);

  return (
    <h3 className={cx(heading({ size: "lg", ...recipeProps }), css(cssProps))}>
      {children}
    </h3>
  );
}

export function Heading4(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = heading.splitVariantProps(rest);

  return (
    <h4 className={cx(heading({ size: "md", ...recipeProps }), css(cssProps))}>
      {children}
    </h4>
  );
}

export function Heading5(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = heading.splitVariantProps(rest);

  return (
    <h5 className={cx(heading({ size: "sm", ...recipeProps }), css(cssProps))}>
      {children}
    </h5>
  );
}

export function Heading6(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = heading.splitVariantProps(rest);

  return (
    <h6 className={cx(heading({ size: "xs", ...recipeProps }), css(cssProps))}>
      {children}
    </h6>
  );
}
