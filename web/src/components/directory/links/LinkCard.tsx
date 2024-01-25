import { ArrowTopRightOnSquareIcon } from "@heroicons/react/24/solid";

import { Link as LinkSchema } from "src/api/openapi/schemas";
import { Empty } from "src/components/site/Empty";
import { Link } from "src/theme/components/Link";

import {
  Box,
  Flex,
  LinkBox,
  LinkOverlay,
  VStack,
  styled,
} from "@/styled-system/jsx";

export function LinkCard(props: LinkSchema) {
  const asset = props.assets?.[0] ?? undefined;
  const domainSearch = `/l?q=${props.domain}`;

  return (
    <styled.article
      display="flex"
      flexDir="column"
      w="full"
      overflow="hidden"
      boxShadow="md"
      borderRadius="lg"
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
      <Flex bgColor="gray.50" p="1" w="full" justify="space-between">
        <styled.h2 color="gray.500" pl="1">
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

        <Link flexShrink="0" kind="ghost" size="xs" href={domainSearch}>
          {props.domain}
        </Link>
      </Flex>

      <LinkBox display="flex" gap="0" maxH="24">
        <Box flexGrow="1" flexShrink="0" width="32">
          {asset ? (
            <styled.img
              src={asset.url}
              height="full"
              width="full"
              objectPosition="center"
              objectFit="cover"
            />
          ) : (
            <VStack justify="center" w="full" h="full">
              <Empty />
            </VStack>
          )}
        </Box>

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
            {props.description || "(No description)"} <br />
            <br />
          </styled.p>
        </styled.div>
      </LinkBox>
    </styled.article>
  );
}
