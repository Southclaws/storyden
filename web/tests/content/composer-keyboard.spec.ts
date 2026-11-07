import { expect, test } from "@playwright/test";

import { registerUser } from "../access_key_admin_assignment";

test.use({ browserName: "webkit" });

test.describe("Keyboard submission in WebKit", () => {
  for (const mode of ["richtext", "markdown"] as const) {
    test(`submits a ${mode} thread with Cmd+Enter and reply with Ctrl+Enter`, async ({
      page,
    }) => {
      await registerUser(page, `keyboard-${mode}`, "TestPassword123!");

      if (mode === "markdown") {
        const preference = await page.request.patch(
          "http://localhost:8001/api/accounts",
          { data: { meta: { editor: { mode } } } },
        );
        expect(preference.ok()).toBe(true);
      }

      await page.goto("/new");
      const editor = page.getByRole("textbox", {
        name: mode === "markdown" ? "Markdown editor" : "Rich text editor",
      });
      await editor.fill("First paragraph");
      await editor.press("Enter");
      await page.keyboard.insertText("Second paragraph");
      await editor.press("Shift+Enter");
      await page.keyboard.insertText("After a line break");

      if (mode === "richtext") {
        await expect(editor.locator("p")).toHaveCount(2);
        await expect(editor.locator("p").nth(1).locator("br")).toHaveCount(1);
      } else {
        await expect(editor).toHaveValue(
          "First paragraph\nSecond paragraph\nAfter a line break",
        );
      }

      await expect(page).toHaveURL(/\/new$/);
      const title = page.locator("#title-input");
      await title.fill(`Keyboard submission ${mode}`);
      await expect(title).toHaveText(`Keyboard submission ${mode}`);
      await editor.fill("A thread posted from the keyboard.");
      await expect(
        page.getByRole("button", { name: "Post", exact: true }),
      ).toBeEnabled();
      await editor.press("Meta+Enter");

      await expect(page).toHaveURL(/\/t\//);
      await expect(
        page.getByRole("heading", { name: `Keyboard submission ${mode}` }),
      ).toBeVisible();

      const form = page.getByRole("main").locator(".reply-composer__form");
      await form.getByRole("button", { name: "Write a reply" }).focus();
      const replyEditor = form.getByRole("textbox");
      await replyEditor.focus();
      await expect(replyEditor).toBeFocused();
      await replyEditor.fill("A reply posted from the keyboard.");
      const post = form.getByRole("button", { name: "Post", exact: true });
      await expect(post).toBeEnabled();

      const posted = page.waitForResponse(
        (response) =>
          /\/api\/threads\/[^/]+\/replies$/.test(
            new URL(response.url()).pathname,
          ) && response.request().method() === "POST",
      );
      await replyEditor.press("Control+Enter");
      const response = await posted;
      expect(response.ok()).toBe(true);
      expect(response.request().postDataJSON().body.trim()).toBe(
        "<p>A reply posted from the keyboard.</p>",
      );
      await expect(
        page
          .locator("li")
          .filter({ hasText: "A reply posted from the keyboard." }),
      ).toBeVisible();
      await expect(post).toBeDisabled();
    });
  }
});
