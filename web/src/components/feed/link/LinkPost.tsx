import NextLink from "next/link";

import { Link, ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";
import { Heading, LinkBox, LinkOverlay } from "src/theme/components";

import { FeedItemMenu } from "../common/FeedItemMenu/FeedItemMenu";

import { Box, Flex, HStack, VStack, styled } from "@/styled-system/jsx";

type Props = {
  thread: ThreadReference;
};

export function LinkPost(props: Props) {
  const session = useSession();

  const permalink = `/t/${props.thread.slug}`;
  const link = props.thread.link as Link;
  const asset = link.assets?.[0] ?? props.thread.assets?.[0];

  return (
    <LinkBox>
      <styled.article
        display="flex"
        flexDir="column"
        width="full"
        boxShadow="md"
        borderRadius="md"
        backgroundColor="white"
        overflow="hidden"
      >
        <Box
          display="flex"
          w="full"
          bgColor="accent.100"
          overflow="hidden"
          height="24"
          p="2"
          gap="2"
        >
          {asset && (
            <Box flexGrow="1" flexShrink="0" width="32">
              <styled.img
                src={asset.url}
                height="full"
                width="full"
                objectPosition="left"
                objectFit="cover"
              />
            </Box>
          )}

          <VStack
            w="full"
            alignItems="start"
            justifyContent="space-evenly"
            gap="0"
            p="2"
          >
            <Flex width="full" justifyContent="space-between">
              <Heading size="sm">
                <LinkOverlay as={NextLink} href={permalink}>
                  {props.thread.title}
                </LinkOverlay>
              </Heading>
            </Flex>

            <styled.p lineClamp={2}>
              <span>{props.thread.short}</span>
              &nbsp;•&nbsp;
              <styled.span color="gray.500">{link.description}</styled.span>
            </styled.p>
          </VStack>
        </Box>

        <Flex justifyContent="space-between" p="2">
          <Byline
            href={permalink}
            author={props.thread.author}
            time={new Date(props.thread.createdAt)}
            updated={new Date(props.thread.updatedAt)}
          />

          <HStack>
            {session && <CollectionMenu thread={props.thread} />}
            <FeedItemMenu {...props.thread} />
          </HStack>
        </Flex>
      </styled.article>
    </LinkBox>
  );
}
