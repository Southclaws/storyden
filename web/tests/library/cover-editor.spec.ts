import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  getAdminAccessKey,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("Library keeps its cover visible while the cropper loads, and saves cropping", async ({
  page,
  browser,
}) => {
  const seed = randomUUID().slice(0, 8);
  // A non-square image makes moving the fixed 4:1 cover crop meaningful.
  const image = await page.evaluate(() => {
    const canvas = document.createElement("canvas");
    canvas.width = 1800;
    canvas.height = 900;
    const ctx = canvas.getContext("2d")!;
    ctx.fillStyle = "#2255aa";
    ctx.fillRect(0, 0, 1800, 450);
    ctx.fillStyle = "#cc5533";
    ctx.fillRect(0, 450, 1800, 450);
    return canvas.toDataURL("image/png").split(",")[1]!;
  });
  const uploaded = await page.request.post(
    "http://localhost:8001/api/assets?filename=library-cover.png",
    {
      headers: {
        Authorization: `Bearer ${getAdminAccessKey()}`,
        "Content-Type": "application/octet-stream",
      },
      data: Buffer.from(image, "base64"),
    },
  );
  expect(uploaded.ok()).toBe(true);
  const asset = await uploaded.json();
  const node = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: `Readable Library ${seed}`,
      slug: `read-edit-${seed}`,
      visibility: "published",
      content: "Content stays readable before editing.",
      primary_image_asset_id: asset.id,
      meta: {
        layout: {
          blocks: [{ type: "title" }, { type: "cover" }, { type: "content" }],
        },
      },
    }),
  );
  const admin = `read-edit-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  await page.mouse.move(0, 0);
  await page.goto(`/l/${node.slug}`, { waitUntil: "networkidle" });
  const main = page.getByRole("main");
  await expect(main.locator('[data-outline-block="cover"] img')).toBeVisible();
  await expect(main.locator(".advanced-cropper")).toHaveCount(0);

  await page.getByRole("button", { name: "Edit", exact: true }).click();
  const quickEdit = page.getByRole("menuitem", {
    name: "Quick edit",
    exact: true,
  });
  await expect(quickEdit).toBeEnabled();
  let release!: () => void;
  const ready = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("**/_next/static/**/*.js", async (route) => {
    await ready;
    await route.continue();
  });
  try {
    await quickEdit.click();
    await expect(
      page.locator("#block-cover_container:visible img"),
    ).toBeVisible();
  } finally {
    release();
  }
  await page.unrouteAll({ behavior: "wait" });
  const cropper = page.locator(
    "#block-cover_container:visible .advanced-cropper",
  );
  await expect(cropper).toBeVisible();
  await expect(cropper.locator("img").first()).toBeVisible();
  const box = await cropper.boundingBox();
  expect(box).not.toBeNull();
  const cropped = page.waitForResponse(
    (response) =>
      new URL(response.url()).pathname === "/api/assets" &&
      new URL(response.url()).searchParams.get("filename") ===
        "cropped-cover" &&
      response.request().method() === "POST",
  );
  await page.mouse.move(box!.x + box!.width / 2, box!.y + box!.height / 2);
  await page.mouse.down();
  await page.mouse.move(
    box!.x + box!.width / 2,
    box!.y + box!.height / 2 + 40,
    { steps: 8 },
  );
  await page.mouse.up();
  const cropResponse = await cropped;
  expect(cropResponse.ok()).toBe(true);
  const croppedAsset = await cropResponse.json();
  expect(croppedAsset.parent.id).toBe(asset.id);
  const saved = () =>
    page.request
      .get(`http://localhost:8001/api/nodes/${node.slug}`)
      .then((response) => response.json());
  await expect
    .poll(async () => (await saved()).primary_image?.id)
    .toBe(croppedAsset.id);
  await page.getByRole("button", { name: "View", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "Edit", exact: true }),
  ).toBeVisible();
  await expect
    .poll(() => new URL(page.url()).searchParams.get("edit"))
    .toBeNull();

  const noJS = await browser.newContext({ javaScriptEnabled: false });
  try {
    const reader = await noJS.newPage();
    await reader.goto(page.url());
    const cover = reader.locator('[data-outline-block="cover"] img');
    // The existing loading.tsx boundary can keep streamed content hidden without JS.
    // Assert the actual SSR elements rather than embedded RSC JSON.
    await expect(cover).toHaveCount(1);
    await expect(cover).toHaveAttribute("src", new RegExp(croppedAsset.id));
    await expect(reader.locator("body")).toContainText(
      "Content stays readable before editing.",
    );
  } finally {
    await noJS.close();
  }
});
