import { Node } from "src/api/openapi-schema";
import { Timestamp } from "src/components/site/Timestamp";

import { Card } from "@/components/ui/rich-card";
import { LibraryPath, joinLibraryPath } from "@/screens/library/library-path";
import { HStack } from "@/styled-system/jsx";
import { RichCardVariantProps } from "@/styled-system/recipes";
import { getAssetURL } from "@/utils/asset";

import { LibraryPageMenu } from "./LibraryPageMenu/LibraryPageMenu";

export type NodeCardContext = "library" | "generic";

export type Props = {
  node: Node;
  libraryPath: LibraryPath;
  context: NodeCardContext;
} & RichCardVariantProps;

export function NodeCard({ node, libraryPath, context, ...rest }: Props) {
  const slug = joinLibraryPath(libraryPath, node.slug);
  const url = `/l/${slug}`;
  const image = getAssetURL(node.primary_image?.path);

  return (
    <Card
      id={node.id}
      title={node.name}
      text={node.description}
      url={url}
      image={image}
      // menu={<LibraryPageMenu node={node} />}
      controls={
        <HStack w="full" justify="space-between">
          <Timestamp created={node.createdAt} href={url} large />

          <LibraryPageMenu node={node} />
        </HStack>
      }
      {...rest}
    ></Card>
  );
}
