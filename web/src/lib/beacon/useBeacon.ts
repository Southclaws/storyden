import { useEffect, useEffectEvent } from "react";

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
 * @param onQueued - Called when the browser queues the beacon; does not confirm delivery
 */
export function useBeacon(
  kind: DatagraphItemKind,
  id: string | undefined,
  onQueued?: () => void | Promise<void>,
) {
  const onBeaconQueued = useEffectEvent(async () => {
    await onQueued?.();
  });

  useEffect(() => {
    if (!id) {
      return;
    }

    queueMicrotask(async () => {
      try {
        if (sendBeacon(kind, id)) {
          await onBeaconQueued();
        }
      } catch (error) {
        console.warn("failed to send beacon:", error);
      }
    });
  }, [kind, id]);
}
