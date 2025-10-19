import Link from "next/link";

import {
  DatagraphItem,
  DatagraphItemKind,
  DatagraphItemNode,
  DatagraphItemPost,
  DatagraphItemProfile,
  DatagraphItemReply,
  DatagraphItemThread,
  Report,
} from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

function getItemPath(item: DatagraphItem): string | undefined {
  switch (item.kind) {
    case DatagraphItemKind.post:
    case DatagraphItemKind.thread:
      return `/t/${item.ref.slug}`;
    case DatagraphItemKind.reply:
      return `/t/${item.ref.root_slug}#${item.ref.id}`;
    case DatagraphItemKind.node:
      return `/l/${item.ref.slug}`;
    case DatagraphItemKind.profile:
      return `/m/${item.ref.handle}`;
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
      return <PostContent item={item} />;
    case DatagraphItemKind.node:
      return <NodeContent item={item} />;
    case DatagraphItemKind.profile:
      return <ProfileContent item={item} />;
    default:
      return (
        <Box p="3" borderRadius="md" bg="bg.muted">
          <styled.p color="fg.muted">Unknown content type</styled.p>
        </Box>
      );
  }
}

type PostContentProps = {
  item: DatagraphItemPost | DatagraphItemThread | DatagraphItemReply;
};

function PostContent({ item }: PostContentProps) {
  const { ref } = item;
  const path = getItemPath(item);

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

        <HStack gap="2" fontSize="sm" color="fg.subtle" flexWrap="wrap">
          <MemberBadge profile={ref.author} size="sm" name="handle" />
        </HStack>
      </LStack>
    </Box>
  );
}

type NodeContentProps = {
  item: DatagraphItemNode;
};

function NodeContent({ item }: NodeContentProps) {
  const { ref } = item;
  const path = getItemPath(item);

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
          <MemberBadge profile={ref.owner} size="sm" name="handle" />
        </LStack>
      </LStack>
    </Box>
  );
}

type ProfileContentProps = {
  item: DatagraphItemProfile;
};

function ProfileContent({ item }: ProfileContentProps) {
  const { ref } = item;
  const path = getItemPath(item);

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

        <MemberBadge profile={ref} size="sm" name="handle" />
      </LStack>
    </Box>
  );
}
