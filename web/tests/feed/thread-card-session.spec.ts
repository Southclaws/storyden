import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  registerUser,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

for (const surface of ["home", "discussion", "category", "search"] as const) {
  test(`${surface} cards receive the server session before client account revalidation`, async ({
    page,
    browser,
  }) => {
    const category =
      surface === "category"
        ? await withAdminAccessKey(({ categoryCreate }) =>
            categoryCreate({
              name: `Card session ${randomUUID().slice(0, 8)}`,
              description: "A category with server-rendered card controls.",
              colour: "#123456",
            }),
          )
        : undefined;
    const thread = await withAdminAccessKey(({ threadCreate }) =>
      threadCreate({
        title: `Server session ${randomUUID().slice(0, 8)}`,
        body: "A thread whose controls are ready in the server response.",
        visibility: "published",
        category: category?.id,
      }),
    );
    await registerUser(page, "card-session", "TestPassword123!");
    const cookies = await page.context().cookies();
    const path =
      surface === "home"
        ? "/"
        : surface === "discussion"
          ? "/d"
          : surface === "category"
            ? `/d/${category!.slug}`
            : `/t?q=${encodeURIComponent(thread.title)}`;
    const context = await browser.newContext({ javaScriptEnabled: false });
    await context.addCookies(cookies);
    try {
      const serverPage = await context.newPage();
      await serverPage.goto(path);
      const serverCard = serverPage.locator(`[id="${thread.id}"]`);
      // Some routes stream through loading.tsx. Inspect their server markup
      // even when a streamed segment still needs JavaScript to reveal it.
      await expect(
        serverCard.getByRole("button", {
          name: "Like",
          exact: true,
          includeHidden: true,
        }),
      ).toBeEnabled();
      await expect(
        serverCard.getByRole("button", {
          name: "Bookmark",
          exact: true,
          includeHidden: true,
        }),
      ).toHaveCount(1);

      let releaseAccount!: () => void;
      const accountGate = new Promise<void>((resolve) => {
        releaseAccount = resolve;
      });
      await page.route("**/api/accounts", async (route) => {
        await accountGate;
        await route.continue();
      });
      try {
        await page.goto(path, { waitUntil: "domcontentloaded" });
        const card = page.locator(`[id="${thread.id}"]`);
        await expect(
          card.getByRole("button", { name: "Like", exact: true }),
        ).toBeEnabled();
        await expect(
          card.getByRole("button", { name: "Bookmark", exact: true }),
        ).toBeVisible();
        const like = card.getByRole("button", { name: "Like", exact: true });
        await like.click();
        await expect(
          card.getByRole("button", { name: "Unlike", exact: true }),
        ).toBeEnabled();
      } finally {
        releaseAccount();
        await page.unrouteAll({ behavior: "wait" });
      }
      await expect(
        page
          .locator(`[id="${thread.id}"]`)
          .getByRole("button", { name: "Unlike", exact: true }),
      ).toBeEnabled();
    } finally {
      await context.close();
    }
  });
}
