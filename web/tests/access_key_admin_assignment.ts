/**
 * NOTE: Use this to create an admin account for a Playwright test. Do not rely
 * on first-registration assignment of the administrator permission. Tests run
 * out of order so this won't work. Instead, we create an admin account via the
 * e2e runner and that account creates an access key which is passed/used here.
 */
import { BrowserContext, Page } from "@playwright/test";

const DEFAULT_ROLE_ADMIN_ID = "00000000000000000a00";

function getApiUrl(): string {
  return process.env["PUBLIC_API_ADDRESS"] || "http://localhost:8001";
}

// This key has permission to assign administrator role to other accounts.
function getAdminAccessKey(): string {
  const key = process.env["E2E_ADMIN_ACCESS_KEY"];
  if (!key) {
    throw new Error("E2E_ADMIN_ACCESS_KEY environment variable not set");
  }
  return key;
}

export async function registerUser(
  page: Page,
  username: string,
  password: string,
) {
  await page.goto("/register");
  await page.getByRole("textbox", { name: "username" }).fill(username);
  await page.getByRole("textbox", { name: "password" }).fill(password);
  await page.getByRole("button", { name: "Register" }).click();
  await page.waitForURL("/", { timeout: 10000 });
}

// Registers for an account via the API (no browser navigation) then uses the
// access key to assign the default administrator role to the new account.
export async function createAdmin(
  context: BrowserContext,
  username: string,
  password: string,
): Promise<void> {
  const apiUrl = getApiUrl();
  const adminAccessKey = getAdminAccessKey();

  const registerResp = await fetch(`${apiUrl}/api/auth/password/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      identifier: username,
      token: password,
    }),
  });

  if (!registerResp.ok) {
    throw new Error(
      `Failed to register user: ${registerResp.status} ${await registerResp.text()}`,
    );
  }

  const assignResp = await fetch(
    `${apiUrl}/api/accounts/${username}/roles/${DEFAULT_ROLE_ADMIN_ID}`,
    {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${adminAccessKey}`,
      },
    },
  );

  if (!assignResp.ok) {
    throw new Error(
      `Failed to assign admin role: ${assignResp.status} ${await assignResp.text()}`,
    );
  }
}

export async function createAccountWithRole(
  context: BrowserContext,
  username: string,
  password: string,
  role: "admin" | "member",
): Promise<void> {
  if (role === "admin") {
    return createAdmin(context, username, password);
  }

  const apiUrl = getApiUrl();

  const registerResp = await fetch(`${apiUrl}/api/auth/password/signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      identifier: username,
      token: password,
    }),
  });

  if (!registerResp.ok) {
    throw new Error(
      `Failed to register user: ${registerResp.status} ${await registerResp.text()}`,
    );
  }
}

export async function login(page: Page, username: string, password: string) {
  await page.goto("/login");
  await page.getByRole("textbox", { name: "username" }).fill(username);
  await page.getByRole("textbox", { name: "password" }).fill(password);
  await page.getByRole("button", { name: "login" }).click();
  await page.waitForURL("/", { timeout: 10000 });
}

export async function logout(page: Page) {
  await page.getByRole("button", { name: "Account menu" }).click();
  await page.getByRole("menuitem", { name: "Log out" }).click();
  await page.waitForURL("/login", { timeout: 10000 });
}
