import { type Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

async function openImport(page: Page, seed: string) {
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
  await page.goto(`/l/${parent.slug}?edit=true&editMode=direct`);
  await page
    .getByRole("button", { name: "Create from URL", exact: true })
    .click();
  await page
    .getByPlaceholder("Enter URL to import...")
    .fill(`https://example.com/${seed}`);
}

test("URL import shows progress, honours reduced motion, and saves a child page", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
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
  await openImport(page, seed);
  try {
    await page.getByRole("button", { name: "Import", exact: true }).click();
    const label = page.getByText("Fetching page...", { exact: true });
    await expect(label).toBeVisible();
    const spinner = label.locator("../..").locator(":scope > div").first();
    const tokenOverride = await page.addStyleTag({
      content:
        ":root { --colors-background-control-track: rgb(12, 34, 56); --colors-text-muted: rgb(65, 43, 21); }",
    });
    await expect(spinner).toHaveCSS("width", "20px");
    await expect(spinner).toHaveCSS("height", "20px");
    await expect(spinner).toHaveCSS("border-top-width", "3px");
    await expect(spinner).toHaveCSS("border-right-color", "rgb(12, 34, 56)");
    await expect(spinner).toHaveCSS("border-top-color", "rgb(65, 43, 21)");
    await expect(spinner).toHaveCSS("animation-name", "spin");
    await expect(spinner).toHaveCSS("animation-duration", "1s");
    await tokenOverride.evaluate((element) =>
      element.parentNode?.removeChild(element),
    );
    await test.info().attach("Import progress spinner", {
      body: await label
        .locator("../../..")
        .screenshot({ animations: "disabled" }),
      contentType: "image/png",
    });
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
  const complete = page.getByText("Complete", { exact: true });
  await expect(complete).toBeVisible();
  await expect(complete.locator("../..").locator("svg")).toHaveCSS(
    "fill",
    "none",
  );
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

test("URL import failures use the same unfilled danger icon as error tooltips", async ({
  page,
}) => {
  await page.route("http://localhost:8001/api/links", (route) =>
    route.fulfill({
      status: 500,
      json: { message: "Import service unavailable" },
    }),
  );
  await openImport(page, randomUUID().slice(0, 8));
  const tokenOverride = await page.addStyleTag({
    content:
      ":root { --colors-status-danger-content: rgb(128, 20, 30); --colors-text-subtle: rgb(77, 92, 103); }",
  });
  await page.getByRole("button", { name: "Import", exact: true }).click();
  const label = page.getByText("Failed to import", { exact: true });
  await expect(label).toBeVisible();
  const icon = label.locator("../..").locator("svg");
  await expect(icon).toHaveCSS("width", "20px");
  await expect(icon).toHaveCSS("height", "20px");
  await expect(icon).toHaveCSS("color", "rgb(128, 20, 30)");
  await expect(icon).toHaveCSS("fill", "none");
  await expect(label).toHaveCSS("color", "rgb(77, 92, 103)");
  await tokenOverride.evaluate((element) =>
    element.parentNode?.removeChild(element),
  );
  await test.info().attach("Import failure indicator", {
    body: await label
      .locator("../../..")
      .screenshot({ animations: "disabled" }),
    contentType: "image/png",
  });
});
