import { LinkListResult } from "src/api/openapi-schema";

import { EmptyState } from "@/components/site/EmptyState";

import { LinkCardRows } from "../LinkCardList";

type Props = {
  links: LinkListResult;
  show?: number;
};

export function LinkResultList({ links, show }: Props) {
  if (links.links.length === 0) {
    return <EmptyState hideContributionLabel>No links were found.</EmptyState>;
  }

  const shown = show ? links.links.slice(0, show) : links.links;

  return <LinkCardRows links={shown} />;
}
