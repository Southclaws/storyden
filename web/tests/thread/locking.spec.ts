import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import { API_ADDRESS } from "../../src/config";
import {
  getAdminAccessKey,
  registerUser,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

for (const viewport of [
  { width: 1280, height: 900 },
  { width: 390, height: 844 },
]) {
  test(`thread lock controls and reply permissions at ${viewport.width}px`, async ({
    page,
    request,
  }, testInfo) => {
    await page.setViewportSize(viewport);
    const username = await registerUser(
      page,
      "thread-lock",
      "TestPassword123!",
    );
    const created = await page.request.post(`${API_ADDRESS}/api/threads`, {
      data: {
        title: `Thread locking ${randomUUID().slice(0, 8)}`,
        body: "A discussion that can be resolved and locked independently.",
        visibility: "published",
      },
    });
    expect(created.ok()).toBe(true);
    const thread = await created.json();
    await withAdminAccessKey(({ replyCreate }) =>
      replyCreate(thread.slug, { body: "A reply before locking." }),
    );
    const locked = await request.patch(
      `${API_ADDRESS}/api/threads/${thread.slug}`,
      {
        headers: { Authorization: `Bearer ${getAdminAccessKey()}` },
        data: { locked: true, lifecycle_state: "resolved" },
      },
    );
    expect(locked.ok()).toBe(true);

    await page.goto(`/t/${thread.slug}`);
    await expect(
      page.getByRole("heading", { name: "Thread locked", exact: true }),
    ).toBeVisible();
    const disabledComposer = page.locator(".reply-composer__form:visible");
    await expect(disabledComposer).toHaveAttribute("aria-disabled", "true");
    await expect(
      disabledComposer.getByRole("button", { name: "Post reply", exact: true }),
    ).toBeDisabled();
    await expect(
      disabledComposer.getByRole("button", {
        name: "Write a reply",
        exact: true,
      }),
    ).toBeDisabled();
    const alert = await page
      .getByRole("heading", { name: "Thread locked", exact: true })
      .boundingBox();
    const title = await page.getByRole("heading", { level: 1 }).boundingBox();
    expect(alert).not.toBeNull();
    expect(title).not.toBeNull();
    expect(alert!.y).toBeLessThan(title!.y);
    await expect(
      page.getByRole("button", { name: "Reply to this", exact: true }),
    ).toHaveCount(0);
    const trigger = page.getByRole("button", {
      name: `More options for ${thread.title}`,
      exact: true,
    });
    await trigger.click();
    await expect(
      page.getByRole("menuitem", { name: "Lock thread", exact: true }),
    ).toHaveCount(0);
    await expect(
      page.getByRole("menuitem", { name: "Unlock thread", exact: true }),
    ).toHaveCount(0);
    await page.keyboard.press("Escape");
    await page.screenshot({
      path: testInfo.outputPath("locked-author.png"),
      fullPage: true,
    });

    await withAdminAccessKey(({ accountAddRole }) =>
      accountAddRole(username, "00000000000000000a00"),
    );
    await page.reload();
    await expect(page.locator(".reply-composer__form:visible")).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Reply to this", exact: true }),
    ).toBeVisible();
    await trigger.click();
    await page
      .getByRole("menuitem", { name: "Unlock thread", exact: true })
      .click();
    await expect(
      page.getByRole("heading", { name: "Thread locked", exact: true }),
    ).toHaveCount(0);
    await expect(page.getByRole("heading", { level: 1 })).toHaveAccessibleName(
      `${thread.title} Resolved`,
    );

    await trigger.click();
    await page
      .getByRole("menuitem", { name: "Lock thread", exact: true })
      .click();
    await expect(
      page.getByRole("heading", { name: "Thread locked", exact: true }),
    ).toBeVisible();
    const composer = page.locator(".reply-composer__form:visible");
    await composer.hover();
    const editor = composer.locator(".content-composer__editor");
    await expect(editor).toHaveAttribute("contenteditable", "true");
    await editor.fill("An administrator can still reply while locked.");
    await composer
      .getByRole("button", { name: "Post reply", exact: true })
      .click();
    await expect(
      page
        .locator(".thread-reply:visible")
        .getByText("An administrator can still reply while locked.", {
          exact: true,
        }),
    ).toBeVisible();
    await expect(editor).toHaveText("");
    await page.screenshot({
      path: testInfo.outputPath("locked-admin.png"),
      fullPage: true,
    });
  });
}
