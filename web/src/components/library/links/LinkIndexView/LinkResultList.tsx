import { LinkListResult } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";

import { LinkCardGrid, LinkCardRows } from "../LinkCardList";

type Props = {
  links: LinkListResult;
  show?: number;
  layout?: "rows" | "grid";
};

export function LinkResultList({ links, show, layout = "rows" }: Props) {
  if (links.links.length === 0) {
    return <EmptyState hideContributionLabel>No links were found.</EmptyState>;
  }

  const shown = show ? links.links.slice(0, show) : links.links;

  if (layout === "grid") {
    return <LinkCardGrid links={shown} />;
  }

  return <LinkCardRows links={shown} />;
}
