import { Page, expect, test } from "@playwright/test";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

const PASSWORD = "TestPassword123!";

async function clickOutsideOpenMenu(page: Page) {
  await page.locator("#block-content_content:visible").click();
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
  const sourceBox = await source.boundingBox();
  const targetBox = await target.boundingBox();

  if (!sourceBox || !targetBox) {
    throw new Error("Block drag geometry is unavailable");
  }

  await page.mouse.move(
    sourceBox.x + sourceBox.width / 2,
    sourceBox.y + sourceBox.height / 2,
  );
  await page.mouse.down();
  await page.mouse.move(
    sourceBox.x + sourceBox.width / 2,
    sourceBox.y + sourceBox.height / 2 + 4,
    { steps: 2 },
  );
  await expect(source).toHaveAttribute("data-dragging", "");

  const liveTargetBox = await target.boundingBox();
  if (!liveTargetBox) {
    throw new Error("Live block drag geometry is unavailable");
  }

  await page.mouse.move(
    liveTargetBox.x + liveTargetBox.width / 2,
    liveTargetBox.y + liveTargetBox.height - 2,
    { steps: 12 },
  );
  await page.mouse.up();
}

test.describe("Library page block editor", () => {
  test("keeps menu chrome inert, dismisses outside, configures, and reorders blocks", async ({
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
