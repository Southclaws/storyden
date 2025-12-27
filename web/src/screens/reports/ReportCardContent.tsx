import Link from "next/link";

import {
  DatagraphItem,
  DatagraphItemKind,
  DatagraphItemNode,
  DatagraphItemPost,
  DatagraphItemProfile,
  DatagraphItemReply,
  DatagraphItemThread,
  Identifier,
  Report,
} from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { Badge } from "@/components/ui/badge";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

function getItemPath(
  reportID: Identifier,
  item: DatagraphItem,
): string | undefined {
  switch (item.kind) {
    case DatagraphItemKind.post:
      return `/t/locate/${item.ref.id}?ctx-report-id=${reportID}`;
    case DatagraphItemKind.thread:
      return `/t/${item.ref.slug}?ctx-report-id=${reportID}`;
    case DatagraphItemKind.reply:
      return `/t/locate/${item.ref.id}?ctx-report-id=${reportID}`;
    case DatagraphItemKind.node:
      return `/l/${item.ref.slug}?ctx-report-id=${reportID}`;
    case DatagraphItemKind.profile:
      return `/m/${item.ref.handle}?ctx-report-id=${reportID}`;
    default:
      return undefined;
  }
}

type Props = {
  report: Report;
};

export function ReportCardContent({ report }: Props) {
  const item = report.item as DatagraphItem | undefined;

  if (!item) {
    return (
      <Box p="3" borderRadius="md" bg="bg.muted">
        <styled.p color="fg.muted" fontStyle="italic">
          Content no longer available
        </styled.p>
      </Box>
    );
  }

  switch (item.kind) {
    case DatagraphItemKind.post:
    case DatagraphItemKind.thread:
    case DatagraphItemKind.reply:
      return <PostContent item={item} reportID={report.id} />;
    case DatagraphItemKind.node:
      return <NodeContent item={item} reportID={report.id} />;
    case DatagraphItemKind.profile:
      return <ProfileContent item={item} reportID={report.id} />;
    default:
      return (
        <Box p="3" borderRadius="md" bg="bg.muted">
          <styled.p color="fg.muted">Unknown content type</styled.p>
        </Box>
      );
  }
}

type PostContentProps = {
  reportID: Identifier;
  item: DatagraphItemPost | DatagraphItemThread | DatagraphItemReply;
};

function PostContent({ reportID, item }: PostContentProps) {
  const { ref } = item;
  const path = getItemPath(reportID, item);

  return (
    <Box p="3" borderRadius="md" bg="bg.muted">
      <LStack gap="2">
        <WStack>
          {path ? (
            <Link href={path}>
              <styled.h3 fontWeight="medium" fontSize="md" lineClamp={1}>
                {ref.title || "(Untitled post)"}
              </styled.h3>
            </Link>
          ) : (
            <styled.h3 fontWeight="medium" fontSize="md" lineClamp={2}>
              {ref.title || "(Untitled post)"}
            </styled.h3>
          )}
          <Timestamp created={ref.createdAt} color="fg.subtle" large />
        </WStack>

        <WStack gap="2" fontSize="sm" color="fg.subtle" flexWrap="wrap">
          <MemberBadge profile={ref.author} size="sm" name="handle" />

          {ref.deletedAt !== undefined && <ContentDeletedBadge />}
        </WStack>
      </LStack>
    </Box>
  );
}

type NodeContentProps = {
  reportID: Identifier;
  item: DatagraphItemNode;
};

function NodeContent({ reportID, item }: NodeContentProps) {
  const { ref } = item;
  const path = getItemPath(reportID, item);

  return (
    <Box p="3" borderRadius="md" bg="bg.muted">
      <LStack gap="2">
        <WStack>
          {path ? (
            <Link href={path}>
              <styled.h3 fontWeight="medium" fontSize="md" lineClamp={1}>
                {ref.name}
              </styled.h3>
            </Link>
          ) : (
            <styled.h3 fontWeight="medium" fontSize="md" lineClamp={2}>
              {ref.name}
            </styled.h3>
          )}
          <Timestamp created={ref.createdAt} color="fg.subtle" large />
        </WStack>

        <LStack>
          <div>
            <styled.p fontSize="sm" color="fg.subtle" lineClamp={3}>
              {ref.description}
            </styled.p>
          </div>

          <WStack gap="2" fontSize="sm" color="fg.subtle" flexWrap="wrap">
            <MemberBadge profile={ref.owner} size="sm" name="handle" />

            {ref.deletedAt !== undefined && <ContentDeletedBadge />}
          </WStack>
        </LStack>
      </LStack>
    </Box>
  );
}

type ProfileContentProps = {
  reportID: Identifier;
  item: DatagraphItemProfile;
};

function ProfileContent({ reportID, item }: ProfileContentProps) {
  const { ref } = item;
  const path = getItemPath(reportID, item);

  return (
    <Box p="3" borderRadius="md" bg="bg.muted">
      <LStack gap="2">
        <WStack>
          {path ? (
            <Link href={path}>
              <styled.h3 fontWeight="medium" fontSize="md" lineClamp={1}>
                {ref.name}
              </styled.h3>
            </Link>
          ) : (
            <styled.h3 fontWeight="medium" fontSize="md">
              {ref.name}
            </styled.h3>
          )}
          <Timestamp created={ref.createdAt} color="fg.subtle" large />
        </WStack>

        <WStack gap="2" fontSize="sm" color="fg.subtle" flexWrap="wrap">
          <MemberBadge profile={ref} size="sm" name="handle" />

          {ref.deletedAt !== undefined && <ContentDeletedBadge />}
        </WStack>
      </LStack>
    </Box>
  );
}

function ContentDeletedBadge() {
  return (
    <Badge
      bg="bg.destructive"
      color="fg.destructive"
      borderColor="border.destructive"
    >
      Content deleted
    </Badge>
  );
}
