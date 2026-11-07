import {
  pushSubscriptionDelete,
  pushSubscriptionUpsert,
} from "@/api/openapi-worker/notifications";

import type { PushContext } from "./storage";

export function decodePublicKey(publicKey: string) {
  const decoded = atob(publicKey.replace(/-/g, "+").replace(/_/g, "/"));
  return Uint8Array.from(decoded, (character) => character.charCodeAt(0));
}

export async function saveSubscription(
  context: PushContext,
  subscription: PushSubscription,
) {
  const input = subscription.toJSON();
  const endpoint = input.endpoint;
  const p256dh = input.keys?.["p256dh"];
  const auth = input.keys?.["auth"];

  if (!endpoint || !p256dh || !auth) {
    throw new Error(
      "Browser push subscription is missing its endpoint or encryption keys",
    );
  }

  const registered = await pushSubscriptionUpsert(
    {
      endpoint,
      keys: { p256dh, auth },
    },
    { apiAddress: context.apiAddress },
  );

  if (registered.status !== 200) {
    throw new Error(
      `Push subscription registration failed (${registered.status})`,
    );
  }

  return { ...context, subscriptionId: registered.data.id };
}

export async function deleteSubscription(context: PushContext) {
  if (!context.subscriptionId) {
    return;
  }

  const response = await pushSubscriptionDelete(context.subscriptionId, {
    apiAddress: context.apiAddress,
  });

  if (
    response.status !== 204 &&
    response.status !== 404 &&
    response.status !== 401
  ) {
    throw new Error(`Push subscription removal failed (${response.status})`);
  }
}
