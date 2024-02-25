import { Cluster } from "src/api/openapi/schemas";

import { Grid, LStack } from "@/styled-system/jsx";

import { ClusterCard } from "./ClusterCard";

type Props = {
  directoryPath: string[];
  clusters: Cluster[];
};

export function ClusterCardRows({ directoryPath, clusters }: Props) {
  return (
    <LStack>
      {clusters.map((c) => (
        <ClusterCard
          key={c.id}
          shape="row"
          directoryPath={directoryPath}
          cluster={c}
        />
      ))}
    </LStack>
  );
}

export function ClusterCardGrid({ directoryPath, clusters }: Props) {
  return (
    <Grid
      w="full"
      gridTemplateColumns={{
        base: "2",
        sm: "4",
        lg: "6",
      }}
    >
      {clusters.map((c) => (
        <ClusterCard
          key={c.id}
          shape="box"
          directoryPath={directoryPath}
          cluster={c}
        />
      ))}
    </Grid>
  );
}
