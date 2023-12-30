import { Link as LinkSchema, ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";
import { Anchor } from "src/components/site/Anchor";
import { Heading1 } from "src/theme/components/Heading/Index";
import { Link } from "src/theme/components/Link";

import { FeedItemByline } from "../common/FeedItemByline/FeedItemByline";
import { FeedItemMenu } from "../common/FeedItemMenu/FeedItemMenu";
import { Empty } from "../common/PostRef/Empty";

import { Box, Flex, HStack, VStack, styled } from "@/styled-system/jsx";
import { Card } from "@/styled-system/patterns";

type Props = {
  thread: ThreadReference;
  onDelete: () => void;
};

export function LinkPost(props: Props) {
  const session = useSession();

  const permalink = `/t/${props.thread.slug}`;
  const link = props.thread.link as LinkSchema;
  const asset = link.assets?.[0] ?? props.thread.assets?.[0];

  return (
    <styled.article className={Card({ kind: "edge" })}>
      <Box display="flex" w="full" height="16">
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

        <VStack
          w="full"
          minW="0"
          alignItems="start"
          justifyContent="space-evenly"
          gap="0"
          p="2"
        >
          <Flex width="full" justifyContent="space-between">
            <Heading1 size="sm" lineClamp={1}>
              <Anchor href={permalink}>{props.thread.title}</Anchor>
            </Heading1>
          </Flex>

          <styled.p
            w="full"
            color="fg.muted"
            overflow="hidden"
            textOverflow="ellipsis"
            textWrap="nowrap"
          >
            <Anchor href={link.url}>{link.url}</Anchor>
          </styled.p>
        </VStack>
      </Box>

      <Box px="2" pb="2">
        <FeedItemByline thread={props.thread} onDelete={props.onDelete} />
      </Box>
    </styled.article>
  );
}
