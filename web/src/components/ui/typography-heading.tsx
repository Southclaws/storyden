import { DetailedHTMLProps, HTMLAttributes, PropsWithChildren } from "react";

import { css, cx } from "@/styled-system/css";
import {
  TypographyHeadingVariantProps,
  typographyHeading,
} from "@/styled-system/recipes";
import { JsxStyleProps } from "@/styled-system/types";

type HeadingProps = PropsWithChildren<
  TypographyHeadingVariantProps &
    JsxStyleProps &
    DetailedHTMLProps<HTMLAttributes<HTMLHeadingElement>, HTMLHeadingElement>
>;

export function Heading1(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = typographyHeading.splitVariantProps(rest);

  return (
    <h1
      className={cx(
        typographyHeading({ size: "2xl", ...recipeProps }),
        css(cssProps),
      )}
    >
      {children}
    </h1>
  );
}

export function Heading2(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = typographyHeading.splitVariantProps(rest);

  return (
    <h2
      className={cx(
        typographyHeading({ size: "xl", ...recipeProps }),
        css(cssProps),
      )}
    >
      {children}
    </h2>
  );
}

export function Heading3(props: HeadingProps) {
  const { children, className, ...rest } = props;
  const [recipeProps, cssProps] = typographyHeading.splitVariantProps(rest);

  return (
    <h3
      className={cx(
        typographyHeading({ size: "lg", ...recipeProps }),
        css(cssProps),
        className,
      )}
    >
      {children}
    </h3>
  );
}

export function Heading4(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = typographyHeading.splitVariantProps(rest);

  return (
    <h4
      className={cx(
        typographyHeading({ size: "md", ...recipeProps }),
        css(cssProps),
      )}
    >
      {children}
    </h4>
  );
}

export function Heading5(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = typographyHeading.splitVariantProps(rest);

  return (
    <h5
      className={cx(
        typographyHeading({ size: "sm", ...recipeProps }),
        css(cssProps),
      )}
    >
      {children}
    </h5>
  );
}

export function Heading6(props: HeadingProps) {
  const { children, ...rest } = props;
  const [recipeProps, cssProps] = typographyHeading.splitVariantProps(rest);

  return (
    <h6
      className={cx(
        typographyHeading({ size: "xs", ...recipeProps }),
        css(cssProps),
      )}
    >
      {children}
    </h6>
  );
}
