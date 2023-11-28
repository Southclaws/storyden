import NextLink from "next/link";

import { ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";
import { Heading, LinkBox, LinkOverlay } from "src/theme/components";

import { FeedItem } from "../common/FeedItem/FeedItem";
import { FeedItemMenu } from "../common/FeedItemMenu/FeedItemMenu";

import { Flex, HStack, styled } from "@/styled-system/jsx";

type Props = {
  thread: ThreadReference;
  onDelete?: () => void;
};

export function TextPost(props: Props) {
  const session = useSession();
  const permalink = `/t/${props.thread.slug}`;

  return (
    <LinkBox>
      <FeedItem>
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
            {props.onDelete && (
              <FeedItemMenu thread={props.thread} onDelete={props.onDelete} />
            )}
          </HStack>
        </Flex>
      </FeedItem>
    </LinkBox>
  );
}
