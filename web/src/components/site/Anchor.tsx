import { Link as ChakraLink, LinkProps } from "@chakra-ui/react";
import NextLink from "next/link";
import { useCallback } from "react";

export function Anchor({ children, href, onClick, ...rest }: LinkProps) {
  const handleClick = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      if (onClick) {
        e.preventDefault();
        return onClick?.(e);
      }
    },
    [onClick]
  );

  return (
    <NextLink
      // we use legacy behaviour since it's the only thing that
      // seems to make it compatible with Chakra's link component.
      legacyBehavior
      href={href ?? ""}
    >
      <ChakraLink onClick={handleClick} {...rest}>
        {children}
      </ChakraLink>
    </NextLink>
  );
}
