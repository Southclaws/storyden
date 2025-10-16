import { useEffect } from "react";

import { DatagraphItemKind } from "@/api/openapi-schema";

import { sendBeacon } from "./beacon";

/**
 * Hook that sends a beacon when the component mounts or when kind/id changes.
 * Useful for tracking page views or item interactions.
 *
 * It's safe to use immediately after a useSWR call where data is undefined,
 * simply pass the data?.id in and it'll fire once the data is available.
 *
 * @param kind - The type of item being tracked
 * @param id - The unique identifier of the item
 */
export function useBeacon(kind: DatagraphItemKind, id: string | undefined) {
  useEffect(() => {
    if (!id) {
      return;
    }

    queueMicrotask(() => {
      try {
        sendBeacon(kind, id);
      } catch (error) {
        console.warn("failed to send beacon:", error);
      }
    });
  }, [kind, id]);
}
