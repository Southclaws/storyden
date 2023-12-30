import { KeyboardEvent, useState } from "react";

import { useCollectionList } from "src/api/openapi/collections";
import { ThreadReference } from "src/api/openapi/schemas";
import { useSession } from "src/auth";
import { useDisclosure } from "src/utils/useDisclosure";

export type Props = {
  thread: ThreadReference;
};

export function useCollectionMenu({ thread }: Props) {
  const account = useSession();
  const { data: collections, error } = useCollectionList();
  const [multiSelect, setMultiSelect] = useState(false);
  const [selected, setSelected] = useState(0);

  const handleReset = () => {
    setMultiSelect(false);
    setSelected(0);
  };

  const { onOpenChange, onToggle } = useDisclosure({
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
    ready: true as const,
    isAlreadySaved,
    collections: collections.collections,
    multiSelect,
    onKeyDown,
    onKeyUp,
    onOpenChange,
  };
}
