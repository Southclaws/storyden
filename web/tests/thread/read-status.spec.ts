import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import { API_ADDRESS } from "../../src/config";
import {
  registerUser,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

for (const viewport of [
  { width: 1280, height: 900 },
  { width: 390, height: 844 },
]) {
  test(`explicit thread read controls at ${viewport.width}px`, async ({
    page,
  }, testInfo) => {
    await page.setViewportSize(viewport);
    await registerUser(page, "read-status", "TestPassword123!");
    const thread = await withAdminAccessKey(({ threadCreate }) =>
      threadCreate({
        title: `Read controls ${randomUUID().slice(0, 8)}`,
        body: "A discussion to revisit later.",
        visibility: "published",
      }),
    );

    await page.goto("/d");
    const card = page.locator(`[id="${thread.id}"]:visible`).first();
    const unreadHeading = card.getByRole("heading", {
      name: `${thread.title} Unread`,
      exact: true,
    });
    await expect(unreadHeading).toBeVisible();
    const trigger = card.getByRole("button", {
      name: `More options for ${thread.title}`,
      exact: true,
    });

    await trigger.focus();
    await trigger.press("Enter");
    await expect(
      page.getByRole("menuitem", { name: "Mark as unread", exact: true }),
    ).toHaveCount(0);
    await page
      .getByRole("menuitem", { name: "Mark as read", exact: true })
      .click();
    await expect(
      card.getByRole("heading", {
        name: `${thread.title} Caught up`,
        exact: true,
      }),
    ).toBeVisible();

    await trigger.click();
    await expect(
      page.getByRole("menuitem", { name: "Mark as read", exact: true }),
    ).toHaveCount(0);
    await page
      .getByRole("menuitem", { name: "Mark as unread", exact: true })
      .click();
    await expect(unreadHeading).toBeVisible();

    await trigger.click();
    await expect(
      page.getByRole("menuitem", { name: "Mark as read", exact: true }),
    ).toBeVisible();
    await expect(
      page.getByRole("menuitem", { name: "Mark as unread", exact: true }),
    ).toHaveCount(0);
    await page.getByRole("menu").screenshot({
      path: testInfo.outputPath("read-status-menu.png"),
      animations: "disabled",
    });
    await page.keyboard.press("Escape");

    await card.getByRole("link", { name: thread.title, exact: true }).click();
    await expect(page).toHaveURL(new RegExp(`/t/${thread.slug}`));

    const status = async () => {
      const response = await page.request.get(
        `${API_ADDRESS}/api/threads/${thread.slug}`,
      );
      expect(response.ok()).toBe(true);
      const data = await response.json();
      return data.read_status;
    };

    await expect
      .poll(async () => (await status())?.marked_unread_at)
      .toBeUndefined();

    const threadTrigger = page
      .getByRole("button", {
        name: `More options for ${thread.title}`,
        exact: true,
      })
      .first();
    await threadTrigger.click();
    await expect(
      page.getByRole("menuitem", { name: "Mark as read", exact: true }),
    ).toHaveCount(0);
    const marked = page.waitForResponse(
      (response) =>
        response.url().endsWith(`/threads/${thread.slug}/read-status`) &&
        response.request().method() === "PATCH",
    );
    await page
      .getByRole("menuitem", { name: "Mark as unread", exact: true })
      .click();
    expect((await marked).ok()).toBe(true);
    await expect
      .poll(async () => (await status())?.marked_unread_at)
      .toBeTruthy();

    await threadTrigger.click();
    await expect(
      page.getByRole("menuitem", { name: "Mark as read", exact: true }),
    ).toBeVisible();
    await expect(
      page.getByRole("menuitem", { name: "Mark as unread", exact: true }),
    ).toHaveCount(0);
    await page.keyboard.press("Escape");

    const reply = page.locator(".reply-composer__form");
    await reply.hover();
    const editor = reply.locator(".content-composer__editor");
    await expect(editor).toHaveAttribute("contenteditable", "true");
    await editor.click();
    await editor.fill(`Following up at ${viewport.width}px`);
    await reply.getByRole("button", { name: "Post", exact: true }).click();
    await expect(editor).toHaveText("");
    await threadTrigger.click();
    await expect(
      page.getByRole("menuitem", { name: "Mark as unread", exact: true }),
    ).toBeVisible();
    await expect(
      page.getByRole("menuitem", { name: "Mark as read", exact: true }),
    ).toHaveCount(0);
    await expect
      .poll(async () => (await status())?.marked_unread_at)
      .toBeUndefined();
  });
}
