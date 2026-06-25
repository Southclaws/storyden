import { expect, test } from "@playwright/test";
import { unlink, writeFile } from "node:fs/promises";

import { withAdminAccessKey } from "../access_key_admin_assignment";

import {
  DEFAULT_ROBOT_MODEL,
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  switchToRobot,
} from "./helpers";

const ROBOT_SCRIPT_DIR = "../tests/robot/scripts";

test.describe("Robot Chat — tool confirmation", () => {
  test.beforeAll(async () => {
    await setupRobotProviderWithScript();
  });

  test("robot_delete pauses for approval before deleting", async ({ page }) => {
    const suffix = Date.now();
    const victimRobotName = `e2e-delete-victim-${suffix}`;
    const deleterRobotName = `e2e-delete-actor-${suffix}`;
    const deleteScriptName = `e2e-delete-confirmation-${suffix}.yaml`;
    const deleteScriptPath = `${ROBOT_SCRIPT_DIR}/${deleteScriptName}`;
    let victimRobotID = "";

    try {
      await withAdminAccessKey(async ({ robotCreate }) => {
        const victimRobot = await robotCreate({
          name: victimRobotName,
          description: "E2E victim robot for confirmation flow",
          playbook:
            "you are the robot that should only be deleted after approval",
          model: DEFAULT_ROBOT_MODEL,
        });
        victimRobotID = victimRobot.id;

        await writeFile(
          deleteScriptPath,
          `steps:
  - match:
      contains: "delete the victim"
    respond:
      tool_calls:
        - id: call_delete_1
          name: robot_delete
          args:
            id: "${victimRobotID}"
  - match:
      tool_result: robot_delete
    respond:
      text: "Delete flow finished."
      finish: "stop"
`,
        );

        await robotCreate({
          name: deleterRobotName,
          description: "E2E robot that requests robot_delete",
          playbook: "you request robot deletion when asked",
          model: `mock/../robot/scripts/${deleteScriptName}`,
          tools: ["robot_delete"],
        });
      });

      await goToNewChat(page);

      await switchToRobot(page, deleterRobotName);

      await sendMessage(page, "delete the victim");

      await expect(
        page
          .locator("summary")
          .filter({ hasText: "Robot Delete" })
          .filter({ hasText: "Needs approval" }),
      ).toBeVisible({ timeout: 15000 });
      await expect(page.getByRole("button", { name: "Approve" })).toBeVisible();
      await expect(page.getByRole("button", { name: "Deny" })).toBeVisible();

      await withAdminAccessKey(async ({ robotGet }) => {
        await expect(robotGet(victimRobotID)).resolves.toMatchObject({
          id: victimRobotID,
        });
      });

      await page.getByRole("button", { name: "Approve" }).click();

      await expect(
        page
          .locator("summary")
          .filter({ hasText: "Robot Delete" })
          .filter({ hasText: "Tool complete" }),
      ).toBeVisible({ timeout: 15000 });
      await expect(page.getByText("Delete flow finished.")).toBeVisible({
        timeout: 15000,
      });

      await withAdminAccessKey(async ({ robotGet }) => {
        let deleted = false;
        try {
          await robotGet(victimRobotID);
        } catch {
          deleted = true;
        }
        expect(deleted).toBe(true);
      });
    } finally {
      await unlink(deleteScriptPath).catch(() => undefined);
    }
  });

  test("robot_delete confirms after another tool call in the same chat", async ({
    page,
  }) => {
    const suffix = Date.now();
    const victimRobotName = `e2e-mixed-victim-${suffix}`;
    const actorRobotName = `e2e-mixed-actor-${suffix}`;
    const scriptName = `e2e-delete-after-create-${suffix}.yaml`;
    const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;
    let victimRobotID = "";

    try {
      await withAdminAccessKey(async ({ robotCreate }) => {
        const victimRobot = await robotCreate({
          name: victimRobotName,
          description: "E2E victim robot for mixed confirmation flow",
          playbook:
            "you are the robot that should only be deleted after approval",
          model: DEFAULT_ROBOT_MODEL,
        });
        victimRobotID = victimRobot.id;

        await writeFile(
          scriptPath,
          `steps:
  - match:
      contains: "create a robot"
    respond:
      tool_calls:
        - id: call_create_1
          name: robot_create
          args:
            name: "Temporary Test Robot ${suffix}"
            description: "temporary robot created before deletion"
            playbook: "you are temporary"
  - match:
      tool_result: robot_create
    respond:
      text: "Created Robot."
      finish: "stop"
  - match:
      contains: "delete it"
    respond:
      tool_calls:
        - id: call_delete_1
          name: robot_delete
          args:
            id: "${victimRobotID}"
  - match:
      tool_result: robot_delete
    respond:
      text: "Delete flow finished."
      finish: "stop"
`,
        );

        await robotCreate({
          name: actorRobotName,
          description: "E2E robot that creates before deleting",
          playbook: "you create first, then delete later",
          model: `mock/../robot/scripts/${scriptName}`,
          tools: ["robot_create", "robot_delete"],
        });
      });

      await goToNewChat(page);

      await switchToRobot(page, actorRobotName);

      await sendMessage(page, "create a robot");
      await expect(
        page
          .locator("summary")
          .filter({ hasText: "Robot Create" })
          .filter({ hasText: "Tool complete" }),
      ).toBeVisible({ timeout: 15000 });
      await expect(page.getByText("Created Robot.")).toBeVisible({
        timeout: 15000,
      });

      await sendMessage(page, "delete it");

      await expect(
        page
          .locator("summary")
          .filter({ hasText: "Robot Delete" })
          .filter({ hasText: "Needs approval" }),
      ).toBeVisible({ timeout: 15000 });
      await expect(page.getByRole("button", { name: "Approve" })).toBeVisible();

      await page.getByRole("button", { name: "Approve" }).click();

      await expect(
        page
          .locator("summary")
          .filter({ hasText: "Robot Delete" })
          .filter({ hasText: "Tool complete" }),
      ).toBeVisible({ timeout: 15000 });
      await expect(page.getByText("Delete flow finished.")).toBeVisible({
        timeout: 15000,
      });
      await expect(page.getByText("No tool invocation found")).toHaveCount(0);

      await withAdminAccessKey(async ({ robotGet }) => {
        let deleted = false;
        try {
          await robotGet(victimRobotID);
        } catch {
          deleted = true;
        }
        expect(deleted).toBe(true);
      });
    } finally {
      await unlink(scriptPath).catch(() => undefined);
    }
  });

  test("robot_delete confirms multiple deletes from the same assistant turn", async ({
    page,
  }) => {
    const suffix = Date.now();
    const firstVictimName = `e2e-multi-delete-first-${suffix}`;
    const secondVictimName = `e2e-multi-delete-second-${suffix}`;
    const actorRobotName = `e2e-multi-delete-actor-${suffix}`;
    const scriptName = `e2e-multi-delete-confirmation-${suffix}.yaml`;
    const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;
    let firstVictimID = "";
    let secondVictimID = "";

    try {
      await withAdminAccessKey(async ({ robotCreate }) => {
        const firstVictim = await robotCreate({
          name: firstVictimName,
          description: "First E2E victim robot for multi-confirmation flow",
          playbook:
            "you are the first robot that should only be deleted after approval",
          model: DEFAULT_ROBOT_MODEL,
        });
        firstVictimID = firstVictim.id;

        const secondVictim = await robotCreate({
          name: secondVictimName,
          description: "Second E2E victim robot for multi-confirmation flow",
          playbook:
            "you are the second robot that should only be deleted after approval",
          model: DEFAULT_ROBOT_MODEL,
        });
        secondVictimID = secondVictim.id;

        await writeFile(
          scriptPath,
          `steps:
  - match:
      contains: "delete both victims"
    respond:
      tool_calls:
        - id: call_delete_1
          name: robot_delete
          args:
            id: "${firstVictimID}"
        - id: call_delete_2
          name: robot_delete
          args:
            id: "${secondVictimID}"
  - match:
      tool_result: robot_delete
    respond:
      text: "Both delete flows finished."
      finish: "stop"
`,
        );

        await robotCreate({
          name: actorRobotName,
          description: "E2E robot that requests two robot_delete calls",
          playbook: "you request multiple robot deletions when asked",
          model: `mock/../robot/scripts/${scriptName}`,
          tools: ["robot_delete"],
        });
      });

      await goToNewChat(page);

      await switchToRobot(page, actorRobotName);
      await expect(
        page.getByRole("button", { name: actorRobotName }),
      ).toBeVisible({ timeout: 15000 });

      await sendMessage(page, "delete both victims");

      const confirmationBatch = page.getByRole("group", {
        name: "Confirmation batch",
      });
      await expect(confirmationBatch).toContainText(
        "Approve these 2 actions?",
        {
          timeout: 15000,
        },
      );
      await expect(
        confirmationBatch.getByRole("button", {
          name: "Approve all confirmations",
        }),
      ).toBeVisible();
      await expect(
        confirmationBatch.getByRole("button", {
          name: "Deny all confirmations",
        }),
      ).toBeVisible();

      await expect(
        confirmationBatch.getByRole("button", { name: /^Approve Delete / }),
      ).toHaveCount(2);
      await expect(
        confirmationBatch.getByRole("button", { name: /^Deny Delete / }),
      ).toHaveCount(2);

      await confirmationBatch
        .getByRole("button", { name: "Approve all confirmations" })
        .click();

      await expect(confirmationBatch).toContainText("2 actions resolved", {
        timeout: 15000,
      });
      await expect(page.getByText("Both delete flows finished.")).toBeVisible({
        timeout: 15000,
      });

      await withAdminAccessKey(async ({ robotGet }) => {
        await expect(robotGet(firstVictimID)).rejects.toThrow();
        await expect(robotGet(secondVictimID)).rejects.toThrow();
      });
    } finally {
      await unlink(scriptPath).catch(() => undefined);
    }
  });
});
