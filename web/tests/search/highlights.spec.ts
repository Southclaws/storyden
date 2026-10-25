import { type Page, expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import { nodeDelete } from "../../src/api/openapi-client/nodes";
import {
  createAdmin,
  getAdminAccessKey,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

async function highlightedText(page: Page) {
  return page.evaluate(() =>
    [...(CSS.highlights.get("storyden-search-match") ?? [])].map((range) =>
      range.toString(),
    ),
  );
}

test("submitted search highlights survive pagination, selection, and clearing", async ({
  page,
}, testInfo) => {
  const seed = randomUUID().slice(0, 8);
  const word = `garden${seed}`;
  const query = `${word} tools`;
  const ids: string[] = [];

  try {
    await withAdminAccessKey(async ({ nodeCreate }) => {
      for (let index = 0; index < 26; index++) {
        const node = await nodeCreate({
          name: `${query} ${index}`,
          description: `A ${query} notebook. ${"A quiet day beside the harbour. ".repeat(30)}`,
          visibility: "published",
        });
        ids.push(node.id);
      }
    });

    await page.goto("/search");
    const input = page.getByRole("searchbox");
    await input.fill(query);
    await expect.poll(() => highlightedText(page)).toEqual([]);
    await input.press("Enter");
    await expect(
      page.getByRole("status", { name: "Search results" }),
    ).toHaveText("25 search results on this page.");
    await expect
      .poll(async () => (await highlightedText(page)).length)
      .toBe(100);

    await page.getByRole("link", { name: "2", exact: true }).click();
    await expect(
      page.getByRole("status", { name: "Search results" }),
    ).toHaveText("1 search result on this page.");
    await expect
      .poll(() => highlightedText(page))
      .toEqual([word, "tools", word, "tools"]);
    await expect(page.locator("main mark")).toHaveCount(0);

    const card = page.locator("main article");
    await page.screenshot({ path: testInfo.outputPath("search-clipped.png") });
    const selected = await card.getByRole("heading").evaluate((element) => {
      const selection = window.getSelection();
      selection?.selectAllChildren(element);
      const text = selection?.toString();
      selection?.removeAllRanges();
      return text;
    });
    expect(selected).toContain(query);

    for (const colorScheme of ["light", "dark"] as const) {
      await page.emulateMedia({ colorScheme });
      await page.screenshot({
        path: testInfo.outputPath(`search-${colorScheme}.png`),
      });
      const colors = await card
        .locator("[data-search-highlight]")
        .first()
        .evaluate((element) => {
          const highlight = getComputedStyle(
            element,
            "::highlight(storyden-search-match)",
          );
          return {
            background: highlight.backgroundColor,
            color: highlight.color,
          };
        });
      expect(colors.background).not.toBe("rgba(0, 0, 0, 0)");
      expect(colors.background).not.toBe(colors.color);
    }

    await page.locator('button[type="reset"]').click();
    await expect.poll(() => highlightedText(page)).toEqual([]);
  } finally {
    for (const id of ids) {
      await nodeDelete(id, undefined, {
        headers: { Authorization: `Bearer ${getAdminAccessKey()}` },
      });
    }
  }
});

for (const layout of ["table", "grid", "grid-placeholder"] as const) {
  test(`directory ${layout} highlights names and stops during direct editing`, async ({
    page,
  }) => {
    const seed = randomUUID().slice(0, 8);
    const ids: string[] = [];
    const password = "TestPassword123!";

    try {
      const parent = await withAdminAccessKey(async ({ nodeCreate }) => {
        const parent = await nodeCreate({
          name: `Directory ${seed}`,
          visibility: "published",
          meta: {
            layout: {
              blocks: [
                { type: "title" },
                {
                  type: "directory",
                  config: {
                    layout: layout === "table" ? "table" : "grid",
                    columns: [
                      {
                        fid: "fixed:name",
                        hidden: layout === "grid-placeholder",
                      },
                      { fid: "fixed:description", hidden: false },
                      { fid: "fixed:link", hidden: true },
                      {
                        fid: "fixed:primary_image",
                        hidden: layout !== "grid-placeholder",
                      },
                    ],
                  },
                },
              ],
            },
          },
        });
        ids.push(parent.id);
        for (const name of ["Garden tools", "Harbour notes"]) {
          const child = await nodeCreate({
            name,
            description: "Garden tools also appear in this description.",
            parent: parent.slug,
            visibility: "published",
          });
          ids.push(child.id);
        }
        return parent;
      });

      await createAdmin(page.context(), `highlights-${seed}`, password);
      await login(page, `highlights-${seed}`, password);
      await page.goto(`/l/${parent.slug}`);
      const directory = page.locator(
        '[data-outline-block="directory"]:visible',
      );
      const input = directory.getByPlaceholder("Search...");
      await expect(input).toHaveValue("");
      // Keep an unfiltered entry so back/forward changes the query on this page.
      await page.evaluate(() => {
        window.history.pushState(null, "", window.location.href);
      });
      await input.fill("garden tools");
      await expect(input).toHaveValue("garden tools");
      await expect(directory.getByRole("status")).toHaveText(
        "1 page in this directory.",
      );
      await expect
        .poll(() => highlightedText(page))
        .toEqual(["Garden", "tools"]);

      await page.goBack();
      await expect(input).toHaveValue("");
      await expect.poll(() => highlightedText(page)).toEqual([]);
      await expect(
        directory.getByText("Harbour notes", { exact: true }),
      ).toBeVisible();

      await page.goForward();
      await expect(input).toHaveValue("garden tools");
      await expect
        .poll(() => highlightedText(page))
        .toEqual(["Garden", "tools"]);

      // Quick edit replaces names with inputs; neither their values nor the
      // image placeholder should retain highlights while editing.
      await page.getByRole("button", { name: "Edit", exact: true }).click();
      await page
        .getByRole("menuitem", { name: "Quick edit", exact: true })
        .click();
      await expect(
        page.getByRole("button", { name: "View", exact: true }),
      ).toBeVisible();
      await expect.poll(() => highlightedText(page)).toEqual([]);
      await page.getByRole("button", { name: "View", exact: true }).click();
      await expect
        .poll(() => highlightedText(page))
        .toEqual(["Garden", "tools"]);

      await input.fill("");
      await expect.poll(() => highlightedText(page)).toEqual([]);
      await expect(
        directory.getByText("Harbour notes", { exact: true }),
      ).toBeVisible();
    } finally {
      for (const id of ids.reverse()) {
        await nodeDelete(id, undefined, {
          headers: { Authorization: `Bearer ${getAdminAccessKey()}` },
        });
      }
    }
  });
}

test("a body-only search match needs no visible highlight", async ({
  page,
}) => {
  const query = `hidden${randomUUID().slice(0, 8)}`;
  const node = await withAdminAccessKey(({ nodeCreate }) =>
    nodeCreate({
      name: "An unrelated title",
      description: "An unrelated preview",
      content: query,
      visibility: "published",
    }),
  );
  try {
    await page.goto(`/search?q=${query}&kind=node`);
    await expect(
      page.getByRole("heading", { name: "An unrelated title" }),
    ).toBeVisible();
    await expect(
      page.getByRole("status", { name: "Search results" }),
    ).toHaveText("1 search result on this page.");
    await expect.poll(() => highlightedText(page)).toEqual([]);
  } finally {
    await nodeDelete(node.id, undefined, {
      headers: { Authorization: `Bearer ${getAdminAccessKey()}` },
    });
  }
});
