import { Locator, Page, expect, test } from "@playwright/test";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

const PASSWORD = "TestPassword123!";

async function createCategoryAndThreads(
  seed: string,
): Promise<{ categoryThreadTitle: string; uncategorisedThreadTitle: string }> {
  const categoryThreadTitle = `Category thread ${seed}`;
  const uncategorisedThreadTitle = `Uncategorised thread ${seed}`;

  await withAdminAccessKey(async ({ categoryCreate, threadCreate }) => {
    const category = await categoryCreate({
      colour: "#3b82f6",
      description: `E2E category ${seed}`,
      name: `E2E Category ${seed}`,
      slug: `e2e-category-${seed}`,
    });

    await threadCreate({
      title: categoryThreadTitle,
      body: `Thread in category ${seed}`,
      category: category.id,
      visibility: "published",
    });

    await threadCreate({
      title: uncategorisedThreadTitle,
      body: `Thread without category ${seed}`,
      visibility: "published",
    });
  });

  return { categoryThreadTitle, uncategorisedThreadTitle };
}

async function openFeedEditor(page: Page) {
  await expect(page.locator("#navigation__leftbar")).toBeVisible();

  const editButton = page
    .getByRole("toolbar", { name: "Primary navigation" })
    .getByRole("button", { name: "Edit site" });
  await expect(editButton).toBeVisible();
  await editButton.click();

  await expect(page).toHaveURL(/editing=site/);
  await expect(
    page
      .getByRole("toolbar", { name: "Primary navigation" })
      .getByRole("button", { name: "Exit edit mode" }),
  ).toHaveAttribute("aria-pressed", "true");
}

async function chooseBlockMenuItem(page: Page, item: string, value: string) {
  const menuItem = page.getByRole("menuitem", { name: item, exact: true });
  await expect(menuItem).toBeVisible();
  await menuItem.click();

  const valueItem = page.getByRole("menuitem", { name: value, exact: true });
  await expect(valueItem).toBeVisible();
  await valueItem.click();
}

async function clickOutsideOpenMenu(page: Page) {
  await page.locator("main").dispatchEvent("pointerdown");
}

async function dragBlockBelow(
  page: Page,
  sourceBlock: Locator,
  targetBlock: Locator,
) {
  const source = sourceBlock.getByRole("button", {
    name: "Move or configure block",
  });
  const sourceBox = await source.boundingBox();
  const targetBox = await targetBlock.boundingBox();

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

  const liveTargetBox = await targetBlock.boundingBox();
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

test.describe("Feed Editor Settings", () => {
  test("configures category layout and thread source from block menus", async ({
    page,
  }) => {
    const seed = Date.now().toString();
    const adminHandle = `admin_zone_${seed}`;
    const { categoryThreadTitle, uncategorisedThreadTitle } =
      await createCategoryAndThreads(seed);

    await createAdmin(page.context(), adminHandle, PASSWORD);
    await login(page, adminHandle, PASSWORD);
    await page.goto("/");
    await openFeedEditor(page);

    const categoryBlock = page.locator(
      '.block-editor__root[data-block-type="categories"]:visible',
    );
    await categoryBlock
      .getByRole("button", { name: "Move or configure block" })
      .click();

    const categoryMenuLabel = page
      .locator('[data-scope="menu"][data-part="item-group-label"]')
      .getByText("Discussion categories", { exact: true });
    await categoryMenuLabel.click({ force: true });
    await expect(
      page.getByRole("menuitem", { name: "Layout", exact: true }),
    ).toBeVisible();

    await clickOutsideOpenMenu(page);
    await expect(
      page.getByRole("menuitem", { name: "Layout", exact: true }),
    ).toBeHidden();

    await categoryBlock
      .getByRole("button", { name: "Move or configure block" })
      .click();
    await chooseBlockMenuItem(page, "Layout", "Grid");

    const threadBlock = page.locator(
      '.block-editor__root[data-block-type="threads"]:visible',
    );
    await threadBlock
      .getByRole("button", { name: "Move or configure block" })
      .click();
    await chooseBlockMenuItem(page, "Source", "All threads");

    await expect(
      page.getByRole("link", { name: categoryThreadTitle, exact: true }),
    ).toBeVisible();
    await expect(
      page.getByRole("link", { name: uncategorisedThreadTitle, exact: true }),
    ).toBeVisible();

    const categoryBeforeDrag = await categoryBlock.boundingBox();
    const threadBeforeDrag = await threadBlock.boundingBox();
    expect(categoryBeforeDrag?.y).toBeLessThan(threadBeforeDrag?.y ?? 0);

    await dragBlockBelow(page, categoryBlock, threadBlock);

    await expect
      .poll(async () => {
        const categoryAfterDrag = await categoryBlock.boundingBox();
        const threadAfterDrag = await threadBlock.boundingBox();
        return (categoryAfterDrag?.y ?? 0) > (threadAfterDrag?.y ?? 0);
      })
      .toBe(true);
    await expect(
      page.getByRole("menuitem", { name: "Layout", exact: true }),
    ).toBeHidden();
  });
});
