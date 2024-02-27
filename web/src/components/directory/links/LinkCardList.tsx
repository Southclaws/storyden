import { Link } from "src/api/openapi/schemas";
import { CardGrid, CardRows } from "src/theme/components/Card";

import { CardVariantProps } from "@/styled-system/recipes";

import { LinkCard } from "./LinkCard";

type Props = {
  links: Link[];
  size?: CardVariantProps["size"];
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
