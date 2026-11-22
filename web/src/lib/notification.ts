import type { Notification } from "@/api/openapi-schema/notification";

function getNotificationContent(n: Notification) {
  const postURL = n.item?.id
    ? `/t/locate/${encodeURIComponent(n.item.id)}`
    : "/notifications";

  switch (n.event) {
    case "global_thread_published":
      return {
        description: "started a new thread",
        url: postURL,
      };

    case "global_reply_published":
      return {
        description: "replied to a thread",
        url: postURL,
      };

    case "global_post_reacted":
      return {
        description: "reacted to a post",
        url: postURL,
      };

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
        url: n.item?.slug
          ? `/m/${encodeURIComponent(n.item.slug)}`
          : "/notifications",
      };

    case "node_version_created":
      return {
        description: "created a draft for your page",
        url: n.item?.slug
          ? `/l/${encodeURIComponent(n.item.slug)}`
          : "/notifications",
      };

    case "node_version_applied":
      return {
        description: "applied your draft",
        url: n.item?.slug
          ? `/l/${encodeURIComponent(n.item.slug)}`
          : "/notifications",
      };

    case "node_version_deleted":
      return {
        description: "discarded your draft",
        url: n.item?.slug
          ? `/l/${encodeURIComponent(n.item.slug)}`
          : "/notifications",
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
  const itemName = n.item?.kind === "profile" ? undefined : n.item?.title;
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
