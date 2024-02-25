import { Cluster } from "src/api/openapi/schemas";
import { Timestamp } from "src/components/site/Timestamp";
import {
  DirectoryPath,
  joinDirectoryPath,
} from "src/screens/directory/datagraph/useDirectoryPath";
import { Card } from "src/theme/components/Card";

import { DirectoryBadge } from "../DirectoryBadge";

import { HStack } from "@/styled-system/jsx";
import { CardVariantProps } from "@/styled-system/recipes";

export type Props = {
  cluster: Cluster;
  directoryPath: DirectoryPath;
} & CardVariantProps;

export function ClusterCard({ cluster, directoryPath, ...rest }: Props) {
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
      <HStack color="fg.muted">
        <DirectoryBadge />

        <Timestamp
          created={cluster.createdAt}
          updated={cluster.updatedAt}
          href={url}
        />
      </HStack>
    </Card>
  );
}
