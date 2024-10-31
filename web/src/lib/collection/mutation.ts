import { Arguments, MutatorCallback, useSWRConfig } from "swr";

import {
  collectionAddNode,
  collectionAddPost,
  collectionRemoveNode,
  collectionRemovePost,
  getCollectionListKey,
} from "@/api/openapi-client/collections";
import { getThreadListKey } from "@/api/openapi-client/threads";
import {
  Account,
  Collection,
  CollectionListOKResponse,
  ThreadListOKResponse,
} from "@/api/openapi-schema";

import { useFeedMutations } from "../feed/mutation";

export function useCollectionItemMutations(session: Account) {
  const { mutate } = useSWRConfig();
  const { revalidate: revalidateFeed } = useFeedMutations();

  const threadQueryMutationKey = getThreadListKey()[0];
  function threadListKeyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0].startsWith(threadQueryMutationKey);
  }

  const collectionListAnyMutationKey = getCollectionListKey()[0];
  function collectionListAnyKeyFilterFn(key: Arguments) {
    return (
      Array.isArray(key) && key[0].startsWith(collectionListAnyMutationKey)
    );
  }

  const addPost = async (collection: Collection, postID: string) => {
    const collectionQueryMutationKey = getCollectionListKey({
      has_item: postID,
      account_handle: session.handle,
    });

    const threadListMutator: MutatorCallback<ThreadListOKResponse> = (data) => {
      if (!data) return;

      const newThreads = data.threads.map((t) => {
        if (t.id === postID) {
          const newThread = {
            ...t,
            collections: {
              in_collections: t.collections.in_collections + 1,
              has_collected: true,
            },
          };

          return newThread;
        }
        return t;
      });

      const newData: ThreadListOKResponse = {
        ...data,
        threads: newThreads,
      };

      return newData;
    };

    await mutate(threadListKeyFilterFn, threadListMutator, {
      revalidate: false,
    });

    const collectionListMutator: MutatorCallback<CollectionListOKResponse> = (
      data,
    ) => {
      if (!data) return;

      const newCollections = data.collections.map((c) => {
        if (c.id === collection.id) {
          const newCollection = {
            ...c,
            has_queried_item: true,
          };

          return newCollection;
        }
        return c;
      });

      const newData: CollectionListOKResponse = {
        ...data,
        collections: newCollections,
      };

      return newData;
    };

    await mutate(collectionQueryMutationKey, collectionListMutator, {
      revalidate: false,
    });

    await collectionAddPost(collection.id, postID);
  };

  const removePost = async (collectionID: string, postID: string) => {
    const collectionQueryMutationKey = getCollectionListKey({
      has_item: postID,
      account_handle: session.handle,
    });

    const threadListMutator: MutatorCallback<ThreadListOKResponse> = (data) => {
      if (!data) return;

      const newThreads = data.threads.map((t) => {
        if (t.id === postID) {
          const newThread = {
            ...t,
            collections: {
              in_collections: t.collections.in_collections - 1,
              has_collected: false,
            },
          };

          return newThread;
        }
        return t;
      });

      const newData: ThreadListOKResponse = {
        ...data,
        threads: newThreads,
      };

      return newData;
    };

    await mutate(threadListKeyFilterFn, threadListMutator, {
      revalidate: false,
    });

    const collectionListMutator: MutatorCallback<CollectionListOKResponse> = (
      data,
    ) => {
      if (!data) return;

      const newCollections = data.collections.map((c) => {
        if (c.id === collectionID) {
          const newCollection = {
            ...c,
            has_queried_item: false,
          };

          return newCollection;
        }
        return c;
      });

      const newData: CollectionListOKResponse = {
        ...data,
        collections: newCollections,
      };

      return newData;
    };

    await mutate(collectionQueryMutationKey, collectionListMutator, {
      revalidate: false,
    });

    await collectionRemovePost(collectionID, postID);
  };

  const addNode = async (collectionId: string, nodeID: string) => {
    // Optimistically update

    collectionAddNode(collectionId, nodeID);
  };

  const removeNode = async (collectionId: string, nodeID: string) => {
    // Optimistically update

    collectionRemoveNode(collectionId, nodeID);
  };

  const revalidate = async () => {
    revalidateFeed();
    mutate(collectionListAnyKeyFilterFn);
  };

  return {
    addPost,
    removePost,
    addNode,
    removeNode,
    revalidate,
  };
}
