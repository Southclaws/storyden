import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("Library retains tag badges while the picker loads and persists editing", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const tag = `original-${seed}`;
  const addedTag = `added-${seed}`;
  const node = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: `Library tags ${seed}`,
      slug: `library-tags-${seed}`,
      visibility: "published",
      tags: [tag],
      meta: { layout: { blocks: [{ type: "title" }, { type: "tags" }] } },
    }),
  );
  const admin = `library-tags-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  await page.mouse.move(0, 0);
  await page.goto(`/l/${node.slug}`, { waitUntil: "networkidle" });
  const main = page.getByRole("main");
  await expect(
    main.getByRole("link", { name: tag, exact: true }),
  ).toBeVisible();
  await expect(main.getByRole("combobox")).toHaveCount(0);
  await page.getByRole("button", { name: "Edit", exact: true }).click();
  const quickEdit = page.getByRole("menuitem", {
    name: "Quick edit",
    exact: true,
  });
  await expect(quickEdit).toBeEnabled();
  let release!: () => void;
  const ready = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("**/_next/static/**/*.js", async (route) => {
    await ready;
    await route.continue();
  });
  const tags = page.locator("#block-tags_container:visible");
  try {
    await quickEdit.click();
    await expect(
      tags.getByRole("link", { name: tag, exact: true }),
    ).toBeVisible();
    await expect(tags.getByRole("combobox")).toHaveCount(0);
  } finally {
    release();
  }
  await page.unrouteAll({ behavior: "wait" });
  const input = tags.getByRole("combobox", { name: "Search for items" });
  await input.fill(addedTag);
  await input.press("Enter");
  await input.press("Escape");
  await expect
    .poll(async () => {
      const response = await page.request.get(
        `http://localhost:8001/api/nodes/${node.slug}`,
      );
      const saved = await response.json();
      return saved.tags.map((item: { name: string }) => item.name);
    })
    .toEqual([tag, addedTag]);
  await page.getByRole("button", { name: "View", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "Edit", exact: true }),
  ).toBeVisible();
  await expect
    .poll(() => new URL(page.url()).searchParams.get("edit"))
    .toBeNull();
  await page.reload();
  await expect(
    main.getByRole("link", { name: tag, exact: true }),
  ).toBeVisible();
  await expect(
    main.getByRole("link", { name: addedTag, exact: true }),
  ).toBeVisible();
  await expect(main.getByRole("combobox")).toHaveCount(0);
});
