import { LinkReference } from "src/api/openapi-schema";

import { CardGrid, CardRows } from "@/components/ui/rich-card";

import { LinkCard } from "./LinkCard";

type Props = {
  links: LinkReference[];
};

export function LinkCardRows({ links }: Props) {
  return (
    <CardRows>
      {links.map((l) => (
        <LinkCard key={l.slug} shape="row" link={l} />
      ))}
    </CardRows>
  );
}

export function LinkCardGrid({ links }: Props) {
  return (
    <CardGrid>
      {links.map((l) => (
        <LinkCard key={l.slug} shape="row" link={l} />
      ))}
    </CardGrid>
  );
}
