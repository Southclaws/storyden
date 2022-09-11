import { Link as ChakraLink, LinkProps } from "@chakra-ui/react";
import NextLink from "next/link";

export default function Link({ children, href, ...rest }: LinkProps) {
  return (
    <ChakraLink {...rest}>
      <NextLink href={href!}>{children}</NextLink>
    </ChakraLink>
  );
}
