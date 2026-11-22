import { beforeEach, describe, expect, it, vi } from "vitest";

const { readContext, writeContext } = vi.hoisted(() => ({
  readContext: vi.fn(),
  writeContext: vi.fn(),
}));

vi.mock("@/lib/webpush/storage", () => ({
  readPushContext: readContext,
  writePushContext: writeContext,
}));

const context = {
  accountId: "d4v1p9b2t53b3lo5e3ca",
  apiAddress: "https://api.example.com",
  publicKey: "AQID",
  subscriptionId: "d4v1p9b2t53b3lo5e3cb",
};

const id = "d4v1p9b2t53b3lo5e3cg";
const listeners = new Map<string, (event: unknown) => void>();
const showNotification = vi.fn();
const fetchRequest = vi.fn();
const openWindow = vi.fn();
const matchAll = vi.fn();
const subscribe = vi.fn();

beforeEach(async () => {
  vi.resetModules();
  readContext.mockResolvedValue(context);
  writeContext.mockResolvedValue(undefined);
  showNotification.mockResolvedValue(undefined);
  matchAll.mockResolvedValue([]);
  subscribe.mockResolvedValue({
    toJSON: () => ({
      endpoint: "https://push.example.com/subscription",
      keys: { p256dh: "public-key", auth: "auth-secret" },
    }),
  });

  vi.stubGlobal("fetch", fetchRequest);
  vi.stubGlobal("navigator", {
    locks: {
      request: (_name: string, action: () => Promise<unknown>) => action(),
    },
  });
  vi.stubGlobal("self", {
    addEventListener: (name: string, listener: (event: unknown) => void) => {
      listeners.set(name, listener);
    },
    location: { origin: "https://community.example.com" },
    registration: { showNotification, pushManager: { subscribe } },
    clients: { matchAll, openWindow },
  });

  await import("./notifications");
});

describe("subscription renewal", () => {
  it("renews for the authenticated account with the configured key", async () => {
    respondWithAccount();
    fetchRequest.mockResolvedValueOnce(
      Response.json({ id: "subscription-two" }),
    );

    await dispatch("pushsubscriptionchange", {});

    expect(subscribe).toHaveBeenCalledWith({
      userVisibleOnly: true,
      applicationServerKey: new Uint8Array([1, 2, 3]),
    });
    expect(writeContext).toHaveBeenCalledWith({
      ...context,
      subscriptionId: "subscription-two",
    });
  });

  it("does not renew a previous account's subscription", async () => {
    respondWithAccount("another-account");

    await dispatch("pushsubscriptionchange", {});

    expect(subscribe).not.toHaveBeenCalled();
    expect(writeContext).not.toHaveBeenCalled();
  });
});

async function dispatch(name: string, data: object) {
  let work: Promise<unknown> | undefined;
  const listener = listeners.get(name);

  if (!listener) {
    throw new Error(`No listener for ${name}`);
  }

  listener({
    ...data,
    waitUntil: (promise: Promise<unknown>) => {
      work = promise;
    },
  });

  await work;
}

const payload = {
  account_id: context.accountId,
  subscription_id: context.subscriptionId,
  notification: {
    id,
    created_at: "2026-11-07T12:00:00Z",
    event: "trail_run_attention",
    target: "run one",
    status: "unread",
  },
};

function receive(data: unknown = payload) {
  return dispatch("push", { data: { json: () => data } });
}

function respondWithAccount(accountId = context.accountId) {
  fetchRequest.mockResolvedValueOnce(Response.json({ id: accountId }));
}

describe("push notifications", () => {
  it("displays structured content without API queries and includes the instance icon", async () => {
    fetchRequest.mockRejectedValue(
      new TypeError("Client certificate unavailable"),
    );

    await receive();

    expect(fetchRequest).not.toHaveBeenCalled();
    expect(showNotification).toHaveBeenCalledWith("System", {
      body: "A Trail run needs your attention",
      icon: "https://api.example.com/api/info/icon/192x192",
      tag: `storyden:${context.accountId}:${id}`,
      data: {
        url: "/robots/trails?run=run%20one",
        accountId: context.accountId,
      },
    });
  });

  it("formats the sender and item title from the payload", async () => {
    await receive({
      ...payload,
      notification: {
        ...payload.notification,
        event: "thread_reply",
        source: { id: context.accountId, handle: "odin", name: "Odin" },
        item: { id, kind: "post", title: "A little harbour" },
      },
    });

    expect(fetchRequest).not.toHaveBeenCalled();
    expect(showNotification).toHaveBeenCalledWith(
      "odin",
      expect.objectContaining({
        body: 'replied to your post "A little harbour"',
        data: { url: `/t/locate/${id}`, accountId: context.accountId },
      }),
    );
  });

  it.each([
    null,
    {},
    id,
    { ...payload, notification: {} },
    {
      ...payload,
      notification: { ...payload.notification, item: { id, kind: "post" } },
    },
  ])("ignores malformed notification data %j", async (data) => {
    await receive(data);

    expect(fetchRequest).not.toHaveBeenCalled();
    expect(showNotification).not.toHaveBeenCalled();
  });

  it("ignores malformed JSON", async () => {
    await dispatch("push", {
      data: {
        json: () => {
          throw new SyntaxError("Invalid JSON");
        },
      },
    });

    expect(showNotification).not.toHaveBeenCalled();
  });

  it("ignores pushes without local subscription state", async () => {
    readContext.mockResolvedValue(undefined);

    await receive();

    expect(showNotification).not.toHaveBeenCalled();
  });

  it("ignores notifications for a previous account", async () => {
    await receive({ ...payload, account_id: "00000000000000000001" });

    expect(fetchRequest).not.toHaveBeenCalled();
    expect(showNotification).not.toHaveBeenCalled();
  });

  it("ignores notifications for a replaced subscription", async () => {
    await receive({ ...payload, subscription_id: "00000000000000000001" });

    expect(showNotification).not.toHaveBeenCalled();
  });

  it("does not display content after logout while reading local state", async () => {
    readContext.mockResolvedValueOnce(context).mockResolvedValueOnce(undefined);

    await receive();

    expect(showNotification).not.toHaveBeenCalled();
  });

  it("does not display content after a subscription changes", async () => {
    readContext.mockResolvedValueOnce(context).mockResolvedValueOnce({
      ...context,
      subscriptionId: "00000000000000000001",
    });

    await receive();

    expect(showNotification).not.toHaveBeenCalled();
  });
});

describe("notification navigation", () => {
  async function click(url: string, accountId = context.accountId) {
    await dispatch("notificationclick", {
      notification: { close: vi.fn(), data: { url, accountId } },
    });
  }

  it.each(["//evil.example", "https://evil.example", "javascript:alert(1)"])(
    "rejects external target %s",
    async (url) => {
      await click(url);

      expect(openWindow).not.toHaveBeenCalled();
      expect(matchAll).not.toHaveBeenCalled();
    },
  );

  it("does not open a previous account's notification", async () => {
    await click("/notifications", "another-account");

    expect(openWindow).not.toHaveBeenCalled();
  });

  it("opens the same-origin notification destination", async () => {
    await click("/robots/trails?run=one");

    expect(openWindow).toHaveBeenCalledWith(
      "https://community.example.com/robots/trails?run=one",
    );
  });

  it("navigates and focuses an existing window", async () => {
    const navigate = vi.fn();
    const focus = vi.fn();
    matchAll.mockResolvedValueOnce([
      { url: "https://community.example.com/", navigate, focus },
    ]);

    await click("/notifications");

    expect(navigate).toHaveBeenCalledWith(
      "https://community.example.com/notifications",
    );
    expect(focus).toHaveBeenCalledOnce();
    expect(openWindow).not.toHaveBeenCalled();
  });
});
