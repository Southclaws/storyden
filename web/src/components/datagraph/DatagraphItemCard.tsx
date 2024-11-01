import { DatagraphItem, DatagraphItemKind } from "@/api/openapi-schema";
import { HStack } from "@/styled-system/jsx";

import { Badge } from "../ui/badge";
import { Card } from "../ui/rich-card";

type Props = {
  item: DatagraphItem;
};

export function DatagraphItemCard({ item }: Props) {
  const url = buildPermalink(item);

  return (
    <Card
      id={item.id}
      url={url}
      title={item.name}
      text={item.description}
      controls={
        <HStack>
          {/* TODO: We need more info for datagraph items on the API. */}
          {/* <MemberBadge profile={item.owner} /> */}

          <DatagraphItemBadge item={item} />
        </HStack>
      }
    />
  );
}

export function DatagraphItemBadge({ item }: Props) {
  const label = getDatagraphKindLabel(item.kind);
  return <Badge>{label}</Badge>;
}

function buildPermalink(d: DatagraphItem): string {
  switch (d.kind) {
    case "post":
      return `/t/${d.slug}`;
    case "profile":
      return `/t/${d.slug}`;
    case "node":
      return `/l/${d.slug}`;
  }
}

function getDatagraphKindLabel(kind: DatagraphItemKind): string {
  switch (kind) {
    case "post":
      return "Post";
    case "profile":
      return "Profile";
    case "node":
      return "Library";
  }
}
