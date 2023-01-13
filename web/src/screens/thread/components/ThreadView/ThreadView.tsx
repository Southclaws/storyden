import { Heading, Spinner, VStack } from "@chakra-ui/react";
import { Thread } from "src/api/openapi/schemas";
import { CategoryPill } from "src/components/CategoryPill";
import { PostListView } from "../PostList";
import { ReplyBox } from "../ReplyBox";
import { useThread } from "./useThread";

export function ThreadView(props: Thread) {
  const { loggedIn, onReply, isLoading } = useThread(props);
  // NOTE: isLoading is a hack to easily reset the ReplyBox + provide feedback.
  return (
    <VStack alignItems="start" gap={2} py={4} width="full">
      <Heading>{props.title}</Heading>
      <CategoryPill category={props.category} />

      <PostListView {...props} />

      {loggedIn &&
        (isLoading ? (
          <VStack width="full" py={6}>
            <Spinner />
          </VStack>
        ) : (
          <ReplyBox onSave={onReply} />
        ))}
    </VStack>
  );
}
