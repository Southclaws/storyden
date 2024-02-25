import { map } from "lodash/fp";

import { Link, LinkListResult } from "src/api/openapi/schemas";
import { LinkCard } from "src/components/directory/links/LinkCard";
import { Empty } from "src/components/site/Empty";
import { CardGrid, CardItem, CardRows } from "src/theme/components/Card";

import { styled } from "@/styled-system/jsx";

type Props = {
  links: LinkListResult;
  show?: number;
};

export function LinkCardList({ links, show }: Props) {
  if (links.links.length === 0) {
    return <Empty>no links were found</Empty>;
  }

  const shown = show ? links.links.slice(0, show) : links.links;

  return (
    <styled.ol w="full" display="flex" flexDir="column" gap="4">
      {shown.map((v) => (
        <styled.li key={v.url}>
          <LinkCard {...v} />
        </styled.li>
      ))}
    </styled.ol>
  );
}

const toCardItems = map<Link, CardItem>((l) => ({
  id: l.slug,
  title: l.title || l.url,
  text: l.description,
  image: l.assets[0]?.url,
  url: l.url,
}));

export function LinkCardRows({ links }: { links: Link[] }) {
  const items = toCardItems(links);
  return <CardRows items={items} />;
}

export function LinkCardGrid({ links }: { links: Link[] }) {
  const items = toCardItems(links);
  return <CardGrid items={items} />;
}
