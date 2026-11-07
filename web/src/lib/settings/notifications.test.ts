import { describe, expect, it } from "vitest";

import { getNotificationGroupState } from "./notifications";

describe("getNotificationGroupState", () => {
  it("enables events without overrides", () => {
    expect(
      getNotificationGroupState(undefined, ["thread_reply", "reply_to_reply"]),
    ).toBe(true);
  });

  it("shows mixed state for independently configured events", () => {
    expect(
      getNotificationGroupState({ thread_reply: false }, [
        "thread_reply",
        "reply_to_reply",
      ]),
    ).toBe("indeterminate");
  });

  it("disables the group only when all its events are disabled", () => {
    expect(
      getNotificationGroupState(
        { thread_reply: false, reply_to_reply: false },
        ["thread_reply", "reply_to_reply"],
      ),
    ).toBe(false);
  });
});
