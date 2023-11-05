"use client";

import { intersection } from "lodash";
import { contains, filter, find, map } from "lodash/fp";
import { MouseEvent, useState } from "react";

import {
  collectionAddPost,
  collectionRemovePost,
} from "src/api/openapi/collections";
import {
  Account,
  Collection,
  CollectionList,
  ThreadReference,
} from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { useFeed } from "src/components/feed/useFeed";
import { useFeedState } from "src/components/feed/useFeedState";

export type Props = {
  thread: ThreadReference;
  initialCollections: CollectionList;
  multiSelect: boolean;
};

type CollectionState = Collection & {
  hasPost: boolean;
};

const hasCollection = (collections: CollectionList, account?: Account) => {
  const f = filter((c: Collection) => c.owner.id === account?.id);
  const l = f(collections);
  return l.length > 0;
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

  const onSelect =
    (c: CollectionState) => async (e: MouseEvent<HTMLDivElement>) => {
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
