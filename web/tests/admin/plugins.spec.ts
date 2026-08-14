import { Page, expect, test } from "@playwright/test";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

const PASSWORD = "TestPassword123!";

function uniqueSeed() {
  return `${Date.now().toString(36)}${Math.random().toString(36).slice(2, 6)}`;
}

async function addExternalPlugin(seed: string) {
  const name = `E2E Plugin ${seed}`;

  const plugin = await withAdminAccessKey(({ pluginAdd }) =>
    pluginAdd({
      mode: "external",
      manifest: {
        id: `e2e-plugin-${seed}`,
        name,
        author: "e2e",
        description: `External plugin for e2e run ${seed}`,
        version: "1.0.0",
        command: "./plugin",
      },
    }),
  );

  return { id: plugin.id, name };
}

async function loginAsAdmin(page: Page, seed: string) {
  const handle = `pluginadmin${seed}`;
  await createAdmin(page.context(), handle, PASSWORD);
  await login(page, handle, PASSWORD);
}

test.describe("admin plugins", () => {
  test("opens a plugin from the list and navigates back", async ({ page }) => {
    const seed = uniqueSeed();
    const plugin = await addExternalPlugin(seed);

    await loginAsAdmin(page, seed);

    await page.goto("/admin/plugins");

    const card = page.getByRole("link", { name: plugin.name, exact: true });
    await expect(card).toBeVisible();
    await expect(card).toHaveAttribute(
      "href",
      new RegExp(`plugin=${plugin.id}`),
    );

    await card.click();

    await expect(page).toHaveURL(new RegExp(`\\?plugin=${plugin.id}`));
    await expect(
      page.getByRole("heading", { name: plugin.name }),
    ).toBeVisible();
    await expect(page.getByRole("tab", { name: "Overview" })).toBeVisible();
    await expect(
      page.getByRole("tab", { name: "Configuration" }),
    ).toBeVisible();

    await page.getByRole("button", { name: "Back" }).click();

    await expect(page).toHaveURL(/\/admin\/plugins$/);
    await expect(card).toBeVisible();

    await withAdminAccessKey(({ pluginDelete }) => pluginDelete(plugin.id));
  });

  test("renders the plugin page when linked to directly", async ({ page }) => {
    const seed = uniqueSeed();
    const plugin = await addExternalPlugin(seed);

    await loginAsAdmin(page, seed);

    await page.goto(`/admin/plugins?plugin=${plugin.id}`);

    await expect(
      page.getByRole("heading", { name: plugin.name }),
    ).toBeVisible();
    await expect(page.getByRole("tab", { name: "Overview" })).toBeVisible();

    await withAdminAccessKey(({ pluginDelete }) => pluginDelete(plugin.id));
  });

  test("deletes a plugin without opening it", async ({ page }) => {
    const seed = uniqueSeed();
    const plugin = await addExternalPlugin(seed);

    await loginAsAdmin(page, seed);

    await page.goto("/admin/plugins");

    const card = page
      .getByRole("listitem")
      .filter({ has: page.getByRole("link", { name: plugin.name }) });
    await expect(card).toBeVisible();

    await card.getByRole("button", { name: "Delete" }).click();
    await card.getByRole("button", { name: "Confirm Delete" }).click();

    await expect(page).toHaveURL(/\/admin\/plugins$/);
    await expect(card).toHaveCount(0);
  });
});
