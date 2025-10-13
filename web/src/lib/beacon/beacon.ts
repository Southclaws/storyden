import { getSendBeaconMutationKey } from "@/api/openapi-client/misc";
import { DatagraphItemKind } from "@/api/openapi-schema";
import { API_ADDRESS } from "@/config";

/**
 * Send a beacon to track user activity for a specific item.
 * This uses the navigator.sendBeacon API for fire-and-forget tracking.
 *
 * @param kind - The type of item being tracked (thread, post, etc.)
 * @param id - The unique identifier of the item
 */
export function sendBeacon(kind: DatagraphItemKind, id: string): void {
  if (typeof navigator === "undefined" || !navigator.sendBeacon) {
    return;
  }

  const beaconData = JSON.stringify({
    k: kind,
    id,
  });

  const endpoint = `${API_ADDRESS}/api${getSendBeaconMutationKey()[0]}`;
  navigator.sendBeacon(endpoint, beaconData);
}
