import chroma from "chroma-js";

import {
  DatagraphItem,
  DatagraphItemKind,
  DatagraphItemNode,
  DatagraphItemPost,
  DatagraphItemProfile,
  DatagraphItemReply,
  DatagraphItemThread,
} from "@/api/openapi-schema";
import { HStack, WStack } from "@/styled-system/jsx";
import { ColorPalette } from "@/styled-system/tokens";
import { getAssetURL } from "@/utils/asset";
import { htmlToMarkdown } from "@/utils/markdown";

import { MemberBadge } from "../member/MemberBadge/MemberBadge";
import { Timestamp } from "../site/Timestamp";
import { Badge } from "../ui/badge";
import { Card } from "../ui/rich-card";

type Props = {
  item: DatagraphItem;
};

export function DatagraphItemCard({ item }: Props) {
  switch (item.kind) {
    case DatagraphItemKind.post:
      return <DatagraphItemPostGenericCard item={item} />;

    case DatagraphItemKind.thread:
      return <DatagraphItemPostGenericCard item={item} />;

    case DatagraphItemKind.reply:
      return <DatagraphItemReplyCard item={item} />;

    case DatagraphItemKind.node:
      return <DatagraphItemNodeCard item={item} />;

    // case DatagraphItemKind.collection:
    //   return null;

    case DatagraphItemKind.profile:
      return <DatagraphItemProfileCard item={item} />;

    // case DatagraphItemKind.event:
    //   return null;
  }
}

export function DatagraphItemPostGenericCard({
  item,
}: {
  item: DatagraphItemPost | DatagraphItemThread;
}) {
  const { ref } = item;
  const url = `/t/${ref.slug}`;

  return (
    <Card
      id={ref.id}
      url={url}
      title={ref.title || "(untitled post)"}
      text={ref.description ?? htmlToMarkdown(ref.body)}
      controls={
        <WStack>
          <HStack gap="1" minWidth="0" color="fg.subtle">
            <MemberBadge
              profile={ref.author}
              size="sm"
              name="full-horizontal"
            />
            <Timestamp created={ref.createdAt} />
          </HStack>

          <DatagraphItemBadge kind={item.kind} />
        </WStack>
      }
    />
  );
}

export function DatagraphItemReplyCard({ item }: { item: DatagraphItemReply }) {
  const { ref } = item;
  const url = `/t/locate/${ref.id}`;

  return (
    <Card
      id={ref.id}
      url={url}
      title={ref.title ? `Thread: ${ref.title}` : "(untitled thread)"}
      text={ref.description ?? htmlToMarkdown(ref.body)}
      controls={
        <WStack>
          <HStack gap="1" minWidth="0" color="fg.subtle">
            <MemberBadge
              profile={ref.author}
              size="sm"
              name="full-horizontal"
            />
            <Timestamp created={ref.createdAt} />
          </HStack>

          <DatagraphItemBadge kind={item.kind} />
        </WStack>
      }
    />
  );
}

export function DatagraphItemNodeCard({ item }: { item: DatagraphItemNode }) {
  const { ref } = item;
  const url = `/l/${ref.slug}`;

  return (
    <Card
      id={ref.id}
      url={url}
      title={ref.name}
      text={ref.description}
      image={getAssetURL(ref.primary_image?.path)}
      controls={
        <WStack>
          <HStack gap="1" minWidth="0" color="fg.subtle">
            <MemberBadge profile={ref.owner} size="sm" name="full-horizontal" />
            <Timestamp created={ref.createdAt} />
          </HStack>

          <DatagraphItemBadge kind={item.kind} />
        </WStack>
      }
    />
  );
}

export function DatagraphItemProfileCard({
  item,
}: {
  item: DatagraphItemProfile;
}) {
  const { ref } = item;
  const url = `/m/${ref.handle}`;

  return (
    <Card
      id={ref.id}
      url={url}
      title={ref.name}
      text={ref.bio}
      controls={
        <WStack>
          <HStack gap="1" minWidth="0" color="fg.subtle">
            <MemberBadge profile={ref} size="sm" name="full-horizontal" />
            <Timestamp created={ref.createdAt} />
          </HStack>

          <DatagraphItemBadge kind={item.kind} />
        </WStack>
      }
    />
  );
}

export function DatagraphItemBadge({ kind }: { kind: DatagraphItemKind }) {
  const label = getDatagraphKindLabel(kind);
  const colour = getDatagraphKindColour(kind);

  const cssVars = badgeColourCSS(colour);

  return (
    <Badge
      style={cssVars}
      backgroundColor="var(--colors-color-palette-bg)"
      borderColor="var(--colors-color-palette-bo)"
      color="var(--colors-color-palette-fg)"
    >
      {label}
    </Badge>
  );
}

export function badgeColourCSS(c: string) {
  const { bg, bo, fg } = badgeColours(c);

  return {
    "--colors-color-palette-fg": fg,
    "--colors-color-palette-bo": bo,
    "--colors-color-palette-bg": bg,
  } as React.CSSProperties;
}

export function badgeColours(c: string) {
  const colour = chroma(c);

  const bg = colour.luminance(0.8).css();
  const bo = colour.luminance(0.6).saturate(1.3).css();
  const fg = colour.darken(1.5).saturate(2).css();

  return { bg, bo, fg };
}

export function getDatagraphKindLabel(kind: DatagraphItemKind): string {
  switch (kind) {
    case DatagraphItemKind.post:
      return "Post";
    case DatagraphItemKind.thread:
      return "Thread";
    case DatagraphItemKind.reply:
      return "Reply";
    case DatagraphItemKind.node:
      return "Library";
    case DatagraphItemKind.collection:
      return "Collection";
    case DatagraphItemKind.profile:
      return "Profile";
    case DatagraphItemKind.event:
      return "Event";
  }
}

export function getDatagraphKindColour(kind: DatagraphItemKind): ColorPalette {
  switch (kind) {
    case DatagraphItemKind.post:
      return "pink";
    case DatagraphItemKind.thread:
      return "pink";
    case DatagraphItemKind.reply:
      return "pink";
    case DatagraphItemKind.node:
      return "green";
    case DatagraphItemKind.collection:
      return "blue";
    case DatagraphItemKind.profile:
      return "red";
    case DatagraphItemKind.event:
      return "amber";
  }
}
