import { expect, test } from "@playwright/test";

test.describe("Login", () => {
  test("should login successfully with correct credentials", async ({
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

    await page.getByRole("button", { name: "Account menu" }).click();
    await page.getByRole("link", { name: "Logout" }).click();

    await page.getByRole("link", { name: "Login" }).click();
    await page.getByRole("textbox", { name: "username" }).fill(username);
    await page.getByRole("textbox", { name: "password" }).fill(password);
    await page.getByRole("button", { name: "Login" }).click();

    await expect(page).toHaveURL("/", { timeout: 10000 });
    await expect(
      page.getByRole("button", { name: `Account menu` }),
    ).toBeVisible();
  });

  test("should fail to login with wrong password", async ({ page }) => {
    const timestamp = Date.now();
    const username = `testuser${timestamp}`;
    const password = "TestPassword123!";
    const wrongPassword = "WrongPassword456!";

    await page.goto("/register");
    await page.getByRole("textbox", { name: "username" }).fill(username);
    await page.getByRole("textbox", { name: "password" }).fill(password);
    await page.getByRole("button", { name: "Register" }).click();

    await expect(page).toHaveURL("/", { timeout: 10000 });

    await page.getByRole("button", { name: "Account menu" }).click();
    await page.getByRole("link", { name: "Logout" }).click();

    await page.getByRole("link", { name: "Login" }).click();
    await page.getByRole("textbox", { name: "username" }).fill(username);
    await page.getByRole("textbox", { name: "password" }).fill(wrongPassword);
    await page.getByRole("button", { name: "Login" }).click();

    await expect(page).not.toHaveURL("/", { timeout: 5000 });
  });
});
