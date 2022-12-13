import { Link as ChakraLink, LinkProps } from "@chakra-ui/react";
import NextLink from "next/link";

export default function Link({ children, href, ...rest }: LinkProps) {
  return (
    <ChakraLink {...rest}>
      <NextLink
        // we use legacy behaviour since it's the only thing that
        // seems to make it compatible with Chakra's link component.
        legacyBehavior
        href={href!}
      >
        {children}
      </NextLink>
    </ChakraLink>
  );
}
