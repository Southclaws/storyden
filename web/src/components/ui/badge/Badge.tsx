import { forwardRef } from "react";

import { StyledBadge } from "./Badge.internal";
import type { StyledBadgeProps } from "./Badge.internal";

export type BadgeProps = StyledBadgeProps;

export const Badge = forwardRef<HTMLDivElement, BadgeProps>((props, ref) => {
  const colorPalette =
    props.colorPalette ?? (props.variant === "solid" ? "accent" : "gray");

  return <StyledBadge ref={ref} {...props} colorPalette={colorPalette} />;
});

Badge.displayName = "Badge";
