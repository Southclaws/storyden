import { Thread } from "src/api/openapi-schema";

import { Byline } from "@/components/content/Byline";
import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";
import { HStack, VStack } from "@/styled-system/jsx";

import { CategoryPill } from "../CategoryPill";
import { PostListView } from "../PostList";
import { Reply } from "../Reply/Reply";
import { Title } from "../Title/Title";
import { ThreadScreenContext } from "../context/context";
import { useThreadScreenState } from "../context/state";

export function ThreadView(props: Thread) {
  const state = useThreadScreenState(props);

  return (
    <ThreadScreenContext.Provider value={state}>
      <VStack alignItems="start" gap="4" width="full">
        <Title {...props} />

        <HStack w="full">
          <Byline
            href={`#${props.id}`}
            author={props.author}
            time={new Date(props.createdAt)}
            updated={new Date(props.updatedAt)}
          />

          <CategoryPill category={props.category} />
        </HStack>
        {state.editingContent ? (
          <>
            {/* <ContentComposer
            onChange={onContentChange}
            initialValue={editingContent}
          />
          <HStack>
            <Button onClick={onPublishEdit}>Update</Button>
            <Button variant="ghost" onClick={onCancelEdit}>
              Cancel
            </Button>
          </HStack> */}
          </>
        ) : (
          <>
            <ContentComposer initialValue={props.body} disabled />
          </>
        )}

        <PostListView posts={props.replies} />

        <Reply {...props} />
      </VStack>
    </ThreadScreenContext.Provider>
  );
}
