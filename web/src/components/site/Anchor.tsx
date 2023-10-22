import Link from "next/link";
import { AnchorHTMLAttributes, PropsWithRef, useCallback } from "react";

import { css, cx } from "@/styled-system/css";

export function Anchor({
  children,
  className,
  onClick,
  href,
  ...rest
}: PropsWithRef<AnchorHTMLAttributes<HTMLAnchorElement>>) {
  // This allows us to progressively enhance features on the application by
  // treating important buttons as links to fallback pages. For example, there
  // may be a button that triggers the opening of a modal dialogue but if the
  // user has JavaScript disabled due to device constraints or privacy reasons,
  // the functionality must also be implemented by a normal page.
  const handleClick = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      if (onClick) {
        e.preventDefault();
        return onClick?.(e);
      }
    },
    [onClick],
  );

  return (
    <Link
      className={cx(
        css({
          _hover: { textDecoration: "underline" },
        }),
        className,
      )}
      href={href!}
      onClick={handleClick}
      {...rest}
    >
      {children}
    </Link>
  );
}
