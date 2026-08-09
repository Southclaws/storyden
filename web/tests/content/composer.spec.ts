import { Page, expect, test } from "@playwright/test";

const PASSWORD = "TestPassword123!";
const IMAGE_BYTES = Buffer.from(
  "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQIHWP4z8DwHwAFgAI/ScL3WQAAAABJRU5ErkJggg==",
  "base64",
);

function uniqueUsername(prefix: string) {
  return `${prefix}-${Date.now().toString(36)}${Math.random()
    .toString(36)
    .slice(2, 8)}`;
}

async function openComposer(page: Page) {
  await page.goto("/register");
  await page
    .getByRole("textbox", { name: "username" })
    .fill(uniqueUsername("composer"));
  await page.getByRole("textbox", { name: "password" }).fill(PASSWORD);
  await page.getByRole("button", { name: "Register" }).click();
  await expect(page).toHaveURL("/");

  await page.goto("/new");
  const editor = page.locator(".ProseMirror[contenteditable='true']");
  await expect(editor).toBeVisible();

  return editor;
}

async function pasteText(page: Page, selector: string, text: string) {
  await page.locator(selector).evaluate((element, value) => {
    const clipboard = new DataTransfer();
    clipboard.setData("text/plain", value);
    element.dispatchEvent(
      new ClipboardEvent("paste", {
        bubbles: true,
        cancelable: true,
        clipboardData: clipboard,
      }),
    );
  }, text);
}

test.describe("Content composer", () => {
  test("formats selected text through the selection menu", async ({ page }) => {
    const editor = await openComposer(page);

    await editor.fill("Format this text");
    await editor.press(process.platform === "darwin" ? "Meta+A" : "Control+A");

    const boldButton = page.locator('button[title="Toggle bold text"]:visible');
    await expect(boldButton).toBeVisible();
    await boldButton.click();

    await expect(editor.locator("strong")).toHaveText("Format this text");
  });

  test("turns a pasted URL into a link through the link paste menu", async ({
    page,
  }) => {
    const editor = await openComposer(page);
    const url = "https://example.com/storyden";

    await editor.click();
    await pasteText(page, ".ProseMirror[contenteditable='true']", url);

    await expect(page.getByText("Show link as")).toBeVisible();
    await page.getByRole("button", { name: "Text", exact: true }).click();

    await expect(editor.locator(`a[href="${url}"]`)).toHaveText(url);
    await expect(page.getByText("Show link as")).toBeHidden();
  });

  test("renders successful image uploads through the custom image node view", async ({
    page,
  }) => {
    const editor = await openComposer(page);

    await page.route("**/api/assets*", async (route) => {
      if (route.request().method() !== "POST") {
        await route.fallback();
        return;
      }

      await route.fulfill({
        contentType: "application/json",
        body: JSON.stringify({
          id: "01HZZZZZZZZZZZZZZZZZZZZZZZ",
          filename: "composer-image.png",
          height: 1,
          mime_type: "image/png",
          path: "/api/assets/composer-image.png",
          width: 1,
        }),
      });
    });

    await page.locator('input[type="file"]').setInputFiles({
      name: "composer-image.png",
      mimeType: "image/png",
      buffer: IMAGE_BYTES,
    });

    const image = editor.locator('img[alt="composer-image.png"]');
    await expect(image).toHaveAttribute("src", /composer-image\.png$/);
    await expect(image.locator("xpath=..")).toHaveAttribute(
      "data-node-view-wrapper",
    );
  });

  test("offers removal when an image upload fails", async ({ page }) => {
    const editor = await openComposer(page);

    await page.route("**/api/assets*", async (route) => {
      if (route.request().method() !== "POST") {
        await route.fallback();
        return;
      }

      await route.fulfill({
        status: 500,
        contentType: "application/json",
        body: JSON.stringify({ message: "Upload failed" }),
      });
    });

    await page.locator('input[type="file"]').setInputFiles({
      name: "failed-image.png",
      mimeType: "image/png",
      buffer: IMAGE_BYTES,
    });

    await expect(page.getByText("Upload failed")).toBeVisible();
    await page.getByRole("button", { name: "Remove" }).click();
    await expect(editor.locator('img[alt="failed-image.png"]')).toHaveCount(0);
  });
});
