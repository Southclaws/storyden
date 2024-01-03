import { DatagraphItem, DatagraphSearchResult } from "src/api/openapi/schemas";
import { EmptyState } from "src/components/feed/EmptyState";
import { Heading1 } from "src/theme/components/Heading/Index";

import { FeedItem } from "../feed/common/FeedItem/FeedItem";

import { Flex, LinkBox, LinkOverlay, styled } from "@/styled-system/jsx";

type Props = {
  result: DatagraphSearchResult;
};

export function DatagraphSearchResults({ result }: Props) {
  if (!result.items?.length) {
    return <EmptyState />;
  }

  console.log(result);

  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap="4">
      {result.items.map((v) => (
        <DatagraphResultItem key={v.id} {...v} />
      ))}
    </styled.ol>
  );
}

export function DatagraphResultItem(props: DatagraphItem) {
  const permalink = `/TODO`;

  return (
    <LinkBox>
      <FeedItem>
        <Flex justifyContent="space-between">
          <Heading1 size="sm">
            <LinkOverlay
              //as={NextLink} // TODO
              href={permalink}
            >
              {props.name}
            </LinkOverlay>
          </Heading1>
        </Flex>

        <styled.p lineClamp={3}>{props.description}</styled.p>
        <styled.p lineClamp={3}>{props.type}</styled.p>
      </FeedItem>
    </LinkBox>
  );
}
