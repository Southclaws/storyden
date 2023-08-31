import NextLink from "next/link";
import { PropsWithChildren, useCallback } from "react";

export function Anchor({ children, onClick, ...rest }: PropsWithChildren<any>) {
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
    [onClick]
  );

  return <NextLink {...rest}>{children}</NextLink>;
}
