import { type Locator, expect, test } from "@playwright/test";

import {
  DEFAULT_ROBOT_MODEL,
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  waitForPersistedChatRoute,
} from "./helpers";

const INTENT = "I’ll make it a focused Library curator Robot.";
const CREATED_ROBOT_NAME = "React Library Directory Curator";

test.describe("Robot Chat — message part order", () => {
  test.beforeAll(async () => {
    await setupRobotProviderWithScript(
      "mock/../robot/scripts/e2e-message-part-order.yaml",
    );
  });

  test.afterAll(async () => {
    await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
  });

  test.afterEach(async ({ page }) => {
    const response = await page.request.get("http://localhost:8001/api/robots");
    if (!response.ok()) {
      return;
    }

    const payload = (await response.json()) as {
      robots?: Array<{ id: string; name: string }>;
    };
    for (const robot of payload.robots ?? []) {
      if (robot.name === CREATED_ROBOT_NAME) {
        await page.request.delete(
          `http://localhost:8001/api/robots/${robot.id}`,
        );
      }
    }
  });

  test("keeps intent text before its tool call while streaming and after hydration", async ({
    page,
  }) => {
    await goToNewChat(page);
    await sendMessage(page, "create the curator robot");

    const messageLog = page.getByRole("log", { name: "Robot chat messages" });
    await expect(messageLog.getByText(INTENT, { exact: true })).toBeVisible({
      timeout: 15000,
    });
    await expect(
      messageLog.getByRole("group", {
        name: "Robot Create tool call",
        exact: true,
      }),
    ).toBeVisible({ timeout: 15000 });

    // The mock pauses before its post-tool response so this assertion observes
    // the live SSE projection, before the persisted session is hydrated.
    await expectMessagePartOrder(messageLog);

    await expect(
      messageLog.getByText("Created Robot: React Library Directory Curator", {
        exact: true,
      }),
    ).toBeVisible({ timeout: 15000 });

    await waitForPersistedChatRoute(page);

    await page.reload();

    const hydratedMessageLog = page.getByRole("log", {
      name: "Robot chat messages",
    });
    await expect(
      hydratedMessageLog.getByText(INTENT, { exact: true }),
    ).toBeVisible({
      timeout: 15000,
    });
    await expectMessagePartOrder(hydratedMessageLog);
  });
});

async function expectMessagePartOrder(messageLog: Locator) {
  const intent = messageLog.getByText(INTENT, { exact: true });
  const tool = messageLog.getByRole("group", {
    name: "Robot Create tool call",
    exact: true,
  });

  await expect(intent).toBeVisible({ timeout: 15000 });
  await expect(tool).toBeVisible({ timeout: 15000 });

  const toolElement = await tool.elementHandle();
  expect(toolElement).not.toBeNull();

  const intentPrecedesTool = await intent.evaluate(
    (intentNode, toolNode) =>
      toolNode !== null &&
      Boolean(
        intentNode.compareDocumentPosition(toolNode) &
        Node.DOCUMENT_POSITION_FOLLOWING,
      ),
    toolElement,
  );

  expect(intentPrecedesTool).toBe(true);
}
