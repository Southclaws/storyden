import { expect, test } from "@playwright/test";

import { API_ADDRESS, WEB_ADDRESS } from "../../src/config";
import { withAdminAccessKey } from "../access_key_admin_assignment";

test("includes a reply in SSR immediately after posting", async ({
  request,
}) => {
  const suffix = Date.now();
  const signup = await request.post(`${API_ADDRESS}/api/auth/password/signup`, {
    data: {
      identifier: `ssr-reply-${suffix}`,
      token: "MetadataTestPassword123!",
    },
  });
  expect(signup.status()).toBe(200);

  const created = await request.post(`${API_ADDRESS}/api/threads`, {
    data: {
      title: `SSR reply freshness ${suffix}`,
      body: "This thread will receive a reply.",
      visibility: "published",
    },
  });
  expect(created.status()).toBe(200);
  const thread = await created.json();
  const path = `/t/${thread.slug}`;
  const initial = await request.get(path);
  expect(initial.status()).toBe(200);

  const replyText = `A reply posted after the SSR cache was primed ${suffix}`;
  const reply = await request.post(
    `${API_ADDRESS}/api/threads/${thread.slug}/replies`,
    { data: { body: replyText } },
  );
  expect(reply.status()).toBe(200);

  const refreshed = await request.get(path);
  expect(refreshed.status()).toBe(200);
  expect(await refreshed.text()).toContain(replyText);
});

test("returns 404 for a missing Library preview rather than a server error", async ({
  request,
}) => {
  const response = await request.get(
    `/l/og?slug=missing-metadata-${Date.now()}`,
  );

  expect(response.status()).toBe(404);
  expect(response.headers()["cache-control"]).toBe("no-store");
});

test("renders guest content without client hydration", async ({ browser }) => {
  const suffix = Date.now();
  const threadTitle = `SSR discussion ${suffix}`;
  const thread = await withAdminAccessKey(({ threadCreate }) =>
    threadCreate({
      title: threadTitle,
      body: "This discussion was rendered on the server.",
      visibility: "published",
    }),
  );
  const nodeTitle = `SSR Library page ${suffix}`;
  const node = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: nodeTitle,
      slug: `ssr-metadata-${suffix}`,
      content: "<p>This Library page was rendered on the server.</p>",
      visibility: "published",
    }),
  );
  const context = await browser.newContext({
    baseURL: WEB_ADDRESS,
    userAgent: "Twitterbot/1.0",
  });

  const browserAPIRequests: string[] = [];
  await context.route("**/*", async (route) => {
    const request = route.request();
    if (request.resourceType() === "script") {
      await route.abort();
      return;
    }

    if (
      request.resourceType() === "fetch" ||
      request.resourceType() === "xhr"
    ) {
      browserAPIRequests.push(request.url());
    }

    await route.continue();
  });

  try {
    const page = await context.newPage();
    for (const [path, title, body] of [
      [
        `/t/${thread.slug}`,
        threadTitle,
        "This discussion was rendered on the server.",
      ],
      [
        `/l/${node.slug}`,
        nodeTitle,
        "This Library page was rendered on the server.",
      ],
    ] as const) {
      const response = await page.goto(path);

      expect(response?.status()).toBe(200);
      await expect(page).toHaveTitle(new RegExp(title));
      await expect(page.getByText(body, { exact: true })).toBeVisible();
      expect(
        await page.locator('script[type="application/ld+json"]').count(),
      ).toBeGreaterThan(0);
    }

    expect(browserAPIRequests).toEqual([]);
  } finally {
    await context.close();
  }
});

test("keeps authorized private titles without public SEO content", async ({
  request,
}) => {
  const suffix = Date.now();
  const handle = `metadata-${suffix}`;
  const signup = await request.post(`${API_ADDRESS}/api/auth/password/signup`, {
    data: { identifier: handle, token: "MetadataTestPassword123!" },
  });
  expect(signup.status()).toBe(200);
  await withAdminAccessKey(({ accountAddRole }) =>
    accountAddRole(handle, "00000000000000000a00"),
  );

  const threadTitle = `In review discussion ${suffix}`;
  const thread = await withAdminAccessKey(({ threadCreate }) =>
    threadCreate({
      title: threadTitle,
      body: "Awaiting moderation",
      visibility: "review",
    }),
  );
  const nodeTitle = `Draft Library page ${suffix}`;
  const nodeResponse = await request.post(`${API_ADDRESS}/api/nodes`, {
    data: {
      name: nodeTitle,
      slug: `draft-metadata-${suffix}`,
      visibility: "draft",
    },
  });
  expect(nodeResponse.status()).toBe(200);
  const node = await nodeResponse.json();

  for (const [path, title] of [
    [`/t/${thread.slug}`, threadTitle],
    [`/l/${node.slug}`, nodeTitle],
  ] as const) {
    const response = await request.get(path);
    expect(response.status(), path).toBe(200);
    const html = await response.text();

    expect(html).toContain(`<title>${title} | `);
    expect(html).toContain('name="robots" content="noindex, nofollow"');
    expect(html).not.toContain(`property="og:title" content="${title}"`);
    expect(html).not.toMatch(/<script type="application\/ld\+json">/);
  }
});
