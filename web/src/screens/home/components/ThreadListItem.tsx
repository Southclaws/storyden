import Link from "next/link";

import { ThreadReference } from "src/api/openapi/schemas";
import { CollectionMenu } from "src/components/CollectionMenu/CollectionMenu";
import { Byline } from "src/screens/thread/components/Byline";

import { Flex, HStack, styled } from "@/styled-system/jsx";

import { ThreadMenu } from "./ThreadMenu/ThreadMenu";

export function ThreadListItem(props: { thread: ThreadReference }) {
  const permalink = `/t/${props.thread.slug}`;

  return (
    <styled.section display="flex" flexDir="column" py={2} width="full" gap={2}>
      <article>
        <Flex justifyContent="space-between">
          <styled.h1 fontSize="sm">
            <Link href={permalink}>{props.thread.title}</Link>
          </styled.h1>
        </Flex>

        <styled.p //</article>noOfLines={3}
        >
          {props.thread.short}
        </styled.p>
      </article>

      <Flex justifyContent="space-between">
        <Byline
          href={permalink}
          author={props.thread.author}
          time={new Date(props.thread.createdAt)}
          updated={new Date(props.thread.updatedAt)}
        />

        {/* Tags list */}

        <HStack>
          <CollectionMenu thread={props.thread} />
          <ThreadMenu {...props.thread} />
        </HStack>
      </Flex>
    </styled.section>
  );
}
