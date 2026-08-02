"use client";

import {
  NodeVersion,
  NodeVersionStatus,
  NodeWithChildren,
} from "@/api/openapi-schema";
import { Breadcrumbs } from "@/components/library/Breadcrumbs";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { CardBox } from "@/components/ui/card-box";
import { LinkButton } from "@/components/ui/link-button";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

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
        <LinkButton href={pageHref} variant="subtle" flexShrink="0">
          View page
        </LinkButton>
      </WStack>

      <WStack alignItems="start" gap="3">
        <LStack gap="1">
          <PageHeading>{node.name}</PageHeading>
          <Text variant="supporting">Version history</Text>
        </LStack>
      </WStack>

      {versions.length === 0 ? (
        <Text variant="supporting">No versions or drafts yet.</Text>
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
          <Text as="span" variant="metadata">
            <Timestamp created={version.updated_at} /> ago
          </Text>
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

          <LinkButton href={versionUrl} variant="subtle">
            {buttonLabel}
          </LinkButton>
        </WStack>
      </LStack>
    </CardBox>
  );
}
