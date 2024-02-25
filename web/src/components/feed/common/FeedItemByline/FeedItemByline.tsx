import { ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { Byline } from "src/components/content/Byline";
import { CollectionMenu } from "src/components/content/CollectionMenu/CollectionMenu";

import { FeedItemMenu } from "../FeedItemMenu/FeedItemMenu";

import { Flex, HStack } from "@/styled-system/jsx";

type Props = {
  thread: ThreadReference;
  onDelete?: () => void;
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
        {props.onDelete && (
          <FeedItemMenu thread={props.thread} onDelete={props.onDelete} />
        )}
      </HStack>
    </Flex>
  );
}
