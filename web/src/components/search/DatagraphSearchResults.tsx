import { DatagraphItem, DatagraphSearchResult } from "src/api/openapi/schemas";
import { EmptyState } from "src/components/feed/EmptyState";

import { Heading } from "@/components/ui/heading";
import { Box, Flex, LinkOverlay, styled } from "@/styled-system/jsx";

import { FeedItem } from "../feed/common/FeedItem/FeedItem";

type Props = {
  result: DatagraphSearchResult;
};

export function DatagraphSearchResults({ result }: Props) {
  if (!result.items?.length) {
    return <EmptyState />;
  }

  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap="4">
      {result.items.map((v) => (
        <DatagraphResultItem key={v.id} {...v} />
      ))}
    </styled.ol>
  );
}

export function DatagraphResultItem(props: DatagraphItem) {
  const permalink = buildPermalink(props);

  return (
    <Box position="relative">
      <FeedItem>
        <Flex justifyContent="space-between">
          <Heading size="sm">
            <LinkOverlay
              //as={NextLink} // TODO
              href={permalink}
            >
              {props.name}
            </LinkOverlay>
          </Heading>
        </Flex>

        <styled.p lineClamp={3}>{props.description}</styled.p>
        <styled.p lineClamp={3}>{props.kind}</styled.p>
      </FeedItem>
    </Box>
  );
}

function buildPermalink(d: DatagraphItem): string {
  switch (d.kind) {
    case "post":
      return `/t/${d.slug}`;
    case "profile":
      return `/t/${d.slug}`;
    case "node":
      return `/directory/${d.slug}`;
  }
}
