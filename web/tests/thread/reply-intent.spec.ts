import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  registerUser,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

for (const { mode, input } of [
  { mode: "richtext", input: "hover" },
  { mode: "richtext", input: "click" },
  { mode: "markdown", input: "keyboard" },
  { mode: "richtext", input: "touch" },
] as const) {
  test.describe(`Reply editor ${input}`, () => {
    if (input === "touch") {
      test.use({
        hasTouch: true,
        isMobile: true,
        viewport: { width: 390, height: 844 },
      });
    }

    test(`loads ${mode} on intent and preserves posting`, async ({ page }) => {
      const seed = randomUUID().slice(0, 8);
      const thread = await withAdminAccessKey(({ threadCreate }) =>
        threadCreate({
          title: `Reply intent ${seed}`,
          body: "A thread to read before replying.",
          visibility: "published",
        }),
      );
      await registerUser(page, `reply-${input}`, "TestPassword123!");
      const preference = await page.request.patch(
        "http://localhost:8001/api/accounts",
        {
          data: { meta: { editor: { mode } } },
        },
      );
      expect(preference.ok()).toBe(true);
      await page.mouse.move(0, 0);
      const response = await page.goto(`/t/${thread.slug}`, {
        waitUntil: "networkidle",
      });
      const initialHTML = await response!.text();
      const serverTrigger = await page.evaluate((html) => {
        const document = new DOMParser().parseFromString(html, "text/html");
        return document.querySelector(".reply-composer__trigger")?.textContent;
      }, initialHTML);
      expect(serverTrigger).toBe("Write your heart out...");
      const form = page.getByRole("main").locator(".reply-composer__form");
      await expect(form).toBeVisible();
      await expect(form.locator(".content-composer__editor")).toHaveCount(0);
      const post = form.getByRole("button", { name: "Post", exact: true });
      await expect(post).toBeDisabled();
      const formElement = await form.elementHandle();
      const postElement = await post.elementHandle();
      const trigger = form.getByRole("button", { name: "Write a reply" });
      const editor = form.getByRole("textbox", {
        name: mode === "markdown" ? "Markdown editor" : "Rich text editor",
      });
      if (input === "click") {
        let release!: () => void;
        const ready = new Promise<void>((resolve) => {
          release = resolve;
        });
        await page.route("**/_next/static/**/*.js", async (route) => {
          await ready;
          await route.continue();
        });
        try {
          await form.hover();
          const pending = form.locator(".content-composer--readonly");
          await expect(pending).toBeVisible();
          await pending.click();
        } finally {
          release();
        }
        await expect(editor).toBeFocused();
      } else if (input === "hover") {
        await form.hover();
        await expect(editor).toBeVisible();
        await expect(editor).not.toBeFocused();
        await editor.click();
      } else if (input === "keyboard") {
        await trigger.focus();
        await expect(editor).toBeFocused();
      } else {
        await trigger.tap();
        await expect(editor).toBeFocused();
      }
      expect(
        await formElement?.evaluate((element) => element.isConnected),
      ).toBe(true);
      expect(
        await postElement?.evaluate((element) => element.isConnected),
      ).toBe(true);
      const body = `Reply loaded through ${input} ${seed}`;
      await editor.fill(body);
      await expect(post).toBeEnabled();
      const posted = page.waitForResponse(
        (response) =>
          new URL(response.url()).pathname ===
            `/api/threads/${thread.slug}/replies` &&
          response.request().method() === "POST",
      );
      await post.click();
      const result = await posted;
      expect(result.ok()).toBe(true);
      const payload = result.request().postDataJSON();
      expect(payload.body.trim()).toBe(`<p>${body}</p>`);
      expect(payload.reply_to).toBeUndefined();
      await expect(page.locator("li").filter({ hasText: body })).toBeVisible();
      if (mode === "markdown") await expect(editor).toHaveValue("");
      else await expect(editor).toHaveText("");
      await expect(post).toBeDisabled();
      await expect(
        form.getByRole("button", { name: "Write a reply" }),
      ).toHaveCount(0);

      await page.context().clearCookies();
      await page.goto(`/t/${thread.slug}`, { waitUntil: "networkidle" });
      await expect(page.getByText(/Please .* to reply/)).toBeVisible();
      await expect(page.locator(".reply-composer")).toHaveCount(0);
      await expect(page.locator(".content-composer__editor")).toHaveCount(0);
    });
  });
}

for (const mode of ["richtext", "markdown"] as const) {
  test(`reply-to activates ${mode}, and clearing the target preserves the draft`, async ({
    page,
  }) => {
    const seed = randomUUID().slice(0, 8);
    const { thread, reply } = await withAdminAccessKey(
      async ({ threadCreate, replyCreate }) => {
        const thread = await threadCreate({
          title: `Reply target ${seed}`,
          body: "Thread with an existing reply.",
          visibility: "published",
        });
        const reply = await replyCreate(thread.slug, {
          body: `Existing reply ${seed}`,
        });
        return { thread, reply };
      },
    );
    await registerUser(page, `reply-target`, "TestPassword123!");
    const preference = await page.request.patch(
      "http://localhost:8001/api/accounts",
      {
        data: { meta: { editor: { mode } } },
      },
    );
    expect(preference.ok()).toBe(true);
    await page.mouse.move(0, 0);
    await page.goto(`/t/${thread.slug}`, { waitUntil: "networkidle" });
    const form = page.getByRole("main").locator(".reply-composer__form");
    await expect(form.locator(".content-composer__editor")).toHaveCount(0);
    const target = page.getByRole("main").locator(`[id="${reply.id}"]`);
    await target.getByRole("button", { name: "Reply to this" }).click();
    await expect(form.getByText("Replying to", { exact: true })).toBeVisible();
    const editor = form.getByRole("textbox", {
      name: mode === "markdown" ? "Markdown editor" : "Rich text editor",
    });
    await expect(editor).toBeFocused();
    await editor.fill("A draft to keep");
    const editorElement = await editor.elementHandle();
    await form.getByRole("button", { name: "Clear reply-to" }).click();
    await expect(form.getByText("Replying to", { exact: true })).toHaveCount(0);
    expect(
      await editorElement?.evaluate((element) => element.isConnected),
    ).toBe(true);
    if (mode === "markdown")
      await expect(editor).toHaveValue("A draft to keep");
    else await expect(editor).toHaveText("A draft to keep");
    await target.getByRole("button", { name: "Reply to this" }).click();
    const posted = page.waitForResponse(
      (response) =>
        new URL(response.url()).pathname ===
          `/api/threads/${thread.slug}/replies` &&
        response.request().method() === "POST",
    );
    await form.getByRole("button", { name: "Post", exact: true }).click();
    const result = await posted;
    expect(result.ok()).toBe(true);
    const payload = result.request().postDataJSON();
    expect(payload.body.trim()).toBe("<p>A draft to keep</p>");
    expect(payload.reply_to).toBe(reply.id);
    await expect(
      page.locator("li").filter({ hasText: "A draft to keep" }),
    ).toBeVisible();
    await expect(form.getByText("Replying to", { exact: true })).toHaveCount(0);
    await expect(editor).toBeVisible();
    await target.getByRole("button", { name: "Reply to this" }).click();
    await expect(editor).toBeFocused();
  });
}
