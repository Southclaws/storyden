import { type ElementType, forwardRef } from "react";

import { Text, type TextProps } from "../text";

export type HeadingSize = "xs" | "sm" | "md" | "lg" | "xl";

export interface HeadingProps extends Omit<TextProps<ElementType>, "size"> {
  size?: HeadingSize;
}

export const Heading = forwardRef<HTMLHeadingElement, HeadingProps>(
  ({ size = "md", ...props }, ref) => (
    <Text as="h2" variant="heading" size={size} ref={ref} {...props} />
  ),
);

Heading.displayName = "Heading";
