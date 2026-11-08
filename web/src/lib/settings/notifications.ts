import {
  type Account,
  NotificationEvent,
  type NotificationSettings,
  Permission,
} from "@/api/openapi-schema";
import { hasPermission } from "@/utils/permissions";

import type { Settings } from "./settings";

type NotificationGroup = {
  id: string;
  label: string;
  description: string;
  events: NotificationEvent[];
};

export function getNotificationGroups(
  account: Account,
  settings: Settings,
): NotificationGroup[] {
  const groups: NotificationGroup[] = [
    {
      id: "replies",
      label: "Replies",
      description: "Replies to your threads and replies to you.",
      events: [
        NotificationEvent.thread_reply,
        NotificationEvent.reply_to_reply,
      ],
    },
  ];

  const usesLibrary = [
    Permission.MANAGE_LIBRARY,
    Permission.SUBMIT_LIBRARY_NODE,
    Permission.SUBMIT_LIBRARY_NODE_CHANGES,
  ].some((permission) => hasPermission(account, permission));

  if (usesLibrary) {
    groups.push({
      id: "library",
      label: "Library contributions",
      description: "Drafts submitted to your pages, applied, or discarded.",
      events: [
        NotificationEvent.node_version_created,
        NotificationEvent.node_version_applied,
        NotificationEvent.node_version_deleted,
      ],
    });
  }

  const usesTrails =
    settings.capabilities.includes("robots") &&
    hasPermission(account, Permission.MANAGE_TRAILS);

  if (usesTrails) {
    groups.push({
      id: "trails",
      label: "Trail alerts",
      description: "Trail runs that need your attention.",
      events: [NotificationEvent.trail_run_attention],
    });
  }

  groups.push({
    id: "moderation",
    label: "Reports and moderation",
    description: "Updates about reports and warnings.",
    events: [
      NotificationEvent.report_submitted,
      NotificationEvent.report_updated,
      NotificationEvent.warning_issued,
    ],
  });

  groups.push(
    {
      id: "new_threads",
      label: "All new threads",
      description: "New threads published across the community.",
      events: [NotificationEvent.global_thread_published],
    },
    {
      id: "all_replies",
      label: "All replies",
      description: "Replies to any thread, including those you haven't joined.",
      events: [NotificationEvent.global_reply_published],
    },
    {
      id: "all_reactions",
      label: "All reactions",
      description:
        "Emoji reactions to threads and replies across the community.",
      events: [NotificationEvent.global_post_reacted],
    },
  );

  return groups;
}

export function getNotificationGroupState(
  overrides: NotificationSettings["events"],
  events: NotificationEvent[],
  channel: "inbox" | "push",
): boolean | "mixed" {
  const enabledCount = events.filter(
    (event) => overrides[event]?.[channel] === true,
  ).length;

  if (enabledCount === events.length) {
    return true;
  }

  if (enabledCount === 0) {
    return false;
  }

  return "mixed";
}
