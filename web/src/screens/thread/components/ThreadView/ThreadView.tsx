import { Heading, VStack } from "@chakra-ui/react";

import { Thread } from "src/api/openapi/schemas";
import { CategoryPill } from "src/components/CategoryPill";

import { PostListView } from "../PostList";
import { Reply } from "../Reply";

export function ThreadView(props: Thread) {
  return (
    <VStack alignItems="start" gap={2} py={4} width="full">
      <Heading>{props.title}</Heading>
      <CategoryPill category={props.category} />

      <PostListView {...props} />

      <Reply {...props} />
    </VStack>
  );
}
