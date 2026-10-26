import { type Locator, type Page, expect } from "@playwright/test";

/** Real two-finger input, without dispatching synthetic DOM events. */
export async function pinchOut(page: Page, target: Locator) {
  const bounds = await target.boundingBox();
  expect(bounds).not.toBeNull();
  // Keep both fingers away from floating tools at the top right.
  const x = bounds!.x + bounds!.width * 0.3;
  const y = bounds!.y + bounds!.height * 0.65;
  const session = await page.context().newCDPSession(page);
  try {
    await session.send("Emulation.setTouchEmulationEnabled", {
      enabled: true,
      maxTouchPoints: 2,
    });
    await session.send("Input.dispatchTouchEvent", {
      type: "touchStart",
      touchPoints: [
        { x: x - 20, y, id: 1 },
        { x: x + 20, y, id: 2 },
      ],
    });
    for (let distance = 24; distance <= 40; distance += 4) {
      await session.send("Input.dispatchTouchEvent", {
        type: "touchMove",
        touchPoints: [
          { x: x - distance, y, id: 1 },
          { x: x + distance, y, id: 2 },
        ],
      });
    }
    await session.send("Input.dispatchTouchEvent", {
      type: "touchEnd",
      touchPoints: [],
    });
  } finally {
    await session.send("Emulation.setTouchEmulationEnabled", {
      enabled: false,
    });
    await session.detach();
  }
}
