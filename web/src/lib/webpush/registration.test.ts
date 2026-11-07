import { expect, it, vi } from "vitest";

const { readContext, writeContext } = vi.hoisted(() => ({
  readContext: vi.fn(),
  writeContext: vi.fn(),
}));

vi.mock("./storage", () => ({
  readPushContext: readContext,
  writePushContext: writeContext,
}));

import { enableWebPush } from "./client";

it("requests permission before asynchronous setup and saves a first subscription through the generated API", async () => {
  const context = {
    accountId: "account-one",
    apiAddress: "https://api.example.com",
    publicKey: "AQID",
  };
  const input = {
    endpoint: "https://push.example.com/subscription",
    keys: { p256dh: "receiver-public-key", auth: "receiver-auth-secret" },
  };
  const subscription = { toJSON: () => input };
  const subscribe = vi.fn().mockResolvedValue(subscription);
  const registration = {
    pushManager: {
      getSubscription: vi.fn().mockResolvedValue(null),
      subscribe,
    },
  };
  const notifications = {
    permission: "default",
    requestPermission: vi.fn(async () => {
      notifications.permission = "granted";
      return "granted";
    }),
  };
  const fetchRequest = vi
    .fn()
    .mockResolvedValue(
      Response.json({ id: "subscription-one", endpoint: input.endpoint }),
    );

  readContext.mockResolvedValue(undefined);
  writeContext.mockResolvedValue(undefined);
  vi.stubGlobal("Notification", notifications);
  vi.stubGlobal("window", {
    isSecureContext: true,
    PushManager: {},
    Notification: notifications,
    indexedDB: {},
  });
  vi.stubGlobal("navigator", {
    serviceWorker: {
      register: vi.fn().mockResolvedValue(registration),
      ready: Promise.resolve(registration),
    },
    locks: {
      request: (_name: string, action: () => Promise<unknown>) => action(),
    },
  });
  vi.stubGlobal("fetch", fetchRequest);

  const enabled = enableWebPush(context);

  expect(notifications.requestPermission).toHaveBeenCalledOnce();
  expect(readContext).not.toHaveBeenCalled();
  expect(await enabled).toBe(true);

  const request = fetchRequest.mock.lastCall?.[0] as Request;

  expect(request.url).toBe(
    "https://api.example.com/api/notifications/push/subscriptions",
  );
  expect(request.method).toBe("PUT");
  expect(request.credentials).toBe("include");
  expect(await request.json()).toEqual(input);
  expect(writeContext).toHaveBeenLastCalledWith({
    ...context,
    subscriptionId: "subscription-one",
  });
});
