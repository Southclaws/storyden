import { Flex, Heading, LinkBox, LinkOverlay, Text } from "@chakra-ui/react";
import { ThreadReference } from "src/api/openapi/schemas";
import { Byline } from "src/screens/thread/components/Byline";
import { ThreadMenu } from "./ThreadMenu/ThreadMenu";

export function ThreadListItem(props: { thread: ThreadReference }) {
  const permalink = `/t/${props.thread.slug}`;

  return (
    <Flex as="section" flexDir="column" py={2} width="full">
      <LinkBox>
        <Flex justifyContent="space-between">
          <Heading size="sm">
            <LinkOverlay href={permalink}>{props.thread.title}</LinkOverlay>
          </Heading>
        </Flex>

        <Text noOfLines={3}>{props.thread.short}</Text>
      </LinkBox>

      <Flex justifyContent="space-between">
        <Byline
          href={permalink}
          author={props.thread.author.handle}
          time={new Date(props.thread.createdAt)}
          updated={new Date(props.thread.updatedAt)}
          more={<ThreadMenu {...props.thread} />}
        />

        {/* Tags list */}
      </Flex>
    </Flex>
  );
}
