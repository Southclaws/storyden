import { Box, Image } from "@chakra-ui/react";
import { PropsWithChildren } from "react";
import { RenderElementProps } from "slate-react";

export function Element({
  attributes,
  children,
  element,
}: PropsWithChildren<RenderElementProps>) {
  switch (element.type) {
    case "paragraph":
      return (
        <Box as="p" mb={3}>
          {children}
        </Box>
      );

    case "image":
      return (
        <Box>
          <Image src={element.link} alt="" />
          {children}
        </Box>
      );

    default:
      return (
        <Box as="p" mb={3} {...attributes}>
          {children}
        </Box>
      );
  }
}
