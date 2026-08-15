import { Page, expect, test } from "@playwright/test";
import { unlink, writeFile } from "node:fs/promises";

import {
  DEFAULT_ROBOT_MODEL,
  goToNewChat,
  sendMessage,
  setupRobotProviderWithScript,
  waitForPersistedChatRoute,
} from "./helpers";

const ROBOT_SCRIPT_DIR = "../tests/robot/scripts";

async function getRobotSessionIDs(page: Page) {
  const response = await page.request.get(
    "http://localhost:8001/api/robots/sessions",
  );
  if (!response.ok()) {
    throw new Error(`failed to list robot sessions: ${response.status()}`);
  }

  const payload = (await response.json()) as {
    sessions?: { id: string }[];
  };
  return payload.sessions?.map((session) => session.id) ?? [];
}

async function waitForPersistedRobotText(
  page: Page,
  sessionID: string,
  text: string,
) {
  await expect
    .poll(
      async () => {
        const response = await page.request.get(
          `http://localhost:8001/api/robots/sessions/${sessionID}`,
        );
        if (!response.ok()) return false;

        const payload = (await response.json()) as {
          message_list?: {
            messages?: { parts?: { text?: string }[] }[];
          };
        };
        return Boolean(
          payload.message_list?.messages?.some((message) =>
            message.parts?.some((part) => part.text === text),
          ),
        );
      },
      { timeout: 15000 },
    )
    .toBe(true);
}

// The mock LLM uses ../tests/robot/scripts/e2e-default.yaml, selected through
// the Robots provider settings API in the suite setup.
// It routes:
//   - any message            → static text response
//   - "list pages"           → library_page_list tool call then text
//   - "trigger error"        → LLM-level error (no text response)

test.describe("Robot Chat — mock LLM stream", () => {
  test.beforeAll(async () => {
    await setupRobotProviderWithScript();
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

  test("an existing session ignores a stale root Robot query", async ({
    page,
  }) => {
    await sendMessage(page, "start the session");
    await expect(
      page.getByText("This is a mock response from the test provider."),
    ).toBeVisible({ timeout: 15000 });

    await expect(page).toHaveURL(
      /\/robots\/chats\/(?!new(?:[/?#]|$))[^/?#]+$/,
      { timeout: 15000 },
    );
    const sessionURL = new URL(page.url());

    await page.goto(`${sessionURL.pathname}?robot=robot_builder`);

    await sendMessage(page, "continue the existing session");

    await expect(
      page.getByText(
        "robotId can only select the root Robot when starting a session",
      ),
    ).toHaveCount(0);
    await expect(
      page.getByText("This is a mock response from the test provider."),
    ).toHaveCount(2, { timeout: 15000 });
  });

  test("tool call renders with name and completion status", async ({
    page,
  }) => {
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

  test("LLM error shows error message in the UI", async ({ page }) => {
    await sendMessage(page, "trigger error");

    await expect(
      page.locator("p").filter({ hasText: "simulated LLM provider failure" }),
    ).toBeVisible({ timeout: 15000 });
  });

  test("model thought is hidden until expanded", async ({ page }) => {
    await sendMessage(page, "show thought");

    await expect(page.getByText("The answer is ready.")).toBeVisible({
      timeout: 15000,
    });

    const thought = page.getByText(
      "I should verify the evidence before answering.",
    );
    await expect(thought).toBeHidden();

    await page.locator("summary").filter({ hasText: "Thought" }).click();
    await expect(thought).toBeVisible();
  });

  test("a new chat updates its URL before the first turn completes", async ({
    page,
  }) => {
    const suffix = Date.now();
    const scriptName = `e2e-session-route-${suffix}.yaml`;
    const scriptPath = `${ROBOT_SCRIPT_DIR}/${scriptName}`;

    try {
      await writeFile(
        scriptPath,
        `steps:
  - match:
      tool_result: library_page_list
    respond:
      delay_ms: 5000
      text: "The delayed first turn completed."
      finish: "stop"
  - match:
      contains: "route before completion"
    respond:
      tool_calls:
        - id: call_route_list
          name: library_page_list
          args: {}
`,
      );

      await setupRobotProviderWithScript(`mock/../robot/scripts/${scriptName}`);
      await page.goto("/robots/chats/new");

      await sendMessage(page, "route before completion");
      await waitForPersistedChatRoute(page);
      await expect(page.getByText("Denbot is responding...")).toBeVisible();
      await expect(
        page.getByText("The delayed first turn completed."),
      ).toBeVisible({ timeout: 15000 });
    } finally {
      await unlink(scriptPath).catch(() => undefined);
      await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
    }
  });

  test("a detached response completes before a later message is sent", async ({
    page,
  }) => {
    const suffix = Date.now();
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
      text: "This detached response completed in the background."
      finish: "stop"
  - match:
      contains: "after detachment"
    respond:
      text: "Follow-up after detachment rendered."
      finish: "stop"
`,
      );

      await setupRobotProviderWithScript(`mock/../robot/scripts/${scriptName}`);

      await page.goto("/robots/chats/new");
      const existingSessionIDs = new Set(await getRobotSessionIDs(page));

      await sendMessage(page, "slow response");
      await expect(page.getByText("Denbot is responding...")).toBeVisible({
        timeout: 15000,
      });

      let sessionID: string | undefined;
      await expect
        .poll(
          async () => {
            const sessionIDs = await getRobotSessionIDs(page);
            sessionID = sessionIDs.find((id) => !existingSessionIDs.has(id));
            return sessionID;
          },
          { timeout: 15000 },
        )
        .toBeTruthy();
      if (!sessionID) throw new Error("new robot session was not created");

      await page.goto("/");
      await waitForPersistedRobotText(
        page,
        sessionID,
        "This detached response completed in the background.",
      );

      await page.goto(`/robots/chats/${sessionID}`);
      await expect(page.getByText("slow response")).toBeVisible({
        timeout: 15000,
      });
      await expect(
        page.getByText("This detached response completed in the background."),
      ).toBeVisible({ timeout: 15000 });

      await sendMessage(page, "after detachment");
      await expect(
        page.getByText("Follow-up after detachment rendered."),
      ).toBeVisible({ timeout: 15000 });

      await page.reload();

      await expect(page.getByText("slow response")).toBeVisible({
        timeout: 15000,
      });
      await expect(
        page.getByText("after detachment", { exact: true }),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        page.getByText("Follow-up after detachment rendered."),
      ).toBeVisible({ timeout: 15000 });
      await expect(
        page.getByText("This detached response completed in the background."),
      ).toBeVisible({ timeout: 15000 });
    } finally {
      await unlink(scriptPath).catch(() => undefined);
      await setupRobotProviderWithScript(DEFAULT_ROBOT_MODEL);
    }
  });
});
