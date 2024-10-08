import { contains, map } from "lodash/fp";
import { KeyboardEvent, useState } from "react";

import {
  collectionAddPost,
  collectionRemovePost,
  useCollectionList,
} from "src/api/openapi-client/collections";
import {
  Account,
  Collection,
  CollectionList,
  PostReference,
} from "src/api/openapi-schema";
import { useDisclosure } from "src/utils/useDisclosure";

import { useFeedMutations } from "@/lib/feed/mutation";

export type Props = {
  account: Account;
  thread: PostReference;
};

export type CollectionWithHasPost = Collection & {
  hasPost: boolean;
};

const hydrateState = (
  collections: CollectionList,
  threadCollections: CollectionList,
): CollectionWithHasPost[] => {
  const cids = threadCollections.map((v) => v.id);

  const m = map((c: Collection) => ({
    ...c,
    hasPost: contains(c.id)(cids),
  }));

  return m(collections);
};

export function useCollectionMenu({ account, thread }: Props) {
  const { data, error } = useCollectionList({ account_handle: account.handle });
  const { revalidate } = useFeedMutations();

  const [multiSelect, setMultiSelect] = useState(false);
  const [selected, setSelected] = useState(0);

  const collections = hydrateState(data?.collections ?? [], thread.collections);

  const handleReset = () => {
    setMultiSelect(false);
    setSelected(0);
  };

  const { onOpenChange: handleOpenChange, onToggle } = useDisclosure({
    onClose: handleReset,
  });

  if (!collections) {
    return {
      ready: false as const,
      error,
    };
  }

  const isAlreadySaved = Boolean(
    thread.collections.filter((c) => c.owner.id === account.id).length,
  );

  const handleKeyDown = (e: KeyboardEvent<HTMLDivElement>) => {
    if (e.shiftKey) setMultiSelect(true);
  };

  const handleKeyUp = (e: KeyboardEvent<HTMLDivElement>) => {
    if (!e.shiftKey && multiSelect) {
      setMultiSelect(false);
      if (selected > 0) {
        onToggle();
      }
    }
  };

  const handleSelect = async ({ value: id }: { value: string }) => {
    switch (id) {
      case "create-collection":
        // handled by <Button />
        return;

      default:
        if (thread.collections.find((v) => v.id === id)) {
          await collectionRemovePost(id, thread.id);
        } else {
          await collectionAddPost(id, thread.id);
        }

        // TODO: Optimistic mutation for collection changes.
        await revalidate();
    }
  };

  return {
    ready: true as const,
    isAlreadySaved,
    collections,
    multiSelect,
    handlers: {
      handleKeyDown,
      handleKeyUp,
      handleOpenChange,
      handleSelect,
    },
  };
}
