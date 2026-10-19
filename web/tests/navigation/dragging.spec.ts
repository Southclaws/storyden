import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("library sidebar separates disclosure, destination, and page actions for keyboard users", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const { parent, child, sibling } = await withAdminAccessKey(
    async ({ nodeCreate }) => {
      const parent = await nodeCreate({
        name: `About${seed} handbook`,
        visibility: "published",
      });
      const child = await nodeCreate({
        name: `Getting started ${seed}`,
        parent: parent.slug,
        visibility: "published",
      });
      const sibling = await nodeCreate({
        name: `Reference ${seed}`,
        parent: parent.slug,
        visibility: "published",
      });
      return { parent, child, sibling };
    },
  );
  const admin = `tree-keyboard-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  await page.goto(`/l/${child.slug}`);

  const navigation = page.getByRole("tree", { name: "Library pages" });
  const toggle = navigation.getByRole("button", {
    name: `Toggle ${parent.name}`,
  });
  const destination = navigation.getByRole("link", {
    name: parent.name,
    exact: true,
  });
  const childLink = navigation.getByRole("treeitem", {
    name: child.name,
    exact: true,
  });
  await expect(childLink).toHaveAttribute("aria-current", "page");
  await expect(toggle).toHaveAttribute("aria-expanded", "true");
  const siblingLink = navigation.getByRole("treeitem", {
    name: sibling.name,
    exact: true,
  });
  const currentRow = childLink.locator("..");
  const parentRow = destination.locator("..");
  const currentBackground = await currentRow.evaluate(
    (element) => getComputedStyle(element).backgroundColor,
  );
  await page.evaluate(() =>
    document.documentElement.style.setProperty(
      "--sd-color-focus-ring",
      "rgb(171, 137, 255)",
    ),
  );
  await childLink.focus();
  await page.keyboard.press("ArrowLeft");
  await expect(destination).toBeFocused();
  await expect(parentRow).toHaveCSS("outline-style", "solid");
  await expect(parentRow).toHaveCSS("outline-width", "2px");
  await expect(parentRow).toHaveCSS("outline-color", "rgb(171, 137, 255)");
  await expect(destination).toHaveCSS("outline-style", "none");
  await expect(currentRow).toHaveCSS("background-color", currentBackground);
  await page.keyboard.press("ArrowLeft");
  await expect(toggle).toHaveAttribute("aria-expanded", "false");
  await expect(childLink).toBeHidden();
  await expect(destination).toBeFocused();
  await page.keyboard.press("ArrowRight");
  await expect(toggle).toHaveAttribute("aria-expanded", "true");
  await expect(destination).toBeFocused();
  await page.keyboard.press("ArrowRight");
  await expect(childLink).toBeFocused();
  await page.keyboard.press("ArrowDown");
  await expect(siblingLink).toBeFocused();
  await page.keyboard.press("ArrowUp");
  await expect(childLink).toBeFocused();
  await page.keyboard.press("Home");
  await expect(
    navigation.locator("a[data-tree-node-id]:visible").first(),
  ).toBeFocused();
  await page.keyboard.press("End");
  await expect(
    navigation.locator("a[data-tree-node-id]:visible").last(),
  ).toBeFocused();
  await page.keyboard.type(`About${seed}`);
  await expect(destination).toBeFocused();
  await page.keyboard.press("ArrowRight");
  await expect(childLink).toBeFocused();
  await expect(childLink).toHaveAttribute("aria-current", "page");
  const moveHint = page
    .getByRole("tooltip")
    .filter({ hasText: "Move or reorder:" });
  await expect(moveHint).toBeVisible();
  await page.keyboard.press("Shift+Tab");
  await expect
    .poll(() =>
      navigation.evaluate((element) =>
        element.contains(document.activeElement),
      ),
    )
    .toBe(false);
  await expect(moveHint).toBeHidden();
  await page.keyboard.press("Tab");
  await expect(childLink).toBeFocused();
  await expect(moveHint).toBeVisible();
  await page.keyboard.press("Tab");
  const create = navigation.getByRole("button", {
    name: `Create page under ${child.name}`,
  });
  await expect(create).toBeFocused();
  await expect(create).toBeVisible();
  await expect(moveHint).toBeHidden();
  await expect(create).toHaveCSS("outline-style", "solid");
  await expect(create).toHaveCSS("outline-width", "2px");
  await expect(currentRow).toHaveCSS("outline-style", "none");
  await expect(parentRow).toHaveCSS("outline-style", "none");
  await expect(currentRow).toHaveCSS("background-color", currentBackground);
  await page.keyboard.press("Tab");
  const options = navigation.getByRole("button", {
    name: `Page options for ${child.name}`,
  });
  await expect(options).toBeFocused();
  await page.keyboard.press("Enter");
  await expect(page.getByRole("menu")).toBeVisible();
  await page.keyboard.press("ArrowDown");
  await expect(page.getByRole("menu")).toBeFocused();
  await expect(page.getByRole("menu")).toHaveAttribute(
    "aria-activedescendant",
    /.+/,
  );
  await page.keyboard.press("Escape");
  await expect(page.getByRole("menu")).toBeHidden();
  await expect(options).toBeFocused();
  await page.keyboard.press("Escape");
  await expect(childLink).toBeFocused();
  await destination.focus();
  await page.keyboard.press("Enter");
  await page.waitForURL(`/l/${parent.slug}`);
  await expect(destination).toHaveAttribute("aria-current", "page");
  await expect(childLink).toBeVisible();
});

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
    const source = tree.locator("a").filter({ hasText: item.source });
    const target = tree.locator("a").filter({ hasText: item.target });
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
    if ((await toggle.getAttribute("aria-expanded")) === "true") {
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
