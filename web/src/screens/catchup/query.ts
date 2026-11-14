import type { ThreadListParams } from "@/api/openapi-schema";

export const catchupViews = ["unread", "unresolved"] as const;
export type CatchupView = (typeof catchupViews)[number];

export function getCatchupThreadParams(
  page: number,
  view: CatchupView,
  category: string,
): ThreadListParams {
  const params: ThreadListParams = {
    page: page.toString(),
    ignore_pinned: true,
  };

  if (view === "unresolved") {
    params.lifecycle_state = "open";
  } else {
    params.read_status = "unread";
  }

  if (category) {
    params.categories = [category];
  }

  return params;
}
