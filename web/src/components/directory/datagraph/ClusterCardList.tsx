import { Cluster } from "src/api/openapi/schemas";
import { CardGrid, CardRows } from "src/theme/components/Card";

import { CardVariantProps } from "@/styled-system/recipes";

import { ClusterCard, ClusterCardContext } from "./ClusterCard";

type Props = {
  directoryPath: string[];
  clusters: Cluster[];
  size?: CardVariantProps["size"];
  context: ClusterCardContext;
};

export function ClusterCardRows({
  directoryPath,
  clusters,
  size,
  context,
}: Props) {
  return (
    <CardRows>
      {clusters.map((c) => (
        <ClusterCard
          key={c.id}
          shape="row"
          size={size}
          context={context}
          directoryPath={directoryPath}
          cluster={c}
        />
      ))}
    </CardRows>
  );
}

export function ClusterCardGrid({ directoryPath, clusters, context }: Props) {
  return (
    <CardGrid>
      {clusters.map((c) => (
        <ClusterCard
          key={c.id}
          shape="box"
          context={context}
          directoryPath={directoryPath}
          cluster={c}
        />
      ))}
    </CardGrid>
  );
}
