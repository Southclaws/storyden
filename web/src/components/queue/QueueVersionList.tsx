import { NodeDraft } from "@/api/openapi-schema";
import { LibraryPageBadge, LibraryBadge } from "@/components/library/LibraryBadge";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { DraftIcon } from "@/components/ui/icons/Draft";
import { LinkButton } from "@/components/ui/link-button";
import { Card, CardRows } from "@/components/ui/rich-card";
import { HStack, WStack } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

export function QueueVersionList({ drafts }: { drafts: NodeDraft[] }) {
  if (drafts.length === 0) {
    return <p>Edit proposals appear here.</p>;
  }

  return (
    <CardRows>
      {drafts.map((draft) => (
        <QueueVersionListItem key={draft.id} draft={draft} />
      ))}
    </CardRows>
  );
}

function QueueVersionListItem({ draft }: { draft: NodeDraft }) {
  const node = draft.node;
  const url = node.parent
    ? `/l/${node.parent.slug}/${node.slug}`
    : `/l/${node.slug}`;

  return (
    <Card
      key={draft.id}
      id={draft.id}
      shape="responsive"
      url={url}
      image={getAssetURL(node.primary_image?.path)}
      title={node.name}
      text="Draft edit proposal"
      titleIcon={<DraftIcon width="4" height="4" />}
      controls={
        <WStack>
          <HStack gap="2">
            <MemberBadge profile={draft.author} size="sm" />
            <Timestamp created={draft.updated_at} large />

            {node.parent ? (
              <LibraryPageBadge {...node.parent} />
            ) : (
              <LibraryBadge />
            )}
          </HStack>

          <LinkButton
            href={`${url}?version=${draft.id}`}
            size="xs"
            variant="subtle"
          >
            Review
          </LinkButton>
        </WStack>
      }
    />
  );
}
