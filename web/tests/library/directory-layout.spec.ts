import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import { API_ADDRESS } from "../../src/config";
import {
  createAdmin,
  getAdminAccessKey,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("Directory Gallery keeps mixed image dimensions bounded, preserves source order and renders before hydration", async ({
  page,
  browser,
}) => {
  const seed = randomUUID().slice(0, 8);
  await page.goto("/");
  const assets: string[] = [];
  for (const [width, height] of [
    [120, 360],
    [480, 120],
    [1, 1000],
    [1000, 1],
    [1, 1],
  ]) {
    const png = await page.evaluate(
      ({ width, height }) => {
        const canvas = document.createElement("canvas");
        canvas.width = width!;
        canvas.height = height!;
        canvas.getContext("2d")!.fillRect(0, 0, width!, height!);
        return canvas.toDataURL("image/png").split(",")[1]!;
      },
      { width, height },
    );
    const upload = await page.request.post(
      `${API_ADDRESS}/api/assets?filename=packing-${width}-${seed}.png`,
      {
        headers: {
          Authorization: `Bearer ${getAdminAccessKey()}`,
          "Content-Type": "application/octet-stream",
        },
        data: Buffer.from(png, "base64"),
      },
    );
    expect(upload.ok()).toBe(true);
    assets.push((await upload.json()).id);
  }
  const parent = await withAdminAccessKey(async ({ nodeCreate }) => {
    const parent = await nodeCreate({
      name: `Packing experiment ${seed}`,
      slug: `packing-${seed}`,
      hide_child_tree: true,
      visibility: "published",
      content:
        "<h2>Reference notes</h2>" +
        "<p>Browse these entries in their original sequence.</p>".repeat(30),
      meta: {
        layout: {
          blocks: [
            { type: "title" },
            { type: "content" },
            {
              type: "directory",
              config: {
                layout: "gallery",
                columns: [
                  { fid: "fixed:name", hidden: false },
                  { fid: "fixed:link", hidden: true },
                  { fid: "fixed:description", hidden: true },
                  { fid: "fixed:primary_image", hidden: false },
                ],
              },
            },
          ],
        },
      },
    });
    for (let index = 0; index < 8; index++) {
      await nodeCreate({
        name: `Reference ${index + 1} ${seed}`,
        slug: `packing-${index}-${seed}`,
        parent: parent.id,
        visibility: "published",
        primary_image_asset_id:
          index < 6 ? assets[index % assets.length] : undefined,
      });
    }
    return parent;
  });
  await page.setViewportSize({ width: 1440, height: 1000 });
  await page.route("**/api/nodes/*/children*", async (route) => {
    const response = await route.fetch();
    const body = await response.json();
    for (const node of body.nodes) {
      if (node.primary_image) {
        node.primary_image.width = 0;
        node.primary_image.height = 0;
      }
    }
    await route.fulfill({ response, json: body });
  });
  await page.goto(`/l/${parent.slug}`);
  const updates: string[] = [];
  page.on("request", (request) => {
    if (request.method() === "PATCH") updates.push(request.url());
  });
  const grid = page.getByRole("list", { name: "Library pages" });
  await expect(grid).toHaveAttribute("data-layout", "gallery");
  await expect(grid.locator(":scope > li")).toHaveCount(8);
  for (const image of await grid.locator("img").all()) {
    await image.scrollIntoViewIfNeeded();
    await expect
      .poll(() =>
        image.evaluate((image) => (image as HTMLImageElement).naturalWidth),
      )
      .toBeGreaterThan(0);
  }
  await expect(grid.locator("img[width], img[height]")).toHaveCount(0);
  await expect
    .poll(() =>
      grid.locator("img").evaluateAll((images) =>
        images.map((image) => {
          const box = image.getBoundingClientRect();
          return {
            bounded:
              box.height >= 31 &&
              box.height <= Math.min(window.innerHeight * 0.7, 640) + 1,
            fit: getComputedStyle(image).objectFit,
          };
        }),
      ),
    )
    .toEqual(
      Array.from({ length: 6 }, () => ({ bounded: true, fit: "cover" })),
    );
  const originalNames = await grid.locator(":scope > li").allTextContents();
  await expect(grid.locator("img[width], img[height]")).toHaveCount(0);
  await expect(grid.locator("img")).toHaveCount(6);
  expect(await grid.locator(":scope > li").allTextContents()).toEqual(
    originalNames,
  );
  await expect
    .poll(() =>
      grid
        .locator("img")
        .evaluateAll((images) =>
          images.every(
            (image) =>
              image.getBoundingClientRect().height <=
              Math.min(window.innerHeight * 0.7, 640) + 1,
          ),
        ),
    )
    .toBe(true);
  // The real background job may still be processing an upload. Wait for the
  // API's processed metadata, then reload the source data before comparing it.
  await expect
    .poll(
      async () => {
        const response = await page.request.get(
          `${API_ADDRESS}/api/nodes/${parent.slug}/children`,
        );
        const result = await response.json();
        const images = result.nodes.flatMap(
          (node: { primary_image?: { width: number; height: number } }) =>
            node.primary_image ? [node.primary_image] : [],
        );
        return (
          images.length === 6 &&
          images.every(
            (image: { width: number; height: number }) =>
              image.width > 0 && image.height > 0,
          )
        );
      },
      { timeout: 30_000 },
    )
    .toBe(true);
  await page.unroute("**/api/nodes/*/children*");
  await page.reload();
  await expect(grid.locator("img[width][height]")).toHaveCount(6);
  await expect(grid.locator("img[srcset]")).toHaveCount(6);
  await expect
    .poll(() =>
      grid
        .locator("img")
        .evaluateAll((images) =>
          images.map((image) => getComputedStyle(image).objectFit),
        ),
    )
    .toEqual(["cover", "cover", "contain", "cover", "cover", "cover"]);
  for (const width of [1440, 900, 390, 1440]) {
    await page.setViewportSize({ width, height: 1000 });
    if (width === 390) {
      await expect
        .poll(() =>
          grid.evaluate((grid) => {
            const cards = [...grid.children];
            return new Set(
              cards.map((card) => Math.round(card.getBoundingClientRect().x)),
            ).size;
          }),
        )
        .toBe(1);
    }
    await expect
      .poll(() =>
        grid.locator(":scope > li").evaluateAll((items) => {
          const rects = items.map((item) => item.getBoundingClientRect());
          return rects.every((first, a) =>
            rects
              .slice(a + 1)
              .every(
                (second) =>
                  first.right <= second.x + 1 ||
                  second.right <= first.x + 1 ||
                  first.bottom <= second.y + 1 ||
                  second.bottom <= first.y + 1,
              ),
          );
        }),
      )
      .toBe(true);
    expect(
      await page.evaluate(
        () =>
          document.documentElement.scrollWidth <=
          document.documentElement.clientWidth,
      ),
    ).toBe(true);
  }
  await page.setViewportSize({ width: 1440, height: 500 });
  const outlineRoot = page.locator(".content-outline__root:visible");
  await expect(outlineRoot).toHaveAttribute("data-visible", "true");
  const outlineFrames = await outlineRoot.evaluate(async (root) => {
    const states: (string | null)[] = [];
    for (let frame = 0; frame < 60; frame++) {
      await new Promise(requestAnimationFrame);
      states.push(root.getAttribute("data-visible"));
    }
    return [...new Set(states)];
  });
  expect(outlineFrames).toEqual(["true"]);
  const handle = `directory_outline_${seed}`;
  await createAdmin(page.context(), handle, "TestPassword123!");
  await login(page, handle, "TestPassword123!");
  await page.goto(`/l/${parent.slug}`);
  const outlineTrigger = page.getByRole("button", {
    name: "On this page",
    exact: true,
  });
  for (const width of [1440, 1280, 1100, 390]) {
    await page.setViewportSize({ width, height: 500 });
    await page.locator("#navigation__scroll").evaluate(async (shell) => {
      await Promise.all(
        shell.getAnimations().map((animation) => animation.finished),
      );
    });
    const view = await grid.evaluate((list) => ({
      width: list.getBoundingClientRect().width,
      columns: new Set(
        [...list.children].map((card) =>
          Math.round(card.getBoundingClientRect().x),
        ),
      ).size,
    }));
    if (width === 1440) {
      await expect(outlineTrigger).toBeVisible();
      const triggerBox = await outlineTrigger.boundingBox();
      const pageBox = await outlineRoot.boundingBox();
      expect(triggerBox!.x).toBeGreaterThanOrEqual(pageBox!.x + pageBox!.width);
      expect(triggerBox!.x + triggerBox!.width).toBeLessThanOrEqual(width);
      const center = triggerBox!.y + triggerBox!.height / 2;
      expect(center).toBeCloseTo(250, 0);
      await page.mouse.wheel(0, 300);
      await expect
        .poll(async () => {
          const box = await outlineTrigger.boundingBox();
          return Math.round(box!.y + box!.height / 2);
        })
        .toBe(250);
      await outlineTrigger.click();
      await expect(
        page.getByRole("navigation", { name: "On this page" }),
      ).toBeVisible();
      await page.keyboard.press("Escape");
    } else {
      await expect(outlineTrigger).toBeHidden();
    }
    await page.getByRole("button", { name: "Edit", exact: true }).click();
    await page
      .getByRole("menuitem", { name: "Quick edit", exact: true })
      .click();
    await expect(
      page.getByRole("button", { name: "View", exact: true }),
    ).toBeVisible();
    await expect(outlineTrigger).toBeHidden();
    await expect
      .poll(() =>
        grid.evaluate((list) => ({
          width: list.getBoundingClientRect().width,
          columns: new Set(
            [...list.children].map((card) =>
              Math.round(card.getBoundingClientRect().x),
            ),
          ).size,
        })),
      )
      .toEqual(view);
    await page.getByRole("button", { name: "View", exact: true }).click();
    await page.waitForURL((url) => !url.searchParams.has("edit"));
  }
  expect(await grid.locator(":scope > li").allTextContents()).toEqual(
    originalNames,
  );
  expect(updates).toEqual([]);

  const ssrContext = await browser.newContext({
    viewport: { width: 1440, height: 1000 },
  });
  try {
    // Next's inline stream bootstrap inserts the server HTML. Block application
    // scripts so this verifies the CSS layout before React can hydrate it.
    await ssrContext.route("**/_next/static/**", (route) =>
      route.request().resourceType() === "script"
        ? route.abort()
        : route.continue(),
    );
    const ssrPage = await ssrContext.newPage();
    await ssrPage.goto(`/l/${parent.slug}`);
    const ssrGrid = ssrPage.getByRole("list", { name: "Library pages" });
    await expect(ssrGrid).toHaveAttribute("data-layout", "gallery");
    await expect(ssrGrid.locator(":scope > li")).toHaveCount(8);
    await expect(ssrGrid).toBeVisible();
    await expect(ssrGrid).not.toHaveAttribute("data-masonry");
    for (const image of await ssrGrid.locator("img").all()) {
      await image.scrollIntoViewIfNeeded();
      await expect(image).toBeVisible();
      await expect(image).not.toHaveAttribute("width", "0");
      await expect(image).not.toHaveAttribute("height", "0");
      await expect
        .poll(() =>
          image.evaluate((image) => (image as HTMLImageElement).naturalWidth),
        )
        .toBeGreaterThan(0);
      const box = await image.boundingBox();
      expect(box!.height).toBeGreaterThanOrEqual(31);
      expect(box!.height).toBeLessThanOrEqual(641);
    }
  } finally {
    await ssrContext.close();
  }

  await page.setViewportSize({ width: 1440, height: 1000 });
  await page.getByRole("button", { name: "Edit", exact: true }).click();
  await page.getByRole("menuitem", { name: "Quick edit", exact: true }).click();
  await page.locator('[id="block-directory_gutter-drag-handle"]').click();
  await page.getByRole("menuitem", { name: "Layout", exact: true }).hover();
  await page.getByRole("menuitem", { name: "Grid", exact: true }).click();
  await page.getByRole("button", { name: "View", exact: true }).click();
  await page.waitForURL((url) => !url.searchParams.has("edit"));
  await expect(grid).toHaveAttribute("data-layout", "grid");
  await expect
    .poll(() =>
      grid.locator("img").evaluateAll((images) =>
        images.every((image) => {
          const box = image.getBoundingClientRect();
          return (
            getComputedStyle(image).objectFit === "cover" &&
            Math.abs(box.width / box.height - 4 / 3) < 0.01
          );
        }),
      ),
    )
    .toBe(true);
});
