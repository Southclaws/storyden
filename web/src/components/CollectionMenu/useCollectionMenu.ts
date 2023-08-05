import { useToast } from "@chakra-ui/react";
import { mutate } from "swr";

import {
  collectionAddPost,
  collectionRemovePost,
  useCollectionList,
} from "src/api/openapi/collections";
import { ThreadReference } from "src/api/openapi/schemas";
import { getThreadListKey } from "src/api/openapi/threads";
import { useSession } from "src/auth";

export type Props = {
  thread: ThreadReference;
};

type CollectionState = {
  id: string;
  name: string;
  hasPost: boolean;
};

export function useCollectionMenu(props: Props) {
  const account = useSession();
  const toast = useToast();
  const collectionList = useCollectionList();

  const postCollections = new Set(props.thread.collections.map((c) => c.id));
  const isAlreadySaved = Boolean(
    props.thread.collections.filter((c) => c.owner.id === account?.id).length
  );

  const collections: CollectionState[] =
    collectionList.data?.collections.map((c) => ({
      id: c.id,
      name: c.name,
      hasPost: postCollections.has(c.id),
    })) ?? [];

  const onSelect = (c: CollectionState) => async () => {
    if (postCollections.has(c.id)) {
      await collectionRemovePost(c.id, props.thread.id);
      toast({ title: `Removed from ${c.name}` });
    } else {
      await collectionAddPost(c.id, props.thread.id);
      toast({ title: `Added to ${c.name}` });
    }
    console.log(getThreadListKey());
    await mutate(getThreadListKey({}));
  };

  return {
    error: collectionList.error,
    collections: collections,
    isAlreadySaved,
    onSelect,
  };
}
