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

async function dismissOnboarding(page: Page) {
  const skipButton = page.getByRole("button", { name: "Skip" });
  if (await skipButton.isVisible({ timeout: 1000 }).catch(() => false)) {
    await skipButton.click();
  }
}

async function openSidebar(page: Page) {
  const sidebarToggle = page.getByRole("button", {
    name: /navigation sidebar/i,
  });
  await expect(sidebarToggle).toBeVisible();

  if ((await sidebarToggle.getAttribute("aria-expanded")) !== "true") {
    await sidebarToggle.click();
  }

  await expect(sidebarToggle).toHaveAttribute("aria-expanded", "true");
  await expect(page.locator("#navigation__leftbar")).toBeVisible();
}

async function closeSidebar(page: Page) {
  const sidebarToggle = page.getByRole("button", {
    name: /navigation sidebar/i,
  });
  await expect(sidebarToggle).toBeVisible();

  if ((await sidebarToggle.getAttribute("aria-expanded")) !== "false") {
    await sidebarToggle.click();
  }

  await expect(sidebarToggle).toHaveAttribute("aria-expanded", "false");
  await expect(page.locator("#navigation__leftbar")).not.toBeVisible();
}

async function openFeedEditorFromSidebar(page: Page) {
  await openSidebar(page);

  const editButton = page.getByRole("button", { name: "Open feed editor" });
  await expect(editButton).toBeVisible();
  await editButton.click();

  await expect(page).toHaveURL(/editing=feed/);
  await expect(page.getByText("Editing Home")).toBeVisible();
}

async function chooseSelectOption(
  page: Page,
  select: Locator,
  optionValue: string,
) {
  await expect(select).toBeVisible();
  const contentId = await select.getAttribute("aria-controls");
  if (!contentId) {
    throw new Error("select trigger did not expose aria-controls");
  }
  const content = page.locator(`[id="${contentId}"]`);
  await select.focus();
  await select.press("ArrowDown");
  if ((await content.getAttribute("data-state")) !== "open") {
    await select.click({ force: true });
  }
  await expect(content).toHaveAttribute("data-state", "open");

  const option = content.locator(
    `[data-scope='select'][data-part='item'][data-value='${optionValue}']`,
  );
  await expect(option).toBeVisible();
  await option.click();
  await expect(content).toHaveAttribute("data-state", "closed");
}

test.describe("AdminZone Feed Settings", () => {
  test("switches source to categories and back to threads with dependent UI updates", async ({
    page,
  }) => {
    const seed = Date.now().toString();
    const adminHandle = `admin_zone_${seed}`;
    const { categoryThreadTitle, uncategorisedThreadTitle } =
      await createCategoryAndThreads(seed);

    await createAdmin(page.context(), adminHandle, PASSWORD);
    await login(page, adminHandle, PASSWORD);
    await dismissOnboarding(page);
    await page.goto("/");
    await openFeedEditorFromSidebar(page);

    const sourceSelect = page.getByRole("combobox", { name: "Source" });
    await expect(sourceSelect).toContainText("Threads");

    await chooseSelectOption(page, sourceSelect, "categories");
    await expect(sourceSelect).toContainText("Categories");

    const layoutSelect = page.getByRole("combobox", { name: "Layout" });
    await expect(layoutSelect).toBeVisible();
    await expect(
      page.getByText("Thread list display", { exact: true }),
    ).toBeVisible();
    await expect(
      page.getByRole("heading", { name: "Discussion categories" }),
    ).toBeVisible();

    await chooseSelectOption(page, layoutSelect, "grid");
    await expect(layoutSelect).toContainText("Grid");

    const listDisplaySelect = page
      .locator("[data-scope='select'][data-part='trigger'][role='combobox']")
      .filter({ hasText: /Uncategorised only|All threads|None/ })
      .first();

    await chooseSelectOption(page, listDisplaySelect, "all");
    await expect(listDisplaySelect).toContainText("All threads");

    await expect(
      page.getByRole("link", { name: categoryThreadTitle, exact: true }),
    ).toBeVisible();
    await expect(
      page.getByRole("link", { name: uncategorisedThreadTitle, exact: true }),
    ).toBeVisible();

    await chooseSelectOption(page, sourceSelect, "threads");

    await expect(page.getByRole("combobox", { name: "Layout" })).toHaveCount(0);
    await expect(
      page.getByText("Thread list display", { exact: true }),
    ).toHaveCount(0);

    const listDisplayComboboxes = page
      .locator("[data-scope='select'][data-part='trigger'][role='combobox']")
      .filter({ hasText: /Uncategorised only|All threads|None/ });
    await expect(listDisplayComboboxes).toHaveCount(0);
    await expect(sourceSelect).toContainText("Threads");

    await closeSidebar(page);
  });
});
