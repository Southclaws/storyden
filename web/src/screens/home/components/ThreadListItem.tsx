import { Flex, Heading, LinkBox, LinkOverlay, Text } from "@chakra-ui/react";
import { formatDistanceToNow } from "date-fns";
import { ThreadReference } from "src/api/openapi/schemas";

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
        <Flex gap={2}>
          <Text>{props.thread.author.handle ?? "Unknown"}</Text>
          <Text>•</Text>
          <Text>
            {formatDistanceToNow(new Date(props.thread.createdAt))} ago
          </Text>
        </Flex>

        {/* Tags list */}
      </Flex>
    </Flex>
  );
}
