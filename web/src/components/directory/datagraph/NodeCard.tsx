import { Node } from "src/api/openapi-schema";
import { Timestamp } from "src/components/site/Timestamp";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/directory-path";

import { Card } from "@/components/ui/rich-card";
import { HStack } from "@/styled-system/jsx";
import { RichCardVariantProps } from "@/styled-system/recipes";
import { getAssetURL } from "@/utils/asset";

import { DirectoryBadge } from "../DirectoryBadge";

export type NodeCardContext = "directory" | "generic";

export type Props = {
  node: Node;
  directoryPath: DirectoryPath;
  context: NodeCardContext;
} & RichCardVariantProps;

export function NodeCard({ node, directoryPath, context, ...rest }: Props) {
  const slug = joinDirectoryPath(directoryPath, node.slug);
  const asset = node.assets?.[0];
  const url = `/directory/${slug}`;

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
          <DirectoryBadge />

          <Timestamp
            created={node.createdAt}
            updated={node.updatedAt}
            href={url}
          />
        </HStack>
      ) : (
        <Timestamp
          created={node.createdAt}
          updated={node.updatedAt}
          large
          href={url}
        />
      )}
    </Card>
  );
}
