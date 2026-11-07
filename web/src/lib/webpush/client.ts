import { decodePublicKey, deleteSubscription, saveSubscription } from "./api";
import { withPushLock } from "./lock";
import { type PushContext, readPushContext, writePushContext } from "./storage";

export function supportsWebPush() {
  return (
    window.isSecureContext &&
    "serviceWorker" in navigator &&
    "PushManager" in window &&
    "Notification" in window &&
    "indexedDB" in window
  );
}

async function getRegistration() {
  const registration = await navigator.serviceWorker.getRegistration("/");
  const worker =
    registration?.active ?? registration?.waiting ?? registration?.installing;

  if (worker && new URL(worker.scriptURL).pathname === "/sw.js") {
    return registration;
  }

  return undefined;
}

async function clearSubscription() {
  const context = await readPushContext();
  await writePushContext();

  const registration = await getRegistration();
  const notifications = await registration?.getNotifications();

  for (const notification of notifications ?? []) {
    notification.close();
  }

  const subscription = await registration?.pushManager.getSubscription();

  if (subscription) {
    await subscription.unsubscribe();
  }

  if (context) {
    await deleteSubscription(context);
  }
}

export async function clearWebPush() {
  if (!supportsWebPush()) {
    return;
  }

  await withPushLock(clearSubscription);
}

export async function reconcileWebPush(context: PushContext) {
  if (!supportsWebPush()) {
    return;
  }

  await withPushLock(async () => {
    const previous = await readPushContext();
    const changed =
      previous &&
      (previous.accountId !== context.accountId ||
        previous.apiAddress !== context.apiAddress ||
        previous.publicKey !== context.publicKey);

    if (changed) {
      await clearSubscription();
    }

    await navigator.serviceWorker.register("/sw.js", {
      scope: "/",
      updateViaCache: "none",
    });

    const registration = await navigator.serviceWorker.ready;
    const current = { ...context };

    if (!changed && previous?.subscriptionId) {
      current.subscriptionId = previous.subscriptionId;
    }

    await writePushContext(current);

    if (Notification.permission !== "granted") {
      return;
    }

    let subscription = await registration.pushManager.getSubscription();
    const key = decodePublicKey(context.publicKey);

    if (subscription && !matchesPublicKey(subscription, key)) {
      await subscription.unsubscribe();
      subscription = null;
    }

    if (!subscription) {
      subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: key,
      });
    }

    await writePushContext(await saveSubscription(current, subscription));
  });
}

function matchesPublicKey(subscription: PushSubscription, key: Uint8Array) {
  const existing = subscription.options.applicationServerKey;

  if (!existing) {
    return false;
  }

  const bytes = new Uint8Array(existing);
  return (
    bytes.length === key.length &&
    bytes.every((value, index) => value === key[index])
  );
}

export async function enableWebPush(context: PushContext) {
  if (!supportsWebPush()) {
    return false;
  }

  let permission = Notification.permission;

  if (permission === "default") {
    permission = await Notification.requestPermission();
  }

  if (permission !== "granted") {
    return false;
  }

  await reconcileWebPush(context);
  return true;
}
