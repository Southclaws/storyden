import { useEffect } from "react";

import { FeedBlockType } from "@/lib/settings/feed";
import { createEmitter } from "@/utils/emitter";

export type FeedBlockEvents = {
  "feed:reorder-block": {
    activeId: FeedBlockType;
    overId: FeedBlockType;
  };
};

export const feedBlockBus = createEmitter<FeedBlockEvents>();

export function useFeedBlockEvent<K extends keyof FeedBlockEvents>(
  type: K,
  handler: (event: FeedBlockEvents[K]) => void,
) {
  useEffect(() => {
    feedBlockBus.on(type, handler);
    return () => feedBlockBus.off(type, handler);
  }, [handler, type]);
}

export function useEmitFeedBlockEvent() {
  return <K extends keyof FeedBlockEvents>(
    type: K,
    payload: FeedBlockEvents[K],
  ) => feedBlockBus.emit(type, payload);
}
