import { expect, test } from "@playwright/test";

import { getInfo } from "../../src/api/openapi-client/misc";
import { resolveDiscoverySettings } from "../../src/lib/settings/discovery";
import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

const PASSWORD = "TestPassword123!";

test("discovery presets, handwritten robots.txt and reset persist through the real API", async ({
  page,
  context,
}) => {
  const original = resolveDiscoverySettings((await getInfo()).discovery);
  const username = `discovery-${Date.now()}`;
  await createAdmin(context, username, PASSWORD);
  await login(page, username, PASSWORD);

  try {
    await withAdminAccessKey(async ({ adminSettingsUpdate }) => {
      await adminSettingsUpdate({
        services: {
          discovery: {
            search_indexing: true,
            training_crawling: false,
            llms_txt: false,
            auth_md: false,
            ai_catalog: false,
            api_catalog: false,
            mcp_server_card: false,
            robots_txt: "",
          },
        },
      });
    });
    await page.goto("/admin/discovery");
    const textarea = page.getByRole("textbox", {
      name: "robots.txt",
      exact: true,
    });
    const save = page.getByRole("button", { name: "Save", exact: true });
    const search = page.getByRole("checkbox", {
      name: "Enable search indexing",
    });
    const training = page.getByRole("checkbox", {
      name: "Allow AI training crawlers",
    });

    await expect(textarea).toHaveValue(/User-agent: GPTBot\nDisallow: \//);
    await expect(textarea).toHaveValue(/User-agent: CCBot\nDisallow: \//);
    await expect(
      page.getByRole("checkbox", { name: "llms.txt" }),
    ).not.toBeChecked();
    expect((await page.request.get("/llms.txt")).status()).toBe(404);
    expect((await page.request.get("/auth.md")).status()).toBe(404);

    await page.getByText("Enable search indexing", { exact: true }).click();
    await expect(search).not.toBeChecked();
    await page.getByText("Allow AI training crawlers", { exact: true }).click();
    await expect(training).toBeChecked();
    await save.click();
    await expect(textarea).toBeEnabled();
    await expect(textarea).toHaveValue("User-agent: *\nDisallow: /\n");
    expect((await getInfo()).discovery?.robots_txt).toBe("");
    expect((await page.request.get("/sitemap.xml")).status()).toBe(404);

    const custom = "User-agent: *\nDisallow: /handwritten-e2e/\n";
    await textarea.fill(custom);
    await page.getByText("Enable search indexing", { exact: true }).click();
    await expect(search).toBeChecked();
    await save.click();
    await expect(textarea).toBeEnabled();
    expect(await (await page.request.get("/robots.txt")).text()).toBe(custom);
    await page.reload();
    await expect(textarea).toHaveValue(custom);

    await page.getByText("Allow AI training crawlers", { exact: true }).click();
    await expect(training).not.toBeChecked();
    await page.getByRole("button", { name: "Use default robots.txt" }).click();
    await expect(textarea).toBeEnabled();
    await expect(textarea).toHaveValue(/User-agent: \*\nAllow: \//);
    await expect(textarea).toHaveValue(/User-agent: CCBot\nDisallow: \//);
    expect(await (await page.request.get("/robots.txt")).text()).toBe(custom);

    await page.getByText("llms.txt", { exact: true }).click();
    await expect(
      page.getByRole("checkbox", { name: "llms.txt", exact: true }),
    ).toBeChecked();
    await save.click();
    await expect(textarea).toBeEnabled();
    expect((await page.request.get("/llms.txt")).status()).toBe(200);
    expect((await page.request.get("/auth.md")).status()).toBe(404);
    const published = await (await page.request.get("/robots.txt")).text();
    await expect(textarea).toHaveValue(published);
    expect(published).not.toContain("handwritten-e2e");
    expect((await getInfo()).discovery?.robots_txt).toBe("");
    await page.reload();
    await expect(textarea).toHaveValue(published);
  } finally {
    await withAdminAccessKey(async ({ adminSettingsUpdate }) => {
      await adminSettingsUpdate({ services: { discovery: original } });
    });
    expect((await getInfo()).discovery).toEqual(original);
  }
});
