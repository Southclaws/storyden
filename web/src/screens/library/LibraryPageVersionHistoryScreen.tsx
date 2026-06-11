"use client";

import {
  NodeVersion,
  NodeVersionStatus,
  NodeWithChildren,
} from "@/api/openapi-schema";
import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import { PageVersionStatusBadge } from "./LibraryPageScreen/PageVersionStatusBadge";

type Props = {
  node: NodeWithChildren;
  versions: NodeVersion[];
  libraryPath: string[];
};

export function LibraryPageVersionHistoryScreen({
  node,
  versions,
  libraryPath,
}: Props) {
  const pageHref = `/l/${libraryPath.join("/")}`;

  return (
    <LStack gap="3">
      <WStack>
        <Breadcrumbs
          libraryPath={libraryPath}
          visibility={node.visibility}
          create="show"
        />
        <LinkButton href={pageHref} size="xs" variant="subtle" flexShrink="0">
          View page
        </LinkButton>
      </WStack>

      <WStack alignItems="start" gap="3">
        <LStack gap="1">
          <Heading fontSize="heading.2" fontWeight="bold">
            {node.name}
          </Heading>
          <styled.p color="fg.muted" fontSize="sm">
            Version history
          </styled.p>
        </LStack>
      </WStack>

      {versions.length === 0 ? (
        <styled.p color="fg.muted" fontSize="sm">
          No versions or drafts yet.
        </styled.p>
      ) : (
        <LStack gap="2">
          {versions.map((version) => (
            <VersionHistoryItem
              key={version.id}
              version={version}
              pageHref={pageHref}
            />
          ))}
        </LStack>
      )}
    </LStack>
  );
}

function VersionHistoryItem({
  version,
  pageHref,
}: {
  version: NodeVersion;
  pageHref: string;
}) {
  const versionUrl = `${pageHref}?version=${version.id}`;
  const isApplied = version.status === NodeVersionStatus.applied;
  const buttonLabel = isApplied ? "View changes" : "Review";

  return (
    <CardBox>
      <LStack>
        <WStack>
          <PageVersionStatusBadge status={version.status} />
          <styled.span color="fg.muted" fontSize="xs">
            <Timestamp created={version.updated_at} /> ago
          </styled.span>
        </WStack>

        <WStack alignItems="end">
          <HStack>
            <MemberBadge
              profile={version.author}
              size="xs"
              name="handle"
              avatar="visible"
            />
          </HStack>

          <LinkButton href={versionUrl} size="xs" variant="subtle">
            {buttonLabel}
          </LinkButton>
        </WStack>
      </LStack>
    </CardBox>
  );
}
