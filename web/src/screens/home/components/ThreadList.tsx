import { Box, Flex, Heading, LinkOverlay, Text } from "@chakra-ui/react";
import { formatDistanceToNow } from "date-fns";
import { ThreadReference } from "src/api/openapi/schemas";

type Props = { threads: ThreadReference[] };
export function ThreadList(props: Props) {
  const children = props.threads.map((t) => (
    <ThreadListItem key={t.id} thread={t} />
  ));

  return <Box as="main">{children}</Box>;
}

export function ThreadListItem(props: { thread: ThreadReference }) {
  return (
    <Flex as="section" flexDir="column" px={4} py={2} width="full">
      <LinkOverlay href={`/${props.thread.slug}`}>
        <Flex justifyContent="space-between">
          <Heading size="sm">{props.thread.title}</Heading>
          {/* Options menu */}
        </Flex>

        <Text noOfLines={3}>{props.thread.short}</Text>
      </LinkOverlay>

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
