import mitt from "mitt";
import { useEffect } from "react";

import { LibraryPageBlockType } from "./metadata";

export type LibraryBlockEvents = {
  "library:reorder-block": {
    activeId: LibraryPageBlockType;
    overId: LibraryPageBlockType;
  };
  "library:add-block": { type: string };
  "library:remove-block": { type: string };
};

export const libraryBus = mitt<LibraryBlockEvents>();

export function useLibraryBlockEvent<K extends keyof LibraryBlockEvents>(
  type: K,
  handler: (event: LibraryBlockEvents[K]) => void,
) {
  useEffect(() => {
    libraryBus.on(type, handler);
    return () => {
      libraryBus.off(type, handler);
    };
  }, [type, handler]);
}

export function useEmitLibraryBlockEvent() {
  return <K extends keyof LibraryBlockEvents>(
    type: K,
    payload: LibraryBlockEvents[K],
  ) => {
    libraryBus.emit(type, payload);
  };
}
