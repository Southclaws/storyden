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
    const selectedPickerValue = (pickerName: string, value: string) =>
      page
        .getByRole("group")
        .filter({
          has: page.getByRole("combobox", { name: pickerName }),
        })
        .getByText(value, { exact: true });

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
    await page
      .getByRole("combobox", { name: "Select Tools" })
      .fill("List Categories");
    await page
      .getByRole("listbox")
      .getByRole("option", { name: "List Categories", exact: true })
      .click();
    await page.getByRole("button", { name: "Create" }).click();

    await expect(page).toHaveURL(/\/robots\/toolsets\/(?!new$)[^/]+$/);
    await expect(page.getByRole("heading", { name })).toBeVisible();
    await expect(page.getByRole("textbox", { name: "Name" })).toHaveValue(name);
    await expect(
      selectedPickerValue("Select Tools", "List Categories"),
    ).toBeVisible();

    await page.reload();
    await expect(
      selectedPickerValue("Select Tools", "List Categories"),
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
    await page.getByRole("combobox", { name: "Select Toolsets" }).fill(name);
    await page
      .getByRole("listbox")
      .getByRole("option", { name, exact: true })
      .click();
    await page.keyboard.press("Escape");

    await page.getByRole("button", { name: "Select individual tools" }).click();
    await page
      .getByRole("combobox", { name: "Select individual tools" })
      .fill("Member Search");
    await page
      .getByRole("listbox")
      .getByRole("option", { name: "Member Search", exact: true })
      .click();
    await page.getByRole("button", { name: "Create" }).click();

    await expect(page).toHaveURL(/\/robots\/(?!new$)[^/]+$/);
    await expect(page.getByRole("heading", { name: robotName })).toBeVisible();
    await expect(selectedPickerValue("Select Toolsets", name)).toBeVisible();
    await expect(
      selectedPickerValue("Select individual tools", "Member Search"),
    ).toBeVisible();

    await page.reload();
    await expect(selectedPickerValue("Select Toolsets", name)).toBeVisible();
    await expect(
      selectedPickerValue("Select individual tools", "Member Search"),
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
