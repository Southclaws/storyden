import { NodeWithChildren } from "@/api/openapi-schema";
import {
  LibraryBadge,
  LibraryPageBadge,
} from "@/components/library/LibraryBadge";
import { LibraryPageMenu } from "@/components/library/LibraryPageMenu/LibraryPageMenu";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { Card, CardRows } from "@/components/ui/rich-card";
import { HStack, WStack } from "@/styled-system/jsx";
import { getAssetURL } from "@/utils/asset";

export function QueueNodeList({ nodes }: { nodes: NodeWithChildren[] }) {
  if (nodes.length === 0) {
    return <p>Submissions appear here.</p>;
  }

  return (
    <CardRows>
      {nodes.map((node) => (
        <QueueNodeListItem key={node.id} node={node} />
      ))}
    </CardRows>
  );
}

export function QueueNodeListItem({ node }: { node: NodeWithChildren }) {
  const url = node.parent
    ? `/l/${node.parent.slug}/${node.slug}`
    : `/l/${node.slug}`;

  return (
    <Card
      key={node.id}
      id={node.id}
      shape="responsive"
      url={url}
      image={getAssetURL(node.primary_image?.path)}
      title={node.name}
      text={node.description}
      controls={
        <WStack>
          <HStack gap="2">
            <MemberBadge profile={node.owner} size="sm" />

            {node.parent ? (
              <LibraryPageBadge {...node.parent} />
            ) : (
              <LibraryBadge />
            )}

            <Timestamp created={node.createdAt} large />
          </HStack>

          <LibraryPageMenu node={node} />
        </WStack>
      }
    />
  );
}
