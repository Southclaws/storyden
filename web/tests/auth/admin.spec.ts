import { expect, test } from "@playwright/test";

import { createAccountWithRole, login } from "../access_key_admin_assignment";

test.describe("utility", () => {
  test("assign_admin_role", async ({ browser }) => {
    const timestamp = Date.now();
    const adminHandle = `admin_${timestamp}`;
    const memberHandle = `member_${timestamp}`;
    const password = "TestPassword123!";

    const adminctx = await browser.newContext();
    const memberctx = await browser.newContext();

    await createAccountWithRole(adminctx, adminHandle, password, "admin");
    await createAccountWithRole(memberctx, memberHandle, password, "member");

    const adminPage = await adminctx.newPage();
    await login(adminPage, adminHandle, password);
    await adminPage.goto("/settings");
    await expect(
      adminPage.getByRole("heading", { name: "Settings" }),
    ).toBeVisible();

    const memberPage = await memberctx.newPage();
    await login(memberPage, memberHandle, password);
    await expect(
      memberPage.getByRole("button", { name: "Account menu" }),
    ).toBeVisible();

    await adminctx.close();
    await memberctx.close();
  });
});
