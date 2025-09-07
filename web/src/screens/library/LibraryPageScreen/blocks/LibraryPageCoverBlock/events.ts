import mitt from "mitt";
import { useEffect } from "react";

import { Asset } from "@/api/openapi-schema";

export type LibraryCoverEvents = {
  "library-cover:update-from-asset": Asset;
};

export const libraryCoverBus = mitt<LibraryCoverEvents>();

export function useLibraryCoverEvent<K extends keyof LibraryCoverEvents>(
  type: K,
  handler: (event: LibraryCoverEvents[K]) => void,
) {
  useEffect(() => {
    libraryCoverBus.on(type, handler);
    return () => {
      libraryCoverBus.off(type, handler);
    };
  }, [type, handler]);
}

export function useEmitLibraryCoverEvent() {
  return <K extends keyof LibraryCoverEvents>(
    type: K,
    payload: LibraryCoverEvents[K],
  ) => {
    libraryCoverBus.emit(type, payload);
  };
}
