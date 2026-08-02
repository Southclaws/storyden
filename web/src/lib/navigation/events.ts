import { useEffect } from "react";

import { createEmitter } from "@/utils/emitter";

export type NavigationItemEvents = {
  "navigation:reorder-item": {
    activeKey: string;
    overKey: string;
  };
};

export const navigationItemBus = createEmitter<NavigationItemEvents>();

export function useNavigationItemEvent<K extends keyof NavigationItemEvents>(
  type: K,
  handler: (event: NavigationItemEvents[K]) => void,
) {
  useEffect(() => {
    navigationItemBus.on(type, handler);
    return () => navigationItemBus.off(type, handler);
  }, [handler, type]);
}

export function useEmitNavigationItemEvent() {
  return <K extends keyof NavigationItemEvents>(
    type: K,
    payload: NavigationItemEvents[K],
  ) => navigationItemBus.emit(type, payload);
}
