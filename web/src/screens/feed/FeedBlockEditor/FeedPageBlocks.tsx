import dynamic from "next/dynamic";

import { LStack } from "@/styled-system/jsx";

import { useFeedBlockEditor } from "./Context";
import { FeedBlockRender } from "./FeedBlockRender";

const FeedPageEditor = dynamic(
  () => import("./FeedPageEditor").then((module) => module.FeedPageEditor),
  { ssr: true, loading: () => <FeedPageView /> },
);

export function FeedPageBlocks() {
  const { isEditing } = useFeedBlockEditor();
  return isEditing ? <FeedPageEditor /> : <FeedPageView />;
}

function FeedPageView() {
  const { feed } = useFeedBlockEditor();
  return (
    <LStack className="feed-page__blocks" gap="4" width="full">
      {feed.blocks.map((block) => (
        <FeedBlockRender key={block.type} block={block} />
      ))}
    </LStack>
  );
}
