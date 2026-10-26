import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import type { Asset } from "../../src/api/openapi-schema";
import { API_ADDRESS } from "../../src/config";
import { DefaultFeedConfig } from "../../src/lib/settings/feed";
import {
  createAdmin,
  getAdminAccessKey,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";
import { pinchOut } from "../touch";

test("feed cover restores its source and crop without changing the social preview", async ({
  page,
}, testInfo) => {
  const handle = `feed-cover-${randomUUID().slice(0, 8)}`;
  const password = "TestPassword123!";
  await createAdmin(page.context(), handle, password);
  await login(page, handle, password);

  await withAdminAccessKey(async ({ adminSettingsUpdate }) => {
    await adminSettingsUpdate({
      metadata: { feed: { blocks: [{ type: "cover" }] } },
    });
  });

  let releaseUpload = () => {};
  let releaseImageLoad = () => {};
  const uploadGate = new Promise<void>((resolve) => {
    releaseUpload = resolve;
  });
  const assetUploads: URL[] = [];
  page.on("request", (request) => {
    const url = new URL(request.url());
    if (url.pathname === "/api/assets" && request.method() === "POST") {
      assetUploads.push(url);
    }
  });

  try {
    const socialImage = await page.evaluate(() => {
      const canvas = document.createElement("canvas");
      canvas.width = 1200;
      canvas.height = 630;
      canvas.getContext("2d")!.fillRect(0, 0, 1200, 630);
      return canvas.toDataURL("image/png").split(",")[1]!;
    });
    const social = await page.request.post(`${API_ADDRESS}/api/info/banner`, {
      headers: {
        Authorization: `Bearer ${getAdminAccessKey()}`,
        "Content-Type": "application/octet-stream",
      },
      data: Buffer.from(socialImage, "base64"),
    });
    expect(social.ok()).toBe(true);

    await page.goto("/");
    await page
      .getByRole("toolbar", { name: "Primary navigation" })
      .getByRole("button", { name: "Edit site", exact: true })
      .click();

    const image = await page.evaluate(() => {
      const canvas = document.createElement("canvas");
      canvas.width = 1920;
      canvas.height = 768;
      const context = canvas.getContext("2d")!;
      context.fillStyle = "#20a060";
      context.fillRect(0, 0, 1920, 768);
      context.fillStyle = "#e04040";
      context.fillRect(0, 0, 1920, 144);
      context.fillStyle = "#4060e0";
      context.fillRect(0, 624, 1920, 144);
      return canvas.toDataURL("image/png").split(",")[1]!;
    });
    const block = page.locator('.block-editor__root[data-block-type="cover"]');
    const initialBounds = await block.boundingBox();
    await page.route(
      (url) =>
        url.pathname === "/api/assets" &&
        url.searchParams.get("filename") === "feed-cover",
      async (route) => {
        await uploadGate;
        await route.continue();
      },
      { times: 1 },
    );
    const uploaded = page.waitForResponse((response) => {
      return (
        new URL(response.url()).pathname === "/api/assets" &&
        new URL(response.url()).searchParams.get("filename") === "feed-cover" &&
        response.request().method() === "POST"
      );
    });
    const sourceUploaded = page.waitForResponse((response) => {
      return (
        new URL(response.url()).searchParams.get("filename") ===
          "source-cover.png" && response.request().method() === "POST"
      );
    });
    const chooser = page.waitForEvent("filechooser");
    await block
      .getByRole("button", { name: "Upload image", exact: true })
      .click();
    await (
      await chooser
    ).setFiles({
      name: "source-cover.png",
      mimeType: "image/png",
      buffer: Buffer.from(image, "base64"),
    });
    await expect
      .poll(() =>
        block
          .locator(".advanced-cropper img")
          .first()
          .evaluate((image: HTMLImageElement) => image.naturalWidth),
      )
      .toBe(1920);

    const cropArea = block.getByRole("group", {
      name: "Crop area",
      exact: true,
    });
    await expect(cropArea).toHaveAccessibleDescription(/arrow keys/);
    await expect(
      block.getByRole("group", { name: "Image cropper", exact: true }),
    ).toBeVisible();
    await expect(cropArea).toHaveAttribute("aria-busy", "true");
    await expect(cropArea).toHaveAttribute("aria-disabled", "true");
    await expect(cropArea).toHaveCSS("opacity", "0.6");
    await expect(
      block.getByRole("status", { name: "Saving image" }),
    ).toBeVisible();
    const spinner = block.locator(".image-cropper__status .spinner");
    const initialTransform = await spinner.evaluate(
      (element) => getComputedStyle(element).transform,
    );
    await expect
      .poll(() =>
        spinner.evaluate((element) => getComputedStyle(element).transform),
      )
      .not.toBe(initialTransform);
    const savingBounds = await block.boundingBox();
    expect(savingBounds!.height).toBe(initialBounds!.height);
    await block.screenshot({
      path: testInfo.outputPath("cover-editor-desktop.png"),
    });
    const viewport = page.viewportSize()!;
    await page.setViewportSize({ width: 375, height: 812 });
    const editBoundsBefore = (await block
      .getByRole("button", { name: "Image tools" })
      .boundingBox())!;
    await block.getByRole("button", { name: "Image tools" }).hover();
    await expect(
      block.getByRole("button", { name: "Replace image" }),
    ).toBeInViewport({ ratio: 1 });
    await expect(
      block.getByRole("button", { name: "Remove image" }),
    ).toBeInViewport({ ratio: 1 });
    const mobileBounds = (await block.boundingBox())!;
    await expect(
      block.getByRole("button", { name: "Replace image" }),
    ).toBeDisabled();
    await expect(
      block.getByRole("button", { name: "Remove image" }),
    ).toBeDisabled();
    const replaceBounds = (await block
      .getByRole("button", { name: "Replace image" })
      .boundingBox())!;
    expect(replaceBounds.x).toBeGreaterThan(
      mobileBounds.x + mobileBounds.width / 2,
    );
    expect(replaceBounds.x + replaceBounds.width).toBeLessThan(
      mobileBounds.x + mobileBounds.width,
    );
    expect(replaceBounds.y + replaceBounds.height).toBeLessThan(
      mobileBounds.y + mobileBounds.height / 2,
    );
    const statusBounds = (await block
      .locator(".image-cropper__status")
      .boundingBox())!;
    expect(statusBounds.height).toBe(replaceBounds.height);
    const editBoundsAfter = (await block
      .getByRole("button", { name: "Image tools" })
      .boundingBox())!;
    const removeBounds = (await block
      .getByRole("button", { name: "Remove image" })
      .boundingBox())!;
    expect(editBoundsAfter.x).toBeCloseTo(editBoundsBefore.x, 0);
    expect(replaceBounds.x + replaceBounds.width).toBeLessThanOrEqual(
      removeBounds.x,
    );
    expect(removeBounds.x + removeBounds.width).toBeLessThanOrEqual(
      editBoundsAfter.x,
    );
    expect(replaceBounds.y).toBeCloseTo(editBoundsAfter.y, 0);
    expect(removeBounds.y).toBeCloseTo(editBoundsAfter.y, 0);
    await block.screenshot({
      path: testInfo.outputPath("cover-editor-mobile.png"),
    });
    await page.setViewportSize(viewport);

    const sourceResponse = await sourceUploaded;
    expect(sourceResponse.ok()).toBe(true);
    const sourceAsset: Asset = await sourceResponse.json();
    const sourceBytes = await page.request.get(
      `${API_ADDRESS}${sourceAsset.path}`,
    );
    expect(await sourceBytes.body()).toEqual(Buffer.from(image, "base64"));

    releaseUpload();

    const response = await uploaded;
    expect(response.ok()).toBe(true);
    let asset: Asset = await response.json();
    const original = await page.request.get(`${API_ADDRESS}${asset.path}`);
    const bytes = await original.body();
    expect(bytes.readUInt32BE(16)).toBe(1920);
    expect(bytes.readUInt32BE(20)).toBe(480);

    await expect
      .poll(async () => {
        const info = await page.request.get(`${API_ADDRESS}/api/info`);
        return (await info.json()).metadata.feed.blocks[0].asset;
      })
      .toEqual({
        id: asset.id,
        path: asset.path,
        parent: { id: sourceAsset.id, path: sourceAsset.path },
      });

    const savedIndicator = block.getByRole("status", { name: "Image saved" });
    await expect(savedIndicator).toBeVisible();
    await expect(savedIndicator).toBeHidden();
    await expect(cropArea).toHaveAttribute("aria-busy", "false");
    await expect(cropArea).toHaveAttribute("aria-disabled", "false");
    await expect(cropArea).toHaveCSS("opacity", "1");

    await page.mouse.move(0, 0);
    const imageTools = block.getByRole("button", { name: "Image tools" });
    await expect(imageTools).toHaveAttribute("aria-expanded", "false");
    await expect
      .poll(async () => {
        const bounds = await block
          .locator(".image-cropper__toolbar")
          .boundingBox();

        return Math.abs(bounds!.width - bounds!.height);
      })
      .toBeLessThan(0.5);

    const collapsedWidth = (await block
      .locator(".image-cropper__toolbar")
      .boundingBox())!.width;
    await imageTools.focus();
    await imageTools.press("Enter");
    await expect(
      block.getByRole("button", { name: "Remove image" }),
    ).toBeVisible();
    await expect
      .poll(
        async () =>
          (await block.locator(".image-cropper__toolbar").boundingBox())!.width,
      )
      .toBeGreaterThan(collapsedWidth);
    await expect(
      block.getByRole("button", { name: "Replace image" }),
    ).toBeEnabled();
    const repositioned = page.waitForResponse((response) => {
      return (
        new URL(response.url()).pathname === "/api/assets" &&
        new URL(response.url()).searchParams.get("filename") === "feed-cover" &&
        response.request().method() === "POST"
      );
    });
    const cropper = await block.locator(".advanced-cropper").boundingBox();
    expect(cropper).not.toBeNull();
    const centerX = cropper!.x + cropper!.width / 2;
    const centerY = cropper!.y + cropper!.height / 2;
    await page.mouse.move(centerX, centerY);
    await page.mouse.down();
    await page.mouse.move(centerX, centerY - 40, { steps: 5 });
    await page.mouse.up();
    await page
      .getByRole("toolbar", { name: "Primary navigation" })
      .getByRole("button", { name: "Exit edit mode", exact: true })
      .click();
    const repositionedResponse = await repositioned;
    expect(repositionedResponse.ok()).toBe(true);
    const previousID = asset.id;
    asset = await repositionedResponse.json();
    expect(asset.id).not.toBe(previousID);
    await expect
      .poll(async () => {
        const info = await page.request.get(`${API_ADDRESS}/api/info`);
        return (await info.json()).metadata.feed.blocks[0].asset.id;
      })
      .toBe(asset.id);

    const movedInfo = await page.request.get(`${API_ADDRESS}/api/info`);
    const savedCrop = (await movedInfo.json()).metadata.feed.blocks[0].crop;
    expect(savedCrop.top).toBeGreaterThan(144);
    expect(savedCrop.width).toBe(1920);
    expect(savedCrop.height).toBe(480);

    const imageLoadGate = new Promise<void>((resolve) => {
      releaseImageLoad = resolve;
    });
    await page.route(
      `${API_ADDRESS}${sourceAsset.path}`,
      async (route) => {
        await imageLoadGate;
        await route.continue();
      },
      { times: 1 },
    );

    await page.reload();
    await page
      .getByRole("toolbar", { name: "Primary navigation" })
      .getByRole("button", { name: "Edit site", exact: true })
      .click();
    const loadingIndicator = block.getByRole("status", {
      name: "Loading image",
    });
    await expect(loadingIndicator).toBeVisible();
    await expect(cropArea).toHaveAttribute("aria-busy", "true");
    await expect(cropArea).toHaveAttribute("aria-disabled", "true");
    releaseImageLoad();
    await expect(loadingIndicator).toBeHidden();
    await expect(cropArea).toHaveAttribute("aria-busy", "false");
    await expect(cropArea).toHaveAttribute("aria-disabled", "false");

    const reopenedImage = block.locator(".advanced-cropper img").first();
    await expect(reopenedImage).toHaveAttribute(
      "src",
      `${API_ADDRESS}${sourceAsset.path}`,
    );
    await expect
      .poll(() =>
        reopenedImage.evaluate(
          (image: HTMLImageElement) => image.naturalHeight,
        ),
      )
      .toBe(768);
    await expect
      .poll(async () => {
        const frame = (await block.locator(".advanced-cropper").boundingBox())!;
        const image = (await reopenedImage.boundingBox())!;
        return ((frame.y - image.y) / image.height) * 768;
      })
      .toBeCloseTo(savedCrop.top, 0);

    const recropped = page.waitForResponse((response) => {
      return (
        new URL(response.url()).searchParams.get("filename") === "feed-cover" &&
        response.request().method() === "POST"
      );
    });
    await cropArea.focus();
    await expect(cropArea).toBeFocused();
    await cropArea.press("Shift+ArrowDown");
    await cropArea.press("Shift+ArrowDown");
    await cropArea.press("Shift+ArrowDown");
    const recroppedResponse = await recropped;
    expect(recroppedResponse.ok()).toBe(true);
    expect(
      new URL(recroppedResponse.url()).searchParams.get("parent_asset_id"),
    ).toBe(sourceAsset.id);
    asset = await recroppedResponse.json();
    await expect
      .poll(async () => {
        const info = await page.request.get(`${API_ADDRESS}/api/info`);
        return (await info.json()).metadata.feed.blocks[0].asset.id;
      })
      .toBe(asset.id);

    const topPixel = await page.evaluate(async (url) => {
      const response = await fetch(url);
      const bitmap = await createImageBitmap(await response.blob());
      const canvas = document.createElement("canvas");
      canvas.width = 1;
      canvas.height = 1;
      const context = canvas.getContext("2d")!;
      context.drawImage(bitmap, 960, 0, 1, 1, 0, 0, 1, 1);
      bitmap.close();
      return Array.from(context.getImageData(0, 0, 1, 1).data);
    }, `${API_ADDRESS}${asset.path}`);
    expect(topPixel).toEqual([224, 64, 64, 255]);

    const sourceUploads = assetUploads.filter(
      (url) => url.searchParams.get("filename") === "source-cover.png",
    );
    const cropUploads = assetUploads.filter(
      (url) => url.searchParams.get("filename") === "feed-cover",
    );
    expect(sourceUploads).toHaveLength(1);
    expect(cropUploads).toHaveLength(3);
    for (const upload of cropUploads) {
      expect(upload.searchParams.get("parent_asset_id")).toBe(sourceAsset.id);
    }

    const zoomed = page.waitForResponse(
      (r) =>
        r.request().method() === "POST" &&
        new URL(r.url()).searchParams.get("filename") === "feed-cover",
    );
    const zoomBounds = (await cropArea.boundingBox())!;
    await page.mouse.move(
      zoomBounds.x + zoomBounds.width * 0.3,
      zoomBounds.y + zoomBounds.height * 0.65,
    );
    await page.mouse.wheel(0, -240);
    asset = await (await zoomed).json();
    await expect
      .poll(
        async () =>
          (await (await page.request.get(`${API_ADDRESS}/api/info`)).json())
            .metadata.feed.blocks[0].asset.id,
      )
      .toBe(asset.id);
    const wheelCrop = (
      await (await page.request.get(`${API_ADDRESS}/api/info`)).json()
    ).metadata.feed.blocks[0].crop;
    expect(wheelCrop.width).toBeLessThan(1920);
    expect(1920 / wheelCrop.width).toBeLessThan(1.04);
    await expect(cropArea).toHaveAttribute("aria-disabled", "false");
    const pinched = page.waitForResponse(
      (r) =>
        r.request().method() === "POST" &&
        new URL(r.url()).searchParams.get("filename") === "feed-cover",
    );
    await pinchOut(page, cropArea);
    asset = await (await pinched).json();
    expect(asset.parent?.id).toBe(sourceAsset.id);
    await expect
      .poll(
        async () =>
          (await (await page.request.get(`${API_ADDRESS}/api/info`)).json())
            .metadata.feed.blocks[0].asset.id,
      )
      .toBe(asset.id);
    const pinchCrop = (
      await (await page.request.get(`${API_ADDRESS}/api/info`)).json()
    ).metadata.feed.blocks[0].crop;
    expect(pinchCrop.width).toBeLessThan(wheelCrop.width);
    const zoomBytes = await (
      await page.request.get(`${API_ADDRESS}${asset.path}`)
    ).body();
    expect(zoomBytes.readUInt32BE(16)).toBe(1920);
    expect(zoomBytes.readUInt32BE(20)).toBe(480);
    expect((await cropArea.boundingBox())!.height).toBeCloseTo(
      zoomBounds.height,
      0,
    );

    const unchangedSocial = await page.request.get(
      `${API_ADDRESS}/api/info/banner?variant=original`,
    );
    expect(await unchangedSocial.body()).toEqual(
      Buffer.from(socialImage, "base64"),
    );

    await page.goto("/admin/brand");
    const bannerChooser = page.waitForEvent("filechooser");
    await page
      .getByRole("button", { name: "Upload banner", exact: true })
      .click();
    await (
      await bannerChooser
    ).setFiles({
      name: "social.png",
      mimeType: "image/png",
      buffer: Buffer.from(image, "base64"),
    });
    await expect
      .poll(() =>
        page
          .locator(".advanced-cropper img")
          .first()
          .evaluate((image: HTMLImageElement) => image.naturalWidth),
      )
      .toBe(1920);
    const banner = page.waitForResponse((response) => {
      return (
        new URL(response.url()).pathname === "/api/info/banner" &&
        response.request().method() === "POST"
      );
    });
    await page
      .getByRole("button", { name: "Save banner", exact: true })
      .click();
    expect((await banner).ok()).toBe(true);
    const changedSocial = await page.request.get(
      `${API_ADDRESS}/api/info/banner?variant=original`,
    );
    const socialBytes = await changedSocial.body();
    expect(socialBytes.readUInt32BE(16)).toBe(1200);
    expect(socialBytes.readUInt32BE(20)).toBe(630);
    expect(socialBytes).not.toEqual(Buffer.from(socialImage, "base64"));

    await page.goto("/");
    const cover = page.locator('[data-sd-block="cover"] img');
    await expect(cover).toBeVisible();
    await expect(cover).toHaveAttribute(
      "src",
      `${API_ADDRESS}${asset.path}?variant=l`,
    );
    for (const width of [1280, 375]) {
      await page.setViewportSize({ width, height: 812 });
      const bounds = await cover.boundingBox();
      expect(bounds).not.toBeNull();
      expect(bounds!.width / bounds!.height).toBeCloseTo(4, 1);
    }

    await page.setViewportSize({ width: 1280, height: 812 });
    await page.route(`${API_ADDRESS}${sourceAsset.path}`, (route) =>
      route.abort("failed"),
    );
    await page
      .getByRole("toolbar", { name: "Primary navigation" })
      .getByRole("button", { name: "Edit site", exact: true })
      .click();
    await expect(
      block.getByRole("status", {
        name: "The image could not be opened. Please choose another.",
      }),
    ).toBeVisible();
    await expect(cropArea).toHaveAttribute("aria-busy", "false");
    await expect(cropArea).toHaveAttribute("aria-disabled", "true");
    await block.getByRole("button", { name: "Image tools" }).focus();
    await block.getByRole("button", { name: "Image tools" }).press("Enter");
    await block.getByRole("button", { name: "Remove image" }).click();
    await expect(
      block.getByRole("button", { name: "Upload image" }),
    ).toBeVisible();
    await expect
      .poll(async () => {
        const info = await page.request.get(`${API_ADDRESS}/api/info`);
        return (await info.json()).metadata.feed.blocks[0];
      })
      .toEqual({ type: "cover" });
    await page.reload();
    await expect(cover).toHaveCount(0);
    const preservedSource = await page.request.get(
      `${API_ADDRESS}${sourceAsset.path}`,
    );
    expect(preservedSource.ok()).toBe(true);
    const preservedSocial = await page.request.get(
      `${API_ADDRESS}/api/info/banner?variant=original`,
    );
    expect(await preservedSocial.body()).toEqual(socialBytes);
  } finally {
    releaseUpload();
    releaseImageLoad();
    await withAdminAccessKey(async ({ adminSettingsUpdate }) => {
      await adminSettingsUpdate({ metadata: { feed: DefaultFeedConfig } });
    });
  }
});
