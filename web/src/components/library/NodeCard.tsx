import { Node } from "src/api/openapi-schema";
import { Timestamp } from "src/components/site/Timestamp";

import { Card } from "@/components/ui/rich-card";
import { LibraryPath, joinLibraryPath } from "@/screens/library/library-path";
import { HStack } from "@/styled-system/jsx";
import { RichCardVariantProps } from "@/styled-system/recipes";
import { getAssetURL } from "@/utils/asset";

import { LibraryBadge } from "./LibraryBadge";

export type NodeCardContext = "library" | "generic";

export type Props = {
  node: Node;
  libraryPath: LibraryPath;
  context: NodeCardContext;
} & RichCardVariantProps;

export function NodeCard({ node, libraryPath, context, ...rest }: Props) {
  const slug = joinLibraryPath(libraryPath, node.slug);
  const asset = node.assets?.[0];
  const url = `/l/${slug}`;

  return (
    <Card
      id={node.id}
      title={node.name}
      text={node.description}
      url={url}
      image={getAssetURL(asset?.path)}
      {...rest}
    >
      {context === "generic" ? (
        <HStack color="fg.muted">
          <LibraryBadge />

          <Timestamp created={node.createdAt} href={url} />
        </HStack>
      ) : (
        <Timestamp created={node.createdAt} large href={url} />
      )}
    </Card>
  );
}
