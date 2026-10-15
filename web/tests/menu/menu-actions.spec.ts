import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAccountWithRole,
  createAdmin,
  login,
} from "../access_key_admin_assignment";

const PASSWORD = "TestPassword123!";
let member: string;

test.beforeEach(async ({ page }) => {
  const seed = randomUUID().slice(0, 8);
  member = `menu-member-${seed}`;
  const admin = `menu-admin-${seed}`;
  await createAccountWithRole(page.context(), member, PASSWORD, "member");
  await createAdmin(page.context(), admin, PASSWORD);
  await login(page, admin, PASSWORD);
  await page.goto(`/m/${member}`, { waitUntil: "networkidle" });
});

test("opens a report while its form loads, retains a draft and submits it", async ({
  page,
}) => {
  const trigger = page.locator("main .member-options-menu__trigger");
  await expect(
    page.getByRole("menuitem", { name: "Report member" }),
  ).toHaveCount(0);
  await trigger.focus();
  await trigger.press("Enter");
  const report = page.getByRole("menuitem", { name: "Report member" });
  await expect(report).toBeVisible();

  // Hold subsequently requested JS so the drawer must remain usable while loading.
  let release!: () => void;
  const ready = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("**/_next/static/**/*.js", async (route) => {
    await ready;
    await route.continue();
  });
  await report.click();
  const dialog = page.getByRole("dialog", { name: `Report ${member}` });
  try {
    await expect(dialog).toBeVisible();
    await expect(dialog.getByRole("status", { name: "Loading" })).toBeVisible();
    await expect(page.getByRole("menu")).toHaveCount(0);
  } finally {
    release();
  }
  const details = dialog.getByPlaceholder("Optional context for moderators");
  await details.fill("Please review this member's conduct.");
  await page.keyboard.press("Escape");
  await expect(dialog).toBeHidden();
  await trigger.click();
  await report.click();
  await expect(details).toHaveValue("Please review this member's conduct.");
  await dialog.getByRole("button", { name: "Cancel", exact: true }).click();
  await expect(dialog).toBeHidden();
  await trigger.click();
  await report.click();
  await expect(details).toHaveValue("");
  await details.fill("Specific report context");
  const submitted = page.waitForResponse(
    (response) =>
      response.url().endsWith("/api/reports") &&
      response.request().method() === "POST",
  );
  await dialog
    .getByRole("button", { name: "Report member", exact: true })
    .click();
  const response = await submitted;
  expect(response.ok()).toBe(true);
  expect(response.request().postDataJSON()).toMatchObject({
    target_kind: "profile",
    comment: "Specific report context",
  });
  await expect(dialog).toBeHidden();
});

test("opens moderation screens and issues a warning after reopening", async ({
  page,
}) => {
  const trigger = page.locator("main .member-options-menu__trigger");
  for (const [action, title, content] of [
    [
      "Suspend",
      `Suspend account ${member}`,
      `Do you want to suspend the account ${member}?`,
    ],
    [
      "Reset password",
      `Reset password for ${member}`,
      `Generate a password reset link for ${member}?`,
    ],
  ] as const) {
    await trigger.click();
    await page.getByRole("menuitem", { name: action, exact: true }).click();
    const dialog = page.getByRole("dialog", { name: title });
    await expect(dialog.getByText(content, { exact: true })).toBeVisible();
    await dialog.getByRole("button", { name: "Cancel", exact: true }).click();
    await expect(dialog).toBeHidden();
  }
  await trigger.click();
  await page
    .getByRole("menuitem", { name: "Warn member", exact: true })
    .click();
  const dialog = page.getByRole("dialog", {
    name: `Issue warning to ${member}`,
  });
  const reason = dialog.getByPlaceholder(
    "Clear, specific reason for this warning",
  );
  await reason.fill("Please follow the community guidelines.");
  await page.keyboard.press("Escape");
  await expect(dialog).toBeHidden();
  await trigger.click();
  await page
    .getByRole("menuitem", { name: "Warn member", exact: true })
    .click();
  await expect(reason).toHaveValue("Please follow the community guidelines.");
  const submitted = page.waitForResponse(
    (response) =>
      response.url().includes("/warnings") &&
      response.request().method() === "POST",
  );
  await dialog
    .getByRole("button", { name: "Issue warning", exact: true })
    .click();
  expect((await submitted).ok()).toBe(true);
  await expect(dialog).toBeHidden();
});

test("creates invitations only on opening, supports retry and creates a fresh link on reopen", async ({
  page,
}) => {
  let requests = 0;
  await page.route("**/api/invitations", async (route) => {
    if (route.request().method() !== "POST") return route.continue();
    requests++;
    if (requests === 1) {
      return route.fulfill({
        status: 503,
        json: { message: "Temporarily unavailable" },
      });
    }
    await route.continue();
  });
  const trigger = page.getByRole("button", {
    name: "Account menu",
    exact: true,
  });
  await trigger.click();
  expect(requests).toBe(0);
  await page.getByRole("menuitem", { name: "Invite", exact: true }).click();
  const dialog = page.getByRole("dialog", { name: "Invite someone" });
  await expect(
    dialog.getByText("Something went wrong while creating the invitation."),
  ).toBeVisible();
  await dialog.getByRole("button", { name: "Try again" }).click();
  const link = dialog.getByRole("textbox", { name: "Invitation link" });
  await expect(link).toHaveValue(/\/invitation\/.+/);
  const first = await link.inputValue();
  expect(requests).toBe(2);
  await dialog.getByRole("button", { name: "Done", exact: true }).click();
  await expect(dialog).toBeHidden();
  await trigger.click();
  await page.getByRole("menuitem", { name: "Invite", exact: true }).click();
  await expect(link).toHaveValue(/\/invitation\/.+/);
  await expect(link).not.toHaveValue(first);
  expect(requests).toBe(3);
});
