import { Page, expect, test } from "@playwright/test";
import { unlink, writeFile } from "node:fs/promises";

import { withAdminAccessKey } from "../access_key_admin_assignment";

import {
  DEFAULT_ROBOT_MODEL,
  dismissOnboarding,
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  switchToRobot,
} from "./helpers";

const ROBOT_SCRIPT_DIR = "../tests/robot/scripts";

async function getLatestRobotSessionID(page: Page) {
  const response = await page.request.get(
    "http://localhost:8001/api/robots/sessions",
  );
  if (!response.ok()) {
    throw new Error(`failed to list robot sessions: ${response.status()}`);
  }

  const payload = (await response.json()) as {
    sessions?: { id: string }[];
  };
  const session = payload.sessions?.[0];
  if (!session) {
    throw new Error("no robot sessions found");
  }

  return session.id;
}

async function openLatestRobotSession(page: Page) {
  const sessionID = await getLatestRobotSessionID(page);

  await page.goto(`/robots/chats/${sessionID}`);
  await dismissOnboarding(page);
}

// The mock LLM uses ../tests/robot/scripts/e2e-default.yaml, selected through
// the Robots provider settings API in the suite setup.
// It routes:
//   - any message            → static text response
//   - "list pages"           → library_page_list tool call then text
//   - "trigger error"        → LLM-level error (no text response)

test.describe("Robot Chat — mock LLM stream", () => {
  // The default global agent does not expose every tool, so this custom Robot
  // declares the library_page_list tool used by the mock script.
  let robotName: string;

  test.beforeAll(async () => {
    robotName = `e2e-robot-${Date.now()}`;
    await setupRobotProviderWithScript();
    await withAdminAccessKey(async ({ robotCreate }) => {
      await robotCreate({
        name: robotName,
        description: "E2E test robot with all tools",
        playbook: "you are a test robot",
        model: DEFAULT_ROBOT_MODEL,
        tools: ["library_page_list"],
      });
    });
  });

  test.beforeEach(async ({ page }) => {
    await goToNewChat(page);
  });

  test("text response arrives and renders in the message list", async ({
    page,
  }) => {
    await sendMessage(page, "hello");

    await expect(
      page.getByText("This is a mock response from the test provider."),
    ).toBeVisible({ timeout: 15000 });
  });

  test("tool call renders with name and completion status", async ({
    page,
  }) => {
    // Switch to the robot that exposes all tools (including library_page_list)
    await switchToRobot(page, robotName);

    await sendMessage(page, "list pages");

    await expect(
      page
        .locator("summary")
        .filter({ hasText: "Library Page List" })
        .filter({ hasText: "Tool complete" }),
    ).toBeVisible({ timeout: 15000 });

    // After the tool completes the mock sends a text follow-up
    await expect(
      page.getByText("I found some pages in the library."),
    ).toBeVisible({ timeout: 15000 });
  });

  test("existing custom robot chat can switch back to the default robot", async ({
    page,
  }) => {
    const suffix = Date.now();
    const customRobotName = `e2e-reset-robot-${suffix}`;
    const scriptName = `e2e-reset-${suffix}.yaml`;
    const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;

    try {
      await writeFile(
        scriptPath,
        `steps:
  - match:
      any: true
    respond:
      text: "Custom reset robot handled this."
      finish: "stop"
`,
      );

      let customRobotID = "";
      await withAdminAccessKey(async ({ robotCreate }) => {
        const customRobot = await robotCreate({
          name: customRobotName,
          description: "E2E reset robot",
          playbook: "you are the reset test robot",
          model: `mock/../robot/scripts/${scriptName}`,
        });
        customRobotID = customRobot.id;
      });

      await page.goto("/robots/chats/new");
      await dismissOnboarding(page);

      await switchToRobot(page, customRobotName);

      await sendMessage(page, "custom first");
      await expect(
        page.getByText("Custom reset robot handled this."),
      ).toBeVisible({ timeout: 15000 });

      const sessionID = await getLatestRobotSessionID(page);

      await page.goto(`/robots/chats/${sessionID}?robot=${customRobotID}`);
      await dismissOnboarding(page);

      await expect(
        page.getByRole("button", { name: customRobotName }),
      ).toBeVisible({ timeout: 15000 });
      await page.getByRole("button", { name: customRobotName }).click();
      await page
        .getByRole("menuitem", { name: "Storyden Robot Builder" })
        .click();
      await expect(
        page.getByRole("button", { name: "Storyden Robot Builder" }),
      ).toBeVisible({ timeout: 15000 });
      await expect(page).toHaveURL(
        new RegExp(`/robots/chats/${sessionID}\\?robot=robot_builder$`),
      );

      await page.reload();
      await dismissOnboarding(page);
      await expect(
        page.getByRole("button", { name: "Storyden Robot Builder" }),
      ).toBeVisible({ timeout: 15000 });

      await sendMessage(page, "default after reset");
      await expect(
        page.getByText("This is a mock response from the test provider."),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        page.getByText("Custom reset robot handled this."),
      ).toHaveCount(1);
    } finally {
      await unlink(scriptPath).catch(() => undefined);
    }
  });

  test("LLM error shows error message in the UI", async ({ page }) => {
    await sendMessage(page, "trigger error");

    await expect(
      page.locator("p").filter({ hasText: "simulated LLM provider failure" }),
    ).toBeVisible({ timeout: 15000 });
  });

  test("cancelled response can be followed by another message and refreshed", async ({
    page,
  }) => {
    const suffix = Date.now();
    const cancellableRobotName = `e2e-cancellable-robot-${suffix}`;
    const scriptName = `e2e-cancellable-${suffix}.yaml`;
    const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;

    try {
      await writeFile(
        scriptPath,
        `steps:
  - match:
      contains: "slow response"
    respond:
      delay_ms: 5000
      text: "This cancelled response should not render."
      finish: "stop"
  - match:
      contains: "after cancel"
    respond:
      text: "Follow-up after cancellation rendered."
      finish: "stop"
`,
      );

      await withAdminAccessKey(async ({ robotCreate }) => {
        await robotCreate({
          name: cancellableRobotName,
          description: "E2E cancellation robot",
          playbook: "you are the cancellation test robot",
          model: `mock/../robot/scripts/${scriptName}`,
        });
      });

      await page.goto("/robots/chats/new");
      await dismissOnboarding(page);

      await switchToRobot(page, cancellableRobotName);

      await sendMessage(page, "slow response");
      await expect(
        page.getByText(`${cancellableRobotName} is responding...`),
      ).toBeVisible({ timeout: 15000 });

      await page.getByRole("button", { name: "Cancel Robot response" }).click();
      await expect(
        page.getByText(`${cancellableRobotName} is responding...`),
      ).toHaveCount(0, { timeout: 15000 });
      await expect(
        page.getByText("This cancelled response should not render."),
      ).toHaveCount(0);

      await sendMessage(page, "after cancel");
      await expect(
        page.getByText("Follow-up after cancellation rendered."),
      ).toBeVisible({ timeout: 15000 });

      await openLatestRobotSession(page);

      await expect(page.getByText("slow response")).toBeVisible({
        timeout: 15000,
      });
      await expect(page.getByText("after cancel", { exact: true })).toBeVisible(
        {
          timeout: 15000,
        },
      );
      await expect(
        page.getByText("Follow-up after cancellation rendered."),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        page.getByText("This cancelled response should not render."),
      ).toHaveCount(0);
    } finally {
      await unlink(scriptPath).catch(() => undefined);
    }
  });

  test("robot switch tool call changes robot and continues the conversation", async ({
    page,
  }) => {
    const suffix = Date.now();
    const targetRobotName = `e2e-target-robot-${suffix}`;
    const switcherRobotName = `e2e-switcher-robot-${suffix}`;
    const switcherScriptName = `e2e-switcher-${suffix}.yaml`;
    const switcherScriptPath = `${ROBOT_SCRIPT_DIR}/${switcherScriptName}`;

    try {
      await withAdminAccessKey(async ({ robotCreate }) => {
        const targetRobot = await robotCreate({
          name: targetRobotName,
          description: "E2E target robot for switch tool",
          playbook: "you are the target test robot",
          model: "mock/../robot/scripts/e2e-switch-target.yaml",
        });

        await writeFile(
          switcherScriptPath,
          `steps:
  - match:
      contains: "start switch flow"
    respond:
      text: "Switcher robot started the conversation."
      finish: "stop"
  - match:
      contains: "switch to target"
    respond:
      tool_calls:
        - id: call_switch_1
          name: robot_switch
          args:
            robot_id: "${targetRobot.id}"
  - match:
      contains: "ROBOT SWITCH"
    respond:
      text: "Switcher robot incorrectly resumed after switch."
      finish: "stop"
  - match:
      any: true
    respond:
      text: "Switcher robot stayed active."
      finish: "stop"
`,
        );

        await robotCreate({
          name: switcherRobotName,
          description: "E2E switcher robot",
          playbook: "you are the switcher test robot",
          model: `mock/../robot/scripts/${switcherScriptName}`,
          tools: ["robot_switch"],
        });
      });

      await page.goto("/robots/chats/new");
      await dismissOnboarding(page);

      await switchToRobot(page, switcherRobotName);

      await sendMessage(page, "start switch flow");
      await expect(
        page.getByText("Switcher robot started the conversation."),
      ).toBeVisible({ timeout: 15000 });

      await sendMessage(page, "switch to target");
      await expect(
        page
          .locator("summary")
          .filter({ hasText: "Robot Switch" })
          .filter({ hasText: "Tool complete" })
          .first(),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        page.getByRole("separator", { name: "Robot switched" }),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        page.getByRole("button", { name: targetRobotName }),
      ).toBeVisible({ timeout: 15000 });

      await sendMessage(page, "target followup one");
      await expect(
        page.getByText("Target robot handled followup one."),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        page.getByText("Switcher robot incorrectly resumed after switch."),
      ).toHaveCount(0);

      await page.getByRole("button", { name: targetRobotName }).click();
      await page
        .getByRole("menuitem", { name: "Storyden Robot Builder" })
        .click();
      await expect(
        page.getByRole("button", { name: "Storyden Robot Builder" }),
      ).toBeVisible({ timeout: 15000 });

      await sendMessage(page, "hello from default robot");
      await expect(
        page.getByText("This is a mock response from the test provider."),
      ).toBeVisible({ timeout: 15000 });

      await switchToRobot(page, targetRobotName);

      await sendMessage(page, "target followup two");
      await expect(
        page.getByText("Target robot handled followup two."),
      ).toBeVisible({ timeout: 15000 });
    } finally {
      await unlink(switcherScriptPath).catch(() => undefined);
    }
  });
});
