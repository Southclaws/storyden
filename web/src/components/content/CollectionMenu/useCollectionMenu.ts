"use client";

import { KeyboardEvent, MouseEvent, useState } from "react";
import { mutate } from "swr";

import {
  collectionAddPost,
  collectionRemovePost,
  useCollectionList,
} from "src/api/openapi/collections";
import { ThreadReference } from "src/api/openapi/schemas";
import { getThreadListKey } from "src/api/openapi/threads";
import { useSession } from "src/auth";
import { useDisclosure } from "src/theme/components";

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
  const collectionList = useCollectionList();
  const [multiSelect, setMultiSelect] = useState(false);
  const [selected, setSelected] = useState(0);

  // called when we want to reset the menu's state.
  const onReset = () => {
    setMultiSelect(false);
    setSelected(0);
  };

  const { isOpen, onOpen, onClose, onToggle } = useDisclosure({
    onClose: onReset,
  });

  const postCollections = new Set(props.thread.collections.map((c) => c.id));
  const isAlreadySaved = Boolean(
    props.thread.collections.filter((c) => c.owner.id === account?.id).length,
  );

  const collections: CollectionState[] =
    collectionList.data?.collections.map((c) => ({
      id: c.id,
      name: c.name,
      hasPost: postCollections.has(c.id),
    })) ?? [];

  const onSelect =
    (c: CollectionState) => async (e: MouseEvent<HTMLButtonElement>) => {
      if (e.shiftKey) {
        setMultiSelect(true);
      }

      if (postCollections.has(c.id)) {
        await collectionRemovePost(c.id, props.thread.id);
      } else {
        await collectionAddPost(c.id, props.thread.id);
      }
      await mutate(getThreadListKey({}));

      setSelected(selected + 1);
    };

  const onKeyDown = (e: KeyboardEvent<HTMLDivElement>) => {
    if (e.shiftKey) setMultiSelect(true);
  };
  const onKeyUp = (e: KeyboardEvent<HTMLDivElement>) => {
    if (!e.shiftKey && multiSelect) {
      setMultiSelect(false);
      if (selected > 0) {
        onToggle();
      }
    }
  };

  return {
    error: collectionList.error,
    collections: collections,
    isAlreadySaved,
    onSelect,
    onKeyDown,
    onKeyUp,
    multiSelect,
    isOpen,
    onOpen,
    onClose,
  };
}
