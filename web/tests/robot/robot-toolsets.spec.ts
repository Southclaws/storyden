import { expect, test } from "@playwright/test";

import { goToNewChat, setupRobotProviderWithScript } from "./helpers";

test.describe("Robot Toolsets", () => {
  test.beforeAll(async () => {
    await setupRobotProviderWithScript();
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
    await page.goto("/robots/toolsets");
  });

  test("creates, edits, and deletes a reusable Toolset", async ({ page }) => {
    const suffix = Date.now().toString(36);
    const name = `Research tools ${suffix}`;
    const robotName = `Research robot ${suffix}`;
    const updatedDescription = "Reusable tools for focused research tasks.";

    await expect(page.getByRole("heading", { name: "Toolsets" })).toBeVisible();
    const newToolsetLink = page
      .getByRole("main")
      .getByRole("link", { name: "New Toolset" });
    await expect(newToolsetLink).toBeVisible();

    await newToolsetLink.click();
    await expect(page).toHaveURL("/robots/toolsets/new");

    await page.getByRole("textbox", { name: "Name" }).fill(name);
    await page
      .getByRole("textbox", { name: "Description" })
      .fill("Research helpers");
    await page
      .getByRole("textbox", { name: "Instructions" })
      .fill("Verify sources before reporting findings.");
    await page.getByRole("button", { name: "Select Tools" }).click();
    const toolsMenu = page.getByRole("menu", { name: "Select Tools" });
    await toolsMenu
      .getByRole("textbox", { name: "Search for items" })
      .fill("List Categories");
    await page.getByRole("menuitem", { name: "List Categories" }).click();
    await page.getByRole("button", { name: "Create" }).click();

    await expect(page).toHaveURL(/\/robots\/toolsets\/(?!new$)[^/]+$/);
    await expect(page.getByRole("heading", { name })).toBeVisible();
    await expect(page.getByRole("textbox", { name: "Name" })).toHaveValue(name);
    await expect(
      page
        .getByRole("button", { name: "Select Tools" })
        .getByText("List Categories"),
    ).toBeVisible();

    await page.reload();
    await expect(
      page
        .getByRole("button", { name: "Select Tools" })
        .getByText("List Categories"),
    ).toBeVisible();

    await page
      .getByRole("textbox", { name: "Description" })
      .fill(updatedDescription);
    await page.getByRole("button", { name: "Save" }).click();
    await expect(page.getByText("Toolset saved")).toBeVisible();

    const toolsetURL = page.url();
    await page.goto("/robots/new");
    await page.getByRole("textbox", { name: "Name" }).fill(robotName);
    await page
      .getByRole("textbox", { name: "Description" })
      .fill("Researches the community using shared and focused tools.");
    await page
      .getByRole("textbox", { name: "Playbook" })
      .fill("Use the shared Toolset first, then search members when needed.");

    await page.getByRole("button", { name: "Select Toolsets" }).click();
    const toolsetMenu = page.getByRole("menu", { name: "Select Toolsets" });
    await toolsetMenu
      .getByRole("textbox", { name: "Search for items" })
      .fill(name);
    await toolsetMenu.getByRole("menuitem", { name }).click();

    await page.getByRole("button", { name: "Select individual tools" }).click();
    const individualToolMenu = page.getByRole("menu", {
      name: "Select individual tools",
    });
    await individualToolMenu
      .getByRole("textbox", { name: "Search for items" })
      .fill("Member Search");
    await individualToolMenu
      .getByRole("menuitem", { name: "Member Search" })
      .click();
    await page.getByRole("button", { name: "Create" }).click();

    await expect(page).toHaveURL(/\/robots\/(?!new$)[^/]+$/);
    await expect(page.getByRole("heading", { name: robotName })).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Select Toolsets" }).getByText(name),
    ).toBeVisible();
    await expect(
      page
        .getByRole("button", { name: "Select individual tools" })
        .getByText("Member Search"),
    ).toBeVisible();

    await page.reload();
    await expect(
      page.getByRole("button", { name: "Select Toolsets" }).getByText(name),
    ).toBeVisible();
    await expect(
      page
        .getByRole("button", { name: "Select individual tools" })
        .getByText("Member Search"),
    ).toBeVisible();

    await page.getByRole("button", { name: "Delete robot" }).click();
    await page
      .getByRole("dialog")
      .getByRole("button", { name: "Delete robot" })
      .click();
    await expect(page).toHaveURL("/robots");

    await page.goto(toolsetURL);

    await page.getByRole("button", { name: "Delete Toolset" }).click();
    await expect(
      page.getByRole("heading", { name: "Delete Toolset" }),
    ).toBeVisible();
    await page
      .getByRole("dialog")
      .getByRole("button", { name: "Delete Toolset" })
      .click();

    await expect(page).toHaveURL("/robots/toolsets");
    await expect(page.getByText(name)).toHaveCount(0);
  });
});
