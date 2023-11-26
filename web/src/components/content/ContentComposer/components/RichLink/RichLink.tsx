import { PropsWithChildren } from "react";
import { RenderElementProps, useFocused, useSelected } from "slate-react";

import { Box, Flex, styled } from "@/styled-system/jsx";

import { Props, useRichLink } from "./useRichLink";

export function RichLink(props: PropsWithChildren<Props & RenderElementProps>) {
  const selected = useSelected();
  const focused = useFocused();
  const { link } = useRichLink(props);

  if (!link) {
    return (
      <styled.a
        bgColor="zinc.100"
        borderRadius="sm"
        color="gray.500"
        lineClamp={1}
        px="1"
        contentEditable={true}
        suppressContentEditableWarning
        href={props.href}
        {...props.attributes}
      >
        {props.href} {props.children}
      </styled.a>
    );
  }

  const asset = link.assets?.[0] ?? undefined;

  return (
    <styled.article
      contentEditable={false}
      data-selected={selected && focused}
      display="flex"
      flexDir="column"
      gap="1"
      w="full"
      bgColor="zinc.100"
      borderRadius="md"
      overflow="hidden"
      outlineStyle="solid"
      outlineWidth="medium"
      outlineColor="gray.100"
      mb="2"
      css={{
        "&[data-selected=true]": {
          outlineStyle: "dashed",
          outlineOffset: "-0.5",
          outlineWidth: "medium",
          outlineColor: "accent.200",
        },
      }}
      {...props.attributes}
    >
      <styled.span
        color="gray.500"
        lineClamp={1}
        px="1"
        contentEditable={true}
        suppressContentEditableWarning
      >
        {props.href} {props.children}
      </styled.span>

      <Flex gap="2">
        {asset && (
          <Box flexGrow="1" flexShrink="0" width="32">
            <styled.img
              src={asset.url}
              height="full"
              width="full"
              objectPosition="center"
              objectFit="cover"
            />
          </Box>
        )}

        <styled.p
          display="flex"
          flexDir="column"
          justifyContent="space-evenly"
          alignItems="start"
          w="full"
          h="full"
          gap="1"
        >
          <styled.span lineClamp={1} fontSize="md" fontWeight="bold">
            <styled.a href={link.url} target="=_blank">
              {link.title || link.url}
            </styled.a>
          </styled.span>

          <styled.span lineClamp={2}>
            {link.description || "(No description)"} <br />
            <br />
          </styled.span>
        </styled.p>
      </Flex>
    </styled.article>
  );
}
