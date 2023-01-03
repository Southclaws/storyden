import { Flex, Heading, LinkBox, LinkOverlay, Text } from "@chakra-ui/react";
import { ThreadReference } from "src/api/openapi/schemas";
import { Byline } from "src/screens/thread/components/Byline";

export function ThreadListItem(props: { thread: ThreadReference }) {
  return (
    <Flex as="section" flexDir="column" py={2} width="full">
      <LinkBox>
        <Flex justifyContent="space-between">
          <Heading size="sm">
            <LinkOverlay href={`/t/${props.thread.slug}`}>
              {props.thread.title}
            </LinkOverlay>
          </Heading>
          {/* Options menu */}
        </Flex>

        <Text noOfLines={3}>{props.thread.short}</Text>
      </LinkBox>

      <Flex justifyContent="space-between">
        <Byline
          author={props.thread.author.handle}
          time={new Date(props.thread.createdAt)}
        />

        {/* Tags list */}
      </Flex>
    </Flex>
  );
}
