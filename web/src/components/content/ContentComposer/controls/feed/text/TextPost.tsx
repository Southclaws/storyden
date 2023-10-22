import { Heading, LinkBox, LinkOverlay } from "@chakra-ui/react";
import NextLink from "next/link";

import { ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";
import { Byline } from "src/components/feed/common/Byline";
import { FeedItemMenu } from "src/components/feed/common/FeedItemMenu/FeedItemMenu";

import { Flex, HStack, styled } from "@/styled-system/jsx";

type Props = {
  thread: ThreadReference;
};

export function TextPost(props: Props) {
  const session = useSession();
  const permalink = `/t/${props.thread.slug}`;

  return (
    <LinkBox>
      <styled.article
        display="flex"
        flexDir="column"
        width="full"
        py={2}
        gap={2}
      >
        <Flex justifyContent="space-between">
          <Heading size="sm">
            <LinkOverlay as={NextLink} href={permalink}>
              {props.thread.title}
            </LinkOverlay>
          </Heading>
        </Flex>

        <styled.p maxLines={3}>{props.thread.short}</styled.p>

        <Flex justifyContent="space-between">
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
