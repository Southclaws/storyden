import {
  type Locator,
  type Page,
  type Response,
  expect,
  test,
} from "@playwright/test";
import { randomUUID } from "node:crypto";

import { API_ADDRESS } from "../../src/config";
import { createAdmin, login } from "../access_key_admin_assignment";
import { pinchOut } from "../touch";

function isBannerSave(response: Response) {
  return (
    new URL(response.url()).pathname === "/api/admin" &&
    response.request().method() === "PATCH" &&
    Boolean(response.request().postDataJSON()?.banner_asset_id)
  );
}

async function waitForImageToSettle(image: Locator) {
  let previous: string | undefined;
  let stable = 0;
  await expect
    .poll(
      async () => {
        const bounds = await image.boundingBox();
        if (!bounds) return 0;
        const current = JSON.stringify(bounds);
        stable = current === previous ? stable + 1 : 0;
        previous = current;
        return stable;
      },
      { intervals: [100] },
    )
    .toBeGreaterThanOrEqual(3);
}

async function openBranding(page: Page) {
  const handle = `banner-${randomUUID().slice(0, 8)}`;
  await createAdmin(page.context(), handle, "TestPassword123!");
  await login(page, handle, "TestPassword123!");
  await page.goto("/admin/brand");
  const editor = page.getByRole("group", { name: "Banner image", exact: true });
  await expect(
    editor.getByRole("group", { name: "Crop area", exact: true }),
  ).toHaveAttribute("aria-disabled", "false");
  return editor;
}

async function sourceImage(page: Page, mimeType = "image/png") {
  return Buffer.from(
    await page.evaluate((mimeType) => {
      const canvas = document.createElement("canvas");
      canvas.width = 1200;
      canvas.height = 1200;
      const context = canvas.getContext("2d")!;
      for (const [index, color] of [
        "#e04040",
        "#20a060",
        "#4060e0",
      ].entries()) {
        context.fillStyle = color;
        context.fillRect(0, index * 400, 1200, 400);
      }
      return canvas.toDataURL(mimeType).split(",")[1]!;
    }, mimeType),
    "base64",
  );
}

async function savedBanner(page: Page) {
  const response = await page.request.get(
    `${API_ADDRESS}/api/info/banner?variant=original`,
  );
  expect(response.ok()).toBe(true);
  const png = await response.body();
  expect(png.subarray(1, 4).toString()).toBe("PNG");
  expect(png.readUInt32BE(16)).toBe(1200);
  expect(png.readUInt32BE(20)).toBe(630);
  return png;
}

test("branding banner saves exact-size crops, pans and zooms with shared controls", async ({
  page,
}) => {
  const editor = await openBranding(page);
  const area = editor.getByRole("group", { name: "Crop area", exact: true });
  const source = await sourceImage(page);
  let releaseSave = () => {};
  const saveGate = new Promise<void>((resolve) => {
    releaseSave = resolve;
  });
  await page.route(`${API_ADDRESS}/api/admin`, async (route) => {
    if (route.request().method() === "PATCH") await saveGate;
    await route.continue();
  });
  try {
    const uploaded = page.waitForResponse(isBannerSave);
    const chooser = page.waitForEvent("filechooser");
    await editor
      .getByRole("button", { name: "Replace image", exact: true })
      .click();
    await (
      await chooser
    ).setFiles({ name: "banner.png", mimeType: "image/png", buffer: source });
    await expect
      .poll(() =>
        editor
          .locator("img")
          .first()
          .evaluate((image: HTMLImageElement) => image.naturalHeight),
      )
      .toBe(1200);
    await expect(area).toHaveAttribute("aria-disabled", "false");
    await page
      .getByRole("button", { name: "Save banner", exact: true })
      .click();
    await expect(
      editor.getByRole("status", { name: "Saving image", exact: true }),
    ).toBeVisible();
    await expect(area).toHaveAttribute("aria-disabled", "true");
    await expect(
      editor.getByRole("button", { name: "Replace image", exact: true }),
    ).toBeDisabled();
    releaseSave();
    expect((await uploaded).ok()).toBe(true);
    await expect(area).toHaveAttribute("aria-disabled", "false");
    const initial = await savedBanner(page);
    const assigned = await (
      await page.request.get(`${API_ADDRESS}/api/admin`)
    ).json();
    const original = assigned.banner_asset.parent;
    expect(original).toBeDefined();
    expect(assigned.metadata.branding.banner.asset_id).toBe(
      assigned.banner_asset.id,
    );
    expect(
      await (await page.request.get(`${API_ADDRESS}${original.path}`)).body(),
    ).toEqual(source);
    const colors = await page.evaluate(async (base64) => {
      const bytes = Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
      const image = await createImageBitmap(
        new Blob([bytes], { type: "image/png" }),
      );
      const canvas = document.createElement("canvas");
      canvas.width = image.width;
      canvas.height = image.height;
      const context = canvas.getContext("2d")!;
      context.drawImage(image, 0, 0);
      image.close();
      return [0, 315, 629].map((y) =>
        Array.from(context.getImageData(600, y, 1, 1).data),
      );
    }, initial.toString("base64"));
    expect(colors).toEqual([
      [224, 64, 64, 255],
      [32, 160, 96, 255],
      [64, 96, 224, 255],
    ]);

    const panned = page.waitForResponse(isBannerSave);
    await area.press("Shift+ArrowDown");
    await page
      .getByRole("button", { name: "Save banner", exact: true })
      .click();
    expect((await panned).ok()).toBe(true);
    await expect(area).toHaveAttribute("aria-disabled", "false");
    const moved = await savedBanner(page);
    expect(moved).not.toEqual(initial);

    const preview = area.locator("img").first();
    await waitForImageToSettle(preview);
    const beforeWheel = (await preview.boundingBox())!.width;
    const wheelBounds = (await area.boundingBox())!;
    await page.mouse.move(
      wheelBounds.x + wheelBounds.width * 0.3,
      wheelBounds.y + wheelBounds.height * 0.65,
    );
    await page.mouse.wheel(0, -120);
    await expect
      .poll(async () => (await preview.boundingBox())!.width)
      .toBeGreaterThan(beforeWheel);
    const scrolled = page.waitForResponse(isBannerSave);
    await page
      .getByRole("button", { name: "Save banner", exact: true })
      .click();
    expect((await scrolled).ok()).toBe(true);
    await expect(area).toHaveAttribute("aria-disabled", "false");
    const zoomed = await savedBanner(page);
    expect(zoomed).not.toEqual(moved);

    await waitForImageToSettle(preview);
    const beforePinch = (await preview.boundingBox())!.width;
    await pinchOut(page, area);
    await expect
      .poll(async () => (await preview.boundingBox())!.width)
      .toBeGreaterThan(beforePinch);
    const pinched = page.waitForResponse(isBannerSave);
    await page
      .getByRole("button", { name: "Save banner", exact: true })
      .click();
    expect((await pinched).ok()).toBe(true);
    await expect(area).toHaveAttribute("aria-disabled", "false");
    const finalCrop = await savedBanner(page);
    expect(finalCrop).not.toEqual(zoomed);

    await page.reload();
    await expect(area).toHaveAttribute("aria-disabled", "false");
    await expect
      .poll(() =>
        preview.evaluate((image: HTMLImageElement) => image.naturalHeight),
      )
      .toBe(1200);
    const reopened = page.waitForResponse(isBannerSave);
    await page
      .getByRole("button", { name: "Save banner", exact: true })
      .click();
    expect((await reopened).ok()).toBe(true);
    await expect(area).toHaveAttribute("aria-disabled", "false");
    expect(await savedBanner(page)).toEqual(finalCrop);
    const recropped = await (
      await page.request.get(`${API_ADDRESS}/api/admin`)
    ).json();
    expect(recropped.banner_asset.parent.id).toBe(original.id);
    expect(
      await (await page.request.get(`${API_ADDRESS}${original.path}`)).body(),
    ).toEqual(source);
    for (const width of [1280, 390]) {
      await page.setViewportSize({ width, height: 900 });
      const bounds = (await area.boundingBox())!;
      expect(bounds.width / bounds.height).toBeCloseTo(1200 / 630, 2);
      expect(bounds.width).toBeLessThan(width);
    }
  } finally {
    releaseSave();
  }
});

test("branding banner recovers from a save failure without losing the editable source", async ({
  page,
}) => {
  const editor = await openBranding(page);
  const area = editor.getByRole("group", { name: "Crop area", exact: true });
  const source = await sourceImage(page, "image/webp");
  await page.route(`${API_ADDRESS}/api/admin`, async (route) => {
    if (route.request().method() === "PATCH") {
      await route.fulfill({
        status: 503,
        json: {
          trace_id: "banner-save-failure",
          title: "Banner save failed",
          detail: "Banner save failed",
          status: 503,
        },
      });
    } else await route.continue();
  });
  const failed = page.waitForResponse(isBannerSave);
  const chooser = page.waitForEvent("filechooser");
  await editor
    .getByRole("button", { name: "Replace image", exact: true })
    .click();
  await (
    await chooser
  ).setFiles({ name: "banner.webp", mimeType: "image/webp", buffer: source });
  await expect
    .poll(() =>
      editor
        .locator("img")
        .first()
        .evaluate((image: HTMLImageElement) => image.naturalHeight),
    )
    .toBe(1200);
  await expect(area).toHaveAttribute("aria-disabled", "false");
  await page.getByRole("button", { name: "Save banner", exact: true }).click();
  expect((await failed).status()).toBe(503);
  await expect(
    page.getByText("Banner save failed", { exact: true }),
  ).toBeVisible();
  await expect(area).toHaveAttribute("aria-disabled", "false");
  await expect(
    editor.getByRole("button", { name: "Replace image", exact: true }),
  ).toBeEnabled();
  await page.unroute(`${API_ADDRESS}/api/admin`);
  const retried = page.waitForResponse(isBannerSave);
  await page.getByRole("button", { name: "Save banner", exact: true }).click();
  expect((await retried).ok()).toBe(true);
  await expect(
    editor.getByRole("status", { name: "Image saved", exact: true }),
  ).toBeVisible();
  await savedBanner(page);
});

test("legacy banner uploads replace the assigned image without restoring an unrelated crop", async ({
  page,
}) => {
  const editor = await openBranding(page);
  const area = editor.getByRole("group", { name: "Crop area", exact: true });
  const source = await sourceImage(page);
  const chooser = page.waitForEvent("filechooser");
  await editor
    .getByRole("button", { name: "Replace image", exact: true })
    .click();
  await (
    await chooser
  ).setFiles({ name: "banner.png", mimeType: "image/png", buffer: source });
  const preview = area.locator("img").first();
  await expect
    .poll(() =>
      preview.evaluate((image: HTMLImageElement) => image.naturalHeight),
    )
    .toBe(1200);
  await expect(area).toHaveAttribute("aria-disabled", "false");
  await waitForImageToSettle(preview);
  const beforeZoom = (await preview.boundingBox())!.width;
  const bounds = (await area.boundingBox())!;
  await page.mouse.move(
    bounds.x + bounds.width / 2,
    bounds.y + bounds.height / 2,
  );
  await page.mouse.wheel(0, -500);
  await expect
    .poll(async () => (await preview.boundingBox())!.width)
    .toBeGreaterThan(beforeZoom);
  const saved = page.waitForResponse(isBannerSave);
  await page.getByRole("button", { name: "Save banner", exact: true }).click();
  expect((await saved).ok()).toBe(true);
  const previous = await (
    await page.request.get(`${API_ADDRESS}/api/admin`)
  ).json();
  expect(previous.metadata.branding.banner.crop.width).toBeLessThan(1);

  const legacy = await page.request.post(`${API_ADDRESS}/api/info/banner`, {
    headers: { "Content-Type": "application/octet-stream" },
    data: source,
  });
  expect(legacy.ok()).toBe(true);
  const replaced = await (
    await page.request.get(`${API_ADDRESS}/api/admin`)
  ).json();
  expect(replaced.banner_asset.id).not.toBe(previous.banner_asset.id);
  expect(replaced.metadata.branding.banner.asset_id).toBe(
    previous.banner_asset.id,
  );
  await page.reload();
  await expect(area).toHaveAttribute("aria-disabled", "false");
  const resaved = page.waitForResponse(isBannerSave);
  await page.getByRole("button", { name: "Save banner", exact: true }).click();
  expect((await resaved).ok()).toBe(true);
  const current = await (
    await page.request.get(`${API_ADDRESS}/api/admin`)
  ).json();
  expect(current.banner_asset.parent.id).toBe(replaced.banner_asset.id);
  expect(current.metadata.branding.banner.asset_id).toBe(
    current.banner_asset.id,
  );
  expect(current.metadata.branding.banner.crop.width).toBeCloseTo(1, 5);
  // Fractional CSS bounds can differ by a fraction of a source pixel.
  expect(
    Math.abs(current.metadata.branding.banner.crop.height * 1200 - 630),
  ).toBeLessThan(1);
});
