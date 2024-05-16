"use client";

import { DatagraphBulkNodeScreen } from "../DatagraphBulkNodeScreen/DatagraphBulkNodeScreen";

import { Props, useNodeCreateManyScreen } from "./useNodeCreateManyScreen";

export function NodeCreateManyScreen(props: Props) {
  const {
    handlers: { handleCreate },
    parent,
  } = useNodeCreateManyScreen(props);

  return (
    <DatagraphBulkNodeScreen
      node={parent}
      onCreateNodeFromLink={handleCreate}
    />
  );
}
