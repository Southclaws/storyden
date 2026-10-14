import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import { DefaultFeedConfig } from "../../src/lib/settings/feed";
import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

const PASSWORD = "TestPassword123!";

test.describe("Homepage Library block layout", () => {
  test.afterEach(async () => {
    await withAdminAccessKey(async ({ adminSettingsUpdate }) => {
      await adminSettingsUpdate({ metadata: { feed: DefaultFeedConfig } });
    });
  });

  for (const input of ["mouse", "keyboard"] as const) {
    test(`saves Grid selected with ${input} and keeps it after reload`, async ({
      page,
    }) => {
      const seed = randomUUID().slice(0, 8);
      const name = `Library layout page ${seed}`;
      const handle = `library_layout_${seed}`;
      await withAdminAccessKey(async ({ nodeCreate, adminSettingsUpdate }) => {
        await nodeCreate({
          name,
          slug: `library-layout-${seed}`,
          content: "Library layout fixture",
          visibility: "published",
        });
        await adminSettingsUpdate({
          metadata: { feed: { blocks: [{ type: "library", layout: "list" }] } },
        });
      });
      await createAdmin(page.context(), handle, PASSWORD);
      await login(page, handle, PASSWORD);
      await page.goto("/");
      await page
        .getByRole("toolbar", { name: "Primary navigation" })
        .getByRole("button", { name: "Edit site", exact: true })
        .click();
      const block = page.locator(
        '.block-editor__root[data-block-type="library"]',
      );
      await expect(
        block.getByRole("link", { name, exact: true }).first(),
      ).toBeVisible();
      await expect(block.locator("[data-balance]")).toHaveCount(0);

      const updates: unknown[] = [];
      page.on("request", (request) => {
        if (
          request.method() === "PATCH" &&
          request.url().endsWith("/api/admin")
        ) {
          updates.push(request.postDataJSON());
        }
      });
      const trigger = block.getByRole("button", {
        name: "Move or configure block",
      });
      const layout = page.getByRole("menuitem", {
        name: "Layout",
        exact: true,
      });
      const grid = page.getByRole("menuitem", { name: "Grid", exact: true });
      if (input === "mouse") {
        await trigger.click();
        await layout.hover();
        await grid.click();
      } else {
        await trigger.focus();
        await page.keyboard.press("Enter");
        const rootMenu = page.locator(
          `[id="${await trigger.getAttribute("aria-controls")}"]`,
        );
        await expect(rootMenu).toBeFocused();
        await page.keyboard.press("Home");
        await page.keyboard.press("ArrowDown");
        await expect(rootMenu).toHaveAttribute(
          "aria-activedescendant",
          await layout.evaluate((element) => element.id),
        );
        await page.keyboard.press("ArrowRight");
        await expect(grid).toBeVisible();
        const submenu = page.locator(
          `[id="${await layout.getAttribute("aria-controls")}"]`,
        );
        await expect(submenu).toBeFocused();
        await page.keyboard.press("Home");
        await page.keyboard.press("ArrowDown");
        await expect(submenu).toHaveAttribute(
          "aria-activedescendant",
          await grid.evaluate((element) => element.id),
        );
        await page.keyboard.press("Enter");
      }

      await expect
        .poll(() => updates)
        .toEqual([
          {
            metadata: {
              feed: { blocks: [{ type: "library", layout: "grid" }] },
            },
          },
        ]);
      await expect(block.locator("[data-balance]")).toHaveCount(1);
      const saved = await page.request.get("http://localhost:8001/api/info");
      expect((await saved.json()).metadata.feed.blocks).toEqual([
        { type: "library", layout: "grid" },
      ]);
      await page.reload();
      await expect(
        page.locator('[data-sd-block="library"] [data-balance]'),
      ).toHaveCount(1);
    });
  }
});
