import { describe, expect, it } from "vitest";

import type { Notification } from "@/api/openapi-schema/notification";
import { NotificationEvent } from "@/api/openapi-schema/notificationEvent";

import { formatNotification } from "./notification";

function notification(event: Notification["event"]): Notification {
  return {
    id: "d4v1p9b2t53b3lo5e3cg",
    created_at: "2026-11-07T12:00:00Z",
    event,
    status: "unread",
  };
}

describe("shared notification formatting", () => {
  it.each(Object.values(NotificationEvent))(
    "formats %s without an attached item",
    (event) => {
      const result = formatNotification(notification(event));

      expect(result.body).not.toBe("");
      expect(result.url).not.toContain("undefined");
    },
  );

  it("encodes Trail run identifiers", () => {
    const result = formatNotification({
      ...notification("trail_run_attention"),
      target: "run & other=true",
    });

    expect(result).toEqual({
      title: "System",
      description: "A Trail run needs your attention",
      body: "A Trail run needs your attention",
      url: "/robots/trails?run=run%20%26%20other%3Dtrue",
    });
  });

  it.each([
    ["thread_reply", "replied to your post"],
    ["global_thread_published", "started a new thread"],
    ["global_reply_published", "replied to a thread"],
    ["global_post_reacted", "reacted to a post"],
  ] as const)(
    "formats %s with the post title and destination",
    (event, description) => {
      const source = {
        id: "member-one",
        joined: "2026-11-07T12:00:00Z",
        handle: "odin",
        name: "Odin",
        roles: [],
      };
      const result = formatNotification({
        ...notification(event),
        source,
        item: {
          kind: "post",
          ref: {
            id: "post-one",
            createdAt: "2026-11-07T12:00:00Z",
            updatedAt: "2026-11-07T12:00:00Z",
            title: "The little harbour is finally playable",
            description:
              "I got the boats moving! The lighting still needs work.",
            slug: "post-one-the-little-harbour-is-finally-playable",
            author: source,
            visibility: "published",
            assets: [],
            reacts: [],
            collections: { in_collections: 0, has_collected: false },
            likes: { likes: 0, liked: false },
            body: "I got the boats moving! The lighting still needs work.",
            body_links: [],
          },
        },
      });

      expect(result.body).toBe(
        `${description} "The little harbour is finally playable"`,
      );
      expect(result.url).toBe("/t/locate/post-one");
    },
  );
});
