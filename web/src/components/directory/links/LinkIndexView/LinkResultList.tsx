import { LinkListResult } from "src/api/openapi-schema";
import { Empty } from "src/components/site/Empty";

import { LinkCardRows } from "../LinkCardList";

type Props = {
  links: LinkListResult;
  show?: number;
};

export function LinkResultList({ links, show }: Props) {
  if (links.links.length === 0) {
    return <Empty>no links were found</Empty>;
  }

  const shown = show ? links.links.slice(0, show) : links.links;

  return <LinkCardRows links={shown} />;
}
