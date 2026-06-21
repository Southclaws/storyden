import { Page, expect, test } from "@playwright/test";
import { unlink, writeFile } from "node:fs/promises";

import { withAdminAccessKey } from "../access_key_admin_assignment";

import {
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  switchToRobot,
} from "./helpers";

const ROBOT_SCRIPT_DIR = "../tests/robot/scripts";

function completedLibraryRequestTools(page: Page) {
  return page
    .locator("summary")
    .filter({ hasText: "Library Request Page" })
    .filter({ hasText: "Tool complete" });
}

async function selectRequestedLibraryPage(
  page: Page,
  name: string,
  expectedCompletedCount: number,
) {
  const pendingTool = page
    .locator("summary")
    .filter({ hasText: "Library Request Page" })
    .filter({ hasText: "Needs selection" });
  const pageOption = page.getByRole("button", { name });
  const completedTools = completedLibraryRequestTools(page);

  await expect
    .poll(
      async () => {
        if ((await completedTools.count()) >= expectedCompletedCount) {
          return "complete";
        }

        if (await pageOption.isVisible().catch(() => false)) {
          return "ready";
        }

        return "waiting";
      },
      { timeout: 15000 },
    )
    .not.toBe("waiting");

  if ((await completedTools.count()) >= expectedCompletedCount) {
    return;
  }

  await expect(pendingTool).toBeVisible({ timeout: 15000 });
  await pageOption.click();
  await expect(completedTools).toHaveCount(expectedCompletedCount, {
    timeout: 15000,
  });
}

async function completedLibraryRequestCardsAfter(
  page: Page,
  afterText: string,
) {
  const afterElement = await page.getByText(afterText).elementHandle();
  expect(afterElement).not.toBeNull();

  return page
    .locator("details")
    .filter({ hasText: "Library Request Page" })
    .filter({ hasText: "Tool complete" })
    .evaluateAll((nodes, anchor) => {
      if (!anchor) {
        return 0;
      }

      return nodes.filter((node) =>
        Boolean(
          anchor.compareDocumentPosition(node) &
          Node.DOCUMENT_POSITION_FOLLOWING,
        ),
      ).length;
    }, afterElement);
}

test.describe("Robot Chat — Library page request tool", () => {
  test.beforeAll(async () => {
    await setupRobotProviderWithScript();
  });

  test("library_request_page pauses for page selection and resumes", async ({
    page,
  }) => {
    const suffix = Date.now();
    const libraryPageName = `E2E Requested Library Page ${suffix}`;
    const libraryPageSlug = `e2e-requested-library-page-${suffix}`;
    const robotName = `e2e-page-request-robot-${suffix}`;
    const scriptName = `e2e-library-request-page-${suffix}.yaml`;
    const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;

    try {
      await withAdminAccessKey(async ({ robotCreate, nodeCreate }) => {
        await nodeCreate({
          name: libraryPageName,
          slug: libraryPageSlug,
          content:
            "This page exists so the Robot can ask the user to choose it.",
          visibility: "published",
        });

        await writeFile(
          scriptPath,
          `steps:
  - match:
      contains: "choose another library page"
    respond:
      tool_calls:
        - id: call_request_page_2
          name: library_request_page
          args: {}
  - match:
      contains: "choose a library page"
    respond:
      tool_calls:
        - id: call_request_page_1
          name: library_request_page
          args: {}
  - match:
      tool_result: library_request_page
    respond:
      text: "I can continue with the selected Library page."
      finish: "stop"
`,
        );

        await robotCreate({
          name: robotName,
          description: "E2E robot that asks the user to select a Library page",
          playbook: "you request a Library page when the task needs one",
          model: `mock/../robot/scripts/${scriptName}`,
          tools: ["library_request_page"],
        });
      });

      await goToNewChat(page);

      await switchToRobot(page, robotName);

      const firstPrompt = "choose a library page";
      await sendMessage(page, firstPrompt);
      await selectRequestedLibraryPage(page, libraryPageName, 1);

      await expect(
        page.getByText("I can continue with the selected Library page."),
      ).toBeVisible({ timeout: 15000 });

      const secondPrompt = "choose another library page";
      await sendMessage(page, secondPrompt);
      await selectRequestedLibraryPage(page, libraryPageName, 2);

      const completedPageRequestTool = completedLibraryRequestTools(page);

      await expect(completedPageRequestTool).toHaveCount(2, {
        timeout: 15000,
      });
      await expect(
        page
          .locator("summary")
          .filter({ hasText: "Library Request Page" })
          .filter({ hasText: "Needs selection" }),
      ).toHaveCount(0);
      await expect(
        page.getByText("I can continue with the selected Library page.").last(),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        page
          .locator("details")
          .filter({ hasText: "Library Request Page" })
          .filter({ hasText: "Tool complete" }),
      ).toHaveCount(2);
      await expect
        .poll(() => completedLibraryRequestCardsAfter(page, secondPrompt), {
          timeout: 15000,
        })
        .toBe(1);
      await expect(page.getByText("No tool invocation found")).toHaveCount(0);
    } finally {
      await unlink(scriptPath).catch(() => undefined);
    }
  });
});
