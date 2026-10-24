import AxeBuilder from "@axe-core/playwright";
import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test.use({ viewport: { width: 1440, height: 1000 } });

for (const audience of ["guest", "administrator"] as const) {
  test(`desktop sidebar trees expose named rows and their controls to ${audience}s`, async ({
    page,
  }) => {
    const seed = randomUUID().slice(0, 8);
    const fixtures = await withAdminAccessKey(
      async ({ categoryCreate, nodeCreate }) => {
        const category = await categoryCreate({
          name: `Category branch ${seed}`,
          colour: "#2255aa",
          description: "Accessible tree fixture",
        });
        const categoryChild = await categoryCreate({
          name: `Category leaf ${seed}`,
          parent: category.id,
          colour: "#2255aa",
          description: "Accessible tree fixture",
        });
        const node = await nodeCreate({
          name: `Library branch ${seed}`,
          visibility: "published",
        });
        const nodeChild = await nodeCreate({
          name: `Library leaf ${seed}`,
          parent: node.slug,
          visibility: "published",
        });
        return [
          {
            tree: "Discussion categories",
            branch: category.name,
            leaf: categoryChild.name,
            options: `Category options for ${categoryChild.name}`,
          },
          {
            tree: "Library pages",
            branch: node.name,
            leaf: nodeChild.name,
            options: `Page options for ${nodeChild.name}`,
          },
        ];
      },
    );

    if (audience === "administrator") {
      const admin = `tree-aria-${seed}`;
      await createAdmin(page.context(), admin, "TestPassword123!");
      await login(page, admin, "TestPassword123!");
    }
    await page.goto("/");

    for (const fixture of fixtures) {
      const tree = page.getByRole("tree", { name: fixture.tree });
      await expect(tree).toBeVisible();
      const toggle = tree.getByRole("button", {
        name: `Toggle ${fixture.branch}`,
      });
      if ((await toggle.getAttribute("aria-expanded")) !== "true") {
        await toggle.click();
      }
      const item = tree.getByRole("treeitem", {
        name: fixture.leaf,
        exact: true,
      });
      const link = item.getByRole("link", {
        name: fixture.leaf,
        exact: true,
      });
      await link.focus();
      await expect(link).toBeFocused();
      await expect(item).toHaveAttribute("aria-level", "2");
      await expect(item).not.toHaveAttribute("aria-expanded");
      if (audience === "administrator" || fixture.tree === "Library pages") {
        await expect(
          item.getByRole("button", { name: fixture.options }),
        ).toBeVisible();
      }
      const results = await new AxeBuilder({ page })
        .include('[role="tree"]')
        .withRules([
          "aria-required-children",
          "aria-required-parent",
          "aria-allowed-attr",
          "aria-valid-attr-value",
          "nested-interactive",
        ])
        .analyze();
      expect(results.violations).toEqual([]);
    }
  });
}
