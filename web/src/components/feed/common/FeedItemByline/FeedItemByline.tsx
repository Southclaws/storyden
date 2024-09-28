import { ThreadReference } from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";

import { Flex, HStack } from "@/styled-system/jsx";

import { FeedItemMenu } from "../FeedItemMenu/FeedItemMenu";

type Props = {
  thread: ThreadReference;
};

export function FeedItemByline(props: Props) {
  const session = useSession();
  const permalink = `/t/${props.thread.slug}`;

  return (
    <Flex w="full" justifyContent="space-between" gap="2">
      <Byline
        href={permalink}
        author={props.thread.author}
        time={new Date(props.thread.createdAt)}
        updated={new Date(props.thread.updatedAt)}
      />

      <HStack>
        {session && <CollectionMenu thread={props.thread} />}
        <FeedItemMenu thread={props.thread} />
      </HStack>
    </Flex>
  );
}
