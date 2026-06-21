import { expect, test } from "@playwright/test";
import { unlink, writeFile } from "node:fs/promises";

import { withAdminAccessKey } from "../access_key_admin_assignment";

import {
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  switchToRobot,
  waitForPersistedChatRoute,
} from "./helpers";

const ROBOT_SCRIPT_DIR = "../tests/robot/scripts";

test.describe("Robot Chat — presentation markup", () => {
  test.beforeAll(async () => {
    await setupRobotProviderWithScript();
  });

  test("library page card renders, hydrates, and does not use tool continuation", async ({
    page,
  }) => {
    const suffix = Date.now();
    const libraryPageName = `E2E Presented Library Page ${suffix}`;
    const libraryPageSlug = `e2e-presented-library-page-${suffix}`;
    const robotName = `e2e-presentation-robot-${suffix}`;
    const scriptName = `e2e-presentation-markup-${suffix}.yaml`;
    const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;
    let libraryPageID = "";

    try {
      await withAdminAccessKey(async ({ robotCreate, nodeCreate }) => {
        const libraryPage = await nodeCreate({
          name: libraryPageName,
          slug: libraryPageSlug,
          description:
            "A Library page rendered from Robot presentation markup.",
          content: "This page exists so the Robot can present it as a card.",
          visibility: "published",
        });
        libraryPageID = libraryPage.id;

        await writeFile(
          scriptPath,
          `steps:
  - match:
      contains: "show the rich page card"
    respond:
      text: |
        Here's the page I found:

        [${libraryPageName}](sdr:node/${libraryPage.id})

        Let me know if you'd like a summary.
      finish: "stop"
  - match:
      contains: "next message"
    respond:
      text: "Next message handled normally."
      finish: "stop"
`,
        );

        await robotCreate({
          name: robotName,
          description: "E2E robot that emits presentation markup",
          playbook: "you present Library pages with presentation markup",
          model: `mock/../robot/scripts/${scriptName}`,
        });
      });

      await goToNewChat(page);

      await switchToRobot(page, robotName);

      await sendMessage(page, "show the rich page card");

      await expect(page.getByText("Here's the page I found:")).toBeVisible({
        timeout: 15000,
      });
      const card = page.locator(`article[id="${libraryPageID}"]`);
      await expect(card).toBeVisible({ timeout: 15000 });
      await expect(
        card.getByRole("link", { name: libraryPageName }),
      ).toBeVisible();
      await expect(page.getByText("sdr:node")).toHaveCount(0);
      await waitForPersistedChatRoute(page);

      await page.reload();
      const hydratedCard = page.locator(`article[id="${libraryPageID}"]`);
      await expect(hydratedCard).toBeVisible({ timeout: 15000 });
      await expect(
        hydratedCard.getByRole("link", { name: libraryPageName }),
      ).toBeVisible();
      await expect(page.getByText("sdr:node")).toHaveCount(0);

      await sendMessage(page, "next message");
      await expect(
        page.getByText("Next message handled normally."),
      ).toBeVisible({
        timeout: 15000,
      });
      await waitForPersistedChatRoute(page);

      await page.reload();
      await expect(page.getByText("show the rich page card")).toBeVisible({
        timeout: 15000,
      });
      await expect(page.getByText("Here's the page I found:")).toBeVisible();
      const finalCard = page.locator(`article[id="${libraryPageID}"]`);
      await expect(finalCard).toBeVisible();
      await expect(
        finalCard.getByRole("link", { name: libraryPageName }),
      ).toBeVisible();
      await expect(
        page.getByText("Let me know if you'd like a summary."),
      ).toBeVisible();
      await expect(
        page.getByText("next message", { exact: true }),
      ).toBeVisible();
      await expect(
        page.getByText("Next message handled normally."),
      ).toBeVisible();
    } finally {
      await unlink(scriptPath).catch(() => undefined);
    }
  });
});
