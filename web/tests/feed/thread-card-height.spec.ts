import { expect, test } from "@playwright/test";
import { randomUUID } from "node:crypto";

import {
  registerUser,
  withAdminAccessKey,
} from "../access_key_admin_assignment";

test("the bookmark wrapper keeps the same height as the other card controls", async ({
  page,
  browser,
}) => {
  const thread = await withAdminAccessKey(({ threadCreate }) =>
    threadCreate({
      title: `Card height ${randomUUID().slice(0, 8)}`,
      body: "A thread with compact controls.",
      visibility: "published",
    }),
  );
  await registerUser(page, "card-height", "TestPassword123!");
  const card = page.locator(`[id="${thread.id}"]`);
  const bookmark = card.getByRole("button", { name: "Bookmark", exact: true });
  await expect(bookmark).toBeVisible();
  const likeHeight = await card
    .getByRole("button", { name: `Like ${thread.title}`, exact: true })
    .evaluate((element) => element.getBoundingClientRect().height);
  const wrapperHeight = await bookmark
    .locator("..")
    .evaluate((element) => element.getBoundingClientRect().height);
  expect(wrapperHeight).toBe(likeHeight);
  await page.evaluate(() => document.fonts.ready);
  const signedInHeight = await card.evaluate(
    (element) => element.getBoundingClientRect().height,
  );
  const anonymous = await browser.newContext();
  try {
    const anonymousPage = await anonymous.newPage();
    await anonymousPage.goto("/");
    await anonymousPage.evaluate(() => document.fonts.ready);
    const anonymousCard = anonymousPage.locator(`[id="${thread.id}"]`);
    await expect(
      anonymousCard.getByRole("button", {
        name: `Like ${thread.title}`,
        exact: true,
      }),
    ).toHaveCount(0);
    await expect(
      anonymousCard.getByText("Sign in to like discussions.", { exact: false }),
    ).toBeAttached();
    await expect(
      anonymousCard.getByRole("button", { name: "Bookmark", exact: true }),
    ).toHaveCount(0);
    const signedOutHeight = await anonymousCard.evaluate(
      (element) => element.getBoundingClientRect().height,
    );
    expect(signedOutHeight).toBe(signedInHeight);
  } finally {
    await anonymous.close();
  }
});
