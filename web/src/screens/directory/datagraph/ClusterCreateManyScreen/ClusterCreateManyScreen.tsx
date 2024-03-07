"use client";

import { DatagraphBulkNodeScreen } from "../DatagraphBulkNodeScreen/DatagraphBulkNodeScreen";

import {
  Props,
  useClusterCreateManyScreen,
} from "./useClusterCreateManyScreen";

export function ClusterCreateManyScreen(props: Props) {
  const {
    handlers: { handleCreate },
    parent,
  } = useClusterCreateManyScreen(props);

  return (
    <DatagraphBulkNodeScreen
      node={parent && { type: "cluster", ...parent }}
      onCreateNodeFromLink={handleCreate}
    />
  );
}
