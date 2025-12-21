import { Page, expect, test } from "@playwright/test";

const PASSWORD = "TestPassword123!";

async function dismissOnboarding(page: Page) {
  const skipButton = page.getByRole("button", { name: "Skip" });
  if (await skipButton.isVisible({ timeout: 1000 }).catch(() => false)) {
    await skipButton.click();
  }
}

async function registerUser(page: Page, username: string) {
  await page.goto("/register");
  await page.getByRole("textbox", { name: "username" }).fill(username);
  await page.getByRole("textbox", { name: "password" }).fill(PASSWORD);
  await page.getByRole("button", { name: "Register" }).click();
  await expect(page).toHaveURL("/", { timeout: 10000 });
  await expect(
    page.getByRole("button", { name: "Account menu" }),
  ).toBeVisible();
  await dismissOnboarding(page);
}

async function logout(page: Page) {
  await page.getByRole("button", { name: "Account menu" }).click();
  await page.getByRole("link", { name: "Logout" }).click();
  await expect(page.getByRole("link", { name: "Login" })).toBeVisible();
}

async function login(page: Page, username: string) {
  await page.goto("/login");
  await page.getByRole("textbox", { name: "username" }).fill(username);
  await page.getByRole("textbox", { name: "password" }).fill(PASSWORD);
  await page.getByRole("button", { name: "Login" }).click();
  await expect(page).toHaveURL("/", { timeout: 10000 });
  await expect(
    page.getByRole("button", { name: "Account menu" }),
  ).toBeVisible();
  await dismissOnboarding(page);
}

async function createThread(
  page: Page,
  title: string,
  body: string,
): Promise<string> {
  await page.getByRole("link", { name: "Post" }).click();
  await expect(page).toHaveURL("/new", { timeout: 5000 });

  await page.locator("#title-input").fill(title);

  const editor = page.locator(".ProseMirror").first();
  await editor.click();
  await editor.fill(body);

  await page.getByRole("button", { name: "Post" }).click();

  await expect(page).toHaveURL(/\/t\//, { timeout: 10000 });
  await dismissOnboarding(page);

  const url = page.url();
  return url;
}

async function postReply(page: Page, body: string) {
  const replyForm = page.locator("form").filter({ hasText: "Reply to" });
  const replyEditor = replyForm.locator(".ProseMirror[contenteditable='true']");
  await replyEditor.click();
  await replyEditor.fill(body);

  await replyForm.getByRole("button", { name: "Post" }).click();

  await page
    .locator("li")
    .filter({ hasText: body })
    .waitFor({ timeout: 10000 });
}

async function navigateToThread(page: Page, threadUrl: string) {
  await page.goto(threadUrl);
  await dismissOnboarding(page);
}

test.describe("Thread Creation", () => {
  test("should create a thread with title and body", async ({ page }) => {
    await registerUser(page, "thread-creator-01");

    const threadUrl = await createThread(
      page,
      "My First Thread",
      "This is the body of my first thread.",
    );

    await expect(
      page.getByRole("heading", { name: "My First Thread" }),
    ).toBeVisible();
    await expect(
      page.locator("main").getByText("This is the body of my first thread."),
    ).toBeVisible();
    expect(threadUrl).toMatch(/\/t\//);
  });
});

test.describe("Thread Replies", () => {
  test("should reply to a thread from a different account", async ({
    page,
  }) => {
    await registerUser(page, "thread-author-01");
    const threadUrl = await createThread(
      page,
      "Thread for Replies",
      "This thread will receive replies.",
    );

    await logout(page);
    await registerUser(page, "replier-01");

    await navigateToThread(page, threadUrl);
    await postReply(page, "This is a reply from another user.");

    await expect(
      page
        .locator("li")
        .filter({ hasText: "This is a reply from another user." }),
    ).toBeVisible();
  });

  test("should reply to a specific reply", async ({ page }) => {
    await registerUser(page, "thread-author-02");
    const threadUrl = await createThread(
      page,
      "Thread for Reply Threading",
      "This thread tests reply-to-reply functionality.",
    );

    await postReply(page, "First reply to the thread.");

    await logout(page);
    await registerUser(page, "replier-02");

    await navigateToThread(page, threadUrl);

    const firstReply = page
      .locator("li")
      .filter({ hasText: "First reply to the thread." });
    await firstReply.getByRole("button", { name: "Reply to this" }).click();

    await expect(page.getByText("Replying to")).toBeVisible();

    const replyBox = page
      .locator("form")
      .filter({ has: page.getByRole("button", { name: "Post" }) })
      .last();
    const replyEditor = replyBox.locator(
      ".ProseMirror[contenteditable='true']",
    );
    await replyEditor.click();
    await replyEditor.fill("This is a reply to the first reply.");

    await replyBox.getByRole("button", { name: "Post" }).click();

    const newReply = page
      .locator("li")
      .filter({ hasText: "This is a reply to the first reply." });
    await expect(newReply).toBeVisible({ timeout: 10000 });
    await expect(
      newReply.locator("a").filter({ hasText: "First reply to the thread." }),
    ).toBeVisible();
  });

  test("should clear reply-to selection", async ({ page }) => {
    await registerUser(page, "thread-author-03");
    const threadUrl = await createThread(
      page,
      "Thread for Clear Reply-To",
      "Testing the clear reply-to button.",
    );

    await postReply(page, "A reply to set as target.");

    await navigateToThread(page, threadUrl);

    const reply = page
      .locator("li")
      .filter({ hasText: "A reply to set as target." });
    await reply.getByRole("button", { name: "Reply to this" }).click();

    await expect(page.getByText("Replying to")).toBeVisible();

    await page.getByRole("button", { name: "Clear reply-to" }).click();

    await expect(page.getByText("Replying to")).not.toBeVisible();
  });
});

test.describe("Thread Notifications", () => {
  test("should show notification when someone replies to your thread", async ({
    page,
  }) => {
    await registerUser(page, "notification-author-01");
    const threadUrl = await createThread(
      page,
      "Thread for Notification Test",
      "Waiting for someone to reply.",
    );

    await logout(page);
    await registerUser(page, "notification-replier-01");

    await navigateToThread(page, threadUrl);
    await postReply(page, "A reply that should trigger a notification.");

    await logout(page);
    await login(page, "notification-author-01");

    await expect(
      page.getByRole("button", { name: "Notifications" }),
    ).toBeVisible();
    await page.getByRole("button", { name: "Notifications" }).click();

    await expect(page.getByText("notification-replier-01")).toBeVisible({
      timeout: 10000,
    });
    await expect(page.getByText("replied to")).toBeVisible();
  });
});

test.describe("Thread Reactions", () => {
  test("should add a reaction to a reply", async ({ page }) => {
    await registerUser(page, "reaction-author-01");
    const threadUrl = await createThread(
      page,
      "Thread for Reaction Test",
      "This thread will test reactions.",
    );

    await postReply(page, "A reply to react to.");

    const reply = page
      .locator("li")
      .filter({ hasText: "A reply to react to." });
    const reactionButton = reply.getByRole("button", { name: "Add reaction" });
    await expect(reactionButton).toBeVisible({ timeout: 10000 });
    await reactionButton.click();

    const emojiPicker = page.locator(".EmojiPickerReact");
    await emojiPicker.waitFor({ timeout: 5000 });
    await emojiPicker.locator("button.epr-emoji").first().click();

    await expect(
      reply.locator("button").filter({ hasText: /^\p{Emoji}.*\d$/u }),
    ).toBeVisible({ timeout: 5000 });
  });
});

test.describe("Thread Editing", () => {
  test("should edit thread title and body", async ({ page }) => {
    await registerUser(page, "edit-author-01");
    const threadUrl = await createThread(
      page,
      "Original Thread Title",
      "Original thread body content.",
    );

    await page.goto(threadUrl + "?edit=true");
    await dismissOnboarding(page);

    const titleInput = page
      .locator("main span[contenteditable='true']")
      .first();
    await expect(titleInput).toBeVisible({ timeout: 5000 });

    await titleInput.click();
    await page.keyboard.press("Meta+a");
    await page.keyboard.type("Updated Thread Title");

    const editor = page
      .locator("main .ProseMirror[contenteditable='true']")
      .first();
    await editor.click();
    await page.keyboard.press("Meta+a");
    await editor.fill("Updated thread body content.");

    await page.getByRole("button", { name: "Save" }).click();

    await expect(
      page.getByRole("heading", { name: "Updated Thread Title" }),
    ).toBeVisible({ timeout: 10000 });
    await expect(
      page.locator("main").getByText("Updated thread body content."),
    ).toBeVisible();
  });

  test("should discard thread edits", async ({ page }) => {
    await registerUser(page, "edit-author-02");
    const threadUrl = await createThread(
      page,
      "Thread to Not Edit",
      "Original content that should remain.",
    );

    await page.goto(threadUrl + "?edit=true");
    await dismissOnboarding(page);

    const titleInput = page
      .locator("main span[contenteditable='true']")
      .first();
    await titleInput.click();
    await page.keyboard.press("Meta+a");
    await page.keyboard.type("This change should be discarded");

    await page.getByRole("button", { name: "Discard" }).click();

    await expect(
      page.getByRole("heading", { name: "Thread to Not Edit" }),
    ).toBeVisible();
    await expect(
      page.locator("main").getByText("Original content that should remain."),
    ).toBeVisible();
  });
});

test.describe("Reply Editing", () => {
  test("should edit a reply", async ({ page }) => {
    await registerUser(page, "reply-edit-author-01");
    await createThread(
      page,
      "Thread for Reply Editing",
      "This thread tests reply editing.",
    );

    await postReply(page, "Original reply content.");

    const reply = page
      .locator("li")
      .filter({ hasText: "Original reply content." });
    const moreButton = reply.getByRole("button", { name: "More options" });
    await expect(moreButton).toBeVisible({ timeout: 10000 });
    await moreButton.click();

    const editMenuItem = page.getByRole("menuitem", { name: /Edit/ });
    await expect(editMenuItem).toBeVisible({ timeout: 5000 });
    await editMenuItem.click();

    const saveButton = page.getByRole("button", { name: "Save" });
    await expect(saveButton).toBeVisible({ timeout: 10000 });

    const replyEditor = page.locator("li .ProseMirror[contenteditable='true']");
    await expect(replyEditor).toBeVisible({ timeout: 5000 });
    await replyEditor.click();
    await page.keyboard.press("Meta+a");
    await replyEditor.fill("Updated reply content.");

    await saveButton.click();

    await expect(
      page.locator("li").filter({ hasText: "Updated reply content." }),
    ).toBeVisible({ timeout: 10000 });
  });

  test("should discard reply edits", async ({ page }) => {
    await registerUser(page, "reply-edit-author-02");
    await createThread(
      page,
      "Thread for Reply Discard",
      "Testing reply edit discard.",
    );

    await postReply(page, "Reply that should not change.");

    const reply = page
      .locator("li")
      .filter({ hasText: "Reply that should not change." });
    const moreButton = reply.getByRole("button", { name: "More options" });
    await expect(moreButton).toBeVisible({ timeout: 10000 });
    await moreButton.click();

    const editMenuItem = page.getByRole("menuitem", { name: /Edit/ });
    await expect(editMenuItem).toBeVisible({ timeout: 5000 });
    await editMenuItem.click();

    const discardButton = page.getByRole("button", { name: "Discard" });
    await expect(discardButton).toBeVisible({ timeout: 10000 });

    const replyEditor = page.locator("li .ProseMirror[contenteditable='true']");
    await expect(replyEditor).toBeVisible({ timeout: 5000 });
    await replyEditor.click();
    await page.keyboard.press("Meta+a");
    await replyEditor.fill("This should be discarded.");

    await discardButton.click();

    await expect(
      page.locator("li").filter({ hasText: "Reply that should not change." }),
    ).toBeVisible();
  });
});

test.describe("Thread Pagination", () => {
  test("should paginate replies when exceeding page size", async ({ page }) => {
    test.setTimeout(300000);

    await registerUser(page, "pagination-author-01");
    const threadUrl = await createThread(
      page,
      "Thread for Pagination Test",
      "This thread will have many replies.",
    );

    for (let i = 1; i <= 53; i++) {
      const replyForm = page
        .locator("form")
        .filter({ has: page.getByRole("button", { name: "Post" }) })
        .last();
      const replyEditor = replyForm.locator(
        ".ProseMirror[contenteditable='true']",
      );
      await replyEditor.click();
      await replyEditor.fill(`Reply number ${i} for pagination test.`);
      await replyForm.getByRole("button", { name: "Post" }).click();
      await page.waitForTimeout(200);
    }

    await navigateToThread(page, threadUrl);

    await expect(
      page
        .locator("li")
        .filter({ hasText: "Reply number 1 for pagination test." }),
    ).toBeVisible();
    await expect(
      page
        .locator("li")
        .filter({ hasText: "Reply number 50 for pagination test." }),
    ).toBeVisible();
    await expect(
      page
        .locator("li")
        .filter({ hasText: "Reply number 51 for pagination test." }),
    ).not.toBeVisible();

    await expect(page.getByRole("link", { name: "2" }).first()).toBeVisible();

    await page.getByRole("link", { name: "2" }).first().click();

    await expect(page).toHaveURL(/page=2/, { timeout: 5000 });

    await expect(
      page
        .locator("li")
        .filter({ hasText: "Reply number 51 for pagination test." }),
    ).toBeVisible();
    await expect(
      page
        .locator("li")
        .filter({ hasText: "Reply number 52 for pagination test." }),
    ).toBeVisible();
    await expect(
      page
        .locator("li")
        .filter({ hasText: "Reply number 53 for pagination test." }),
    ).toBeVisible();
    await expect(
      page
        .locator("li")
        .filter({ hasText: "Reply number 1 for pagination test." }),
    ).not.toBeVisible();
  });
});
