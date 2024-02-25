import { map } from "lodash/fp";

import { Link } from "src/api/openapi/schemas";
import { CardItem, CardRows } from "src/theme/components/Card";

const toCardItems = map<Link, CardItem>((l) => ({
  title: l.title || l.url,
  text: l.description,
  image: l.assets[0]?.url,
  url: l.url,
}));

export function LinkListRows({ links }: { links: Link[] }) {
  const items = toCardItems(links);
  return <CardRows items={items} />;
}
