import { Cluster } from "src/api/openapi/schemas";
import { CardGrid, CardRows } from "src/theme/components/Card";

import { ClusterCard } from "./ClusterCard";

type Props = {
  directoryPath: string[];
  clusters: Cluster[];
};

export function ClusterCardRows({ directoryPath, clusters }: Props) {
  return (
    <CardRows>
      {clusters.map((c) => (
        <ClusterCard
          key={c.id}
          shape="row"
          directoryPath={directoryPath}
          cluster={c}
        />
      ))}
    </CardRows>
  );
}

export function ClusterCardGrid({ directoryPath, clusters }: Props) {
  return (
    <CardGrid>
      {clusters.map((c) => (
        <ClusterCard
          key={c.id}
          shape="box"
          directoryPath={directoryPath}
          cluster={c}
        />
      ))}
    </CardGrid>
  );
}
