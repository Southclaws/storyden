import { LinkReference } from "src/api/openapi-schema";

import { CardGrid, CardRows } from "@/components/ui/rich-card";
import { RichCardVariantProps } from "@/styled-system/recipes";

import { LinkCard } from "./LinkCard";

type Props = {
  links: LinkReference[];
  size?: RichCardVariantProps["size"];
};

export function LinkCardRows({ links, size }: Props) {
  return (
    <CardRows>
      {links.map((l) => (
        <LinkCard key={l.slug} shape="row" size={size} link={l} />
      ))}
    </CardRows>
  );
}

export function LinkCardGrid({ links, size }: Props) {
  return (
    <CardGrid>
      {links.map((l) => (
        <LinkCard key={l.slug} shape="row" size={size} link={l} />
      ))}
    </CardGrid>
  );
}
