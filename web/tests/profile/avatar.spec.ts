import { expect, test } from "@playwright/test";

import { registerUser } from "../access_key_admin_assignment";

test("opens the avatar editor on demand and saves an uploaded image", async ({
  page,
}) => {
  const handle = await registerUser(page, "avatar", "TestPassword123!");

  await page.getByRole("button", { name: "Account menu" }).click();
  await page.getByRole("menuitem", { name: "Profile", exact: true }).click();
  await expect(page).toHaveURL(`/m/${handle}`);
  await page.getByRole("button", { name: "Edit", exact: true }).click();

  const trigger = page.getByRole("button", { name: "Change avatar" });
  const dialog = page.getByRole("dialog", { name: "Edit avatar" });
  await expect(trigger).toBeVisible();
  await expect(dialog).toHaveCount(0);

  await trigger.focus();
  await trigger.press("Enter");
  await expect(dialog).toBeVisible();
  const canvas = dialog.locator("canvas");
  await expect(canvas).toBeVisible();

  await dialog.locator('input[type="file"]').setInputFiles({
    name: "avatar.svg",
    mimeType: "image/svg+xml",
    buffer: Buffer.from(
      '<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256"><rect width="256" height="256" fill="red"/></svg>',
    ),
  });
  await expect
    .poll(() =>
      canvas.evaluate((element) =>
        Array.from(
          (element as HTMLCanvasElement)
            .getContext("2d")!
            .getImageData(125, 125, 1, 1).data,
        ),
      ),
    )
    .toEqual([255, 0, 0, 255]);

  const upload = page.waitForResponse(
    (response) =>
      new URL(response.url()).pathname === "/api/accounts/self/avatar" &&
      response.request().method() === "POST",
  );
  await dialog.getByRole("button", { name: "Save", exact: true }).click();
  expect((await upload).ok()).toBe(true);
  await expect(dialog).toBeHidden();
  await expect(
    page.getByText("Avatar updated!", { exact: true }),
  ).toBeVisible();

  await trigger.click();
  await expect(dialog).toBeVisible();
  await expect(dialog.locator("canvas")).toBeVisible();
  await dialog.getByRole("button", { name: "Close", exact: true }).click();
  await expect(dialog).toBeHidden();
});
