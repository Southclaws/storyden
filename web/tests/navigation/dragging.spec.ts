import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("sidebar category and Library drags preserve nesting and disclosure", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  await withAdminAccessKey(async ({ nodeCreate }) => {
    for (let index = 0; index < 24; index++) {
      await nodeCreate({
        name: `Sidebar overflow ${seed} ${index}`,
        visibility: "published",
      });
    }
  });
  const { category, targetCategory, node, targetNode } =
    await withAdminAccessKey(async ({ categoryCreate, nodeCreate }) => ({
      category: await categoryCreate({
        name: `Source category ${seed}`,
        colour: "#2255aa",
        description: "Drag source",
      }),
      targetCategory: await categoryCreate({
        name: `Target category ${seed}`,
        colour: "#aa5522",
        description: "Drag target",
      }),
      node: await nodeCreate({
        name: `Source page ${seed}`,
        visibility: "published",
      }),
      targetNode: await nodeCreate({
        name: `Target page ${seed}`,
        visibility: "published",
      }),
    }));
  const admin = `tree-drag-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");

  for (const item of [
    {
      tree: "Discussion categories",
      source: category.name,
      target: targetCategory.name,
      path: `/api/categories/${category.slug}/position`,
      targetId: targetCategory.id,
    },
    {
      tree: "Library pages",
      source: node.name,
      target: targetNode.name,
      path: `/api/nodes/${node.id}/position`,
      targetId: targetNode.id,
    },
  ]) {
    const tree = page.getByRole("tree", { name: item.tree });
    const source = tree.getByRole("link", { name: item.source, exact: true });
    const target = tree.getByRole("link", { name: item.target, exact: true });
    // Keep the drop target away from the auto-scroll edge so it does not move
    // out from under the measured mouse coordinates during the drag.
    await target.evaluate((element) =>
      element.scrollIntoView({ block: "center", inline: "nearest" }),
    );
    await source.scrollIntoViewIfNeeded();
    const sourceBox = await source.boundingBox();
    const targetBox = await target.boundingBox();
    expect(sourceBox).not.toBeNull();
    expect(targetBox).not.toBeNull();
    const updated = page.waitForResponse(
      (response) =>
        new URL(response.url()).pathname === item.path &&
        response.request().method() === "PATCH",
    );
    await page.mouse.move(
      sourceBox!.x + 20,
      sourceBox!.y + sourceBox!.height / 2,
    );
    await page.mouse.down();
    await page.mouse.move(
      sourceBox!.x + 25,
      sourceBox!.y + sourceBox!.height / 2,
      { steps: 4 },
    );
    await page.mouse.move(
      targetBox!.x + 30,
      targetBox!.y + targetBox!.height / 2,
      { steps: 12 },
    );
    await page.mouse.up();
    const response = await updated;
    expect(response.ok()).toBe(true);
    expect(response.request().postDataJSON()).toMatchObject({
      parent: item.targetId,
    });
    // Check the persisted hierarchy and let the drag sensor's post-drop click
    // suppression finish before exercising the disclosure controls.
    await page.reload();
    const toggle = tree.getByRole("button", { name: `Toggle ${item.target}` });
    await expect(toggle).toBeVisible();
    // Disclosure and navigation stay separate after the hierarchy changes.
    const targetRow = target.locator("..");
    const sourceRow = source.locator("..");
    if (
      (await targetRow.locator("..").getAttribute("aria-expanded")) === "true"
    ) {
      await toggle.click();
      await expect(source).toBeHidden();
    }
    await toggle.click();
    await expect(source).toBeVisible();
    await target.hover({ position: { x: 10, y: 10 } });
    await expect(targetRow.and(page.locator(":hover"))).toHaveCount(1);
    await expect(sourceRow.and(page.locator(":hover"))).toHaveCount(0);
    await source.hover({ position: { x: 10, y: 10 } });
    await expect(sourceRow.and(page.locator(":hover"))).toHaveCount(1);
    await expect(targetRow.and(page.locator(":hover"))).toHaveCount(0);
    await target.focus();
    await expect(target).toBeFocused();
  }
});
