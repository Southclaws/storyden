import { Link as ChakraLink } from "@chakra-ui/react";
import NextLink from "next/link";

export default function Link({ children, href, ...rest }) {
  return (
    <ChakraLink {...rest}>
      <NextLink href={href} {...rest}>
        {children}
      </NextLink>
    </ChakraLink>
  );
}
