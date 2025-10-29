import { uniqueId } from "lodash";
import { Arguments, MutatorCallback, useSWRConfig } from "swr";

import {
  collectionAddNode,
  collectionAddPost,
  collectionCreate,
  collectionDelete,
  collectionRemoveNode,
  collectionRemovePost,
  collectionUpdate,
  getCollectionListKey,
} from "@/api/openapi-client/collections";
import { getThreadListKey } from "@/api/openapi-client/threads";
import {
  Account,
  Collection,
  CollectionInitialProps,
  CollectionListOKResponse,
  CollectionMutableProps,
  Identifier,
  ThreadListOKResponse,
} from "@/api/openapi-schema";
import { slugify } from "@/utils/slugify";

import { useFeedMutations } from "../feed/mutation";

export function useCollectionMutations(session?: Account) {
  const { mutate } = useSWRConfig();

  const collectionListAnyMutationKey = getCollectionListKey();
  function collectionListAnyKeyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0] === collectionListAnyMutationKey[0];
  }

  const create = async (create: CollectionInitialProps) => {
    const mutator: MutatorCallback<CollectionListOKResponse> = (data) => {
      if (!data) return;
      if (!session) return;

      const newCollection = {
        id: uniqueId("optimistic_collection_"),
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        owner: session,
        has_queried_item: false,
        item_count: 0,
        slug: create.slug ?? slugify(create.name),
        ...create,
      } satisfies Collection;

      const newCollections = [...data.collections, newCollection];

      return {
        ...data,
        collections: newCollections,
      };
    };

    await mutate(collectionListAnyKeyFilterFn, mutator, {
      revalidate: false,
    });

    await collectionCreate(create);
  };

  const update = async (id: Identifier, update: CollectionMutableProps) => {
    const mutator: MutatorCallback<CollectionListOKResponse> = (data) => {
      if (!data) return;

      const newCollections = data.collections.map((c) => {
        if (c.id === id) {
          return {
            ...c,
            ...update,
          };
        }
        return c;
      });

      return {
        ...data,
        collections: newCollections,
      };
    };

    await mutate(collectionListAnyKeyFilterFn, mutator, {
      revalidate: false,
    });

    await collectionUpdate(id, update);
  };

  const deleteCollection = async (id: Identifier) => {
    const mutator: MutatorCallback<CollectionListOKResponse> = (data) => {
      if (!data) return;

      const newCollections = data.collections.filter((c) => c.id !== id);

      return {
        ...data,
        collections: newCollections,
      };
    };

    await mutate(collectionListAnyKeyFilterFn, mutator, {
      revalidate: false,
    });

    await collectionDelete(id);
  };

  const revalidate = async () => {
    mutate(collectionListAnyKeyFilterFn);
  };

  return {
    create,
    update,
    deleteCollection,
    revalidate,
  };
}

export function useCollectionItemMutations(session: Account) {
  const { mutate } = useSWRConfig();
  const { revalidate: revalidateFeed } = useFeedMutations();

  const threadQueryMutationKey = getThreadListKey()[0];
  function threadListKeyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0] === threadQueryMutationKey;
  }

  const collectionListAnyMutationKey = getCollectionListKey()[0];
  function collectionListAnyKeyFilterFn(key: Arguments) {
    return Array.isArray(key) && key[0] === collectionListAnyMutationKey;
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
              has_collected: t.collections.in_collections - 1 > 0,
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
