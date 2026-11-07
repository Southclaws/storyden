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
  accountId: "account-one",
  apiAddress: "https://api.example.com",
  publicKey: "AQID",
  subscriptionId: "subscription-one",
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

function receive(payload = id) {
  return dispatch("push", { data: { text: () => payload } });
}

function respondWithAccount(accountId = context.accountId) {
  fetchRequest.mockResolvedValueOnce(Response.json({ id: accountId }));
}

describe("push notifications", () => {
  it("retrieves structured data using cookies and displays shared formatting", async () => {
    respondWithAccount();
    fetchRequest.mockResolvedValueOnce(
      Response.json({
        notifications: [
          {
            id,
            created_at: "2026-11-07T12:00:00Z",
            event: "trail_run_attention",
            target: "run one",
            status: "unread",
          },
        ],
      }),
    );

    await receive();

    const request = fetchRequest.mock.lastCall?.[0] as Request;

    expect(request.url).toBe("https://api.example.com/api/notifications");
    expect(request.credentials).toBe("include");
    expect(request.cache).toBe("no-store");
    expect(request.redirect).toBe("error");
    expect(showNotification).toHaveBeenCalledWith("System", {
      body: "A Trail run needs your attention",
      icon: "https://api.example.com/api/info/icon/192x192",
      tag: `storyden:account-one:${id}`,
      data: {
        url: "/robots/trails?run=run%20one",
        accountId: context.accountId,
      },
    });
  });

  it.each(["", "../secrets", "<script>", "z".repeat(20), "a".repeat(21)])(
    "ignores invalid payload %s",
    async (payload) => {
      await receive(payload);

      expect(fetchRequest).not.toHaveBeenCalled();
      expect(showNotification).not.toHaveBeenCalled();
    },
  );

  it("ignores pushes without local subscription state", async () => {
    readContext.mockResolvedValue(undefined);

    await receive();

    expect(showNotification).not.toHaveBeenCalled();
  });

  it.each([401, 403])(
    "suppresses personal content when the session returns %i",
    async (status) => {
      fetchRequest.mockResolvedValueOnce(new Response(null, { status }));

      await receive();

      expect(showNotification).not.toHaveBeenCalled();
    },
  );

  it("ignores notifications for a previous account", async () => {
    respondWithAccount("another-account");

    await receive();

    expect(fetchRequest).toHaveBeenCalledTimes(1);
    expect(showNotification).not.toHaveBeenCalled();
  });

  it("ignores identifiers missing from the authenticated inbox", async () => {
    respondWithAccount();
    fetchRequest.mockResolvedValueOnce(Response.json({ notifications: [] }));

    await receive();

    expect(showNotification).not.toHaveBeenCalled();
  });

  it("uses generic frontend copy when offline", async () => {
    fetchRequest.mockRejectedValueOnce(new TypeError("Offline"));
    vi.spyOn(console, "warn").mockImplementation(() => {});

    await receive();

    expect(showNotification).toHaveBeenCalledWith(
      "System",
      expect.objectContaining({
        body: "You have a new notification",
        data: { url: "/notifications", accountId: context.accountId },
      }),
    );
  });

  it("does not display data after logout during a fetch", async () => {
    respondWithAccount();
    fetchRequest.mockResolvedValueOnce(new Response(null, { status: 503 }));
    readContext.mockResolvedValueOnce(context).mockResolvedValueOnce(undefined);

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
