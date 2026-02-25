/**
 * NOTE: Use this to create an admin account for a Playwright test. Do not rely
 * on first-registration assignment of the administrator permission. Tests run
 * out of order so this won't work. Instead, we create an admin account via the
 * e2e runner and that account creates an access key which is passed/used here.
 */
import { BrowserContext, Page } from "@playwright/test";

import { accountAddRole } from "../src/api/openapi-client/accounts";
import { authPasswordSignup } from "../src/api/openapi-client/auth";

const DEFAULT_ROLE_ADMIN_ID = "00000000000000000a00";

// This key has permission to assign administrator role to other accounts.
export function getAdminAccessKey(): string {
  const key = process.env["E2E_ADMIN_ACCESS_KEY"];
  if (!key) {
    throw new Error("E2E_ADMIN_ACCESS_KEY environment variable not set");
  }
  return key;
}

export async function withAccessKey<T>(
  accessKey: string,
  fn: () => Promise<T>,
): Promise<T> {
  const baseFetch = globalThis.fetch;

  const authorizedFetch: typeof fetch = async (
    ...args: Parameters<typeof fetch>
  ) => {
    const [input, init] = args;
    const request = input instanceof Request ? input : new Request(input, init);
    const headers = new Headers(request.headers);
    headers.set("Authorization", `Bearer ${accessKey}`);
    return baseFetch(new Request(request, { headers }));
  };

  globalThis.fetch = authorizedFetch;
  try {
    return await fn();
  } finally {
    globalThis.fetch = baseFetch;
  }
}

export async function withAdminAccessKey<T>(fn: () => Promise<T>): Promise<T> {
  return withAccessKey(getAdminAccessKey(), fn);
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
  await authPasswordSignup({
    identifier: username,
    token: password,
  });

  await withAdminAccessKey(async () => {
    await accountAddRole(username, DEFAULT_ROLE_ADMIN_ID);
  });
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

  await authPasswordSignup({
    identifier: username,
    token: password,
  });
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
