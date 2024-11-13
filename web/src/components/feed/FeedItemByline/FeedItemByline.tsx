import { ThreadReference } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { Byline } from "@/components/content/Byline";
import { CollectionMenu } from "@/components/content/CollectionMenu/CollectionMenu";
import { ThreadMenu } from "@/components/thread/ThreadMenu/ThreadMenu";
import { Flex, HStack } from "@/styled-system/jsx";

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
        {session && <CollectionMenu account={session} thread={props.thread} />}
        <ThreadMenu thread={props.thread} />
      </HStack>
    </Flex>
  );
}
