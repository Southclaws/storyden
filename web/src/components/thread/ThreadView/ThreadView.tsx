import { formatDate } from "date-fns";

import { Thread } from "src/api/openapi-schema";

import { RightNavPortal } from "@/components/site/Navigation/Right/context";
import { VStack, styled } from "@/styled-system/jsx";

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

      <RightNavPortal>
        <styled.h1 fontWeight="bold">{props.title}</styled.h1>
        <p>posted {formatDate(props.createdAt, "PP")}</p>
        <p>{props.replies.length} replies</p>
      </RightNavPortal>
    </ThreadScreenContext.Provider>
  );
}
