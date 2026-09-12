import { useEffect } from "react";

import { createEmitter } from "@/utils/emitter";

import { LibraryPageBlockType } from "./metadata";

export type LibraryEvents = {
  "library:reorder-block": {
    activeId: LibraryPageBlockType;
    overId: LibraryPageBlockType;
  };
  "library:add-block": { type: LibraryPageBlockType; index?: number };
  "library:remove-block": { type: LibraryPageBlockType };
  "library:revalidate": Record<string, never>;
};

export const libraryBus = createEmitter<LibraryEvents>();

export function useLibraryEvent<K extends keyof LibraryEvents>(
  type: K,
  handler: (event: LibraryEvents[K]) => void,
) {
  useEffect(() => {
    libraryBus.on(type, handler);
    return () => {
      libraryBus.off(type, handler);
    };
  }, [type, handler]);
}

export function useEmitLibraryBlockEvent() {
  return <K extends keyof LibraryEvents>(
    type: K,
    payload: LibraryEvents[K],
  ) => {
    libraryBus.emit(type, payload);
  };
}
