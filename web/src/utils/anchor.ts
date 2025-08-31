import Link from "next/link";
import { ComponentProps, MouseEvent } from "react";

/**
 * Generate props to disable a Next.js Link component.
 */
export function linkDisabledProps(
  disabled: boolean,
): Omit<ComponentProps<typeof Link>, "href"> {
  return disabled
    ? {
        onClick: (e: MouseEvent<HTMLAnchorElement>) => {
          e.preventDefault();
          e.stopPropagation();
        },
        onAuxClick: (e: MouseEvent<HTMLAnchorElement>) => {
          e.preventDefault();
          e.stopPropagation();
        },
        prefetch: false,
        style: {
          pointerEvents: "inherit",
          // Cascade the default cursor (e.g. caret for text) in children.
          cursor: "inherit",
          textDecoration: "none",
          color: "inherit",
        },
        tabIndex: -1,
        "aria-disabled": true,
      }
    : {};
}
