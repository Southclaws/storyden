import { expect, test } from "@playwright/test";

import { withAdminAccessKey } from "../access_key_admin_assignment";

import {
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  switchToRobot,
} from "./helpers";

test.beforeAll(async () => {
  await setupRobotProviderWithScript();
});

test.describe("Robot Chat — content_search tool", () => {
  let robotName: string;

  test.beforeAll(async () => {
    robotName = `e2e-content-search-robot-${Date.now()}`;
    await withAdminAccessKey(
      async ({ robotCreate, nodeCreate, categoryCreate, threadCreate }) => {
        await nodeCreate({
          name: `Magnolia Library Page ${Date.now()}`,
          visibility: "published",
        });

        const category = await categoryCreate({
          colour: "#3b82f6",
          description: "search test category",
          name: `Search Category ${Date.now()}`,
          slug: `search-category-${Date.now()}`,
        });

        await threadCreate({
          title: "Magnolia Forum Thread",
          body: "<p>discussion about magnolia</p>",
          category: category.id,
          visibility: "published",
        });

        await robotCreate({
          name: robotName,
          description: "E2E robot for content_search tool",
          playbook: "you are a test robot",
          model: "mock/../robot/scripts/robot-chat-content-search.yaml",
        });
      },
    );
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
    await switchToRobot(page, robotName);
  });

  test("content_search tool completes and returns results", async ({
    page,
  }) => {
    await sendMessage(page, "search for magnolia content");

    await expect(
      page
        .locator("summary")
        .filter({ hasText: "Content Search" })
        .filter({ hasText: "Tool complete" }),
    ).toBeVisible({ timeout: 15000 });

    await expect(page.getByText("Content search complete.")).toBeVisible({
      timeout: 15000,
    });
  });
});

test.describe("Robot Chat — library_search_pages tool", () => {
  let robotName: string;

  test.beforeAll(async () => {
    robotName = `e2e-library-search-robot-${Date.now()}`;
    await withAdminAccessKey(async ({ robotCreate, nodeCreate }) => {
      await nodeCreate({
        name: `Magnolia Library Page ${Date.now()}`,
        visibility: "published",
      });

      await robotCreate({
        name: robotName,
        description: "E2E robot for library_search_pages tool",
        playbook: "you are a test robot",
        model: "mock/../robot/scripts/robot-chat-library-search-pages.yaml",
      });
    });
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
    await switchToRobot(page, robotName);
  });

  test("library_search_pages tool is invoked", async ({ page }) => {
    await sendMessage(page, "search for magnolia pages");

    await expect(
      page
        .locator("summary")
        .filter({ hasText: "Library Search Pages" })
        .filter({ hasText: "Tool complete" }),
    ).toBeVisible({ timeout: 15000 });

    await expect(
      page.getByText("Search for library pages complete."),
    ).toBeVisible({ timeout: 15000 });
  });
});

test.describe("Robot Chat — thread_search tool", () => {
  let robotName: string;

  test.beforeAll(async () => {
    robotName = `e2e-thread-search-robot-${Date.now()}`;
    await withAdminAccessKey(
      async ({ robotCreate, categoryCreate, threadCreate }) => {
        const category = await categoryCreate({
          colour: "#3b82f6",
          description: "thread search test category",
          name: `Thread Search Category ${Date.now()}`,
          slug: `thread-search-cat-${Date.now()}`,
        });

        await threadCreate({
          title: "Magnolia Thread Discussion",
          body: "<p>talk about magnolia trees</p>",
          category: category.id,
          visibility: "published",
        });

        await robotCreate({
          name: robotName,
          description: "E2E robot for thread_search tool",
          playbook: "you are a test robot",
          model: "mock/../robot/scripts/robot-chat-thread-search.yaml",
        });
      },
    );
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
    await switchToRobot(page, robotName);
  });

  test("thread_search tool is invoked", async ({ page }) => {
    await sendMessage(page, "search threads about magnolia");

    await expect(
      page
        .locator("summary")
        .filter({ hasText: "Thread Search" })
        .filter({ hasText: "Tool complete" }),
    ).toBeVisible({ timeout: 15000 });

    await expect(page.getByText("Thread search complete.")).toBeVisible({
      timeout: 15000,
    });
  });
});

test.describe("Robot Chat — reply_search tool", () => {
  let robotName: string;

  test.beforeAll(async () => {
    robotName = `e2e-reply-search-robot-${Date.now()}`;
    await withAdminAccessKey(
      async ({ robotCreate, categoryCreate, threadCreate, replyCreate }) => {
        const category = await categoryCreate({
          colour: "#3b82f6",
          description: "reply search test category",
          name: `Reply Search Category ${Date.now()}`,
          slug: `reply-search-cat-${Date.now()}`,
        });

        const thread = await threadCreate({
          title: "Reply Search Test Thread",
          body: "<p>base thread for reply tests</p>",
          category: category.id,
          visibility: "published",
        });

        await replyCreate(thread.slug, {
          body: "<p>magnolia blossom reply content</p>",
        });

        await robotCreate({
          name: robotName,
          description: "E2E robot for reply_search tool",
          playbook: "you are a test robot",
          model: "mock/../robot/scripts/robot-chat-reply-search.yaml",
        });
      },
    );
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
    await switchToRobot(page, robotName);
  });

  test("reply_search tool is invoked", async ({ page }) => {
    await sendMessage(page, "search replies about magnolia");

    await expect(
      page
        .locator("summary")
        .filter({ hasText: "Reply Search" })
        .filter({ hasText: "Tool complete" }),
    ).toBeVisible({ timeout: 15000 });

    await expect(page.getByText("Reply search complete.")).toBeVisible({
      timeout: 15000,
    });
  });
});

test.describe("Robot Chat — post_search tool", () => {
  let robotName: string;

  test.beforeAll(async () => {
    robotName = `e2e-post-search-robot-${Date.now()}`;
    await withAdminAccessKey(
      async ({ robotCreate, categoryCreate, threadCreate }) => {
        const category = await categoryCreate({
          colour: "#3b82f6",
          description: "post search test category",
          name: `Post Search Category ${Date.now()}`,
          slug: `post-search-cat-${Date.now()}`,
        });

        await threadCreate({
          title: "Magnolia Post Thread",
          body: "<p>magnolia post content</p>",
          category: category.id,
          visibility: "published",
        });

        await robotCreate({
          name: robotName,
          description: "E2E robot for post_search tool",
          playbook: "you are a test robot",
          model: "mock/../robot/scripts/robot-chat-post-search.yaml",
        });
      },
    );
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
    await switchToRobot(page, robotName);
  });

  test("post_search tool is invoked", async ({ page }) => {
    await sendMessage(page, "search posts about magnolia");

    await expect(
      page
        .locator("summary")
        .filter({ hasText: "Post Search" })
        .filter({ hasText: "Tool complete" }),
    ).toBeVisible({ timeout: 15000 });

    await expect(page.getByText("Post search complete.")).toBeVisible({
      timeout: 15000,
    });
  });
});

test.describe("Robot Chat — member_search tool", () => {
  let robotName: string;

  test.beforeAll(async () => {
    robotName = `e2e-member-search-robot-${Date.now()}`;
    await withAdminAccessKey(async ({ robotCreate }) => {
      await robotCreate({
        name: robotName,
        description: "E2E robot for member_search tool",
        playbook: "you are a test robot",
        model: "mock/../robot/scripts/robot-chat-member-search.yaml",
      });
    });
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
    await switchToRobot(page, robotName);
  });

  test("member_search tool is invoked", async ({ page }) => {
    await sendMessage(page, "find member odin");

    await expect(
      page
        .locator("summary")
        .filter({ hasText: "Member Search" })
        .filter({ hasText: "Tool complete" }),
    ).toBeVisible({ timeout: 15000 });

    await expect(page.getByText("Member search complete.")).toBeVisible({
      timeout: 15000,
    });
  });
});
