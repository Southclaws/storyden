import { DatagraphItemKind } from "@/api/openapi-schema/datagraphItemKind";
import { NotificationEvent } from "@/api/openapi-schema/notificationEvent";
import { NotificationStatus } from "@/api/openapi-schema/notificationStatus";
import { accountGet } from "@/api/openapi-worker/accounts";
import { formatNotification } from "@/lib/notification";
import * as z from "@/lib/validation";
import { decodePublicKey, saveSubscription } from "@/lib/webpush/api";
import { withPushLock } from "@/lib/webpush/lock";
import {
  type PushContext,
  readPushContext,
  writePushContext,
} from "@/lib/webpush/storage";

const IdentifierSchema = z.string().check(z.regex(/^[0-9a-v]{20}$/));

// NOTE: Hand built from the notification schema - gotta match the OpenAPI def
// otherwise it'll break. If that def changes, it has to change here too!
const PushNotificationPayloadSchema = z.object({
  account_id: IdentifierSchema,
  subscription_id: IdentifierSchema,
  notification: z.object({
    id: IdentifierSchema,
    created_at: z.string(),
    event: z.enum(NotificationEvent),
    status: z.enum(NotificationStatus),
    target: z.optional(z.string()),
    item: z.optional(
      z.object({
        id: IdentifierSchema,
        kind: z.enum(DatagraphItemKind),
        title: z.string(),
        slug: z.optional(z.string()),
      }),
    ),
    source: z.optional(
      z.object({
        id: IdentifierSchema,
        handle: z.string(),
        name: z.string(),
      }),
    ),
  }),
});

type PushNotificationPayload = z.infer<typeof PushNotificationPayloadSchema>;

declare const self: ServiceWorkerGlobalScope;

self.addEventListener("install", (event) => {
  event.waitUntil(self.skipWaiting());
});

self.addEventListener("activate", (event) => {
  event.waitUntil(self.clients.claim());
});

self.addEventListener("push", (event) => {
  event.waitUntil(receivePush(event));
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();
  event.waitUntil(openNotification(event.notification));
});

self.addEventListener("pushsubscriptionchange", (event) => {
  const change = event as ExtendableEvent;
  change.waitUntil(withPushLock(renewSubscription));
});

async function isCurrentAccount(context: PushContext) {
  const response = await accountGet({ apiAddress: context.apiAddress });

  if (response.status !== 200) {
    return false;
  }

  return response.data.id === context.accountId;
}

async function receivePush(event: PushEvent) {
  let payload: PushNotificationPayload;
  try {
    const data: unknown = event.data?.json();
    const result = PushNotificationPayloadSchema.safeParse(data);
    if (!result.success) {
      return;
    }

    payload = result.data;
  } catch (e) {
    console.warn("push notification payload did not validate", e);
    return;
  }

  const context = await readPushContext();
  if (
    context?.accountId !== payload.account_id ||
    context.subscriptionId !== payload.subscription_id
  ) {
    return;
  }

  const notification = payload.notification;
  const content = formatNotification(notification);

  const current = await readPushContext();

  if (
    current?.accountId !== context.accountId ||
    current.subscriptionId !== context.subscriptionId
  ) {
    return;
  }

  try {
    await self.registration.showNotification(content.title, {
      body: content.body,
      icon: `${context.apiAddress.replace(/\/$/, "")}/api/info/icon/192x192`,
      tag: `storyden:${context.accountId}:${notification.id}`,
      data: { url: content.url, accountId: context.accountId },
    });
  } catch (error) {
    console.error("[Web Push worker] Notification display failed", error);
    throw error;
  }
}

async function openNotification(notification: Notification) {
  const context = await readPushContext();
  const data = notification.data;

  if (!context || data?.accountId !== context.accountId) {
    return;
  }

  if (typeof data.url !== "string" || !data.url.startsWith("/")) {
    return;
  }

  const url = new URL(data.url, self.location.origin);

  if (url.origin !== self.location.origin) {
    return;
  }

  const windows = await self.clients.matchAll({
    type: "window",
    includeUncontrolled: true,
  });

  for (const client of windows) {
    if (new URL(client.url).origin === url.origin) {
      await client.navigate(url.href);
      await client.focus();

      return;
    }
  }

  await self.clients.openWindow(url.href);
}

async function renewSubscription() {
  const context = await readPushContext();

  if (!context) {
    return;
  }

  if (!(await isCurrentAccount(context))) {
    return;
  }

  const subscription = await self.registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: decodePublicKey(context.publicKey),
  });

  const updated = await saveSubscription(context, subscription);

  const current = await readPushContext();

  if (
    current?.accountId === context.accountId &&
    current.publicKey === context.publicKey
  ) {
    await writePushContext(updated);
  }
}
