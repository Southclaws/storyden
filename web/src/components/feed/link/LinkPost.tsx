import { Link as LinkSchema, ThreadReference } from "src/api/openapi-schema";
import { Anchor } from "src/components/site/Anchor";

import { Empty } from "@/components/site/Empty";
import { Heading } from "@/components/ui/heading";
import { Box, Flex, VStack, styled } from "@/styled-system/jsx";
import { CardBox } from "@/styled-system/patterns";
import { getAssetURL } from "@/utils/asset";

import { FeedItemByline } from "../FeedItemByline/FeedItemByline";

type Props = {
  thread: ThreadReference;
};

export function LinkPost(props: Props) {
  const permalink = `/t/${props.thread.slug}`;
  const link = props.thread.link as LinkSchema;
  const asset = link.assets?.[0] ?? props.thread.assets?.[0];

  return (
    <styled.article className={CardBox({ kind: "edge" })}>
      <Box display="flex" w="full" height="16">
        <Box flexGrow="1" flexShrink="0" width="32">
          {asset ? (
            <styled.img
              src={getAssetURL(asset.path)}
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

        <VStack
          w="full"
          minW="0"
          alignItems="start"
          justifyContent="space-evenly"
          gap="0"
          p="2"
        >
          <Flex width="full" justifyContent="space-between">
            <Heading size="sm" lineClamp={1}>
              <Anchor href={permalink}>{props.thread.title}</Anchor>
            </Heading>
          </Flex>

          <Box>
            <styled.p lineClamp={1} wordBreak="break-all">
              <Anchor href={link.url}>{link.url}</Anchor>
            </styled.p>
          </Box>
        </VStack>
      </Box>

      <Box px="2" pb="2">
        <FeedItemByline thread={props.thread} />
      </Box>
    </styled.article>
  );
}
