import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  getAdminAccessKey,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("member decorations and navigation bullets resolve Panda tokens in CSS", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const admin = `decoration-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  const created = await page.request.post("http://localhost:8001/api/threads", {
    data: {
      title: `Panda decoration ${seed}`,
      body: "A member with theme-aware decorations.",
      visibility: "published",
    },
  });
  expect(created.ok()).toBe(true);
  const thread = await created.json();
  const node = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: `Bullet token ${seed}`,
      visibility: "published",
    }),
  );
  await page.goto(`/t/${thread.slug}`);
  const name = page
    .getByRole("main")
    .locator(".member-name__show-handle-handle")
    .filter({ hasText: `@${admin}` })
    .first();
  await expect(name).toHaveCSS("font-weight", "400");
  await expect(name).toHaveCSS("font-style", "normal");
  await page.addStyleTag({ content: ":root { --font-weights-normal: 450; }" });
  await expect(name).toHaveCSS("font-weight", "450");
  await page.goto(`/m/${admin}`);
  const profileName = page
    .getByRole("main")
    .locator(".member-name__show-vertical-display-name");
  await expect(profileName).toHaveCSS("font-weight", "700");
  await expect(profileName).toHaveCSS("font-style", "normal");

  const roleResponse = await page.request.post(
    "http://localhost:8001/api/roles",
    {
      headers: { Authorization: `Bearer ${getAdminAccessKey()}` },
      data: {
        name: `Decorated ${seed}`,
        colour: "#123456",
        permissions: [],
        meta: { coloured: true, bold: true, italic: true },
      },
    },
  );
  expect(roleResponse.ok()).toBe(true);
  const role = await roleResponse.json();
  expect(role.meta).toMatchObject({ coloured: true, bold: true, italic: true });
  await withAdminAccessKey(({ accountAddRole }) =>
    accountAddRole(admin, role.id),
  );
  // A new thread avoids the existing cached author reference on the first thread.
  const decoratedResponse = await page.request.post(
    "http://localhost:8001/api/threads",
    {
      data: {
        title: `Decorated author ${seed}`,
        body: "Updated role decoration.",
        visibility: "published",
      },
    },
  );
  expect(decoratedResponse.ok()).toBe(true);
  const decoratedThread = await decoratedResponse.json();
  await page.goto(`/t/${decoratedThread.slug}`);
  await expect(name).toHaveCSS("font-weight", "600");
  await expect(name).toHaveCSS("font-style", "italic");
  await expect(name).toHaveCSS("color", "rgb(18, 52, 86)");
  await page.addStyleTag({
    content:
      ":root { --font-weights-semibold: 650; --colors-text-subtle: rgb(12, 34, 56); }",
  });
  await expect(name).toHaveCSS("font-weight", "650");
  const bullet = page
    .getByRole("tree", { name: "Library pages" })
    .getByRole("link", { name: node.name, exact: true })
    .locator("..")
    .locator('svg:has(circle[r="2.5"])');
  await expect(bullet).toHaveCSS("fill", "rgb(12, 34, 56)");
  await expect(bullet).toHaveCSS("width", "16px");
  await expect(bullet).toHaveCSS("height", "16px");
  await page.goto(`/m/${admin}`);
  await expect(profileName).toHaveCSS("font-weight", "800");
  await expect(profileName).toHaveCSS("font-style", "italic");
  await expect(profileName).toHaveCSS("color", "rgb(18, 52, 86)");
  await page.addStyleTag({
    content: ":root { --font-weights-extrabold: 850; }",
  });
  await expect(profileName).toHaveCSS("font-weight", "850");
});
