import { expect, test } from "@playwright/test";
import { unlink, writeFile } from "node:fs/promises";

import { login, withAdminAccessKey } from "../access_key_admin_assignment";
import { setupRobotProviderWithScript } from "../robot/helpers";

const ADMIN_USERNAME = "e2e_admin";
const ADMIN_PASSWORD = "E2EAdminPassword123!";
const ROBOT_SCRIPT_DIR = "../tests/robot/scripts";

test.describe("Trails", () => {
  const suffix = Date.now();
  const robotName = `Trail E2E Robot ${suffix}`;
  const scriptName = `e2e-trails-${suffix}.yaml`;
  const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;
  let robotID = "";

  test.beforeAll(async () => {
    await writeFile(
      scriptPath,
      `steps:
  - match:
      contains: "complete trail action"
    respond:
      tool_calls:
        - id: call_trail_complete
          name: robot_run_finish
          args:
            status: completed
            summary: "The scheduled community prompt was prepared."
  - match:
      contains: "block trail action"
    respond:
      tool_calls:
        - id: call_trail_blocked
          name: robot_run_finish
          args:
            status: blocked
            summary: "The moderation action needs approval."
            attention:
              reason: needs_approval
              message: "Approve the moderation action before continuing."
`,
    );

    await setupRobotProviderWithScript();
    await withAdminAccessKey(async ({ robotCreate }) => {
      const robot = await robotCreate({
        name: robotName,
        description: "Exercises completed and blocked Trail action results.",
        playbook:
          "Follow the unattended instruction and finish with the structured result tool.",
        model: `mock/../robot/scripts/${scriptName}`,
        tools: [],
      });
      robotID = robot.id;
    });
  });

  test.afterAll(async () => {
    if (robotID) {
      await withAdminAccessKey(async ({ robotDelete }) => robotDelete(robotID));
    }
    await unlink(scriptPath).catch(() => undefined);
  });

  test("creates, previews, runs, records attention, and pauses a Trail", async ({
    page,
  }) => {
    const trailName = `Weekly community prompt ${suffix}`;

    await login(page, ADMIN_USERNAME, ADMIN_PASSWORD);
    await page.goto("/robots/trails/new");

    await page.getByLabel("Name").fill(trailName);
    await page
      .getByLabel("Description")
      .fill("A durable scheduled workflow exercised end to end.");
    const firstRobot = page.getByRole("combobox", {
      name: "Robot action 1",
    });
    await firstRobot.fill(robotName);
    await firstRobot.press("ArrowDown");
    await firstRobot.press("Enter");
    await page.getByRole("button", { name: "Preview next five" }).click();
    await expect(page.locator("form time").first()).toBeVisible();

    await page.getByRole("button", { name: "Add Robot action" }).click();
    const secondRobot = page.getByRole("combobox", {
      name: "Robot action 2",
    });
    await secondRobot.fill(robotName);
    await secondRobot.press("ArrowDown");
    await secondRobot.press("Enter");
    await page
      .getByLabel("Unattended instruction 1")
      .fill("complete trail action");
    await page
      .getByLabel("Unattended instruction 2")
      .fill("block trail action");
    await page.getByRole("button", { name: "Create Trail" }).click();
    await expect(page).toHaveURL(/\/robots\/trails\/(?!new$)[a-z0-9]+$/);

    const trailID = page.url().split("/").at(-1);
    expect(trailID).toBeTruthy();

    const robotsNavigation = page.getByRole("navigation", {
      name: "Robots navigation",
    });
    await expect(
      robotsNavigation.getByRole("link", { name: trailName }),
    ).toBeVisible();
    const trailResponse = await page.request.get(
      `http://localhost:8001/api/trails/${trailID}`,
    );
    expect(trailResponse.ok()).toBe(true);
    const createdTrail = (await trailResponse.json()) as {
      trigger: {
        type: string;
        schedule: {
          start: string;
          timezone: string;
          rule: {
            frequency: string;
            interval: number;
            by_weekday?: string[];
          };
        };
      };
    };
    expect(createdTrail.trigger).toEqual({
      type: "schedule",
      schedule: {
        start: expect.stringMatching(/^\d{4}-\d{2}-\d{2}T09:00:00$/),
        timezone: expect.any(String),
        rule: {
          frequency: "weekly",
          interval: 1,
          by_weekday: [expect.any(String)],
        },
      },
    });

    await page.getByRole("link", { name: "Edit" }).click();
    await expect(page).toHaveURL(`/robots/trails/${trailID}/edit`);
    await expect(
      page.getByRole("heading", { name: "Edit Trail" }),
    ).toBeVisible();
    await page.getByRole("button", { name: "Back" }).click();
    await expect(page).toHaveURL(`/robots/trails/${trailID}`);

    await page.getByRole("button", { name: "Run now" }).click();

    await expect(
      page.getByText("Needs attention", { exact: true }),
    ).toBeVisible({ timeout: 30_000 });
    await expect(page.getByText(/^Manual run at /)).toBeVisible();
    await expect(
      page.getByText("The scheduled community prompt was prepared."),
    ).toBeVisible();
    await expect(
      page.getByText("The moderation action needs approval."),
    ).toBeVisible();
    await expect(
      page.getByText("Needs approval", { exact: true }),
    ).toBeVisible();
    await expect(
      page.getByText("Approve the moderation action before continuing."),
    ).toBeVisible();

    const runsResponse = await page.request.get(
      `http://localhost:8001/api/trails/${trailID}/runs`,
    );
    expect(runsResponse.ok()).toBe(true);
    const runs = (await runsResponse.json()) as {
      runs: {
        id: string;
        actions: {
          target?: {
            type: "robot_run";
            robot_session_id: string;
            output?: { status: string; summary: string };
          };
        }[];
      }[];
    };
    await expect(
      page.getByText(runs.runs[0]?.id ?? "missing run", { exact: true }),
    ).toHaveCount(0);
    expect(
      runs.runs[0]?.actions.map((action) => action.target?.output?.summary),
    ).toEqual([
      "The scheduled community prompt was prepared.",
      "The moderation action needs approval.",
    ]);
    const sessionIDs = runs.runs[0]?.actions.flatMap((action) =>
      action.target === undefined ? [] : [action.target.robot_session_id],
    );
    expect(sessionIDs).toHaveLength(2);
    for (const sessionID of sessionIDs ?? []) {
      const sessionResponse = await page.request.get(
        `http://localhost:8001/api/robots/sessions/${sessionID}`,
      );
      expect(sessionResponse.ok()).toBe(true);
      const session = (await sessionResponse.json()) as { name: string };
      expect(session.name).toBe(`${trailName} (Run 1)`);
    }

    await page.getByRole("button", { name: "Pause" }).click();
    await expect(page.getByText("Paused", { exact: true })).toBeVisible();

    await expect
      .poll(
        async () => {
          const response = await page.request.get(
            "http://localhost:8001/api/notifications",
          );
          if (!response.ok()) return [];
          const body = (await response.json()) as {
            notifications?: { event?: string; target?: string }[];
          };
          return body.notifications ?? [];
        },
        { timeout: 15_000 },
      )
      .toContainEqual(
        expect.objectContaining({
          event: "trail_run_attention",
          target: expect.any(String),
        }),
      );

    const notificationsResponse = await page.request.get(
      "http://localhost:8001/api/notifications",
    );
    const notificationsBody = (await notificationsResponse.json()) as {
      notifications: { event: string; target?: string }[];
    };
    const target = notificationsBody.notifications.find(
      (item) => item.event === "trail_run_attention",
    )?.target;
    if (target === undefined) {
      throw new Error("Trail attention notification has no target");
    }

    await page.goto(`/robots/trails?run=${target}`);
    await expect(page).toHaveURL(`/robots/trails/${trailID}#run-${target}`);
    await expect(page.locator(`#run-${target}`)).toBeVisible();
  });
});
