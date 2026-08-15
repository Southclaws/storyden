import { expect, test } from "@playwright/test";

import {
  DEFAULT_ROBOT_MODEL,
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  waitForPersistedChatRoute,
} from "./helpers";

test.describe("Robot Chat — progressive tool discovery", () => {
  test.afterAll(async () => {
    await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
  });

  test("Denbot searches, inspects, loads, and invokes one tool", async ({
    page,
  }) => {
    await setupRobotProviderWithScript(
      "mock/../robot/scripts/e2e-tool-discovery.yaml",
    );
    await goToNewChat(page);

    await sendMessage(page, "find and load one error tool");

    for (const toolName of [
      "Tool Search",
      "Tool Get",
      "Tool Load",
      "Throw An Error",
    ]) {
      await expect(
        page
          .locator("summary")
          .filter({ hasText: toolName })
          .filter({ hasText: "Tool complete" }),
      ).toBeVisible({ timeout: 15000 });
    }
    await expect(page.getByText("Loaded 1 tools")).toBeVisible();
    await expect(
      page.getByText(
        "The dynamically loaded tool ran and reported its expected test error.",
      ),
    ).toBeVisible({ timeout: 15000 });
    await waitForPersistedChatRoute(page);

    await page.reload();
    await expect(
      page.getByText(
        "The dynamically loaded tool ran and reported its expected test error.",
      ),
    ).toBeVisible({ timeout: 15000 });
    await expect(
      page
        .locator("summary")
        .filter({ hasText: "Tool Load" })
        .filter({ hasText: "Tool complete" }),
    ).toBeVisible();
  });
});
