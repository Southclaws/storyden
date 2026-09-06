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
  await valueItem.hover();
  await expect(valueItem).toHaveAttribute("data-highlighted", "");
  await page.keyboard.press("Enter");
}

async function clickOutsideOpenMenu(page: Page) {
  await page.locator("main").dispatchEvent("pointerdown");
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

async function dragBlockAcross(
  page: Page,
  sourceBlock: Locator,
  targetBlock: Locator,
  direction: "above" | "below",
) {
  const source = sourceBlock.getByRole("button", {
    name: "Move or configure block",
  });
  await activateDrag(page, source);

  const liveTargetBox = await targetBlock.boundingBox();
  if (!liveTargetBox) {
    throw new Error("Live block drag geometry is unavailable");
  }

  await page.mouse.move(
    liveTargetBox.x + liveTargetBox.width / 2,
    direction === "below"
      ? liveTargetBox.y + 2
      : liveTargetBox.y + liveTargetBox.height - 2,
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
    const categoryHandle = categoryBlock.getByRole("button", {
      name: "Move or configure block",
    });
    await categoryHandle.click();
    await expectMenuAnchoredToTrigger(page, categoryHandle);

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

    await categoryHandle.click();
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
    if (!categoryBeforeDrag || !threadBeforeDrag) {
      throw new Error("Block order geometry is unavailable");
    }

    const categoryStartedAbove = categoryBeforeDrag.y < threadBeforeDrag.y;

    await dragBlockAcross(
      page,
      categoryBlock,
      threadBlock,
      categoryStartedAbove ? "below" : "above",
    );

    await expect
      .poll(async () => {
        const categoryAfterDrag = await categoryBlock.boundingBox();
        const threadAfterDrag = await threadBlock.boundingBox();
        if (!categoryAfterDrag || !threadAfterDrag) {
          return false;
        }

        return categoryAfterDrag.y < threadAfterDrag.y !== categoryStartedAbove;
      })
      .toBe(true);
    await expect(
      page.getByRole("menuitem", { name: "Layout", exact: true }),
    ).toBeHidden();
  });
});
