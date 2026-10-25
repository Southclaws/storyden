import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

for (const layout of ["grid", "gallery"] as const) {
  test(`Directory ${layout} filters entries and persists a drag reorder`, async ({
    page,
  }) => {
    const seed = randomUUID().slice(0, 8);
    const fixture = await withAdminAccessKey(async ({ nodeCreate }) => {
      const parent = await nodeCreate({
        name: `Directory ordering ${seed}`,
        slug: `directory-ordering-${seed}`,
        visibility: "published",
        hide_child_tree: true,
        meta: {
          layout: {
            blocks: [
              { type: "title" },
              {
                type: "directory",
                config: {
                  layout,
                  columns: [
                    { fid: "fixed:name", hidden: false },
                    { fid: "fixed:description", hidden: false },
                    { fid: "fixed:link", hidden: true },
                    { fid: "fixed:primary_image", hidden: false },
                  ],
                },
              },
            ],
          },
        },
      });
      for (const name of ["Alpha", "Bravo", "Charlie", "Delta"]) {
        await nodeCreate({
          name: `${name} ${seed}`,
          slug: `directory-${name.toLowerCase()}-${seed}`,
          parent: parent.id,
          visibility: "published",
          description:
            name === "Bravo"
              ? "A description to vary card content."
              : undefined,
        });
      }
      return parent;
    });
    const handle = `directory_order_${seed}`;
    await createAdmin(page.context(), handle, "TestPassword123!");
    await login(page, handle, "TestPassword123!");
    await page.setViewportSize({ width: 1440, height: 1000 });
    await page.goto(`/l/${fixture.slug}`);
    const list = page.getByRole("list", { name: "Library pages" });
    await expect(list.locator(":scope > li")).toHaveCount(4);
    await page
      .getByRole("textbox", { name: "Search...", exact: true })
      .fill("Bravo");
    await expect(list.locator(":scope > li")).toHaveCount(1);
    await expect(
      list.getByRole("link", { name: `Bravo ${seed}`, exact: true }).first(),
    ).toBeVisible();
    await page
      .getByRole("textbox", { name: "Search...", exact: true })
      .fill("");
    await expect(list.locator(":scope > li")).toHaveCount(4);

    await page.getByRole("button", { name: "Edit", exact: true }).click();
    await page
      .getByRole("menuitem", { name: "Quick edit", exact: true })
      .click();
    const source = list.getByRole("button", {
      name: `Move Bravo ${seed}`,
      exact: true,
    });
    const target = list.getByRole("button", {
      name: `Move Alpha ${seed}`,
      exact: true,
    });
    const sourceBox = await source.boundingBox();
    const targetBox = await target.boundingBox();
    expect(sourceBox).not.toBeNull();
    expect(targetBox).not.toBeNull();
    const moved = page.waitForResponse(
      (response) =>
        response.request().method() === "PATCH" &&
        response.url().includes("/position") &&
        response.ok(),
    );
    await page.mouse.move(
      sourceBox!.x + sourceBox!.width / 2,
      sourceBox!.y + sourceBox!.height / 2,
    );
    await page.mouse.down();
    await page.mouse.move(
      sourceBox!.x + sourceBox!.width / 2 + 8,
      sourceBox!.y + sourceBox!.height / 2,
      { steps: 4 },
    );
    await page.mouse.move(
      targetBox!.x + targetBox!.width / 2,
      targetBox!.y + targetBox!.height / 2,
      { steps: 16 },
    );
    await page.mouse.up();
    await moved;
    await page.goto(`/l/${fixture.slug}`);
    await expect(list.locator(":scope > li")).toHaveCount(4);
    const names = await list
      .locator(":scope > li")
      .evaluateAll((cards) =>
        cards.map((card) =>
          card.querySelector("a[aria-label]")?.getAttribute("aria-label"),
        ),
      );
    expect(names).toEqual([
      `Bravo ${seed}`,
      `Alpha ${seed}`,
      `Charlie ${seed}`,
      `Delta ${seed}`,
    ]);
  });
}
