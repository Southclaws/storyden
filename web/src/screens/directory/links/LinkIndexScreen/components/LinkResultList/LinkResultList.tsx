import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/solid";

import { LinkList, Link as LinkSchema } from "src/api/openapi/schemas";

import { Box, LinkBox, LinkOverlay, styled } from "@/styled-system/jsx";

export function LinkResultList(props: { links: LinkList }) {
  return (
    <styled.ol display="flex" flexDir="column" gap="4">
      {props.links.map((v) => (
        <LinkResultListItem key={v.url} {...v} />
      ))}
    </styled.ol>
  );
}

function LinkResultListItem(props: LinkSchema) {
  const asset = props.assets?.[0] ?? undefined;

  return (
    <styled.article
      display="flex"
      flexDir="column"
      w="full"
      overflow="hidden"
      boxShadow="md"
      borderRadius="md"
      backgroundColor="white"
      css={{
        "&[data-selected=true]": {
          outlineStyle: "dashed",
          outlineOffset: "-0.5",
          outlineWidth: "medium",
          outlineColor: "accent.200",
        },
      }}
    >
      <styled.h2 color="gray.500" px="1">
        <styled.a
          display="flex"
          alignItems="center"
          flexWrap="nowrap"
          gap="1"
          _hover={{ textDecoration: "underline" }}
          href={props.url}
          target="=_blank"
        >
          <ArrowTopRightOnSquareIcon height="1em" />
          <styled.span lineClamp={1} wordBreak="break-all">
            {props.url}
          </styled.span>
        </styled.a>
      </styled.h2>

      <LinkBox display="flex" gap="0" maxH="24">
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

        <styled.div
          display="flex"
          flexDir="column"
          justifyContent="space-evenly"
          alignItems="start"
          w="full"
          h="full"
          gap="1"
          p="2"
        >
          <styled.h1 fontSize="md" fontWeight="bold">
            <LinkOverlay
              lineClamp={1}
              wordBreak="break-all"
              href={`/l/${props.slug}`}
            >
              {props.title || props.url}
            </LinkOverlay>
          </styled.h1>

          <styled.p lineClamp={2}>
            <styled.a
              color="gray.500"
              _hover={{ textDecoration: "underline" }}
              href={`/l?q=${props.domain}`}
            >
              {props.domain}
            </styled.a>
            <styled.span>&nbsp;•&nbsp;</styled.span>
            {props.description || "(No description)"} <br />
            <br />
          </styled.p>
        </styled.div>
      </LinkBox>
    </styled.article>
  );
}
