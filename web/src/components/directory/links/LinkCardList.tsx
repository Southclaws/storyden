import { map } from "lodash/fp";

import { Link } from "src/api/openapi/schemas";
import { CardGrid, CardItem, CardRows } from "src/theme/components/Card";

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
