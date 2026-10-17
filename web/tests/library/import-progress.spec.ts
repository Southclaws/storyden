import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("URL import shows progress, honours reduced motion, and saves a child page", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const parent = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: `Import parent ${seed}`,
      visibility: "published",
      meta: {
        layout: {
          blocks: [
            { type: "title" },
            {
              type: "directory",
              config: {
                layout: "table",
                columns: [{ fid: "fixed:name", hidden: false }],
              },
            },
          ],
        },
      },
    }),
  );
  const admin = `import-progress-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  let release!: () => void;
  const ready = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("http://localhost:8001/api/links", async (route) => {
    await ready;
    await route.fulfill({
      json: { title: `Imported page ${seed}`, description: "" },
    });
  });
  await page.goto(`/l/${parent.slug}?edit=true&editMode=direct`);
  await page
    .getByRole("button", { name: "Create from URL", exact: true })
    .click();
  await page
    .getByPlaceholder("Enter URL to import...")
    .fill(`https://example.com/${seed}`);
  try {
    await page.getByRole("button", { name: "Import", exact: true }).click();
    const label = page.getByText("Fetching page...", { exact: true });
    await expect(label).toBeVisible();
    await expect(label.locator("..")).toHaveCSS(
      "animation-name",
      "slideFadeIn",
    );
    await expect(label.locator("..")).toHaveCSS("animation-duration", "0.15s");
    await page.emulateMedia({ reducedMotion: "reduce" });
    await expect(label.locator("..")).toHaveCSS("animation-name", "none");
  } finally {
    release();
  }
  await expect(page.getByText("Complete", { exact: true })).toBeVisible();
  await page.keyboard.press("Escape");
  await page.getByRole("button", { name: "View", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "Edit", exact: true }),
  ).toBeVisible();
  await expect(
    page
      .getByRole("main")
      .getByRole("link", { name: `Imported page ${seed}`, exact: true }),
  ).toBeVisible();
});
