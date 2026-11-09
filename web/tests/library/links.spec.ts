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

test("Links page previews a link once the typed URL is complete", async ({
  page,
  request,
}) => {
  const seed = randomUUID().slice(0, 8);
  const host = `links-${seed}.invalid`;
  const url = `https://${host}/`;

  const input = await openLinksPage(page, seed);
  const linkCreates = recordLinkCreates(page);

  await input.pressSequentially(`https://links-${seed}`, { delay: 20 });
  await expect(page.getByRole("button", { name: "Search" })).toBeVisible();
  await page.waitForTimeout(1000);
  expect(linkCreates).toEqual([]);

  await input.pressSequentially(".invalid", { delay: 20 });
  await expect(page.getByRole("button", { name: "Add link" })).toBeVisible();

  await expect(page.getByRole("main").getByText(url).first()).toBeVisible();
  expect(linkCreates).toEqual([url]);
  expect(await searchLinkURLs(request, seed)).toEqual([url]);
});

test("Links page adds a bare domain immediately on submit", async ({
  page,
  request,
}) => {
  const seed = randomUUID().slice(0, 8);
  const url = `https://links-${seed}.invalid/`;

  const input = await openLinksPage(page, seed);
  const linkCreates = recordLinkCreates(page);

  await input.fill(`links-${seed}.invalid`);
  await input.press("Enter");

  await expect(page.getByRole("main").getByText(url).first()).toBeVisible();
  await page.waitForTimeout(1000);
  expect(linkCreates).toEqual([url]);
  expect(await searchLinkURLs(request, seed)).toEqual([url]);
});
