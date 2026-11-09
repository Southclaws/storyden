import { APIRequestContext, Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import { LinkListOKResponse } from "../../src/api/openapi-schema";
import { API_ADDRESS } from "../../src/config";
import { createAccountWithRole, login } from "../access_key_admin_assignment";

async function searchLinkURLs(request: APIRequestContext, q: string) {
  const response = await request.get(`${API_ADDRESS}/api/links`, {
    params: { q },
  });
  expect(response.ok()).toBe(true);
  const body = (await response.json()) as LinkListOKResponse;
  return body.links.map((link) => link.url);
}

function recordLinkCreates(page: Page) {
  const urls: string[] = [];
  page.on("request", (request) => {
    if (
      request.method() === "POST" &&
      new URL(request.url()).pathname === "/api/links"
    ) {
      urls.push((request.postDataJSON() as { url: string }).url);
    }
  });
  return urls;
}

async function openLinksPage(page: Page, seed: string) {
  const username = `links-${seed}`;
  await createAccountWithRole(
    page.context(),
    username,
    "TestPassword123!",
    "member",
  );
  await login(page, username, "TestPassword123!");
  await page.goto("/links", { waitUntil: "networkidle" });
  return page.getByPlaceholder("Search or paste a new link");
}

test("Links page saves a typed link only after clicking Add link", async ({
  page,
  request,
}) => {
  const seed = randomUUID().slice(0, 8);
  const host = `links-${seed}.invalid`;
  const url = `https://${host}/`;

  const input = await openLinksPage(page, seed);
  const linkCreates = recordLinkCreates(page);

  await input.pressSequentially(`https://links-${seed}`, { delay: 20 });
  await expect(
    page.getByRole("button", { name: "Search", exact: true }),
  ).toBeVisible();
  await page.waitForTimeout(1000);
  expect(linkCreates).toEqual([]);

  await input.pressSequentially(".invalid", { delay: 20 });
  await expect(page.getByRole("button", { name: "Add link" })).toBeVisible();
  await page.waitForTimeout(1000);
  expect(linkCreates).toEqual([]);
  expect(await searchLinkURLs(request, seed)).toEqual([]);

  await page.getByRole("button", { name: "Add link" }).click();

  await expect(
    page.getByRole("main").getByRole("link", { name: url, exact: true }),
  ).toBeVisible();
  expect(linkCreates).toEqual([url]);
  expect(await searchLinkURLs(request, seed)).toEqual([url]);
});

test("Links page adds a bare domain immediately on submit", async ({
  page,
  request,
}) => {
  const seed = randomUUID().slice(0, 8);
  const url = `https://links-${seed}.invalid/`;

  await page.setViewportSize({ width: 390, height: 844 });
  const input = await openLinksPage(page, seed);
  const linkCreates = recordLinkCreates(page);

  await input.fill(`links-${seed}.invalid`);
  await input.press("Enter");

  await expect(
    page.getByRole("main").getByRole("link", { name: url, exact: true }),
  ).toBeVisible();
  await page.waitForTimeout(1000);
  expect(linkCreates).toEqual([url]);
  expect(await searchLinkURLs(request, seed)).toEqual([url]);
  expect(
    await page.evaluate(
      () =>
        document.documentElement.scrollWidth <=
        document.documentElement.clientWidth,
    ),
  ).toBe(true);
});

async function seedLinks(page: Page, urls: string[]) {
  for (const url of urls) {
    const response = await page
      .context()
      .request.post(`${API_ADDRESS}/api/links`, {
        data: { url },
      });
    expect(response.ok()).toBe(true);
  }
}

function linkResults(page: Page) {
  return page.getByRole("region", { name: "Link results" });
}

test("Links search filters text and pasted URLs, survives reload, and clears", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const host = `flow-${seed}.invalid`;
  const alpha = `https://${host}/alpha`;
  const beta = `https://${host}/beta`;
  const input = await openLinksPage(page, seed);
  await seedLinks(page, [alpha, beta]);
  await page.reload({ waitUntil: "networkidle" });
  const creates = recordLinkCreates(page);

  await input.fill(`flow-${seed}`);
  await expect(linkResults(page).getByRole("article")).toHaveCount(2);
  await expect(
    page.getByRole("button", { name: "Search", exact: true }),
  ).toBeVisible();

  await input.fill(alpha);
  await expect(linkResults(page).getByRole("article")).toHaveCount(1);
  await expect(
    linkResults(page).getByRole("link", { name: alpha, exact: true }),
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "Saved", exact: true }),
  ).toBeDisabled();
  await input.press("Enter");
  expect(creates).toEqual([]);

  await page.reload({ waitUntil: "networkidle" });
  await expect(input).toHaveValue(alpha);
  await expect(linkResults(page).getByRole("article")).toHaveCount(1);
  expect(creates).toEqual([]);

  await page.getByRole("button", { name: "Clear search" }).click();
  await expect(input).toHaveValue("");
  await expect(
    linkResults(page).getByRole("link", { name: beta, exact: true }),
  ).toBeVisible();
  await expect(page).toHaveURL(/\/links$/);
});

test("Links pagination follows the active search and browser history", async ({
  page,
  request,
}) => {
  const seed = randomUUID().slice(0, 8);
  const input = await openLinksPage(page, seed);
  const initialResponse = await request.get(`${API_ADDRESS}/api/links`);
  const initial = (await initialResponse.json()) as LinkListOKResponse;
  const host = `pages-${seed}.invalid`;
  const urls = Array.from(
    { length: initial.page_size + 1 },
    (_, index) => `https://${host}/${index.toString().padStart(3, "0")}`,
  );
  await seedLinks(page, urls);
  const creates = recordLinkCreates(page);

  await input.fill(`pages-${seed}`);
  await expect(linkResults(page).getByRole("article")).toHaveCount(
    initial.page_size,
  );
  await linkResults(page).getByRole("link", { name: "2", exact: true }).click();
  await expect(linkResults(page).getByRole("article")).toHaveCount(1);
  await expect(page).toHaveURL(/page=2/);

  await page.goBack();
  await expect(input).toHaveValue(`pages-${seed}`);
  await expect(linkResults(page).getByRole("article")).toHaveCount(
    initial.page_size,
  );
  await page.goForward();
  await expect(linkResults(page).getByRole("article")).toHaveCount(1);

  await input.fill(urls[0]!);
  await expect(page).not.toHaveURL(/page=2/);
  await expect(
    linkResults(page).getByRole("link", { name: urls[0]!, exact: true }),
  ).toBeVisible();
  await expect(linkResults(page).getByRole("article")).toHaveCount(1);
  await expect(
    linkResults(page).getByRole("link", { name: "2", exact: true }),
  ).toHaveCount(0);
  expect(creates).toEqual([]);

  await page.getByRole("button", { name: "Clear search" }).click();
  await expect(page).toHaveURL(/\/links$/);
  await expect(
    linkResults(page).getByRole("link", { name: "2", exact: true }),
  ).toBeVisible();
});

test("Links can retry a failed Add link without losing the search", async ({
  page,
  request,
}) => {
  const seed = randomUUID().slice(0, 8);
  const url = `https://retry-${seed}.invalid/`;
  const input = await openLinksPage(page, seed);
  const creates = recordLinkCreates(page);
  let failNext = true;
  await page.route("**/api/links", async (route) => {
    if (route.request().method() === "POST" && failNext) {
      failNext = false;
      await route.abort("failed");
      return;
    }

    await route.continue();
  });

  await input.fill(url);
  await expect(
    linkResults(page).getByText("No links were found."),
  ).toBeVisible();
  await page.getByRole("button", { name: "Add link", exact: true }).click();
  await expect(
    page.getByRole("alert").filter({ hasText: "Could not add link" }),
  ).toBeVisible();
  await expect(input).toHaveValue(url);
  await page.getByRole("button", { name: "Add link", exact: true }).click();
  await expect(
    linkResults(page).getByRole("link", { name: url, exact: true }),
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "Saved", exact: true }),
  ).toBeDisabled();
  expect(creates).toEqual([url, url]);
  expect(await searchLinkURLs(request, seed)).toEqual([url]);
});

test("Links can recover from a search error without showing old pagination", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const url = `https://search-${seed}.invalid/`;
  const input = await openLinksPage(page, seed);
  await seedLinks(page, [url]);
  let failNext = true;
  await page.route("**/api/links?**", async (route) => {
    const query = new URL(route.request().url()).searchParams.get("q");
    if (query === seed && failNext) {
      failNext = false;
      await route.fulfill({
        status: 503,
        contentType: "application/json",
        body: JSON.stringify({ title: "Temporarily unavailable" }),
      });
      return;
    }

    await route.continue();
  });

  await input.fill(seed);
  await expect(linkResults(page).getByRole("alert")).toHaveText(
    "Could not search links.",
  );
  await expect(linkResults(page).getByRole("article")).toHaveCount(0);
  await expect(
    linkResults(page).getByRole("link", { name: "2", exact: true }),
  ).toHaveCount(0);
  await linkResults(page).getByRole("button", { name: "Try again" }).click();
  await expect(
    linkResults(page).getByRole("link", { name: url, exact: true }),
  ).toBeVisible();
});

test("Links ignores a completed add after the user starts another search", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const addedURL = `https://slow-${seed}.invalid/`;
  const existingURL = `https://existing-${seed}.invalid/`;
  const input = await openLinksPage(page, seed);
  await seedLinks(page, [existingURL]);
  const creates = recordLinkCreates(page);
  let release!: () => void;
  const gate = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("**/api/links", async (route) => {
    if (route.request().method() === "POST") {
      await gate;
    }

    await route.continue();
  });

  await input.fill(addedURL);
  await page.getByRole("button", { name: "Add link", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "Adding…", exact: true }),
  ).toBeDisabled();
  await input.press("Enter");
  expect(creates).toEqual([addedURL]);

  await input.fill(existingURL);
  await expect(
    linkResults(page).getByRole("link", { name: existingURL, exact: true }),
  ).toBeVisible();
  const completed = page.waitForResponse(
    (response) =>
      response.request().method() === "POST" &&
      new URL(response.url()).pathname === "/api/links",
  );
  release();
  expect((await completed).ok()).toBe(true);
  await expect(input).toHaveValue(existingURL);
  await expect(linkResults(page).getByRole("article")).toHaveCount(1);
  await expect(
    linkResults(page).getByRole("link", { name: addedURL, exact: true }),
  ).toHaveCount(0);
});
