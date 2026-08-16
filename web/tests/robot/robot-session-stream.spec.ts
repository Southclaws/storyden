import { Browser, BrowserContext, Page, expect, test } from "@playwright/test";
import { unlink, writeFile } from "node:fs/promises";

import { createAdmin, login } from "../access_key_admin_assignment";

import {
  DEFAULT_ROBOT_MODEL,
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  waitForPersistedChatRoute,
} from "./helpers";

const ROBOT_SCRIPT_DIR = "../tests/robot/scripts";
const VIEWER_PASSWORD = "E2EStreamViewerPassword123!";
const MOCK_RESPONSE = "This is a mock response from the test provider.";

type IndependentPage = {
  context: BrowserContext;
  page: Page;
  username: string;
};

async function createIndependentAdminPage(
  browser: Browser,
  name: string,
): Promise<IndependentPage> {
  const context = await browser.newContext();
  const page = await context.newPage();
  const suffix = Date.now().toString(36);
  const prefix = `stream-${name}`;
  const username = `${prefix.slice(0, 30 - suffix.length - 1)}-${suffix}`;

  await createAdmin(context, username, VIEWER_PASSWORD);
  await login(page, username, VIEWER_PASSWORD);

  return { context, page, username };
}

async function createIndependentCurrentMemberPage(
  browser: Browser,
): Promise<IndependentPage> {
  const context = await browser.newContext();
  const page = await context.newPage();

  await goToNewChat(page);

  return { context, page, username: "e2e_admin" };
}

function messageLog(page: Page) {
  return page.getByRole("log", { name: "Robot chat messages" });
}

test.describe("Robot Chat — session stream", () => {
  test.beforeAll(async () => {
    await setupRobotProviderWithScript();
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
  });

  test("an idle member observes a turn started in another browser", async ({
    browser,
    page: sender,
  }) => {
    await sendMessage(sender, "establish a shared session");
    await expect(sender.getByText(MOCK_RESPONSE)).toBeVisible({
      timeout: 15000,
    });
    await waitForPersistedChatRoute(sender);
    const sessionURL = new URL(sender.url()).pathname;

    const viewer = await createIndependentAdminPage(browser, "idle-viewer");
    try {
      await viewer.page.goto(sessionURL);
      await expect(viewer.page.getByText(MOCK_RESPONSE)).toBeVisible({
        timeout: 15000,
      });
      await expect(
        viewer.page.getByRole("status").filter({ hasText: /is responding/ }),
      ).toHaveCount(0);

      const liveMessage = "sent while another member is watching";
      await sendMessage(sender, liveMessage);

      await expect(
        messageLog(viewer.page)
          .getByRole("article", { name: "@e2e_admin message" })
          .filter({ hasText: liveMessage }),
      ).toBeVisible({ timeout: 15000 });
      await expect(viewer.page.getByText(MOCK_RESPONSE)).toHaveCount(2, {
        timeout: 15000,
      });

      const viewerReply = "reply from the observing member";
      await sendMessage(viewer.page, viewerReply);

      await expect(
        messageLog(sender)
          .getByRole("article", {
            name: `@${viewer.username} message`,
          })
          .filter({ hasText: viewerReply }),
      ).toBeVisible({ timeout: 15000 });
      await expect(sender.getByText(MOCK_RESPONSE)).toHaveCount(3, {
        timeout: 15000,
      });
    } finally {
      await viewer.context.close();
    }
  });

  test("a member attaching mid-turn catches up and follows the live tail", async ({
    browser,
    page: sender,
  }) => {
    const suffix = Date.now();
    const scriptName = `e2e-session-stream-mid-turn-${suffix}.yaml`;
    const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;
    const viewer = await createIndependentAdminPage(browser, "mid-turn-viewer");

    try {
      await writeFile(
        scriptPath,
        `steps:
  - match:
      contains: "attach during this turn"
    respond:
      delay_ms: 5000
      text: "The observer followed the live turn to completion."
      finish: "stop"
`,
      );
      await setupRobotProviderWithScript(`mock/../robot/scripts/${scriptName}`);

      await sender.goto("/robots/chats/new");
      const initiatingMessage = "attach during this turn";
      await sendMessage(sender, initiatingMessage);
      await waitForPersistedChatRoute(sender);
      await expect(sender.getByText("Denbot is responding...")).toBeVisible();

      await viewer.page.goto(new URL(sender.url()).pathname);
      await expect(
        messageLog(viewer.page).getByText(initiatingMessage, { exact: true }),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        viewer.page.getByText("Denbot is responding..."),
      ).toBeVisible({ timeout: 15000 });

      await expect(
        viewer.page.getByText(
          "The observer followed the live turn to completion.",
        ),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        sender.getByText("The observer followed the live turn to completion."),
      ).toBeVisible({ timeout: 15000 });
    } finally {
      await viewer.context.close();
      await unlink(scriptPath).catch(() => undefined);
      await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
    }
  });

  test("a returning member catches up once and resumes observing", async ({
    browser,
    page: sender,
  }) => {
    await sendMessage(sender, "establish a reconnectable session");
    await expect(sender.getByText(MOCK_RESPONSE)).toBeVisible({
      timeout: 15000,
    });
    await waitForPersistedChatRoute(sender);
    const sessionURL = new URL(sender.url()).pathname;

    const viewer = await createIndependentCurrentMemberPage(browser);
    try {
      await viewer.page.goto(sessionURL);
      await expect(viewer.page.getByText(MOCK_RESPONSE)).toBeVisible({
        timeout: 15000,
      });

      await viewer.page.goto("/robots/chats");

      const missedMessage = "sent while the viewer is disconnected";
      await sendMessage(sender, missedMessage);
      await expect(sender.getByText(MOCK_RESPONSE)).toHaveCount(2, {
        timeout: 15000,
      });

      await viewer.page.goto(sessionURL);
      await expect(
        messageLog(viewer.page).getByText(missedMessage, { exact: true }),
      ).toHaveCount(1, { timeout: 15000 });
      await expect(viewer.page.getByText(MOCK_RESPONSE)).toHaveCount(2, {
        timeout: 15000,
      });

      const resumedMessage = "sent after the viewer reconnects";
      await sendMessage(sender, resumedMessage);
      await expect(
        messageLog(viewer.page).getByText(resumedMessage, { exact: true }),
      ).toHaveCount(1, { timeout: 15000 });
      await expect(viewer.page.getByText(MOCK_RESPONSE)).toHaveCount(3, {
        timeout: 15000,
      });
    } finally {
      await viewer.context.close();
    }
  });
});
