import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import { API_ADDRESS } from "../../src/config";
import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("Directory layouts save without changing the visible fields", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const parent = await withAdminAccessKey(async ({ nodeCreate }) => {
    const parent = await nodeCreate({
      name: `Directory ${seed}`,
      slug: `directory-${seed}`,
      visibility: "published",
      hide_child_tree: true,
      meta: {
        layout: {
          blocks: [
            { type: "title" },
            {
              type: "directory",
              config: {
                layout: "table",
                columns: [
                  { fid: "fixed:link", hidden: true },
                  { fid: "fixed:name", hidden: false },
                  { fid: "fixed:description", hidden: false },
                  { fid: "fixed:primary_image", hidden: true },
                ],
              },
            },
          ],
        },
      },
    });
    await nodeCreate({
      name: `Child ${seed}`,
      slug: `child-${seed}`,
      description: "A text-only entry in the same directory.",
      parent: parent.id,
      visibility: "published",
    });
    return parent;
  });
  const handle = `directory_layout_${seed}`;
  await createAdmin(page.context(), handle, "TestPassword123!");
  await login(page, handle, "TestPassword123!");
  await page.goto(`/l/${parent.slug}`);

  for (const mode of ["Grid", "Gallery", "Table"] as const) {
    await page.getByRole("button", { name: "Edit", exact: true }).click();
    await page
      .getByRole("menuitem", { name: "Quick edit", exact: true })
      .click();
    await page.locator('[id="block-directory_gutter-drag-handle"]').click();
    await page.getByRole("menuitem", { name: "Layout", exact: true }).hover();
    await page.getByRole("menuitemradio", { name: mode, exact: true }).click();
    // Exit immediately: the selected layout must survive the pending debounce.
    await page.getByRole("button", { name: "View", exact: true }).click();
    await expect(
      page.getByRole("button", { name: "View", exact: true }),
    ).toBeHidden();
    await page.waitForURL((url) => !url.searchParams.has("edit"));
    await page.reload();
    const response = await page.request.get(
      `${API_ADDRESS}/api/nodes/${parent.slug}`,
    );
    const node = await response.json();
    expect(node.meta.layout.blocks[1].config).toEqual({
      layout: mode.toLowerCase(),
      columns: [
        { fid: "fixed:link", hidden: true },
        { fid: "fixed:name", hidden: false },
        { fid: "fixed:description", hidden: false },
        { fid: "fixed:primary_image", hidden: true },
      ],
    });
    await expect(
      page.getByRole("link", { name: `Child ${seed}`, exact: true }).first(),
    ).toBeVisible();
    if (mode !== "Table") {
      const list = page.getByRole("list", { name: "Library pages" });
      await expect(list).toHaveAttribute("data-layout", mode.toLowerCase());
      await expect(list.locator("img")).toHaveCount(0);
    }
  }
});
