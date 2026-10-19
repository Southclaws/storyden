import { Locator, Page, expect, test } from "@playwright/test";

import { nodeDelete } from "../../src/api/openapi-client/nodes";
import type { Node } from "../../src/api/openapi-schema";
import {
  createAdmin,
  getAdminAccessKey,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

const PASSWORD = "TestPassword123!";

test("library navigation does not speculatively fetch linked routes", async ({
  page,
}) => {
  const slug = `library-navigation-${Date.now()}`;
  const prefetches: string[] = [];

  await withAdminAccessKey(async ({ nodeCreate }) => {
    const node = await nodeCreate({
      name: "Library navigation prefetch check",
      slug,
      content: "A public page with navigation links.",
      visibility: "published",
    });

    try {
      page.on("request", (request) => {
        if (request.headers()["next-router-prefetch"] === "1") {
          prefetches.push(request.url());
        }
      });

      await page.goto(`/l/${slug}`);
      await expect(
        page.getByRole("tree", { name: "Library pages" }),
      ).toBeVisible();
      await page.locator('[data-navigation-anchor="home"]:visible').hover();
      await page.waitForLoadState("networkidle");
      expect(prefetches).toEqual([]);

      await page.locator('[data-navigation-anchor="home"]:visible').click();
      await page.waitForURL("/");
      await expect(
        page.locator('[data-navigation-anchor="home"]:visible'),
      ).toBeVisible();
    } finally {
      await nodeDelete(node.id, undefined, {
        headers: { Authorization: `Bearer ${getAdminAccessKey()}` },
      });
    }
  });
});

test("library cards prefetch only on intent and stop after warming the page", async ({
  page,
}) => {
  const seed = Date.now();

  await withAdminAccessKey(async ({ nodeCreate }) => {
    const nodes: Node[] = [];

    try {
      for (const name of ["Hover", "Keyboard"]) {
        nodes.push(
          await nodeCreate({
            name: `${name} prefetch ${seed}`,
            slug: `${name.toLowerCase()}-prefetch-${seed}`,
            content: `${name} page content`,
            visibility: "published",
          }),
        );
      }

      const prefetches: string[] = [];

      function waitForContentPrefetch(slug: string) {
        return page.waitForResponse((response) => {
          const request = response.request();
          return (
            new URL(response.url()).pathname === `/l/${slug}` &&
            request.headers()["next-router-prefetch"] === "1"
          );
        });
      }

      page.on("request", (request) => {
        const path = new URL(request.url()).pathname;
        if (
          request.headers()["next-router-prefetch"] === "1" &&
          nodes.some((node) => path === `/l/${node.slug}`)
        ) {
          prefetches.push(path);
        }
      });

      await page.goto("/l");
      await page.waitForLoadState("networkidle");
      expect(prefetches).toEqual([]);

      const hoverCard = page.locator(`article[id="${nodes[0]!.id}"]`);
      const keyboardCard = page.locator(`article[id="${nodes[1]!.id}"]`);
      const hoverLink = hoverCard.getByRole("link", {
        name: nodes[0]!.name,
      });

      const hoverPrefetch = waitForContentPrefetch(nodes[0]!.slug);
      await hoverLink.hover();
      await hoverPrefetch;
      await page.waitForLoadState("networkidle");
      expect(new Set(prefetches)).toEqual(new Set([`/l/${nodes[0]!.slug}`]));

      await hoverCard.getByRole("heading").getByRole("link").focus();
      await hoverLink.hover();
      await page.waitForLoadState("networkidle");
      expect(new Set(prefetches)).toEqual(new Set([`/l/${nodes[0]!.slug}`]));

      const keyboardPrefetch = waitForContentPrefetch(nodes[1]!.slug);
      await keyboardCard.getByRole("heading").getByRole("link").focus();
      await keyboardPrefetch;
      await page.waitForLoadState("networkidle");

      const keyboardLink = keyboardCard.getByRole("link", {
        name: nodes[1]!.name,
        exact: true,
      });
      await expect(keyboardLink).toHaveCount(1);
      await expect(keyboardLink).toBeFocused();
      await expect(
        keyboardCard.getByText("Keyboard page content"),
      ).toBeVisible();
      expect(
        await keyboardCard
          .getByText("Keyboard page content")
          .evaluate((element) => element.closest("a")),
      ).toBeNull();
      await page.waitForLoadState("networkidle");
      expect(new Set(prefetches)).toEqual(
        new Set(nodes.map((node) => `/l/${node.slug}`)),
      );

      await hoverLink.click();
      await page.waitForURL(`/l/${nodes[0]!.slug}`);
      await expect(
        page.getByText("Hover page content", { exact: true }),
      ).toBeVisible();
    } finally {
      for (const node of nodes) {
        await nodeDelete(node.id, undefined, {
          headers: { Authorization: `Bearer ${getAdminAccessKey()}` },
        });
      }
    }
  });
});

async function clickOutsideOpenMenu(page: Page) {
  await page.locator("#block-content_content:visible").click();
}

async function expectMenuAnchoredToTrigger(page: Page, trigger: Locator) {
  const menuID = await trigger.getAttribute("aria-controls");
  if (!menuID) {
    throw new Error("Menu trigger does not control a menu");
  }

  const menu = page.locator(`[id="${menuID}"]`);
  await expect(menu).toBeVisible();

  const triggerBox = await trigger.boundingBox();
  const menuBox = await menu.boundingBox();
  if (!triggerBox || !menuBox) {
    throw new Error("Menu positioning geometry is unavailable");
  }

  expect(
    Math.abs(menuBox.x - (triggerBox.x + triggerBox.width)),
  ).toBeLessThanOrEqual(2);
  expect(Math.abs(menuBox.y - triggerBox.y)).toBeLessThanOrEqual(2);
}

async function activateDrag(page: Page, source: Locator) {
  for (let attempt = 0; attempt < 3; attempt++) {
    await source.scrollIntoViewIfNeeded();
    const sourceBox = await source.boundingBox();
    if (!sourceBox) {
      throw new Error("Block drag geometry is unavailable");
    }

    const sourceX = sourceBox.x + sourceBox.width / 2;
    const sourceY = sourceBox.y + sourceBox.height / 2;
    await page.mouse.move(sourceX, sourceY);
    await page.mouse.down();
    await page.mouse.move(sourceX + 8, sourceY, { steps: 4 });

    try {
      await expect(source).toHaveAttribute("data-dragging", "", {
        timeout: 1000,
      });
      return;
    } catch (error) {
      await page.mouse.up();
      await page.keyboard.press("Escape");

      if (attempt === 2) {
        throw error;
      }
    }
  }
}

async function dragBlockBelow(
  page: Page,
  sourceType: string,
  targetType: string,
) {
  const source = page.locator(
    `#block-${sourceType}_gutter-drag-handle:visible`,
  );
  const target = page.locator(`#block-${targetType}_container:visible`);
  await activateDrag(page, source);

  const liveTargetBox = await target.boundingBox();
  if (!liveTargetBox) {
    throw new Error("Live block drag geometry is unavailable");
  }

  await page.mouse.move(
    liveTargetBox.x + liveTargetBox.width / 2,
    liveTargetBox.y + 2,
    { steps: 12 },
  );
  await page.mouse.up();
}

test.describe("Library page block editor", () => {
  test("keeps menu chrome inert, dismisses outside, configures, adds, and reorders blocks", async ({
    page,
  }) => {
    const seed = Date.now().toString();
    const adminHandle = `library_editor_${seed}`;
    const pageName = `Library block editor ${seed}`;
    const pageSlug = `library-block-editor-${seed}`;

    await withAdminAccessKey(async ({ nodeCreate }) => {
      await nodeCreate({
        name: pageName,
        slug: pageSlug,
        content: `Library editor content ${seed}`,
        visibility: "published",
        meta: {
          layout: {
            blocks: [
              { type: "title" },
              {
                type: "directory",
                config: { layout: "table", columns: [] },
              },
              { type: "content" },
            ],
          },
        },
      });
    });

    await createAdmin(page.context(), adminHandle, PASSWORD);
    await login(page, adminHandle, PASSWORD);
    await page.goto(`/l/${pageSlug}?edit=true&editMode=direct`);

    const directoryHandle = page.locator(
      "#block-directory_gutter-drag-handle:visible",
    );
    await expect(directoryHandle).toBeVisible();
    await directoryHandle.click();
    await expectMenuAnchoredToTrigger(page, directoryHandle);

    const directoryMenuLabel = page
      .locator('[data-scope="menu"][data-part="item-group-label"]')
      .getByText("Directory", { exact: true });
    await directoryMenuLabel.dispatchEvent("pointerdown");
    await expect(
      page.getByRole("menuitem", { name: "Layout", exact: true }),
    ).toBeVisible();

    await clickOutsideOpenMenu(page);
    await expect(
      page.getByRole("menuitem", { name: "Layout", exact: true }),
    ).toBeHidden();

    await directoryHandle.click();
    await page.getByRole("menuitem", { name: "Layout", exact: true }).click();
    await page.getByRole("menuitem", { name: "Grid", exact: true }).click();
    await expect(
      page.getByRole("menuitem", { name: "Layout", exact: true }),
    ).toBeHidden();

    await directoryHandle.click();
    await page
      .getByRole("menuitem", { name: "Add Block", exact: true })
      .click();
    await page
      .getByRole("menuitem", { name: "Cover image", exact: true })
      .click();

    await expect(page.locator("#block-cover_container:visible")).toBeVisible();
    await expect
      .poll(() =>
        page
          .locator(".block-editor__root[data-block-type]")
          .evaluateAll((blocks) =>
            blocks.map((block) => block.getAttribute("data-block-type")),
          ),
      )
      .toEqual(["title", "directory", "cover", "content"]);
    await expect(
      page.getByRole("menuitem", { name: "Add Block", exact: true }),
    ).toBeHidden();

    const titleBlock = page.locator("#block-title_container:visible");
    const contentBlock = page.locator("#block-content_container:visible");
    const titleBeforeDrag = await titleBlock.boundingBox();
    const contentBeforeDrag = await contentBlock.boundingBox();
    expect(titleBeforeDrag?.y).toBeLessThan(contentBeforeDrag?.y ?? 0);

    await dragBlockBelow(page, "title", "content");

    await expect
      .poll(async () => {
        const titleAfterDrag = await titleBlock.boundingBox();
        const contentAfterDrag = await contentBlock.boundingBox();
        return (titleAfterDrag?.y ?? 0) > (contentAfterDrag?.y ?? 0);
      })
      .toBe(true);
    await expect(
      page.getByRole("menuitem", { name: "Update URL slug", exact: true }),
    ).toBeHidden();
  });
});
