import { Heading, VStack } from "@chakra-ui/react";
import { Thread } from "src/api/openapi/schemas";
import { CategoryPill } from "src/components/CategoryPill";
import { PostListView } from "../PostList";
import { ReplyBox } from "../ReplyBox";
import { useThread } from "./useThread";

export function ThreadView(props: Thread) {
  const { loggedIn, onReply } = useThread(props);
  return (
    <VStack alignItems="start" gap={2} py={4}>
      <Heading>{props.title}</Heading>
      <CategoryPill category={props.category} />

      <PostListView {...props} />

      {loggedIn && <ReplyBox onSave={onReply} />}
    </VStack>
  );
}
