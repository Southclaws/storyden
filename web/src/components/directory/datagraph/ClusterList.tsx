import { Cluster } from "src/api/openapi/schemas";

import { VStack } from "@/styled-system/jsx";

import { ClusterCard } from "./ClusterCard/ClusterCard";

type Props = {
  directoryPath: string[];
  clusters: Cluster[];
};

export function ClusterList(props: Props) {
  return (
    <VStack w="full">
      {props.clusters.map((cluster) => (
        <ClusterCard
          key={cluster.id}
          directoryPath={props.directoryPath}
          cluster={cluster}
        />
      ))}
    </VStack>
  );
}
