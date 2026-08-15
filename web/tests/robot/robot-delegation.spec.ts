import { expect, test } from "@playwright/test";
import { unlink, writeFile } from "node:fs/promises";

import { withAdminAccessKey } from "../access_key_admin_assignment";

import {
  DEFAULT_ROBOT_MODEL,
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  waitForPersistedChatRoute,
} from "./helpers";

const ROBOT_SCRIPT_DIR = "../tests/robot/scripts";

test.describe("Robot Chat — specialist delegation", () => {
  const suffix = Date.now();
  const specialistName = `E2E research specialist ${suffix}`;
  const specialistScriptName = `e2e-delegation-specialist-${suffix}.yaml`;
  const specialistScriptPath = `${ROBOT_SCRIPT_DIR}/${specialistScriptName}`;
  const coordinatorScriptName = `e2e-delegation-coordinator-${suffix}.yaml`;
  const coordinatorScriptPath = `${ROBOT_SCRIPT_DIR}/${coordinatorScriptName}`;
  let specialistID = "";

  test.beforeAll(async () => {
    const specialistTranscript = Array.from(
      { length: 36 },
      (_, index) =>
        `        - Evidence note ${index + 1}: inspected a bounded candidate and recorded its relevance.`,
    ).join("\n");

    await writeFile(
      specialistScriptPath,
      `steps:
  - match:
      tool_result: content_search
    respond:
      delay_ms: 5000
      text: "The specialist found the delegated evidence."
      finish: "stop"
  - match:
      contains: "test the specialist directly"
    respond:
      text: "The specialist handled this conversation directly."
      finish: "stop"
  - match:
      any: true
    respond:
      text: |
        I’ll inspect the relevant content.
${specialistTranscript}
      tool_calls:
        - id: call_specialist_search
          name: content_search
          args:
            query: "delegated evidence"
            kind: ["thread"]
            max_results: 1
`,
    );

    await withAdminAccessKey(async ({ robotCreate }) => {
      const specialist = await robotCreate({
        name: specialistName,
        description: "Finds a bounded piece of evidence for Denbot.",
        playbook:
          "Complete only the delegated research request and return concise evidence.",
        model: `mock/../robot/scripts/${specialistScriptName}`,
        tools: ["content_search"],
      });
      specialistID = specialist.id;
    });

    const delegatedAgentName = `robot_${specialistID}`;
    await writeFile(
      coordinatorScriptPath,
      `steps:
  - match:
      tool_result: ${delegatedAgentName}
    respond:
      text: "Denbot synthesised the specialist evidence."
      finish: "stop"
  - match:
      contains: "delegate this research"
    respond:
      tool_calls:
        - id: call_delegate_1
          name: ${delegatedAgentName}
          args:
            request: "Find the delegated evidence."
`,
    );
  });

  test.afterAll(async () => {
    if (specialistID) {
      await withAdminAccessKey(async ({ robotDelete }) => {
        await robotDelete(specialistID);
      });
    }
    await unlink(coordinatorScriptPath).catch(() => undefined);
    await unlink(specialistScriptPath).catch(() => undefined);
    await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
  });

  test("an arbitrary custom Robot can be the root conversation", async ({
    page,
  }) => {
    await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
    await goToNewChat(page);
    await page.goto(`/robots/chats/new?robot=${specialistID}`);

    await sendMessage(page, "test the specialist directly");

    await expect(
      page.getByText("The specialist handled this conversation directly."),
    ).toBeVisible({ timeout: 15000 });
    await waitForPersistedChatRoute(page);

    await page.reload();
    await expect(
      page.getByText("The specialist handled this conversation directly."),
    ).toBeVisible({ timeout: 15000 });
  });

  test("Denbot presents delegated work as an expandable Robot transcript", async ({
    page,
  }) => {
    await setupRobotProviderWithScript(
      `mock/../robot/scripts/${coordinatorScriptName}`,
    );
    await goToNewChat(page);

    await sendMessage(page, "delegate this research");

    const delegation = page.getByRole("group", {
      name: `${specialistName} delegation`,
    });
    const delegatedEvidence = delegation.getByText(
      "The specialist found the delegated evidence.",
    );
    await expect(
      page.getByText("Denbot synthesised the specialist evidence."),
    ).toBeVisible({ timeout: 15000 });
    await expect(delegation).toBeVisible();
    await expect(delegation.getByText("Completed")).toBeVisible();
    await expect(
      delegation
        .locator("summary")
        .getByText("Find the delegated evidence.", { exact: true }),
    ).toBeVisible();
    await expect(delegatedEvidence).toBeHidden();
    await expect(
      page.locator(
        '[role="group"][aria-label^="Robot D9"][aria-label$=" tool call"]',
      ),
    ).toHaveCount(0);
    await expect(page.getByText("Thought", { exact: true })).toHaveCount(0);

    await delegation.locator("summary").first().click();
    await expect(delegatedEvidence).toBeVisible();
    await expect(
      delegation
        .getByRole("group", { name: "Content Search tool call" })
        .first(),
    ).toBeVisible();
    await waitForPersistedChatRoute(page);

    await page.reload();
    await expect(
      page.getByText("Denbot synthesised the specialist evidence."),
    ).toBeVisible({ timeout: 15000 });
    await expect(delegatedEvidence).toBeHidden();
    await expect(delegation.getByText("Completed")).toBeVisible();
    await delegation.locator("summary").first().click();
    await expect(delegatedEvidence).toBeVisible();
    await expect(
      delegation
        .getByRole("group", { name: "Content Search tool call" })
        .first(),
    ).toBeVisible();
  });

  test("a long delegated transcript stays inside the chat viewport", async ({
    page,
  }) => {
    await page.setViewportSize({ width: 1280, height: 720 });
    await setupRobotProviderWithScript(
      `mock/../robot/scripts/${coordinatorScriptName}`,
    );
    await goToNewChat(page);

    await sendMessage(page, "delegate this research");

    const delegation = page.getByRole("group", {
      name: `${specialistName} delegation`,
    });
    await expect(delegation).toBeVisible({ timeout: 15000 });
    await delegation.locator("summary").first().click();

    const delegatedToolCalls = delegation.getByRole("group", {
      name: "Content Search tool call",
      exact: true,
    });
    await expect(delegatedToolCalls).toHaveCount(1, { timeout: 15000 });

    const messageViewport = page.getByRole("log", {
      name: "Robot chat messages",
    });
    const composer = page.getByRole("form", {
      name: "Send message to Robot",
    });

    await expect
      .poll(() =>
        messageViewport.evaluate(
          (element) => element.scrollHeight > element.clientHeight,
        ),
      )
      .toBe(true);

    await messageViewport.evaluate((element) => {
      element.scrollTop = 0;
      element.dispatchEvent(new Event("scroll"));
    });

    const newMessages = page.getByRole("button", { name: "New messages" });
    await expect(newMessages).toBeVisible({ timeout: 15000 });

    const viewportBox = await messageViewport.boundingBox();
    const composerBox = await composer.boundingBox();
    const newMessagesBox = await newMessages.boundingBox();
    expect(viewportBox).not.toBeNull();
    expect(composerBox).not.toBeNull();
    expect(newMessagesBox).not.toBeNull();
    expect(viewportBox!.y + viewportBox!.height).toBeLessThanOrEqual(
      composerBox!.y + 1,
    );
    expect(newMessagesBox!.y).toBeGreaterThanOrEqual(viewportBox!.y);
    expect(newMessagesBox!.y + newMessagesBox!.height).toBeLessThanOrEqual(
      viewportBox!.y + viewportBox!.height + 1,
    );

    await newMessages.click();
    await expect
      .poll(() =>
        messageViewport.evaluate(
          (element) =>
            element.scrollHeight - element.scrollTop - element.clientHeight,
        ),
      )
      .toBeLessThan(5);

    const finalResponse = page.getByText(
      "Denbot synthesised the specialist evidence.",
    );
    await expect(finalResponse).toBeVisible({ timeout: 15000 });
    const finalResponseBox = await finalResponse.boundingBox();
    expect(finalResponseBox).not.toBeNull();
    expect(finalResponseBox!.y + finalResponseBox!.height).toBeLessThanOrEqual(
      composerBox!.y + 1,
    );
  });
});
