import type { Notification } from "@/api/openapi-schema/notification";

import { getCommonProperties } from "./datagraph/item";

function getNotificationContent(n: Notification) {
  const p = n.item && getCommonProperties(n.item);
  const postURL = p?.id
    ? `/t/locate/${encodeURIComponent(p.id)}`
    : "/notifications";

  switch (n.event) {
    case "trail_run_attention":
      return {
        description: "A Trail run needs your attention",
        url: n.target
          ? `/robots/trails?run=${encodeURIComponent(n.target)}`
          : "/robots/trails/",
      };

    case "thread_reply":
      return {
        description: "replied to your post",
        url: postURL,
      };

    case "reply_to_reply":
      return {
        description: "replied to you",
        url: postURL,
      };

    case "post_like":
      return {
        description: "liked your post",
        url: postURL,
      };

    case "follow":
      return {
        description: "followed you",
        url: n.source
          ? `/m/${encodeURIComponent(n.source.handle)}`
          : "/notifications",
      };

    case "profile_mention":
      return {
        description: "mentioned you",
        url: postURL,
      };

    case "event_host_added":
      return {
        description: "added you as an event host",
        url: "/notifications",
      };

    case "member_attending_event":
      return {
        description: "is attending your event",
        url: "/notifications",
      };

    case "member_declined_event":
      return {
        description: "declined your event",
        url: "/notifications",
      };

    case "attendee_removed":
      return {
        description: "removed you from their event",
        url: "/notifications",
      };

    case "report_submitted":
      return {
        description: "submitted a report",
        url: `/reports`,
      };

    case "report_updated":
      return {
        description: "report status updated",
        url: `/reports`,
      };

    case "warning_issued":
      return {
        description: "issued you a warning",
        url: p?.slug ? `/m/${p.slug}` : "/notifications",
      };

    case "node_version_created":
      return {
        description: "created a draft for your page",
        url: p?.slug ? `/l/${p.slug}` : "/notifications",
      };

    case "node_version_applied":
      return {
        description: "applied your draft",
        url: p?.slug ? `/l/${p.slug}` : "/notifications",
      };

    case "node_version_deleted":
      return {
        description: "discarded your draft",
        url: p?.slug ? `/l/${p.slug}` : "/notifications",
      };

    default:
      return {
        description: "You have a new notification",
        url: "/notifications",
      };
  }
}

export function formatNotification(n: Notification) {
  const content = getNotificationContent(n);
  const properties = n.item && getCommonProperties(n.item);
  const itemName = n.item?.kind === "profile" ? undefined : properties?.name;
  const body = itemName
    ? `${content.description} "${itemName}"`
    : content.description;

  return {
    title: n.source?.handle ?? "System",
    description: content.description,
    body,
    url: content.url,
  };
}

export function formatNotificationFallback() {
  return {
    title: "System",
    body: "You have a new notification",
    url: "/notifications",
  };
}
