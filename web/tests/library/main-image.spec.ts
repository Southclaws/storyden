import { type Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import type { Asset } from "../../src/api/openapi-schema";
import { API_ADDRESS } from "../../src/config";
import {
  createAdmin,
  getAdminAccessKey,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";
import { pinchOut } from "../touch";

async function upload(page: Page, width: number, height: number) {
  const image = await page.evaluate(
    ({ width, height }) => {
      const canvas = document.createElement("canvas");
      canvas.width = width;
      canvas.height = height;
      const ctx = canvas.getContext("2d")!;
      ctx.fillStyle = "#2255aa";
      ctx.fillRect(0, 0, width, height);
      ctx.fillStyle = "#cc5533";
      ctx.fillRect(0, height / 2, width, height / 2);
      return canvas.toDataURL("image/png").split(",")[1]!;
    },
    { width, height },
  );
  const response = await page.request.post(
    `${API_ADDRESS}/api/assets?filename=main-image-${width}x${height}.png`,
    {
      headers: {
        Authorization: `Bearer ${getAdminAccessKey()}`,
        "Content-Type": "application/octet-stream",
      },
      data: Buffer.from(image, "base64"),
    },
  );
  expect(response.ok(), await response.text()).toBe(true);
  return response.json();
}

async function edit(page: Page) {
  await page.getByRole("button", { name: "Edit", exact: true }).click();
  await page.getByRole("menuitem", { name: "Quick edit", exact: true }).click();
  await expect(
    page.getByRole("group", { name: "Crop area", exact: true }),
  ).toHaveAttribute("aria-disabled", "false");
}

async function presentation(page: Page, name: string) {
  await page
    .getByRole("button", { name: "Main image options", exact: true })
    .click();
  await page
    .getByRole("menuitem", { name: "Page layout", exact: true })
    .hover();
  await page.getByRole("menuitem", { name: new RegExp(`^${name} ·`) }).click();
}

async function view(page: Page) {
  await page.getByRole("button", { name: "View", exact: true }).click();
  await page.waitForURL((url) => !url.searchParams.has("edit"));
}

async function saved(page: Page, slug: string) {
  const response = await page.request.get(`${API_ADDRESS}/api/nodes/${slug}`);
  expect(response.ok()).toBe(true);
  return response.json();
}

test("Main image presets preserve assets, match edit/view, and bound stress dimensions", async ({
  page,
}) => {
  test.setTimeout(120_000);
  const seed = randomUUID().slice(0, 8);
  await createAdmin(page.context(), `main-image-${seed}`, "TestPassword123!");
  await login(page, `main-image-${seed}`, "TestPassword123!");
  for (const [width, height] of [
    [1600, 900],
    [600, 900],
    [2400, 400],
    [400, 1600],
    [1, 1000],
    [1000, 1],
    [1, 1],
  ] as const) {
    const asset = await upload(page, width, height);
    const node = await withAdminAccessKey(({ nodeCreate }) =>
      nodeCreate({
        name: `Main image ${width}x${height} ${seed}`,
        slug: `main-image-${width}-${height}-${seed}`,
        visibility: "published",
        primary_image_asset_id: asset.id,
        meta: {
          coverImage: { left: 12, top: 34 },
          layout: {
            blocks: [{ type: "title" }, { type: "cover" }, { type: "content" }],
          },
        },
      }),
    );
    await page.goto(`/l/${node.slug}`);
    const image = page
      .getByRole("main")
      .locator("[data-main-image-layout]:visible img");
    await expect(image).toBeVisible();
    await expect(
      page.getByRole("main").locator("[data-main-image-layout]:visible"),
    ).toHaveAttribute("data-main-image-layout", "banner");
    for (const mode of ["Image", "Artwork", "Banner"]) {
      await edit(page);
      await presentation(page, mode);
      const frame = page
        .getByRole("main")
        .locator("[data-main-image-layout]:visible");
      await expect(frame).toHaveAttribute(
        "data-main-image-layout",
        mode.toLowerCase(),
      );
      await expect(
        page.getByRole("group", { name: "Crop area", exact: true }),
      ).toBeVisible();
      await expect(
        page.getByRole("button", { name: "Crop image", exact: true }),
      ).toHaveCount(0);
      await view(page);
      for (const viewport of [
        { width: 1440, height: 900 },
        { width: 390, height: 844 },
      ]) {
        await page.setViewportSize(viewport);
        await expect
          .poll(() =>
            image.evaluate((el) => (el as HTMLImageElement).naturalWidth),
          )
          .toBeGreaterThan(0);
        const bounds = (await frame.boundingBox())!;
        const pixels = (await image.boundingBox())!;
        expect(pixels.height).toBeLessThanOrEqual(
          mode === "banner" ? viewport.width / 4 : 640,
        );
        expect(pixels.width).toBeLessThanOrEqual(bounds.width + 1);
        if (mode === "Banner")
          expect(bounds.width / bounds.height).toBeCloseTo(4, 1);
        if (mode === "Artwork") {
          expect(bounds.width / bounds.height).toBeCloseTo(1, 1);
          expect(bounds.width).toBeLessThanOrEqual(320);
        }
        if (mode === "Image") {
          if (width >= height)
            expect(pixels.height).toBeCloseTo(
              (pixels.width * height) / width,
              1,
            );
          else
            expect(pixels.width).toBeCloseTo(
              (pixels.height * width) / height,
              1,
            );
        }
      }
      await page.reload();
      await expect(frame).toHaveAttribute(
        "data-main-image-layout",
        mode.toLowerCase(),
      );
      const persisted = await saved(page, node.slug);
      expect(persisted.primary_image.id).toBe(asset.id);
      expect(persisted.meta.coverImage).toEqual({ left: 12, top: 34 });
      expect(
        persisted.meta.layout.blocks.map((b: { type: string }) => b.type),
      ).toEqual(["title", "cover", "content"]);
      expect(persisted.meta.layout.blocks[1].config.layout).toBe(
        mode.toLowerCase(),
      );
    }
  }
});

async function cropRatio(page: Page, name: string) {
  await page
    .getByRole("button", { name: "Main image options", exact: true })
    .click();
  await page.getByRole("menuitem", { name: "Crop ratio", exact: true }).hover();
  await page.getByRole("menuitem", { name, exact: true }).click();
}

function cropResponse(page: Page) {
  return page.waitForResponse(
    (r) =>
      r.request().method() === "POST" &&
      new URL(r.url()).searchParams.get("filename") === "cropped-cover",
  );
}

test("Main image edits immediately, autosaves adjustments, and retains original pixels after reload", async ({
  page,
}) => {
  await page.setViewportSize({ width: 1280, height: 1000 });
  const seed = randomUUID().slice(0, 8);
  await createAdmin(page.context(), `crop-main-${seed}`, "TestPassword123!");
  await login(page, `crop-main-${seed}`, "TestPassword123!");
  const original = await upload(page, 1600, 900);
  const node = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: `Crop ${seed}`,
      slug: `crop-main-${seed}`,
      visibility: "published",
      primary_image_asset_id: original.id,
      meta: {
        layout: {
          blocks: [
            { type: "title" },
            { type: "cover", config: { layout: "image" } },
          ],
        },
      },
    }),
  );
  await page.goto(`/l/${node.slug}`);
  await edit(page);
  const cropper = page.getByRole("group", {
    name: "Crop area",
    exact: true,
  });
  await expect(cropper).toBeVisible();
  for (const name of [
    "Crop image",
    "Apply crop",
    "Zoom in",
    "Zoom out",
    "Cancel",
  ])
    await expect(page.getByRole("button", { name, exact: true })).toHaveCount(
      0,
    );
  await expect(page.getByRole("combobox", { name: "Crop shape" })).toHaveCount(
    0,
  );
  const toolbar = page.locator(
    "[data-main-image-layout] .image-cropper__controls:visible",
  );
  const frame = page.locator("[data-main-image-layout]:visible");
  const outer = (await frame.boundingBox())!;
  const tools = (await toolbar.boundingBox())!;
  expect(tools.x + tools.width).toBeLessThan(outer.x + outer.width);
  expect(tools.y).toBeGreaterThan(outer.y);
  expect(tools.y - outer.y).toBeLessThan(20);
  expect((await saved(page, node.slug)).primary_image.id).toBe(original.id);

  // Hold a real upload open to prove View cannot discard a pending crop.
  let release!: () => void;
  let started!: () => void;
  const held = new Promise<void>((resolve) => {
    release = resolve;
  });
  const uploading = new Promise<void>((resolve) => {
    started = resolve;
  });
  await page.route("**/api/assets?**", async (route) => {
    if (
      new URL(route.request().url()).searchParams.get("filename") ===
      "cropped-cover"
    ) {
      started();
      await held;
    }
    await route.continue();
  });
  const cropped = cropResponse(page);
  await cropRatio(page, "Square · 1:1");
  await uploading;
  try {
    await expect(
      page.getByRole("button", { name: "View", exact: true }),
    ).toBeDisabled();
    await expect(cropper).toHaveAttribute("aria-busy", "true");
  } finally {
    release();
  }
  const child = await (await cropped).json();
  await page.unrouteAll({ behavior: "wait" });
  expect(child.parent.id).toBe(original.id);
  await expect(
    page.getByRole("button", { name: "View", exact: true }),
  ).toBeEnabled();
  expect((await frame.boundingBox())!.height).toBeCloseTo(outer.height, 0);
  // Move the selection toward the bottom of the original without resetting
  // the workspace or changing the page layout to repair its boundary.
  const selected = (await page
    .locator(".image-cropper__selection")
    .boundingBox())!;
  const dragged = cropResponse(page);
  await page.mouse.move(
    selected.x + selected.width / 2,
    selected.y + selected.height / 2,
  );
  await page.mouse.down();
  await page.mouse.move(
    selected.x + selected.width / 2,
    selected.y + selected.height / 2 + 60,
    { steps: 8 },
  );
  await page.mouse.up();
  const draggedAsset = await (await dragged).json();
  expect(draggedAsset.parent.id).toBe(original.id);
  await expect
    .poll(
      async () =>
        (await saved(page, node.slug)).meta.layout.blocks[1].config.crop.top,
    )
    .toBeGreaterThan(0.1);
  expect((await frame.boundingBox())!.height).toBeCloseTo(outer.height, 0);
  await view(page);
  await page.reload();
  let persisted = await saved(page, node.slug);
  expect(persisted.primary_image.id).toBe(draggedAsset.id);
  const crop = persisted.meta.layout.blocks[1].config.crop;
  expect(crop.width * 1600).toBeCloseTo(crop.height * 900, 0);

  await edit(page);
  await expect(cropper.locator("img").first()).toHaveAttribute(
    "src",
    `${API_ADDRESS}${original.path}`,
  );
  const rectangle = (await page
    .locator(".image-cropper__selection")
    .boundingBox())!;
  expect(rectangle.width / rectangle.height).toBeCloseTo(1, 1);
  const moved = cropResponse(page);
  await cropper.press("ArrowRight");
  const movedAsset = await (await moved).json();
  expect(movedAsset.parent.id).toBe(original.id);
  await expect(
    page.getByRole("button", { name: "View", exact: true }),
  ).toBeEnabled();
  const free = cropResponse(page);
  await cropRatio(page, "Free crop");
  const freeAsset = await (await free).json();
  await expect(
    page.getByRole("button", { name: "View", exact: true }),
  ).toBeEnabled();
  const freeRectangle = (await page
    .locator(".image-cropper__selection")
    .boundingBox())!;
  const resized = cropResponse(page);
  const right = freeRectangle.x + freeRectangle.width;
  const bottom = freeRectangle.y + freeRectangle.height;
  await page.mouse.move(right - 1, bottom - 1);
  await page.mouse.down();
  await page.mouse.move(right - 81, bottom - 1, { steps: 8 });
  await page.mouse.up();
  const resizedAsset = await (await resized).json();
  expect(resizedAsset.id).not.toBe(freeAsset.id);
  await expect(
    page.getByRole("button", { name: "View", exact: true }),
  ).toBeEnabled();
  const resizedRectangle = (await page
    .locator(".image-cropper__selection")
    .boundingBox())!;
  expect(resizedRectangle.width / resizedRectangle.height).toBeLessThan(
    freeRectangle.width / freeRectangle.height - 0.1,
  );
  expect((await frame.boundingBox())!.height).toBeCloseTo(outer.height, 0);
  const expanded = cropResponse(page);
  await cropRatio(page, "Original proportions");
  const expandedAsset = await (await expanded).json();
  expect(expandedAsset.parent.id).toBe(original.id);
  await view(page);
  await page.reload();
  persisted = await saved(page, node.slug);
  expect(persisted.primary_image.id).toBe(expandedAsset.id);
  expect(persisted.meta.layout.blocks[1].config.crop.width).toBeCloseTo(1, 3);
  expect(persisted.meta.layout.blocks[1].config.crop.height).toBeCloseTo(1, 3);
});

for (const mode of ["Banner", "Artwork"] as const) {
  test(`${mode} fills a fixed crop viewport and autosaves dragging, wheel zoom and pinch zoom`, async ({
    page,
  }) => {
    const seed = randomUUID().slice(0, 8);
    await createAdmin(page.context(), `fixed-${seed}`, "TestPassword123!");
    await login(page, `fixed-${seed}`, "TestPassword123!");
    const original = await upload(page, 1600, 900);
    const node = await withAdminAccessKey(({ nodeCreate }) =>
      nodeCreate({
        name: `Fixed ${mode} ${seed}`,
        slug: `fixed-${seed}`,
        visibility: "published",
        primary_image_asset_id: original.id,
        meta: {
          layout: {
            blocks: [
              { type: "title" },
              { type: "cover", config: { layout: mode.toLowerCase() } },
            ],
          },
        },
      }),
    );
    await page.goto(`/l/${node.slug}`);
    await edit(page);
    const area = page.getByRole("group", {
      name: "Crop area",
      exact: true,
    });
    const frame = page.locator("[data-main-image-layout]:visible");
    const stencil = area.locator(".image-cropper__selection");
    for (const viewport of [
      { width: 1280, height: 900 },
      { width: 390, height: 844 },
    ]) {
      await page.setViewportSize(viewport);
      await expect
        .poll(async () => {
          const bounds = (await area.boundingBox())!;
          const selected = (await stencil.boundingBox())!;
          return (
            Math.abs(bounds.width - selected.width) +
            Math.abs(bounds.height - selected.height)
          );
        })
        .toBeLessThan(2);
      const bounds = (await frame.boundingBox())!;
      expect(bounds.width / bounds.height).toBeCloseTo(
        mode === "Banner" ? 4 : 1,
        1,
      );
      await page
        .getByRole("button", { name: "Main image options", exact: true })
        .click();
      await expect(
        page.getByRole("menuitem", { name: "Page layout", exact: true }),
      ).toBeVisible();
      await expect(
        page.getByRole("menuitem", { name: "Crop ratio", exact: true }),
      ).toHaveCount(0);
      await page.keyboard.press("Escape");
    }
    const before = (await frame.boundingBox())!;
    const cropped = cropResponse(page);
    const x = before.x + before.width / 2;
    const y = before.y + before.height / 2;
    await page.mouse.move(x, y);
    await page.mouse.down();
    await page.mouse.move(
      x + (mode === "Artwork" ? 35 : 0),
      y + (mode === "Banner" ? 25 : 0),
      { steps: 8 },
    );
    await page.mouse.up();
    const child = await (await cropped).json();
    expect(child.parent.id).toBe(original.id);
    await expect(
      page.getByRole("button", { name: "View", exact: true }),
    ).toBeEnabled();
    const after = (await frame.boundingBox())!;
    expect(after.height).toBeCloseTo(before.height, 0);
    await view(page);
    await page.reload();
    const persisted = await saved(page, node.slug);
    const crop = persisted.meta.layout.blocks[1].config.crop;
    expect((crop.width * 1600) / (crop.height * 900)).toBeCloseTo(
      mode === "Banner" ? 4 : 1,
      3,
    );
    expect(persisted.primary_image.id).toBe(child.id);
    expect((await frame.boundingBox())!.height).toBeCloseTo(before.height, 0);
    await edit(page);
    await expect(area.locator("img").first()).toHaveAttribute(
      "src",
      `${API_ADDRESS}${original.path}`,
    );
    await expect
      .poll(async () => {
        const bounds = (await area.boundingBox())!;
        const selected = (await stencil.boundingBox())!;
        return (
          Math.abs(bounds.width - selected.width) +
          Math.abs(bounds.height - selected.height)
        );
      })
      .toBeLessThan(2);

    // Fixed frames zoom the image without resizing the viewport or changing the ratio.
    const beforeZoom = await saved(page, node.slug);
    const startingCrop = beforeZoom.meta.layout.blocks[1].config.crop;
    const zoomed = cropResponse(page);
    const zoomFrame = (await area.boundingBox())!;
    await page.mouse.move(
      zoomFrame.x + zoomFrame.width * 0.3,
      zoomFrame.y + zoomFrame.height * 0.65,
    );
    const scrollBefore = await page.evaluate(() => window.scrollY);
    await page.mouse.wheel(0, -240);
    const zoomedAsset: Asset = await (await zoomed).json();
    expect(zoomedAsset.parent?.id).toBe(original.id);
    await expect
      .poll(async () => (await saved(page, node.slug)).primary_image.id)
      .toBe(zoomedAsset.id);
    const zoomedCrop = (await saved(page, node.slug)).meta.layout.blocks[1]
      .config.crop;
    expect(zoomedCrop.width).toBeLessThan(startingCrop.width);
    // One wheel event should make a fine adjustment, including tiny trackpad events.
    expect(startingCrop.width / zoomedCrop.width).toBeLessThan(1.04);
    expect(zoomedCrop.height).toBeLessThan(startingCrop.height);
    expect(await page.evaluate(() => window.scrollY)).toBe(scrollBefore);
    expect((await area.boundingBox())!.height).toBeCloseTo(zoomFrame.height, 0);
    const zoomedPNG = await (
      await page.request.get(`${API_ADDRESS}${zoomedAsset.path}`)
    ).body();
    expect(zoomedPNG.readUInt32BE(16) / zoomedPNG.readUInt32BE(20)).toBeCloseTo(
      mode === "Banner" ? 4 : 1,
      1,
    );
    await expect(
      page.getByRole("button", { name: "View", exact: true }),
    ).toBeEnabled();

    const pinched = cropResponse(page);
    await pinchOut(page, area);
    const pinchedAsset: Asset = await (await pinched).json();
    expect(pinchedAsset.parent?.id).toBe(original.id);
    await expect
      .poll(async () => (await saved(page, node.slug)).primary_image.id)
      .toBe(pinchedAsset.id);
    const pinchedCrop = (await saved(page, node.slug)).meta.layout.blocks[1]
      .config.crop;
    expect(pinchedCrop.width).toBeLessThan(zoomedCrop.width);
    expect((pinchedCrop.width * 1600) / (pinchedCrop.height * 900)).toBeCloseTo(
      mode === "Banner" ? 4 : 1,
      2,
    );
    expect((await area.boundingBox())!.height).toBeCloseTo(zoomFrame.height, 0);
    await view(page);
    await page.reload();
    await edit(page);
    const restored = (await saved(page, node.slug)).meta.layout.blocks[1].config
      .crop;
    expect(restored).toEqual(pinchedCrop);
    // The original can still be revealed by zooming out; no restore button is required.
    const zoomedOut = cropResponse(page);
    await area.press("-");
    const zoomedOutAsset: Asset = await (await zoomedOut).json();
    await expect
      .poll(async () => (await saved(page, node.slug)).primary_image.id)
      .toBe(zoomedOutAsset.id);
    const zoomedOutCrop = (await saved(page, node.slug)).meta.layout.blocks[1]
      .config.crop;
    expect(zoomedOutCrop.width).toBeGreaterThan(pinchedCrop.width);
  });
}

test("Main image replacement keeps the whole original and removal preserves page layout", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  await createAdmin(page.context(), `replace-main-${seed}`, "TestPassword123!");
  await login(page, `replace-main-${seed}`, "TestPassword123!");
  const original = await upload(page, 1600, 900);
  const replacement = await upload(page, 900, 600);
  const bytes = await (
    await page.request.get(`${API_ADDRESS}${replacement.path}`)
  ).body();
  const node = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: `Replace ${seed}`,
      slug: `replace-main-${seed}`,
      visibility: "published",
      primary_image_asset_id: original.id,
      meta: {
        layout: {
          blocks: [
            { type: "title" },
            {
              type: "cover",
              config: {
                layout: "image",
                alt: "A changing landscape",
                crop: { left: 0, top: 0, width: 0.5, height: 0.5 },
              },
            },
            { type: "content" },
          ],
        },
      },
    }),
  );
  await page.goto(`/l/${node.slug}`);
  await edit(page);
  const received = page.waitForResponse(
    (r) =>
      r.request().method() === "POST" &&
      new URL(r.url()).searchParams.get("filename") === "replacement.png",
  );
  await page
    .locator('[data-main-image-layout] input[type="file"]')
    .setInputFiles({
      name: "replacement.png",
      mimeType: "image/png",
      buffer: bytes,
    });
  const asset: Asset = await (await received).json();
  await expect(
    page.getByRole("group", { name: "Crop area", exact: true }),
  ).toHaveAttribute("aria-disabled", "false");
  await view(page);
  let persisted = await saved(page, node.slug);
  expect(persisted.primary_image.id).toBe(asset.id);
  expect(persisted.primary_image.parent).toBeUndefined();
  expect(persisted.meta.layout.blocks[1].config).toEqual({
    layout: "image",
    alt: "A changing landscape",
  });
  // The exported PNG is square, while the source is still available at its full dimensions.
  await edit(page);
  const cropped = cropResponse(page);
  await cropRatio(page, "Square · 1:1");
  const child: Asset = await (await cropped).json();
  const png = await (
    await page.request.get(`${API_ADDRESS}${child.path}`)
  ).body();
  expect(png.readUInt32BE(16)).toBeGreaterThan(0);
  expect(png.readUInt32BE(16)).toBeLessThanOrEqual(600);
  expect(png.readUInt32BE(20)).toBe(png.readUInt32BE(16));
  expect(child.parent?.id).toBe(asset.id);
  await expect(
    page.getByRole("button", { name: "View", exact: true }),
  ).toBeEnabled();
  await page.getByRole("button", { name: "Remove image", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "Upload main image" }),
  ).toBeVisible();
  await view(page);
  await page.reload();
  persisted = await saved(page, node.slug);
  expect(persisted.primary_image).toBeUndefined();
  expect(
    persisted.meta.layout.blocks.map((b: { type: string }) => b.type),
  ).toEqual(["title", "cover", "content"]);
  expect(persisted.meta.layout.blocks[1].config).toMatchObject({
    layout: "image",
    alt: "A changing landscape",
  });
  expect((await page.request.get(`${API_ADDRESS}${asset.path}`)).ok()).toBe(
    true,
  );
});

test("Main image exports preserve extreme proportions and bound large originals", async ({
  page,
}) => {
  const seed = randomUUID().slice(0, 8);
  await createAdmin(page.context(), `export-main-${seed}`, "TestPassword123!");
  await login(page, `export-main-${seed}`, "TestPassword123!");
  for (const [width, height, expectedWidth, expectedHeight] of [
    // Stay below the test instance's upload pixel limit while exceeding each export bound.
    [3000, 400, 1920, 256],
    [400, 3000, 256, 1920],
    [1, 1000, 1, 1000],
    [1000, 1, 1000, 1],
    [1, 1, 1, 1],
  ] as const) {
    const original = await upload(page, width, height);
    const node = await withAdminAccessKey(({ nodeCreate }) =>
      nodeCreate({
        name: `Export ${width}x${height} ${seed}`,
        slug: `export-main-${width}-${height}-${seed}`,
        visibility: "published",
        primary_image_asset_id: original.id,
        meta: {
          layout: { blocks: [{ type: "cover", config: { layout: "image" } }] },
        },
      }),
    );
    await page.goto(`/l/${node.slug}`);
    await edit(page);
    const cropped = cropResponse(page);
    await cropRatio(page, "Original proportions");
    const child: Asset = await (await cropped).json();
    expect(child.parent?.id).toBe(original.id);
    const png = await (
      await page.request.get(`${API_ADDRESS}${child.path}`)
    ).body();
    expect(png.readUInt32BE(16)).toBe(expectedWidth);
    expect(png.readUInt32BE(20)).toBe(expectedHeight);
    await view(page);
    await page.reload();
    const persisted = await saved(page, node.slug);
    expect(persisted.primary_image.id).toBe(child.id);
    expect(persisted.meta.layout.blocks[0].config.crop).toEqual({
      left: 0,
      top: 0,
      width: 1,
      height: 1,
    });
  }
});
