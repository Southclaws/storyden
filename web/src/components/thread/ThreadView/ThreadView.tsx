import { Thread } from "src/api/openapi/schemas";

import { PostListView } from "../PostList";
import { Reply } from "../Reply/Reply";
import { Title } from "../Title/Title";
import { ThreadScreenContext } from "../context/context";
import { useThreadScreenState } from "../context/state";

import { VStack } from "@/styled-system/jsx";

export function ThreadView(props: Thread) {
  const state = useThreadScreenState(props);

  return (
    <ThreadScreenContext.Provider value={state}>
      <VStack alignItems="start" gap={2} width="full">
        <Title {...props} />

        <PostListView {...props} />

        <Reply {...props} />
      </VStack>
    </ThreadScreenContext.Provider>
  );
}
