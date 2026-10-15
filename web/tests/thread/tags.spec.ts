import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("thread tags remain readable and the deferred picker saves edits", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const tag = `original-${seed}`;
  const addedTag = `added-${seed}`;
  const thread = await withAdminAccessKey(({ threadCreate }) =>
    threadCreate({
      title: `Thread tags ${seed}`,
      body: "A readable thread with editable tags.",
      tags: [tag],
      visibility: "published",
    }),
  );
  const admin = `thread-tags-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  await page.mouse.move(0, 0);
  await page.goto(`/t/${thread.slug}`, { waitUntil: "networkidle" });
  const main = page.getByRole("main");
  await expect(
    main.getByRole("link", { name: tag, exact: true }),
  ).toBeVisible();
  await expect(main.getByRole("combobox")).toHaveCount(0);
  await page.goto(`/t/${thread.slug}?edit=true`);
  const input = main.getByRole("combobox", { name: "Search for items" });
  await input.fill(addedTag);
  await input.press("Enter");
  await input.press("Escape");
  await main.getByRole("button", { name: "Save", exact: true }).click();
  await expect(
    main.getByRole("link", { name: addedTag, exact: true }),
  ).toBeVisible();
  await page.reload();
  await expect(
    main.getByRole("link", { name: tag, exact: true }),
  ).toBeVisible();
  await expect(
    main.getByRole("link", { name: addedTag, exact: true }),
  ).toBeVisible();
  await expect(main.getByRole("combobox")).toHaveCount(0);
});
