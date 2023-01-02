import { Heading, VStack } from "@chakra-ui/react";
import { Thread } from "src/api/openapi/schemas";
import { CategoryPill } from "src/components/CategoryPill";
import { PostListView } from "./PostList";

export function ThreadView(props: Thread) {
  return (
    <VStack alignItems="start">
      <Heading>{props.title}</Heading>
      <CategoryPill category={props.category} />

      <PostListView {...props} />
    </VStack>
  );
}
