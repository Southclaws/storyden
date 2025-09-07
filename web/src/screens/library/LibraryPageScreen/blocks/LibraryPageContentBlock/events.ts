import mitt from "mitt";
import { useEffect } from "react";

export type LibraryContentEvents = {
  "library-content:update-generated": string;
};

export const libraryContentBus = mitt<LibraryContentEvents>();

export function useLibraryContentEvent<K extends keyof LibraryContentEvents>(
  type: K,
  handler: (event: LibraryContentEvents[K]) => void,
) {
  useEffect(() => {
    libraryContentBus.on(type, handler);
    return () => {
      libraryContentBus.off(type, handler);
    };
  }, [type, handler]);
}

export function useEmitLibraryContentEvent() {
  return <K extends keyof LibraryContentEvents>(
    type: K,
    payload: LibraryContentEvents[K],
  ) => {
    libraryContentBus.emit(type, payload);
  };
}