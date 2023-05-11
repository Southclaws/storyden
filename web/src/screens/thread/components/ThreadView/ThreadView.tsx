import { Editable, Heading, Input, VStack } from "@chakra-ui/react";

import { Thread } from "src/api/openapi/schemas";
import { CategoryPill } from "src/components/CategoryPill";

import { PostListView } from "../PostList";
import { Reply } from "../Reply/Reply";
import { useThreadView } from "./useThreadView";

export function ThreadView(props: Thread) {
  const { editing, editingTitle, onTitleChange } = useThreadView(props);
  return (
    <VStack alignItems="start" gap={2} py={4} width="full">
      {editing ? (
        <Input value={editingTitle} onChange={onTitleChange} />
      ) : (
        <Heading>{props.title}</Heading>
      )}
      <CategoryPill category={props.category} />

      <PostListView {...props} />

      <Reply {...props} />
    </VStack>
  );
}
