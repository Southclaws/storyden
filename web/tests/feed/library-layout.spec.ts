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
    test(`selects a page outside the menu bounds with ${input}`, async ({
      page,
    }) => {
      const seed = randomUUID().slice(0, 8);
      const name = `Library choice ${seed} with a long page name that must fit within a single option row`;
      const content = `Selected Library page content ${seed}`;
      const node = await withAdminAccessKey(
        async ({ nodeCreate, adminSettingsUpdate }) => {
          for (let index = 0; index < 10; index++) {
            await nodeCreate({
              name: `Picker option ${index} ${seed}`,
              slug: `picker-option-${index}-${seed}`,
              visibility: "published",
            });
          }
          const node = await nodeCreate({
            name,
            slug: `picker-choice-${seed}`,
            content: `<p>${content}</p>`,
            visibility: "published",
          });
          await adminSettingsUpdate({
            metadata: {
              feed: { blocks: [{ type: "library", layout: "list" }] },
            },
          });
          return node;
        },
      );
      const handle = `library_picker_${seed}`;
      await createAdmin(page.context(), handle, PASSWORD);
      await login(page, handle, PASSWORD);
      await page.setViewportSize({ width: 1440, height: 900 });
      await page.goto("/");
      await page
        .getByRole("toolbar", { name: "Primary navigation" })
        .getByRole("button", { name: "Edit site", exact: true })
        .click();
      const block = page.locator(
        '.block-editor__root[data-block-type="library"]',
      );
      await block
        .getByRole("button", { name: "Move or configure block" })
        .click();
      await page.getByRole("menuitem", { name: "Page", exact: true }).hover();
      const picker = page.getByPlaceholder("Select a Page");
      await page.getByRole("button", { name: "open", exact: true }).click();
      const popup = page.locator(".combobox__content:visible");
      await expect(popup).toBeVisible();
      for (const viewport of [
        { width: 1440, height: 900 },
        { width: 390, height: 844 },
        { width: 1440, height: 500 },
      ]) {
        await page.setViewportSize(viewport);
        await expect
          .poll(() =>
            popup.evaluate((element) => {
              const rect = element.getBoundingClientRect();
              return (
                [
                  [rect.left + 4, rect.top + 4],
                  [rect.right - 4, rect.top + 4],
                  [rect.left + 4, rect.bottom - 4],
                  [rect.right - 4, rect.bottom - 4],
                ] as const
              ).every(([x, y]) =>
                element.contains(document.elementFromPoint(x, y)),
              );
            }),
          )
          .toBe(true);
        const bounds = (await popup.boundingBox())!;
        expect(bounds.x).toBeGreaterThanOrEqual(0);
        expect(bounds.y).toBeGreaterThanOrEqual(0);
        expect(bounds.x + bounds.width).toBeLessThanOrEqual(viewport.width);
        expect(bounds.y + bounds.height).toBeLessThanOrEqual(viewport.height);
        const control = (await picker.boundingBox())!;
        expect(control.x).toBeGreaterThanOrEqual(0);
        expect(control.x + control.width).toBeLessThanOrEqual(viewport.width);
      }
      await page.setViewportSize({ width: 1440, height: 900 });
      expect((await popup.boundingBox())!.width).toBeGreaterThan(280);
      const option = page.getByRole("option", { name, exact: true });
      await option.scrollIntoViewIfNeeded();
      const label = option.locator(".combobox__itemText");
      expect(
        await label.evaluate((element) => {
          const text = element.getBoundingClientRect();
          const row = element
            .closest('[role="option"]')!
            .getBoundingClientRect();
          return {
            truncated: element.scrollWidth > element.clientWidth,
            insideRow: text.top >= row.top && text.bottom <= row.bottom,
          };
        }),
      ).toEqual({ truncated: true, insideRow: true });
      const menuBottom = await picker.evaluate(
        (element) =>
          element.closest('[role="menu"]')!.getBoundingClientRect().bottom,
      );
      await expect
        .poll(
          async () =>
            (await popup.boundingBox())!.y +
            (await popup.boundingBox())!.height,
        )
        .toBeGreaterThan(menuBottom + 20);
      const saved = page.waitForResponse(
        (response) =>
          response.request().method() === "PATCH" &&
          response.url().endsWith("/api/admin") &&
          response.ok(),
      );
      if (input === "mouse") {
        expect(
          await option.evaluate((element) => {
            const rect = element.getBoundingClientRect();
            return element.contains(
              document.elementFromPoint(
                rect.x + rect.width / 2,
                rect.y + rect.height / 2,
              ),
            );
          }),
        ).toBe(true);
        await option.click();
      } else {
        await picker.fill(name);
        await picker.press("ArrowDown");
        await picker.press("Enter");
      }
      const response = await saved;
      expect(response.request().postDataJSON().metadata.feed.blocks).toEqual([
        { type: "library", layout: "list", node: node.id },
      ]);
      await expect(
        page.getByText(content, { exact: true }).first(),
      ).toBeVisible();
      await expect(picker).toBeHidden();
      await page.reload();
      await expect(
        page.getByText(content, { exact: true }).first(),
      ).toBeVisible();
    });
  }

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
