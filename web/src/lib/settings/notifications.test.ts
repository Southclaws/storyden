import { describe, expect, it } from "vitest";

import { getNotificationGroupState } from "./notifications";

const events = {
  thread_reply: { inbox: true, push: false },
  reply_to_reply: { inbox: true, push: true },
};

describe("getNotificationGroupState", () => {
  it("resolves each delivery channel independently", () => {
    expect(
      getNotificationGroupState(
        events,
        ["thread_reply", "reply_to_reply"],
        "inbox",
      ),
    ).toBe(true);
    expect(
      getNotificationGroupState(
        events,
        ["thread_reply", "reply_to_reply"],
        "push",
      ),
    ).toBe("mixed");
  });

  it("disables the group only when all its events are disabled", () => {
    expect(
      getNotificationGroupState(
        {
          thread_reply: { inbox: false, push: false },
          reply_to_reply: { inbox: false, push: false },
        },
        ["thread_reply", "reply_to_reply"],
        "inbox",
      ),
    ).toBe(false);
  });

  it("does not invent defaults for missing events", () => {
    expect(
      getNotificationGroupState(
        {},
        ["thread_reply", "reply_to_reply"],
        "inbox",
      ),
    ).toBe(false);
  });
});
