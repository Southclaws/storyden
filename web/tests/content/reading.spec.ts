import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  createAdmin,
  login,
  registerUser,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

const CONTENT =
  '<h2>Field notes</h2><p><strong>Bold note</strong> and <a href="https://example.com/guide">a guide</a>.</p><ol><li><p>Next step</p></li></ol><img src="https://example.com/photo.png" alt="Field photograph">';

test("includes published link-card destinations in server-rendered HTML", async ({
  browser,
}) => {
  const slug = `reading-link-${randomUUID().slice(0, 8)}`;
  await withAdminAccessKey(async ({ nodeCreate }) => {
    await nodeCreate({
      name: "Published link preview",
      slug,
      content:
        '<div data-display="card" data-href="https://example.com/card"></div>',
      visibility: "published",
      meta: { layout: { blocks: [{ type: "title" }, { type: "content" }] } },
    });
  });
  const context = await browser.newContext({ javaScriptEnabled: false });
  try {
    const page = await context.newPage();
    await page.goto(`http://localhost:3001/l/${slug}`);
    const link = page.locator('a[href="https://example.com/card"]');
    await expect(link).toHaveText("https://example.com/card");
    await expect(link).toHaveAttribute("href", "https://example.com/card");
  } finally {
    await context.close();
  }
});

for (const mode of ["richtext", "markdown"] as const) {
  test(`includes stored HTML in the server response and uses the ${mode} editor only when editing`, async ({
    browser,
    page,
  }) => {
    const seed = randomUUID().slice(0, 8);
    const slug = `reading-${seed}`;
    await withAdminAccessKey(async ({ nodeCreate }) => {
      await nodeCreate({
        name: `Reading fixture ${seed}`,
        slug,
        content: CONTENT,
        visibility: "published",
        meta: { layout: { blocks: [{ type: "title" }, { type: "content" }] } },
      });
    });

    const plainContext = await browser.newContext({ javaScriptEnabled: false });
    try {
      const plain = await plainContext.newPage();
      const response = await plain.goto(`http://localhost:3001/l/${slug}`);
      expect(await response!.text()).toContain("<h2>Field notes</h2>");
      // Loading boundaries stream this HTML in a hidden segment until JavaScript
      // reveals it. Assert server output here; hydrated reading is checked below.
      await expect(
        plain.locator("h1").filter({ hasText: `Reading fixture ${seed}` }),
      ).toHaveCount(1);
      await expect(
        plain.getByRole("heading", {
          name: "Field notes",
          exact: true,
          includeHidden: true,
        }),
      ).toHaveText("Field notes");
      await expect(plain.locator("strong")).toHaveText("Bold note");
      await expect(
        plain.locator('a[href="https://example.com/guide"]'),
      ).toHaveAttribute("href", "https://example.com/guide");
      await expect(plain.locator("ol li")).toHaveText("Next step");
      await expect(
        plain.locator('img[alt="Field photograph"]'),
      ).toHaveAttribute("src", "https://example.com/photo.png");
      await expect(plain.locator(".content-composer__editor")).toHaveCount(0);
    } finally {
      await plainContext.close();
    }

    const handle = `reader_${seed}`;
    const password = "TestPassword123!";
    await createAdmin(page.context(), handle, password);
    await login(page, handle, password);
    const preference = await page.request.patch(
      "http://localhost:8001/api/accounts",
      {
        data: { meta: { editor: { mode } } },
      },
    );
    expect(preference.ok()).toBe(true);

    await page.goto(`/l/${slug}`);
    const content = page.locator(
      'main [data-outline-block="content"]:visible, main #block-content_content:visible',
    );
    await expect(
      content.getByRole("heading", { name: "Field notes" }),
    ).toBeVisible();
    await expect(content.locator("strong")).toHaveText("Bold note");
    await expect(
      content.getByRole("img", { name: "Field photograph" }),
    ).toBeAttached();
    await expect(content.locator(".content-composer__editor")).toHaveCount(0);

    await page.getByRole("button", { name: "Edit", exact: true }).click();
    await page
      .getByRole("menuitem", { name: "Quick edit", exact: true })
      .click();
    if (mode === "richtext") {
      const editor = content.locator(".content-composer__editor");
      await expect(editor).toBeVisible();
      await expect(editor.locator("strong")).toHaveText("Bold note");
      await expect(
        editor.getByRole("img", { name: "Field photograph" }),
      ).toBeAttached();
      await expect(
        content.locator("textarea.content-composer__editor"),
      ).toHaveCount(0);
      await editor.locator("p").first().click();
      await page.keyboard.press("End");
      await page.keyboard.type(" Saved paragraph.");
    } else {
      const editor = content.locator("textarea.content-composer__editor");
      await expect(editor).toBeVisible();
      await expect(editor).toHaveValue(/Field notes/);
      await expect(editor).toHaveValue(/https:\/\/example.com\/photo.png/);
      await expect(
        content.locator(
          ".content-composer--richtext .content-composer__editor",
        ),
      ).toHaveCount(0);
      await editor.fill(`${await editor.inputValue()}\n\nSaved paragraph.`);
    }

    await expect
      .poll(async () => {
        const response = await page.request.get(
          `http://localhost:8001/api/nodes/${slug}`,
        );
        return (await response.json()).content;
      })
      .toContain("Saved paragraph.");
    await page.getByRole("button", { name: "View", exact: true }).click();
    await expect(
      content.getByRole("heading", { name: "Field notes" }),
    ).toBeVisible();
    await expect(content).toContainText("Saved paragraph.");
    await expect(content.locator(".content-composer__editor")).toHaveCount(0);
    await page.reload();
    await expect(content).toContainText("Saved paragraph.");
    await expect(
      content.getByRole("img", { name: "Field photograph" }),
    ).toBeAttached();
  });
}

test("server-renders a published thread and its replies for anonymous readers", async ({
  browser,
  page,
}) => {
  const thread = await withAdminAccessKey(
    async ({ threadCreate, replyCreate }) => {
      const thread = await threadCreate({
        title: `Reading thread ${randomUUID().slice(0, 8)}`,
        body: CONTENT,
        visibility: "published",
      });
      await replyCreate(thread.slug, {
        body: "<h3>Reader reply</h3><p>A server-rendered reply.</p>",
      });
      return thread;
    },
  );
  const plainContext = await browser.newContext({ javaScriptEnabled: false });
  try {
    const plain = await plainContext.newPage();
    await plain.goto(`/t/${thread.slug}`);
    await expect(plain.locator(".thread-page h2")).toHaveText("Field notes");
    await expect(plain.locator(".thread-reply h3")).toHaveText("Reader reply");
  } finally {
    await plainContext.close();
  }
  await page.goto(`/t/${thread.slug}`);
  await expect(
    page.getByRole("heading", { name: "Field notes" }),
  ).toBeVisible();
  await expect(
    page.getByRole("heading", { name: "Reader reply" }),
  ).toBeVisible();
  await expect(page.locator(".content-composer__editor")).toHaveCount(0);
});

test("server-renders profile biographies with the same rich content", async ({
  browser,
  page,
}) => {
  const handle = await registerUser(
    page,
    "reading-profile",
    "TestPassword123!",
  );
  const updated = await page.request.patch(
    "http://localhost:8001/api/accounts",
    { data: { bio: CONTENT } },
  );
  expect(updated.ok()).toBe(true);
  const plainContext = await browser.newContext({ javaScriptEnabled: false });
  try {
    const plain = await plainContext.newPage();
    await plain.goto(`/m/${handle}`);
    await expect(plain.locator(".content-composer__content h2")).toHaveText(
      "Field notes",
    );
    await expect(plain.locator(".content-composer__content strong")).toHaveText(
      "Bold note",
    );
  } finally {
    await plainContext.close();
  }
  await page.goto(`/m/${handle}`);
  await expect(
    page.getByRole("heading", { name: "Field notes" }),
  ).toBeVisible();
  await expect(
    page.getByRole("img", { name: "Field photograph" }),
  ).toBeAttached();
  await expect(page.locator(".content-composer__editor")).toHaveCount(0);
});
