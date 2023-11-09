"use client";

import { contains, map } from "lodash/fp";

import {
  collectionAddPost,
  collectionRemovePost,
} from "src/api/openapi/collections";
import {
  Collection,
  CollectionList,
  ThreadReference,
} from "src/api/openapi/schemas";
import { useFeedState } from "src/components/feed/useFeedState";

export type Props = {
  thread: ThreadReference;
  initialCollections: CollectionList;
  multiSelect: boolean;
};

type CollectionState = Collection & {
  hasPost: boolean;
};

const hydrateState = (
  collections: CollectionList,
  threadCollections: CollectionList,
) => {
  const cids = threadCollections.map((v) => v.id);

  const m = map((c: Collection) => ({
    ...c,
    hasPost: contains(c.id)(cids),
  }));

  return m(collections);
};

export function useCollectionMenuItems({
  thread,
  initialCollections,
  multiSelect,
}: Props) {
  const { mutate } = useFeedState();
  const collections = hydrateState(initialCollections, thread.collections);

  const onSelect = (c: CollectionState) => async () => {
    if (thread.collections.find((v) => v.id === c.id)) {
      await collectionRemovePost(c.id, thread.id);
    } else {
      await collectionAddPost(c.id, thread.id);
    }

    await mutate?.();
  };

  return {
    ready: true as const,
    collections,
    onSelect,
    multiSelect,
  };
}
