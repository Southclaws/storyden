import { expect, test } from "@playwright/test";

import { createAdmin, login } from "../access_key_admin_assignment";

const PASSWORD = "TestPassword123!";
const CSS_SOURCE =
  ":root { --e2e-theme-marker: active; --sd-color-accent: rgb(17 85 204); }";
const SCRIPT_SOURCE =
  "document.documentElement.dataset.e2eTheme='loaded';let n=0;document.addEventListener('storyden:navigate',()=>document.documentElement.dataset.e2eNavCount=String(++n));";

test("publishes CSS and JavaScript into anonymous SSR, then disables them", async ({
  browser,
  page,
}) => {
  const seed = Date.now().toString(36);
  const handle = `themeadmin-${seed}`;

  await createAdmin(page.context(), handle, PASSWORD);
  await login(page, handle, PASSWORD);
  await page.goto("/admin/theme");

  await expect(page.locator("[data-sd-theme-editor]")).toHaveCount(0);
  await page
    .getByRole("button", { name: "Enable theme editing", exact: true })
    .click();
  await expect(page.locator("[data-sd-theme-editor]")).toBeVisible();

  await page.locator('a[href="/"]').first().click();
  await expect(page).toHaveURL("/");
  await expect(page.locator("[data-sd-theme-editor]")).toBeVisible();

  await page.getByRole("button", { name: "Add CSS", exact: true }).click();
  await page.getByRole("textbox", { name: "CSS 1 source" }).fill(CSS_SOURCE);
  await page.getByRole("button", { name: "Add JavaScript" }).click();
  await page.getByRole("textbox", { name: "JS 1 source" }).fill(SCRIPT_SOURCE);

  page.once("dialog", (dialog) => dialog.accept());
  await Promise.all([
    page.waitForEvent("load"),
    page.getByRole("button", { name: "Save live theme", exact: true }).click(),
  ]);
  await expect(page.locator("[data-sd-theme-editor]")).toBeVisible();
  await expect(page.getByRole("textbox", { name: "CSS 1 source" })).toHaveValue(
    CSS_SOURCE,
  );
  await page.getByRole("tab", { name: "JS 1", exact: true }).click();
  await expect(page.getByRole("textbox", { name: "JS 1 source" })).toHaveValue(
    SCRIPT_SOURCE,
  );
  await expect(page.locator("html")).toHaveAttribute(
    "data-e2e-theme",
    "loaded",
  );

  const anonymous = await browser.newContext();
  const anonymousPage = await anonymous.newPage();
  const response = await anonymousPage.goto("/");
  expect(response).not.toBeNull();
  const html = await response!.text();
  const head = html.match(/<head[^>]*>([\s\S]*?)<\/head>/)?.[1] ?? "";
  expect(head).toMatch(/<link[^>]+data-sd-theme-asset="stylesheet"/);
  expect(head).toMatch(/<script[^>]+data-sd-theme-asset="script"/);
  await expect(
    anonymousPage.locator('link[data-sd-theme-asset="stylesheet"]'),
  ).toHaveAttribute(
    "href",
    /^http:\/\/localhost:8001\/api\/info\/theme\/assets\//,
  );
  await expect(
    anonymousPage.locator('script[data-sd-theme-asset="script"]'),
  ).toHaveAttribute(
    "src",
    /^http:\/\/localhost:8001\/api\/info\/theme\/assets\//,
  );

  await expect(anonymousPage.locator("html")).toHaveAttribute(
    "data-e2e-theme",
    "loaded",
  );
  await expect
    .poll(() =>
      anonymousPage.evaluate(() =>
        getComputedStyle(document.documentElement)
          .getPropertyValue("--e2e-theme-marker")
          .trim(),
      ),
    )
    .toBe("active");

  const initialNavigationCount = Number(
    (await anonymousPage.locator("html").getAttribute("data-e2e-nav-count")) ??
      "0",
  );
  await anonymousPage.getByRole("link", { name: "Members" }).click();
  await expect
    .poll(async () =>
      Number(
        (await anonymousPage
          .locator("html")
          .getAttribute("data-e2e-nav-count")) ?? "0",
      ),
    )
    .toBeGreaterThan(initialNavigationCount);

  await page.goto("/admin/theme");
  await page
    .locator("[data-sd-theme-editor]")
    .getByRole("button", { name: "Exit theme editing", exact: true })
    .click();
  await expect(page.locator("[data-sd-theme-editor]")).toHaveCount(0);
  await page
    .getByRole("button", { name: "Disable live theme", exact: true })
    .click();

  await anonymous.close();
  const disabledAnonymous = await browser.newContext();
  const disabledPage = await disabledAnonymous.newPage();
  const disabledResponse = await disabledPage.goto("/");
  const disabledHTML = await disabledResponse!.text();
  const disabledHead =
    disabledHTML.match(/<head[^>]*>([\s\S]*?)<\/head>/)?.[1] ?? "";
  expect(disabledHead).not.toMatch(/<(?:link|script)[^>]+data-sd-theme-asset=/);
  await expect(disabledPage.locator("head [data-sd-theme-asset]")).toHaveCount(
    0,
  );
  await expect(disabledPage.locator("html")).not.toHaveAttribute(
    "data-e2e-theme",
    "loaded",
  );
  expect(
    await disabledPage.evaluate(() =>
      getComputedStyle(document.documentElement)
        .getPropertyValue("--e2e-theme-marker")
        .trim(),
    ),
  ).toBe("");

  await disabledAnonymous.close();
});
