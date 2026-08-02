import { Link as LinkSchema, ThreadReference } from "@/api/openapi-schema";
import { Anchor } from "@/components/site/Anchor";
import { Empty } from "@/components/site/Empty";
import { Text } from "@/components/ui/text";
import { Box, Flex, VStack, styled } from "@/styled-system/jsx";
import { cardBox } from "@/styled-system/recipes";
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
    <styled.article className={cardBox({ kind: "edge" })}>
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
            <styled.h1
              color="text.default"
              fontWeight="semibold"
              fontSize="sm"
              lineHeight="normal"
              lineClamp={1}
            >
              <Anchor href={permalink}>{props.thread.title}</Anchor>
            </styled.h1>
          </Flex>

          <Box>
            <Text variant="supporting" lineClamp={1} wordBreak="break-all">
              <Anchor href={link.url}>{link.url}</Anchor>
            </Text>
          </Box>
        </VStack>
      </Box>

      <Box px="2" pb="2">
        <FeedItemByline thread={props.thread} />
      </Box>
    </styled.article>
  );
}
