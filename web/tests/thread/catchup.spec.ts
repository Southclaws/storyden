import { type Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import type { ThreadListResult } from "../../src/api/openapi-schema";
import { API_ADDRESS } from "../../src/config";
import {
  getAdminAccessKey,
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

  await ssrPage.goto("/catchup");
  const content = ssrPage.locator('[data-sd-page="catchup"]:visible').first();
  await expect(
    content.getByRole("heading", { name: "Catch Up", exact: true }),
  ).toBeVisible();
  await expect(
    content.getByRole("link", { name: unread.title, exact: true }),
  ).toBeVisible();
  await expect(
    content.getByRole("link", { name: caughtUp.title, exact: true }),
  ).toHaveCount(0);
  await expect(
    ssrPage.locator('[data-navigation-anchor="catchup"]'),
  ).toHaveAttribute("href", "/catchup");
  expect(browserAPIRequests).toEqual([]);
});

for (const viewport of [
  { width: 1280, height: 900 },
  { width: 390, height: 844 },
]) {
  test(`keeps read threads until refresh and records visits at ${viewport.width}px`, async ({
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

    await page.locator('[data-navigation-anchor="catchup"]:visible').click();
    await expect(page).toHaveURL(/\/catchup$/);
    const content = page.locator('[data-sd-page="catchup"]:visible').first();
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
    await expect(card).toBeVisible();
    await expect(card.getByRole("heading")).toHaveAccessibleName(
      `${fromList.title} Caught up`,
    );
    await expect(
      content.getByRole("link", { name: byVisit.title, exact: true }),
    ).toBeVisible();
    await page.screenshot({
      path: testInfo.outputPath("catchup-page.png"),
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
    await content.getByRole("button", { name: "Refresh", exact: true }).click();
    await expect(
      content.getByText("You're all caught up.", { exact: true }),
    ).toBeVisible();
  });
}

test("keeps the last page stable until refresh then recovers to the remaining page", async ({
  page,
}) => {
  await registerUser(page, "unread-pages", "TestPassword123!");
  await markAllRead(page);
  const category = await withAdminAccessKey(
    async ({ categoryCreate, threadCreate }) => {
      const category = await categoryCreate({
        name: `Pagination ${randomUUID().slice(0, 8)}`,
        description: "Catch Up pagination fixtures.",
        colour: "#336699",
      });
      const prefix = `Unread pagination ${randomUUID().slice(0, 8)}`;
      for (let index = 0; index < 51; index++) {
        await threadCreate({
          title: `${prefix} ${index}`,
          category: category.id,
          body: "A paginated unread discussion.",
          visibility: "published",
        });
      }
      return category;
    },
  );

  await page.goto(`/catchup?category=${category.slug}&page=2`);
  const content = page.locator('[data-sd-page="catchup"]:visible').first();
  await expect(content.locator("ol").first().getByRole("heading")).toHaveCount(
    1,
  );
  await content.getByRole("button", { name: /^More options for / }).click();
  await page
    .getByRole("menuitem", { name: "Mark as read", exact: true })
    .click();
  await expect(page).toHaveURL(
    new RegExp(`/catchup\\?category=${category.slug}&page=2$`),
  );
  await expect(content.locator("ol").first().getByRole("heading")).toHaveCount(
    1,
  );
  await content.getByRole("button", { name: "Refresh", exact: true }).click();
  await expect(page).toHaveURL(
    new RegExp(`/catchup\\?category=${category.slug}$`),
  );
  await expect(content.locator("ol").first().getByRole("heading")).toHaveCount(
    50,
  );
});

test("requires sign-in and hides the Catch Up navigation entry for guests", async ({
  page,
}) => {
  await page.goto("/unread?view=unresolved&category=general&page=2");
  await expect(page).toHaveURL(
    "/catchup?view=unresolved&category=general&page=2",
  );
  await expect(
    page.getByText("Please log in to see this page.", { exact: true }),
  ).toBeVisible();
  await expect(page.locator('[data-navigation-anchor="catchup"]')).toHaveCount(
    0,
  );
});

for (const width of [1280, 390]) {
  test(`bookmarks and switches category presets with the keyboard at ${width}px`, async ({
    page,
    context,
    request,
  }, testInfo) => {
    await page.setViewportSize({ width, height: 900 });
    await registerUser(page, "catchup-filters", "TestPassword123!");
    await markAllRead(page);

    const { category, emptyCategory, openRead, openUnread, resolved, outside } =
      await withAdminAccessKey(
        async ({ categoryCreate, threadCreate, replyCreate }) => {
          const category = await categoryCreate({
            name: `Community${randomUUID().replaceAll("-", "")}${randomUUID().replaceAll("-", "").slice(0, 23)}`,
            description: "Discussions for the Catch Up view.",
            colour: "#336699",
          });
          const emptyCategory = await categoryCreate({
            name: `Empty ${randomUUID().slice(0, 8)}`,
            description: "An empty category.",
            colour: "#336699",
          });
          const openRead = await threadCreate({
            title: `Read but unresolved ${randomUUID().slice(0, 8)}`,
            category: category.id,
            visibility: "published",
          });
          await replyCreate(openRead.slug, {
            body: "This thread has an answer, but is still open.",
          });
          const openUnread = await threadCreate({
            title: `Unread and unresolved ${randomUUID().slice(0, 8)}`,
            category: category.id,
            visibility: "published",
          });
          const resolved = await threadCreate({
            title: `Unread but resolved ${randomUUID().slice(0, 8)}`,
            category: category.id,
            visibility: "published",
          });
          const outside = await threadCreate({
            title: `Another category ${randomUUID().slice(0, 8)}`,
            visibility: "published",
          });
          return {
            category,
            emptyCategory,
            openRead,
            openUnread,
            resolved,
            outside,
          };
        },
      );

    const marked = await page.request.patch(
      `${API_ADDRESS}/api/threads/${openRead.slug}/read-status`,
      { data: { status: "read" } },
    );
    expect(marked.ok()).toBe(true);
    const closed = await request.patch(
      `${API_ADDRESS}/api/threads/${resolved.slug}`,
      {
        headers: { Authorization: `Bearer ${getAdminAccessKey()}` },
        data: { lifecycle_state: "resolved" },
      },
    );
    expect(closed.ok()).toBe(true);

    const bookmark = `/catchup?view=unresolved&category=${category.slug}`;
    // The selected preset and category must also work without hydration.
    const ssrPage = await context.newPage();
    await ssrPage.route("**/*", async (route) => {
      if (route.request().resourceType() === "script") {
        await route.abort();
        return;
      }
      await route.continue();
    });
    await ssrPage.goto(bookmark);
    const ssrContent = ssrPage
      .locator('[data-sd-page="catchup"]:visible')
      .first();
    await expect(
      ssrContent.getByRole("tab", { name: "Unresolved", selected: true }),
    ).toBeVisible();
    await expect(
      ssrContent.getByRole("link", { name: openRead.title, exact: true }),
    ).toBeVisible();
    await expect(
      ssrContent.getByRole("link", { name: resolved.title, exact: true }),
    ).toHaveCount(0);
    await expect(
      ssrContent.getByRole("link", { name: outside.title, exact: true }),
    ).toHaveCount(0);
    await ssrPage.close();

    await page.goto(bookmark);
    const content = page.locator('[data-sd-page="catchup"]:visible').first();
    const categorySelect = content.getByRole("combobox", {
      name: "Filter by category",
    });
    await expect(categorySelect).toContainText(category.name);
    expect(
      await categorySelect.evaluate(
        (element) => element.scrollWidth <= element.clientWidth,
      ),
    ).toBe(true);
    await expect(
      content.getByRole("link", { name: openRead.title, exact: true }),
    ).toBeVisible();
    await expect(
      content.getByRole("link", { name: openUnread.title, exact: true }),
    ).toBeVisible();

    await content.getByRole("tab", { name: "Unresolved" }).focus();
    await page.keyboard.press("ArrowLeft");
    await expect(content.getByRole("tab", { name: "Unread" })).toBeFocused();
    await page.keyboard.press("Enter");
    await expect(
      content.getByRole("tab", { name: "Unread", selected: true }),
    ).toBeVisible();
    await expect(page).toHaveURL(`/catchup?category=${category.slug}`);
    await expect(
      content.getByRole("link", { name: openRead.title, exact: true }),
    ).toHaveCount(0);
    await expect(
      content.getByRole("link", { name: resolved.title, exact: true }),
    ).toBeVisible();

    await categorySelect.focus();
    await page.keyboard.press("Enter");
    await page
      .getByRole("option", { name: emptyCategory.name, exact: true })
      .click();
    await expect(page).toHaveURL(`/catchup?category=${emptyCategory.slug}`);
    await expect(
      content.getByText(
        `There are no unread threads in ${emptyCategory.name}.`,
        { exact: true },
      ),
    ).toBeVisible();
    await content.getByRole("tab", { name: "Unresolved" }).click();
    await expect(
      content.getByText(
        `There are no unresolved threads in ${emptyCategory.name}.`,
        { exact: true },
      ),
    ).toBeVisible();

    await categorySelect.click();
    await page
      .getByRole("option", { name: "All categories", exact: true })
      .click();
    await expect(page).toHaveURL("/catchup?view=unresolved");
    await expect(
      content.getByRole("link", { name: outside.title, exact: true }),
    ).toBeVisible();
    await page.goBack();
    await expect(categorySelect).toContainText(emptyCategory.name);
    await expect(
      content.getByRole("tab", { name: "Unresolved", selected: true }),
    ).toBeVisible();

    await page.goto(bookmark);
    await page.screenshot({
      path: testInfo.outputPath("catchup-unresolved.png"),
      fullPage: true,
      animations: "disabled",
    });
    expect(
      await page.evaluate(
        () => document.documentElement.scrollWidth <= window.innerWidth,
      ),
    ).toBe(true);
  });
}
