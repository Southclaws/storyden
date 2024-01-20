"use client";

import { DatagraphNodeScreen } from "../DatagraphNodeScreen/DatagraphNodeScreen";

import { Props, useClusterCreateScreen } from "./useClusterCreateScreen";

export function ClusterCreateScreen(props: Props) {
  const {
    handlers: { handleCreate },
    initial,
  } = useClusterCreateScreen(props);

  return (
    <DatagraphNodeScreen
      node={{ type: "cluster", ...initial }}
      initialEditingState={true}
      onSave={handleCreate}
    />
  );
}
