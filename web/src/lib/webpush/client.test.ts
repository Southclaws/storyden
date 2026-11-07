import { beforeEach, describe, expect, it, vi } from "vitest";

const { readContext, writeContext, saveSubscription, deleteSubscription } =
  vi.hoisted(() => ({
    readContext: vi.fn(),
    writeContext: vi.fn(),
    saveSubscription: vi.fn(),
    deleteSubscription: vi.fn(),
  }));

vi.mock("./storage", () => ({
  readPushContext: readContext,
  writePushContext: writeContext,
}));

vi.mock("./api", async (importOriginal) => ({
  ...(await importOriginal<typeof import("./api")>()),
  saveSubscription,
  deleteSubscription,
}));

import { clearWebPush, enableWebPush, reconcileWebPush } from "./client";

const context = {
  accountId: "account-one",
  apiAddress: "https://api.example.com",
  publicKey: "AQID",
};

const unsubscribe = vi.fn();
const subscription = {
  options: { applicationServerKey: new Uint8Array([1, 2, 3]).buffer },
  unsubscribe,
};

const getSubscription = vi.fn();
const subscribe = vi.fn();
const register = vi.fn();
const closeNotification = vi.fn();
const requestPermission = vi.fn();
const registration = {
  active: { scriptURL: "https://community.example.com/sw.js" },
  pushManager: { getSubscription, subscribe },
  getNotifications: vi.fn(),
};

beforeEach(() => {
  readContext.mockResolvedValue(undefined);
  writeContext.mockResolvedValue(undefined);
  deleteSubscription.mockResolvedValue(undefined);
  saveSubscription.mockResolvedValue({
    ...context,
    subscriptionId: "subscription-one",
  });
  getSubscription.mockResolvedValue(subscription);
  subscribe.mockResolvedValue(subscription);
  unsubscribe.mockResolvedValue(true);
  register.mockResolvedValue(registration);
  registration.getNotifications.mockResolvedValue([
    { close: closeNotification },
  ]);

  vi.stubGlobal("window", {
    isSecureContext: true,
    PushManager: {},
    Notification: {},
    indexedDB: {},
  });

  vi.stubGlobal("Notification", { permission: "granted", requestPermission });
  vi.stubGlobal("navigator", {
    serviceWorker: {
      register,
      ready: Promise.resolve(registration),
      getRegistration: vi.fn().mockResolvedValue(registration),
    },
    locks: { request: vi.fn((_name, action) => action()) },
  });
});

describe("Web Push subscription lifecycle", () => {
  it("registers and refreshes an existing subscription without prompting", async () => {
    await reconcileWebPush(context);

    expect(register).toHaveBeenCalledWith("/sw.js", {
      scope: "/",
      updateViaCache: "none",
    });
    expect(saveSubscription).toHaveBeenCalledWith(context, subscription);
    expect(subscribe).not.toHaveBeenCalled();
    expect(requestPermission).not.toHaveBeenCalled();
  });

  it("subscribes when permission is already granted", async () => {
    getSubscription.mockResolvedValueOnce(null);

    await reconcileWebPush(context);

    expect(subscribe).toHaveBeenCalledWith({
      userVisibleOnly: true,
      applicationServerKey: new Uint8Array([1, 2, 3]),
    });
    expect(writeContext).toHaveBeenLastCalledWith({
      ...context,
      subscriptionId: "subscription-one",
    });
  });

  it.each(["default", "denied"])(
    "does not request permission or subscribe when %s",
    async (permission) => {
      vi.stubGlobal("Notification", { permission, requestPermission });

      await reconcileWebPush(context);

      expect(writeContext).toHaveBeenCalledWith(context);
      expect(subscribe).not.toHaveBeenCalled();
      expect(requestPermission).not.toHaveBeenCalled();
    },
  );

  it("replaces a subscription made with another application server key", async () => {
    getSubscription.mockResolvedValueOnce({
      ...subscription,
      options: { applicationServerKey: new Uint8Array([4]).buffer },
    });

    await reconcileWebPush(context);

    expect(unsubscribe).toHaveBeenCalledOnce();
    expect(subscribe).toHaveBeenCalledOnce();
  });

  it("clears state and unregisters the previous account's subscription", async () => {
    const previous = {
      ...context,
      accountId: "old-account",
      subscriptionId: "old-subscription",
    };
    readContext.mockResolvedValueOnce(previous);
    readContext.mockResolvedValueOnce(previous);
    getSubscription
      .mockResolvedValueOnce(subscription)
      .mockResolvedValueOnce(null);

    await reconcileWebPush(context);

    expect(writeContext).toHaveBeenNthCalledWith(1);
    expect(deleteSubscription).toHaveBeenCalledWith(previous);
    expect(unsubscribe).toHaveBeenCalledOnce();
    expect(subscribe).toHaveBeenCalledOnce();
  });

  it("clears local state before deleting server registration during logout", async () => {
    readContext.mockResolvedValueOnce({
      ...context,
      subscriptionId: "subscription-one",
    });

    await clearWebPush();

    expect(writeContext).toHaveBeenCalledWith();
    expect(closeNotification).toHaveBeenCalledOnce();
    expect(unsubscribe).toHaveBeenCalledOnce();
    expect(deleteSubscription).toHaveBeenCalledOnce();
  });

  it("starts the permission prompt synchronously when explicitly requested", async () => {
    requestPermission.mockResolvedValueOnce("denied");

    vi.stubGlobal("Notification", { permission: "default", requestPermission });

    const enabled = enableWebPush(context);

    expect(requestPermission).toHaveBeenCalledOnce();
    expect(await enabled).toBe(false);
  });
});
