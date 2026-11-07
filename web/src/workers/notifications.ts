import { accountGet } from "@/api/openapi-worker/accounts";
import { notificationList } from "@/api/openapi-worker/notifications";
import {
  formatNotification,
  formatNotificationFallback,
} from "@/lib/notification";
import { decodePublicKey, saveSubscription } from "@/lib/webpush/api";
import { withPushLock } from "@/lib/webpush/lock";
import {
  type PushContext,
  readPushContext,
  writePushContext,
} from "@/lib/webpush/storage";

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
  const id = event.data?.text();

  if (!id || !/^[0-9a-v]{20}$/.test(id)) {
    return;
  }

  const context = await readPushContext();

  if (!context?.subscriptionId) {
    return;
  }

  let content = formatNotificationFallback();

  try {
    if (!(await isCurrentAccount(context))) {
      return;
    }

    const response = await notificationList(undefined, {
      apiAddress: context.apiAddress,
    });

    if (response.status === 401 || response.status === 403) {
      return;
    }

    if (response.status === 200) {
      const result = response.data;
      const notification = result.notifications.find(
        (notification) => notification.id === id,
      );

      if (!notification) {
        return;
      }

      content = formatNotification(notification);
    } else {
      console.warn(
        "[Web Push worker] Notification lookup failed; using generic notification",
        {
          status: response.status,
        },
      );
    }
  } catch (error) {
    console.warn(
      "[Web Push worker] Notification lookup failed; using generic notification",
      error,
    );
  }

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
      tag: `storyden:${context.accountId}:${id}`,
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
