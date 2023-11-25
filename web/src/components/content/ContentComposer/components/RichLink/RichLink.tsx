import { PropsWithChildren } from "react";
import { useFocused, useSelected } from "slate-react";

import { Box, styled } from "@/styled-system/jsx";

import { Props, useRichLink } from "./useRichLink";

export function RichLink(props: PropsWithChildren<Props>) {
  const selected = useSelected();
  const focused = useFocused();
  const { link } = useRichLink(props);

  if (!link) {
    return (
      <styled.a bgColor="accent.100" borderRadius="sm" px="1" href={props.href}>
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
      w="full"
      bgColor="accent.100"
      borderRadius="md"
      overflow="hidden"
      outlineStyle="solid"
      outlineWidth="medium"
      outlineColor="accent.100"
      mb="2"
      css={{
        "&[data-selected=true]": {
          outlineStyle: "dashed",
          outlineOffset: "-0.5",
          outlineWidth: "medium",
          outlineColor: "accent.200",
        },
      }}
    >
      <styled.span
        color="gray.600"
        lineClamp={1}
        contentEditable={true}
        suppressContentEditableWarning
      >
        {props.children}
      </styled.span>

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
        w="full"
        justifyContent="space-evenly"
        gap="0"
        p="2"
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
    </styled.article>
  );
}
