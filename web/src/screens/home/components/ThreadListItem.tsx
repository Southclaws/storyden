import {
  Flex,
  HStack,
  Heading,
  LinkBox,
  LinkOverlay,
  Text,
} from "@chakra-ui/react";
import NextLink from "next/link";

import { ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";
import { Byline } from "src/screens/thread/components/Byline";

import { ThreadMenu } from "./ThreadMenu/ThreadMenu";

export function ThreadListItem(props: { thread: ThreadReference }) {
  const session = useSession();
  const permalink = `/t/${props.thread.slug}`;

  return (
    <Flex as="section" flexDir="column" py={2} width="full" gap={2}>
      <LinkBox as="article">
        <Flex justifyContent="space-between">
          <Heading size="sm">
            <LinkOverlay as={NextLink} href={permalink}>
              {props.thread.title}
            </LinkOverlay>
          </Heading>
        </Flex>

        <Text noOfLines={3}>{props.thread.short}</Text>
      </LinkBox>

      <Flex justifyContent="space-between">
        <Byline
          href={permalink}
          author={props.thread.author}
          time={new Date(props.thread.createdAt)}
          updated={new Date(props.thread.updatedAt)}
        />

        {/* Tags list */}

        <HStack>
          {session && <CollectionMenu thread={props.thread} />}
          <ThreadMenu {...props.thread} />
        </HStack>
      </Flex>
    </Flex>
  );
}
