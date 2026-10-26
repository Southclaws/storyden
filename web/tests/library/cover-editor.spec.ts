import { Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import type { Asset } from "../../src/api/openapi-schema";
import { API_ADDRESS } from "../../src/config";
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
    canvas.width = 1600;
    canvas.height = 800;
    const ctx = canvas.getContext("2d")!;
    ctx.fillStyle = "#2255aa";
    ctx.fillRect(0, 0, 1600, 400);
    ctx.fillStyle = "#cc5533";
    ctx.fillRect(0, 400, 1600, 400);
    return canvas.toDataURL("image/png").split(",")[1]!;
  });
  const uploaded = await page.request.post(
    `${API_ADDRESS}/api/assets?filename=library-cover.png`,
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
  const variantURL = `${API_ADDRESS}${asset.path}?variant=l`;
  const cold = await page.request.get(variantURL);
  expect(cold.headers()["cache-control"]).toContain("no-cache");
  await expect
    .poll(async () => {
      const response = await page.request.get(variantURL);
      return response.headers()["cache-control"];
    })
    .toContain("86400");
  const warm = await page.request.get(variantURL);
  expect(warm.headers()["content-type"]).toBe("image/webp");
  const originalBytes = Buffer.from(image, "base64").length;
  const warmBytes = (await warm.body()).length;
  expect(warmBytes).toBeLessThan(originalBytes);
  expect(warm.headers()["cache-control"]).toContain("86400");
  await test.info().attach("image-transfer", {
    body: JSON.stringify({ originalBytes, warmBytes }),
    contentType: "application/json",
  });

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
  const coverSource = await main
    .locator('[data-outline-block="cover"] img')
    .evaluate((image) => (image as HTMLImageElement).currentSrc);
  expect(["s", "m", "l"]).toContain(
    new URL(coverSource).searchParams.get("variant"),
  );

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
    '[aria-label="Image cropper"] .advanced-cropper',
  );
  const croppedAsset = await dragCover(page, 40);
  expect(croppedAsset.parent?.id).toBe(asset.id);
  const saved = () =>
    page.request
      .get(`${API_ADDRESS}/api/nodes/${node.slug}`)
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

  await expectVariant(page, asset, "s", 320, 160);
  await expectVariant(page, croppedAsset, "s", 320, 80);

  await page.reload({ waitUntil: "networkidle" });
  await page.getByRole("button", { name: "Edit", exact: true }).click();
  await page.getByRole("menuitem", { name: "Quick edit", exact: true }).click();
  await expect(cropper).toBeVisible();
  await expect(cropper.locator("img").first()).toHaveAttribute(
    "src",
    `${API_ADDRESS}${asset.path}`,
  );
  const recropped = await dragCover(page, -40);
  expect(recropped.id).not.toBe(croppedAsset.id);
  expect(recropped.parent?.id).toBe(asset.id);
  await expect
    .poll(async () => (await saved()).primary_image?.id)
    .toBe(recropped.id);
  const original = await page.request.get(`${API_ADDRESS}${asset.path}`);
  expect(await original.body()).toEqual(Buffer.from(image, "base64"));
  await expectVariant(page, recropped, "s", 320, 80);
  await page.getByRole("button", { name: "View", exact: true }).click();
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
    await expect(cover).toHaveAttribute("src", new RegExp(recropped.id));
    await expect(reader.locator("body")).toContainText(
      "Content stays readable before editing.",
    );
  } finally {
    await noJS.close();
  }
});

test("a cover uploaded in quick edit keeps its original when cropped", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  const admin = `cover-upload-${seed}`;
  await createAdmin(page.context(), admin, "TestPassword123!");
  await login(page, admin, "TestPassword123!");
  const node = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: `Cover upload ${seed}`,
      slug: `cover-upload-${seed}`,
      visibility: "published",
      meta: { layout: { blocks: [{ type: "title" }, { type: "cover" }] } },
    }),
  );
  await page.goto(`/l/${node.slug}`, { waitUntil: "networkidle" });
  await page.getByRole("button", { name: "Edit", exact: true }).click();
  await page.getByRole("menuitem", { name: "Quick edit", exact: true }).click();

  const image = await page.evaluate(() => {
    const canvas = document.createElement("canvas");
    canvas.width = 1600;
    canvas.height = 800;
    const context = canvas.getContext("2d")!;
    context.fillStyle = "#2255aa";
    context.fillRect(0, 0, 1600, 800);
    return canvas.toDataURL("image/png").split(",")[1]!;
  });
  const uploaded = page.waitForResponse(
    (response) =>
      new URL(response.url()).pathname === "/api/assets" &&
      response.request().method() === "POST",
  );
  const chooser = page.waitForEvent("filechooser");
  await page
    .getByRole("button", { name: "Upload main image", exact: true })
    .click();
  await (
    await chooser
  ).setFiles({
    name: "cover.png",
    mimeType: "image/png",
    buffer: Buffer.from(image, "base64"),
  });
  const uploadResponse = await uploaded;
  expect(uploadResponse.ok()).toBe(true);
  const original: Asset = await uploadResponse.json();
  const crop = await dragCover(page, 40);
  expect(crop.parent?.id).toBe(original.id);
  await expect
    .poll(async () => {
      const response = await page.request.get(
        `${API_ADDRESS}/api/nodes/${node.slug}`,
      );
      return (await response.json()).primary_image?.id;
    })
    .toBe(crop.id);
  await expectVariant(page, original, "s", 320, 160);
  await expectVariant(page, crop, "s", 320, 80);
});

async function dragCover(page: Page, distance: number): Promise<Asset> {
  await page
    .getByRole("button", { name: "Main image options", exact: true })
    .click();
  await page
    .getByRole("menuitem", { name: "Page layout", exact: true })
    .hover();
  await page
    .getByRole("menuitem", { name: "Banner · wide page header", exact: true })
    .click();
  const cropper = page.locator(
    '[aria-label="Image cropper"] .advanced-cropper',
  );
  await expect(cropper).toBeVisible();
  await expect(
    page.getByRole("group", { name: "Crop area", exact: true }),
  ).toHaveAttribute("aria-disabled", "false");
  await expect(cropper.locator("img").first()).toBeVisible();
  const box = await cropper.boundingBox();
  expect(box).not.toBeNull();
  const response = page.waitForResponse(
    (response) =>
      new URL(response.url()).pathname === "/api/assets" &&
      new URL(response.url()).searchParams.get("filename") ===
        "cropped-cover" &&
      response.request().method() === "POST",
  );
  const x = box!.x + box!.width / 2;
  const y = box!.y + box!.height / 2;
  await page.mouse.move(x, y);
  await page.mouse.down();
  await page.mouse.move(x, y + distance, { steps: 8 });
  await page.mouse.up();
  const uploaded = await response;
  expect(uploaded.ok()).toBe(true);
  return uploaded.json();
}

async function expectVariant(
  page: Page,
  asset: Asset,
  preset: string,
  width: number,
  height: number,
) {
  const url = `${API_ADDRESS}${asset.path}?variant=${preset}`;
  await expect
    .poll(async () => {
      const response = await page.request.get(url);
      return response.headers()["cache-control"];
    })
    .toContain("86400");
  const response = await page.request.get(url);
  expect(response.headers()["content-type"]).toBe("image/webp");
  const dimensions = await page.evaluate(
    async (bytes) => {
      const bitmap = await createImageBitmap(
        new Blob([new Uint8Array(bytes)], { type: "image/webp" }),
      );
      const size = { width: bitmap.width, height: bitmap.height };
      bitmap.close();
      return size;
    },
    Array.from(await response.body()),
  );
  expect(dimensions).toEqual({ width, height });
}
