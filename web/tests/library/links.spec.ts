import { APIRequestContext, expect, test } from "@playwright/test";
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

test("Links page only creates a link when the user submits it", async ({
  page,
  request,
}) => {
  const seed = randomUUID().slice(0, 8);
  const username = `links-${seed}`;
  const url = `https://links-${seed}.invalid`;

  await createAccountWithRole(
    page.context(),
    username,
    "TestPassword123!",
    "member",
  );
  await login(page, username, "TestPassword123!");

  const linkCreates: string[] = [];
  page.on("request", (request) => {
    if (
      request.method() === "POST" &&
      new URL(request.url()).pathname === "/api/links"
    ) {
      linkCreates.push(request.postData() ?? "");
    }
  });

  await page.goto("/links", { waitUntil: "networkidle" });

  const input = page.getByPlaceholder("Search or paste a new link");
  await input.pressSequentially(url, { delay: 20 });
  await expect(input).toHaveValue(url);

  const submit = page.getByRole("button", { name: "Add link" });
  await expect(submit).toBeVisible();
  await page.waitForLoadState("networkidle");

  expect(linkCreates).toHaveLength(0);
  expect(await searchLinkURLs(request, seed)).toHaveLength(0);

  await submit.click();

  await expect.poll(() => linkCreates.length).toBe(1);
  await expect
    .poll(() => searchLinkURLs(request, seed))
    .toEqual([expect.stringContaining(`links-${seed}.invalid`)]);
});
