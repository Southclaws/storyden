import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

for (const kind of ["Library", "Category"] as const) {
  test(`${kind} keyboard moves nest, sort, return to root, and persist without navigating`, async ({
    page,
  }) => {
    const seed = randomUUID().slice(0, 8);
    const { source, target } = await withAdminAccessKey(
      async ({ nodeCreate, categoryCreate }) => {
        const create = async (name: string) =>
          kind === "Library"
            ? nodeCreate({ name, visibility: "published" })
            : categoryCreate({
                name,
                colour: "#2255aa",
                description: "Keyboard move test",
              });
        return {
          source: await create(`Carried${kind}${seed}`),
          target: await create(`Destination${kind}${seed}`),
        };
      },
    );
    const admin = `keyboard-move-${seed}`;
    await createAdmin(page.context(), admin, "TestPassword123!");
    await login(page, admin, "TestPassword123!");
    await page.goto("/");
    const initialURL = page.url();
    const path =
      kind === "Library"
        ? `/api/nodes/${source.id}/position`
        : `/api/categories/${source.slug}/position`;
    const tree = page.getByRole("tree", {
      name: kind === "Library" ? "Library pages" : "Discussion categories",
    });
    const root = tree.locator("..");
    const carried = tree.locator("a").filter({ hasText: source.name });
    const destination = tree.locator("a").filter({ hasText: target.name });
    const patches: object[] = [];
    page.on("request", (request) => {
      if (
        new URL(request.url()).pathname === path &&
        request.method() === "PATCH"
      )
        patches.push(request.postDataJSON());
    });

    await carried.focus();
    await page.keyboard.press("Alt+Shift+KeyM");
    await expect(root.getByRole("status")).toContainText(
      `Carrying ${source.name}`,
    );
    await page.keyboard.type(target.name);
    await expect(destination).toBeFocused();
    const focusColor = await destination
      .locator("..")
      .evaluate((element) => getComputedStyle(element).outlineColor);
    expect(patches).toEqual([]);
    const nested = page.waitForResponse(
      (response) =>
        new URL(response.url()).pathname === path &&
        response.request().method() === "PATCH",
    );
    await page.keyboard.press("Enter");
    expect((await nested).ok()).toBe(true);
    await expect(root.getByRole("status")).toContainText(
      `Moved ${source.name} inside ${target.name}`,
    );
    await expect(carried).toBeFocused();
    expect(patches).toEqual([{ parent: target.id }]);
    expect(page.url()).toBe(initialURL);

    // The carried child navigates to its parent; placing before it promotes
    // the entire branch to the root rather than leaving the old parent set.
    await page.keyboard.press("Alt+Shift+KeyM");
    await page.keyboard.press("ArrowLeft");
    await expect(destination).toBeFocused();
    await page.keyboard.press("Space");
    await page.keyboard.press("Space");
    await expect(root.getByRole("status")).toContainText(
      `before ${target.name}, at the top level`,
    );
    await expect(destination.locator("..")).toHaveCSS("outline-style", "none");
    await expect(destination.locator("..")).toHaveCSS("box-shadow", "none");
    await expect(
      root.locator('.dragTree__dropIndicatorLine[data-active="true"]'),
    ).toHaveCount(1);
    await expect(
      root.locator('.dragTree__dropIndicatorLine[data-active="true"]'),
    ).toHaveCSS("background-color", focusColor);
    const promoted = page.waitForResponse(
      (response) =>
        new URL(response.url()).pathname === path &&
        response.request().method() === "PATCH",
    );
    await page.keyboard.press("Enter");
    expect((await promoted).ok()).toBe(true);
    await expect(carried).toBeFocused();
    expect(patches[1]).toEqual({ parent: null, before: target.id });

    await page.keyboard.press("Alt+Shift+KeyM");
    await page.keyboard.press("ArrowDown");
    await expect(destination).toBeFocused();
    await page.keyboard.press("Space");
    await expect(destination.locator("..")).toHaveCSS("outline-style", "none");
    const sorted = page.waitForResponse(
      (response) =>
        new URL(response.url()).pathname === path &&
        response.request().method() === "PATCH",
    );
    await page.keyboard.press("Enter");
    expect((await sorted).ok()).toBe(true);
    await expect(carried).toBeFocused();
    expect(patches[2]).toEqual({ parent: null, after: target.id });

    await page.reload();
    await expect(carried).toBeVisible();
    expect(await carried.locator("..").getAttribute("data-depth")).toBe("1");
    const titles = await tree.locator("a[data-tree-node-id]").allTextContents();
    expect(titles.indexOf(source.name)).toBe(titles.indexOf(target.name) + 1);
    await carried.focus();
    await page.keyboard.press("Alt+Shift+KeyM");
    await page.keyboard.press("ArrowUp");
    await page.keyboard.press("Escape");
    await expect(carried).toBeFocused();
    await expect(root.getByRole("status")).toContainText("Move cancelled");
    expect(patches).toHaveLength(3);
  });
}

test("moving into a page that hides children retains focus on the visible parent", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const { source, target } = await withAdminAccessKey(
    async ({ nodeCreate }) => ({
      source: await nodeCreate({
        name: `HiddenSource${seed}`,
        visibility: "published",
      }),
      target: await nodeCreate({
        name: `HiddenDestination${seed}`,
        visibility: "published",
        hide_child_tree: true,
      }),
    }),
  );
  const admin = `hidden-move-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  await page.goto("/");
  const tree = page.getByRole("tree", { name: "Library pages" });
  const carried = tree.locator("a").filter({ hasText: source.name });
  const destination = tree.locator("a").filter({ hasText: target.name });
  await carried.focus();
  await page.keyboard.press("Alt+Shift+KeyM");
  await page.keyboard.type(target.name);
  await expect(destination).toBeFocused();
  await page.keyboard.press("Enter");
  await expect(tree.locator("..").getByRole("status")).toContainText(
    `Moved ${source.name} inside ${target.name}`,
  );
  await expect(carried).toHaveCount(0);
  await expect(destination).toBeFocused();
  await page.keyboard.press("ArrowUp");
  await expect(destination).not.toBeFocused();
});

test("a page menu starts a move and a failed request remains available for retry", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const { source, target } = await withAdminAccessKey(
    async ({ nodeCreate }) => ({
      source: await nodeCreate({
        name: `RetrySource${seed}`,
        visibility: "published",
      }),
      target: await nodeCreate({
        name: `RetryDestination${seed}`,
        visibility: "published",
      }),
    }),
  );
  const admin = `move-retry-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  await page.goto("/");
  const tree = page.getByRole("tree", { name: "Library pages" });
  const root = tree.locator("..");
  const carried = tree.locator("a").filter({ hasText: source.name });
  const destination = tree.locator("a").filter({ hasText: target.name });
  const create = tree.getByRole("button", {
    name: `Create page under ${source.name}`,
  });
  const options = tree.getByRole("button", {
    name: `Page options for ${source.name}`,
  });
  await carried.focus();
  await expect(create).toHaveAttribute("tabindex", "0");
  await page.keyboard.press("Tab");
  await expect(create).toBeFocused();
  await page.keyboard.press("Tab");
  await expect(options).toBeFocused();
  await page.keyboard.press("Enter");
  await page.getByRole("menuitem", { name: "Move or reorder…" }).click();
  await expect(root.getByRole("status")).toContainText(
    `Carrying ${source.name}`,
  );
  await expect(carried).toBeFocused();
  await page.keyboard.type(target.name);
  await expect(destination).toBeFocused();
  await page.route(`**/api/nodes/${source.id}/position`, (route) =>
    route.fulfill({
      status: 503,
      contentType: "application/json",
      body: JSON.stringify({ message: "Temporarily unavailable" }),
    }),
  );
  await page.keyboard.press("Enter");
  await expect(root.getByRole("alert")).toContainText("Could not move");
  await expect(destination).toBeFocused();
  await expect(root.getByRole("status")).toContainText(
    `Carrying ${source.name}`,
  );
  await page.unroute(`**/api/nodes/${source.id}/position`);
  await page.keyboard.press("Enter");
  await expect(root.getByRole("status")).toContainText(`Moved ${source.name}`);
  await expect(carried).toBeFocused();
  await expect(root.getByRole("alert")).toHaveCount(0);
});
