import { Thread } from "src/api/openapi-schema";

import { VStack } from "@/styled-system/jsx";

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

        <PostListView posts={props.replies} />

        <Reply {...props} />
      </VStack>
    </ThreadScreenContext.Provider>
  );
}
