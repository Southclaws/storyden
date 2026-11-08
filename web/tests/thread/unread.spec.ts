import { type Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import type { ThreadListResult } from "../../src/api/openapi-schema";
import { API_ADDRESS } from "../../src/config";
import {
  registerUser,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

async function markAllRead(page: Page) {
  while (true) {
    const response = await page.request.get(
      `${API_ADDRESS}/api/threads?read_status=unread`,
    );
    expect(response.ok()).toBe(true);
    const data: ThreadListResult = await response.json();
    if (data.threads.length === 0) {
      return;
    }

    for (const thread of data.threads) {
      const marked = await page.request.patch(
        `${API_ADDRESS}/api/threads/${thread.slug}/read-status`,
        { data: { status: "read" } },
      );
      expect(marked.ok()).toBe(true);
    }
  }
}

test("server renders unread threads and member navigation without browser API requests", async ({
  page,
  context,
}) => {
  await registerUser(page, "unread-ssr", "TestPassword123!");
  await markAllRead(page);
  const [unread, caughtUp] = await withAdminAccessKey(
    async ({ threadCreate }) => {
      const unread = await threadCreate({
        title: `Unread server render ${randomUUID().slice(0, 8)}`,
        body: "This thread arrives in the server-rendered HTML.",
        visibility: "published",
      });
      const caughtUp = await threadCreate({
        title: `Already read ${randomUUID().slice(0, 8)}`,
        body: "This thread should be excluded.",
        visibility: "published",
      });
      return [unread, caughtUp];
    },
  );
  const marked = await page.request.patch(
    `${API_ADDRESS}/api/threads/${caughtUp.slug}/read-status`,
    { data: { status: "read" } },
  );
  expect(marked.ok()).toBe(true);

  const ssrPage = await context.newPage();
  const browserAPIRequests: string[] = [];
  await ssrPage.route("**/*", async (route) => {
    const request = route.request();
    if (request.resourceType() === "script") {
      await route.abort();
      return;
    }

    const isBrowserFetch = ["fetch", "xhr"].includes(request.resourceType());
    if (isBrowserFetch && request.url().startsWith(`${API_ADDRESS}/api/`)) {
      browserAPIRequests.push(request.url());
      await route.abort();
      return;
    }

    await route.continue();
  });

  await ssrPage.goto("/unread");
  const content = ssrPage.locator('[data-sd-page="unread"]:visible').first();
  await expect(
    content.getByRole("heading", { name: "Unread threads", exact: true }),
  ).toBeVisible();
  await expect(
    content.getByRole("link", { name: unread.title, exact: true }),
  ).toBeVisible();
  await expect(
    content.getByRole("link", { name: caughtUp.title, exact: true }),
  ).toHaveCount(0);
  await expect(
    ssrPage.locator('[data-navigation-anchor="unread"]'),
  ).toHaveAttribute("href", "/unread");
  expect(browserAPIRequests).toEqual([]);
});

for (const viewport of [
  { width: 1280, height: 900 },
  { width: 390, height: 844 },
]) {
  test(`unread list stays current after marking read and visiting at ${viewport.width}px`, async ({
    page,
  }, testInfo) => {
    await page.setViewportSize(viewport);
    await registerUser(page, "unread-flow", "TestPassword123!");
    await markAllRead(page);
    const { fromList, byVisit } = await withAdminAccessKey(
      async ({ threadCreate }) => {
        const fromList = await threadCreate({
          title: `Read from the list ${randomUUID().slice(0, 8)}`,
          body: "A discussion to catch up on.",
          visibility: "published",
        });
        const byVisit = await threadCreate({
          title: `Read by visiting ${randomUUID().slice(0, 8)}`,
          body: "A discussion to catch up on.",
          visibility: "published",
        });
        return { fromList, byVisit };
      },
    );

    if (viewport.width < 768) {
      await page
        .getByRole("button", { name: "Open navigation menu", exact: true })
        .click();
    }

    await page.locator('[data-navigation-anchor="unread"]:visible').click();
    await expect(page).toHaveURL(/\/unread$/);
    const content = page.locator('[data-sd-page="unread"]:visible').first();
    const card = content.locator(`[id="${fromList.id}"]`);
    await expect(card).toBeVisible();
    await card
      .getByRole("button", {
        name: `More options for ${fromList.title}`,
        exact: true,
      })
      .click();
    await page
      .getByRole("menuitem", { name: "Mark as read", exact: true })
      .click();
    await expect(card).toHaveCount(0);
    await expect(
      content.getByRole("link", { name: byVisit.title, exact: true }),
    ).toBeVisible();
    await page.screenshot({
      path: testInfo.outputPath("unread-page.png"),
      animations: "disabled",
      fullPage: true,
    });

    await content
      .getByRole("link", { name: byVisit.title, exact: true })
      .click();
    await expect(page).toHaveURL(new RegExp(`/t/${byVisit.slug}`));
    await expect
      .poll(async () => {
        const response = await page.request.get(
          `${API_ADDRESS}/api/threads/${byVisit.slug}`,
        );
        return (await response.json()).read_status?.last_read_at;
      })
      .toBeTruthy();
    await page.goBack();
    await expect(
      content.getByText("You're all caught up.", { exact: true }),
    ).toBeVisible();
  });
}

test("moves to the remaining page when the last unread item on page two is read", async ({
  page,
}) => {
  await registerUser(page, "unread-pages", "TestPassword123!");
  await markAllRead(page);
  await withAdminAccessKey(async ({ threadCreate }) => {
    const prefix = `Unread pagination ${randomUUID().slice(0, 8)}`;
    for (let index = 0; index < 51; index++) {
      await threadCreate({
        title: `${prefix} ${index}`,
        body: "A paginated unread discussion.",
        visibility: "published",
      });
    }
  });

  await page.goto("/unread?page=2");
  const content = page.locator('[data-sd-page="unread"]:visible').first();
  await expect(content.locator("ol").first().getByRole("heading")).toHaveCount(
    1,
  );
  await content.getByRole("button", { name: /^More options for / }).click();
  await page
    .getByRole("menuitem", { name: "Mark as read", exact: true })
    .click();
  await expect(page).toHaveURL(/\/unread\?page=1$/);
  await expect(content.locator("ol").first().getByRole("heading")).toHaveCount(
    50,
  );
});

test("requires sign-in and hides the unread navigation entry for guests", async ({
  page,
}) => {
  await page.goto("/unread");
  await expect(
    page.getByText("Please log in to see this page.", { exact: true }),
  ).toBeVisible();
  await expect(page.locator('[data-navigation-anchor="unread"]')).toHaveCount(
    0,
  );
});
