import { ThreadReference } from "src/api/openapi/schemas";
import { Anchor } from "src/components/site/Anchor";
import { Heading1 } from "src/theme/components/Heading/Index";

import { FeedItem } from "../common/FeedItem/FeedItem";
import { FeedItemByline } from "../common/FeedItemByline/FeedItemByline";

import { Flex, styled } from "@/styled-system/jsx";

type Props = {
  thread: ThreadReference;
  onDelete?: () => void;
};

export function TextPost(props: Props) {
  const permalink = `/t/${props.thread.slug}`;

  return (
    <FeedItem>
      <Flex justifyContent="space-between">
        <Heading1 size="sm">
          <Anchor href={permalink}>{props.thread.title}</Anchor>
        </Heading1>
      </Flex>

      {/* Suggestion from Jonas: do we actually need a short text preview? */}
      {/* <styled.p lineClamp={1}>{props.thread.short}</styled.p> */}

      <FeedItemByline thread={props.thread} onDelete={props.onDelete} />
    </FeedItem>
  );
}
