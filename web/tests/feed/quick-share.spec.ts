import { expect, test } from "@playwright/test";

import { DefaultFeedConfig } from "../../src/lib/settings/feed";
import {
  registerUser,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

for (const { mode, input } of [
  { mode: "richtext", input: "hover" },
  { mode: "markdown", input: "keyboard" },
  { mode: "richtext", input: "touch" },
] as const) {
  test.describe(`Quick Share ${input}`, () => {
    if (input === "touch") {
      test.use({
        hasTouch: true,
        isMobile: true,
        viewport: { width: 390, height: 844 },
      });
    }

    test.afterEach(async () => {
      await withAdminAccessKey(async ({ adminSettingsUpdate }) => {
        await adminSettingsUpdate({ metadata: { feed: DefaultFeedConfig } });
      });
    });

    test(`loads the ${mode} editor on intent and preserves posting`, async ({
      page,
    }) => {
      await withAdminAccessKey(async ({ adminSettingsUpdate }) => {
        await adminSettingsUpdate({
          metadata: {
            feed: {
              blocks: [
                { type: "quick-share", showCategorySelect: true },
                { type: "threads", source: "all" },
              ],
            },
          },
        });
      });
      const handle = await registerUser(
        page,
        `quick-${input}`,
        "TestPassword123!",
      );
      const preference = await page.request.patch(
        "http://localhost:8001/api/accounts",
        {
          data: { meta: { editor: { mode } } },
        },
      );
      expect(preference.ok()).toBe(true);
      await page.mouse.move(0, 0);
      const response = await page.goto("/", { waitUntil: "networkidle" });
      const serverCategory = await page.evaluate(
        (html) => {
          const document = new DOMParser().parseFromString(html, "text/html");
          return document.querySelector('.quick-share [role="combobox"]')
            ?.textContent;
        },
        await response!.text(),
      );
      expect(serverCategory).toBe("Category");
      const card = page.locator(".quick-share");
      const editor = card.getByRole("textbox", {
        name: mode === "markdown" ? "Markdown editor" : "Rich text editor",
      });
      await expect(card).toBeVisible();
      await expect(card.locator(".content-composer__editor")).toHaveCount(0);
      const prompt = card.getByRole("button", {
        name: "Quickly create a new thread",
      });
      await expect(prompt).toBeVisible();
      const form = await card.locator(".quick-share__form").elementHandle();
      const actions = await card
        .locator(".quick-share__actions")
        .elementHandle();
      const share = card.getByRole("button", { name: "Share", exact: true });
      const category = card.getByRole("combobox");
      await expect(share).toBeVisible();
      await expect(category).toBeVisible();

      if (input === "keyboard") {
        await category.focus();
        await category.press("ArrowDown");
        const noCategory = page.getByRole("option", {
          name: "No category",
          exact: true,
        });
        await expect(noCategory).toBeVisible();
        await page.keyboard.press("Home");
        await expect(noCategory).toHaveAttribute("data-highlighted", "");
        await page.keyboard.press("Enter");
        await expect(category).toContainText("No category");
        await expect(card.locator(".content-composer__editor")).toHaveCount(0);
      }

      if (input === "hover") {
        await card.hover();
        await expect(editor).toBeVisible();
        await expect(editor).not.toBeFocused();
        await editor.click();
      } else if (input === "keyboard") {
        await prompt.focus();
        await expect(editor).toBeFocused();
      } else {
        await prompt.tap();
        await expect(editor).toBeFocused();
      }

      expect(await form?.evaluate((element) => element.isConnected)).toBe(true);
      expect(await actions?.evaluate((element) => element.isConnected)).toBe(
        true,
      );
      if (input === "keyboard") {
        await expect(category).toContainText("No category");
      }

      const title = `Quick Share ${handle}`;
      await editor.fill(title);
      await page.mouse.move(0, 0);
      if (mode === "markdown") {
        await expect(editor).toHaveValue(title);
      } else {
        await expect(editor).toHaveText(title);
      }

      const posted = page.waitForResponse(
        (response) =>
          new URL(response.url()).pathname === "/api/threads" &&
          response.request().method() === "POST",
      );
      await card.getByRole("button", { name: "Share", exact: true }).click();
      expect((await posted).ok()).toBe(true);
      await expect(
        page.getByRole("link", { name: title, exact: true }).first(),
      ).toBeVisible();
      if (mode === "markdown") {
        await expect(editor).toHaveValue("");
      } else {
        await expect(editor).toHaveText("");
      }
    });
  });
}
