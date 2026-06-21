import { Page, expect } from "@playwright/test";

import { login, withAdminAccessKey } from "../access_key_admin_assignment";

const ADMIN_USERNAME = "e2e_admin";
const ADMIN_PASSWORD = "E2EAdminPassword123!";
const MOCK_PROVIDER = "mock";

export const DEFAULT_ROBOT_MODEL = "mock/../robot/scripts/e2e-default.yaml";

export async function setupRobotProviderWithScript(
  modelRef = DEFAULT_ROBOT_MODEL,
) {
  await withAdminAccessKey(
    async ({ adminSettingsUpdate, robotProviderUpdate }) => {
      await robotProviderUpdate(MOCK_PROVIDER, {
        enabled: true,
      });
      await adminSettingsUpdate({
        services: {
          robots: {
            default_model: modelRef,
          },
        },
      });
    },
  );
}

export async function dismissOnboarding(page: Page) {
  const skipButton = page.getByRole("button", { name: "Skip" });
  if (await skipButton.isVisible({ timeout: 1000 }).catch(() => false)) {
    await skipButton.click();
  }
}

export async function goToNewChat(page: Page) {
  await login(page, ADMIN_USERNAME, ADMIN_PASSWORD);
  await dismissOnboarding(page);
  await page.goto("/robots/chats/new");
  await dismissOnboarding(page);
}

export async function sendMessage(page: Page, text: string) {
  const textarea = page.getByPlaceholder("Type a message...");
  const sendButton = page.getByRole("button", {
    name: "Send message",
    exact: true,
  });
  const respondingStatus = page
    .getByRole("status")
    .filter({ hasText: /is responding/ });

  await expect(respondingStatus).toHaveCount(0, { timeout: 15000 });
  await expect(textarea).toBeEnabled({ timeout: 15000 });

  await expect
    .poll(
      async () => {
        await textarea.fill(text);
        return textarea.inputValue();
      },
      { timeout: 15000 },
    )
    .toBe(text);

  await expect(sendButton).toBeEnabled({ timeout: 15000 });
  await sendButton.click();

  await expect(
    page.getByRole("article", { name: "You message" }).filter({
      hasText: text,
    }),
  ).toBeVisible({ timeout: 15000 });
}

export async function switchToRobot(page: Page, robotName: string) {
  await page.getByRole("button", { name: "Storyden Robot Builder" }).click();

  const menuItem = page.getByRole("menuitem", { name: robotName });
  await expect(menuItem).toBeVisible({ timeout: 15000 });
  await menuItem.click();

  await expect(page.getByRole("button", { name: robotName })).toBeVisible({
    timeout: 15000,
  });
}

export async function waitForPersistedChatRoute(page: Page) {
  await expect(page).toHaveURL(/\/robots\/chats\/(?!new(?:[/?#]|$))[^/?#]+/, {
    timeout: 15000,
  });
}
