import { type Locator, type Page, expect } from "@playwright/test";

export async function openImageTools(cropper: Locator | Page) {
  const toggle = cropper.getByRole("button", {
    name: "Image tools",
    exact: true,
  });
  await expect(toggle).toBeVisible();

  if ((await toggle.getAttribute("aria-expanded")) === "false") {
    await toggle.click();
  }

  await expect(toggle).toHaveAttribute("aria-expanded", "true");
}
