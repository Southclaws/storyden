import { expect, test } from "@playwright/test";

test.describe("Registration", () => {
  test("should register successfully with valid credentials", async ({
    page,
  }) => {
    const timestamp = Date.now();
    const username = `testuser${timestamp}`;
    const password = "TestPassword123!";

    await page.goto("/register");

    await page.getByRole("textbox", { name: "username" }).fill(username);
    await page.getByRole("textbox", { name: "password" }).fill(password);

    await page.getByRole("button", { name: "Register" }).click();

    await expect(page).toHaveURL("/", { timeout: 10000 });
    await expect(
      page.getByRole("button", { name: `Account menu` }),
    ).toBeVisible();
  });

  test("should fail to register with username that is too short", async ({
    page,
  }) => {
    await page.goto("/register");

    await page.getByRole("textbox", { name: "username" }).fill("ab");
    await page
      .getByRole("textbox", { name: "password" })
      .fill("TestPassword123!");

    await page.getByRole("button", { name: "Register" }).click();

    await expect(page).not.toHaveURL("/", { timeout: 5000 });
  });
});
