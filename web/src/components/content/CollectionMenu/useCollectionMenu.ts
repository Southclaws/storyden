import { contains, map } from "lodash/fp";
import { KeyboardEvent, useState } from "react";

import {
  collectionAddPost,
  collectionRemovePost,
  useCollectionList,
} from "src/api/openapi/collections";
import {
  Collection,
  CollectionList,
  ThreadReference,
} from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { useFeedState } from "src/components/feed/useFeedState";
import { useDisclosure } from "src/utils/useDisclosure";

export type Props = {
  thread: ThreadReference;
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

export function useCollectionMenu({ thread }: Props) {
  const account = useSession();
  const { data, error } = useCollectionList();
  const { mutate } = useFeedState();

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
    thread.collections.filter((c) => c.owner.id === account?.id).length,
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

        await mutate?.mutateThreads();
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
