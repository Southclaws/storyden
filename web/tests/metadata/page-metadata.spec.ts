import { expect, test } from "@playwright/test";

import { API_ADDRESS, WEB_ADDRESS } from "../../src/config";
import { withAdminAccessKey } from "../access_key_admin_assignment";

const userAgents = ["Googlebot", "Twitterbot/1.0", "Mozilla/5.0"];

for (const userAgent of userAgents) {
  for (const failure of [
    { status: 429, clientAddress: "198.51.100.86" },
    { status: 500, clientAddress: "198.51.100.87" },
  ]) {
    test(`returns retryable HTML for ${userAgent} after an anonymous ${failure.status}`, async ({
      request,
    }) => {
      const suffix = `${Date.now()}-${failure.status}`;
      const signup = await request.post(
        `${API_ADDRESS}/api/auth/password/signup`,
        {
          data: {
            identifier: `metadata-${Date.now()}`,
            token: "MetadataTestPassword123!",
          },
        },
      );
      expect(signup.status()).toBe(200);

      const thread = await withAdminAccessKey(({ threadCreate }) =>
        threadCreate({
          title: `Metadata thread ${suffix}`,
          body: "A published discussion",
          visibility: "published",
        }),
      );
      const node = await withAdminAccessKey(({ nodeCreate }) =>
        nodeCreate({
          name: `Metadata page ${suffix}`,
          slug: `metadata-${suffix}`,
          visibility: "published",
        }),
      );

      for (const path of ["/", `/t/${thread.slug}`, `/l/${node.slug}`]) {
        const healthy = await request.get(path, {
          headers: { "User-Agent": userAgent },
        });
        expect(healthy.status(), path).toBe(200);
        expect(await healthy.text()).toContain("<title>");

        const response = await request.get(path, {
          headers: {
            "User-Agent": userAgent,
            "X-Forwarded-For": failure.clientAddress,
          },
        });

        expect(response.status(), path).toBe(200);
        const html = await response.text();
        expect(html).toContain("<title>");
        if (path.startsWith("/t/")) {
          expect(html).toContain(`<title>${thread.title} | `);
          expect(html).toContain("A published discussion");
        } else if (path.startsWith("/l/")) {
          expect(html).toContain(`<title>${node.name} | `);
        }

        if (path !== "/") {
          expect(html).not.toMatch(/<script type="application\/ld\+json">/);
        }

        expect(html).not.toContain("<title>Thread Not Found</title>");
        expect(html).not.toContain("<title>Page not found</title>");
      }
    });
  }
}

test("returns 404 for a missing Library preview rather than a server error", async ({
  request,
}) => {
  const response = await request.get(
    `/l/og?slug=missing-metadata-${Date.now()}`,
  );

  expect(response.status()).toBe(404);
  expect(response.headers()["cache-control"]).toBe("no-store");
});

test("renders guest content without client hydration or a second public read", async ({
  browser,
}) => {
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
    extraHTTPHeaders: { "X-Forwarded-For": "198.51.100.89" },
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
