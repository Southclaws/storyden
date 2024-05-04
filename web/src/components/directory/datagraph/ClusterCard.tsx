import { Cluster } from "src/api/openapi/schemas";
import { Timestamp } from "src/components/site/Timestamp";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/directory-path";
import { Card } from "src/theme/components/Card";

import { DirectoryBadge } from "../DirectoryBadge";

import { HStack } from "@/styled-system/jsx";
import { CardVariantProps } from "@/styled-system/recipes";

export type ClusterCardContext = "directory" | "generic";

export type Props = {
  cluster: Cluster;
  directoryPath: DirectoryPath;
  context: ClusterCardContext;
} & CardVariantProps;

export function ClusterCard({
  cluster,
  directoryPath,
  context,
  ...rest
}: Props) {
  const slug = joinDirectoryPath(directoryPath, cluster.slug);
  const asset = cluster.assets?.[0];
  const url = `/directory/${slug}`;

  return (
    <Card
      id={cluster.id}
      title={cluster.name}
      text={cluster.description}
      url={url}
      image={asset?.url}
      {...rest}
    >
      {context === "generic" ? (
        <HStack color="fg.muted">
          <DirectoryBadge />

          <Timestamp
            created={cluster.createdAt}
            updated={cluster.updatedAt}
            href={url}
          />
        </HStack>
      ) : (
        <Timestamp
          created={cluster.createdAt}
          updated={cluster.updatedAt}
          large
          href={url}
        />
      )}
    </Card>
  );
}
